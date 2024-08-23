package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	ipv4URL       = "https://4.ident.me/"
	ipv6URL       = "https://6.ident.me/"
	ipHistoryPath = "history.json"
	checkInterval = 30 * time.Second
	maxRetries    = 3
	retryDelay    = 10 * time.Second
)

type IPRecord struct {
	Timestamp string `json:"timestamp"`
	IPv4      string `json:"ipv4"`
	IPv6      string `json:"ipv6"`
}

func getPublicIPs() (string, string, error) {
	var ipv4, ipv6 string
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ipv4, ipv6, err = fetchIPs()
		if err == nil {
			return ipv4, ipv6, nil
		}
		fmt.Printf("Attempt %d: Error fetching IPs: %v\n", attempt, err)
		time.Sleep(retryDelay)
	}
	return "", "", err
}

func fetchIPs() (string, string, error) {
	ipv4, err := fetchIP(ipv4URL)
	if err != nil {
		return "", "", err
	}

	ipv6, err := fetchIP(ipv6URL)
	if err != nil {
		fmt.Println("Warning: Could not fetch IPv6 address:", err)
	}

	return ipv4, ipv6, nil
}

func fetchIP(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

func readHistory() ([]IPRecord, error) {
	file, err := os.Open(ipHistoryPath)
	if os.IsNotExist(err) {
		return []IPRecord{}, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Check if the file is empty
	if fileInfo.Size() == 0 {
		return []IPRecord{}, nil
	}

	var history []IPRecord
	err = json.NewDecoder(file).Decode(&history)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func writeHistory(history []IPRecord) error {
	file, err := os.Create(ipHistoryPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print with indentation
	err = encoder.Encode(history)
	if err != nil {
		return err
	}
	return nil
}

func trackIP(ipChan <-chan IPRecord) {
	var lastLoggedIP IPRecord
	hasNotifiedUpToDate := false

	for ipRecord := range ipChan {
		history, err := readHistory()
		if err != nil {
			fmt.Println("Error reading IP history:", err)
			continue
		}

		if len(history) > 0 {
			lastLoggedIP = history[0]
		}

		// Check if IPs are the same as the last logged one
		if ipRecord.IPv4 == lastLoggedIP.IPv4 && ipRecord.IPv6 == lastLoggedIP.IPv6 {
			if !hasNotifiedUpToDate {
				fmt.Println("🤷 IPs are already up to date")
				hasNotifiedUpToDate = true
			}
			continue
		}

		// Reset the notification flag when IP changes
		hasNotifiedUpToDate = false

		ipRecord.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		history = append([]IPRecord{ipRecord}, history...)
		if err := writeHistory(history); err != nil {
			fmt.Println("Error writing IP history:", err)
			continue
		}

		fmt.Printf("📄 New IP logged: {Timestamp: %s, IPv4: %s, IPv6: %s}\n", ipRecord.Timestamp, ipRecord.IPv4, ipRecord.IPv6)
	}
}


func serveWeb() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("web"))))

	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		history, err := readHistory()
		if err != nil {
			http.Error(w, "Error reading IP history", http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(history)
		if err != nil {
			http.Error(w, "Error converting history to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	ipChan := make(chan IPRecord)

	go func() {
		for {
			ipv4, ipv6, err := getPublicIPs()
			if err != nil {
				fmt.Println("Error fetching public IPs:", err)
			} else {
				ipChan <- IPRecord{IPv4: ipv4, IPv6: ipv6}
			}
			time.Sleep(checkInterval)
		}
	}()

	go trackIP(ipChan)
	serveWeb()
}
