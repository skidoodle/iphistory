package main

import (
	"context"
	"embed"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

//go:generate go tool templ generate

//go:embed assets/*
var assetsFS embed.FS

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	store, err := NewStore("history.db")
	if err != nil {
		logger.Error("db init failed", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		type provider struct {
			server string
			host   string
			isTXT  bool
		}

		providers := []provider{
			{server: "216.239.32.10:53", host: "o-o.myaddr.l.google.com", isTXT: true}, // ns1.google.com
			{server: "193.108.88.1:53", host: "whoami.akamai.net", isTXT: false},       // ns1-1.akamaitech.net
			{server: "208.67.222.222:53", host: "myip.opendns.com", isTXT: false},      // resolver1.opendns.com
		}

		for {
			var detectedIP string

			for _, p := range providers {
				resolver := &net.Resolver{
					PreferGo: true,
					Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
						d := net.Dialer{Timeout: 5 * time.Second}
						return d.DialContext(ctx, "udp4", p.server)
					},
				}

				var raw string
				if p.isTXT {
					txt, err := resolver.LookupTXT(gCtx, p.host)
					if err == nil && len(txt) > 0 {
						raw = strings.Trim(txt[0], "\"")
					}
				} else {
					ips, err := resolver.LookupHost(gCtx, p.host)
					if err == nil && len(ips) > 0 {
						raw = ips[0]
					}
				}

				if ip := net.ParseIP(raw); ip != nil && ip.To4() != nil {
					detectedIP = ip.String()
					break
				}
			}

			if detectedIP != "" {
				last, _ := store.GetLatest()
				if detectedIP != last {
					if err := store.Insert(detectedIP); err != nil {
						logger.Error("failed to save IP", "err", err)
					} else {
						logger.Info("IP change detected", "ip", detectedIP)
					}
				}
			}

			select {
			case <-gCtx.Done():
				return nil
			case <-ticker.C:
			}
		}
	})

	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.FileServer(http.FS(assetsFS)))
	mux.HandleFunc("GET /{$}", handleList(store, logger))
	mux.HandleFunc("GET /p/{page}", handleList(store, logger))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		logger.Info("server started", "url", "http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		sCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(sCtx)
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("application error", "err", err)
	}
}

func handleList(store *Store, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.PathValue("page"))
		if page < 1 {
			page = 1
		}
		query := r.URL.Query().Get("q")

		records, hasMore, err := store.FetchPage(query, page, 50)
		if err != nil {
			http.Error(w, "Internal Error", 500)
			return
		}

		if r.Header.Get("HX-Request") == "true" && r.Header.Get("HX-Target") == "main-content" {
			_ = MainContent(records, query, page, hasMore).Render(r.Context(), w)
			return
		}
		_ = Page(records, query, page, hasMore).Render(r.Context(), w)
	}
}
