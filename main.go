package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	traceURL      = "https://one.one.one.one/cdn-cgi/trace"
	ipHistoryPath = "ip_history.txt"
	checkInterval = 30 * time.Second
	maxRetries    = 3
	retryDelay    = 10 * time.Second
)

func getPublicIP() (string, error) {
	var ip string
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ip, err = fetchIP()
		if err == nil {
			return ip, nil
		}
		fmt.Printf("Attempt %d: Error fetching IP: %v\n", attempt, err)
		time.Sleep(retryDelay)
	}
	return "", err
}

func fetchIP() (string, error) {
	resp, err := http.Get(traceURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(body), "\n") {
		if strings.HasPrefix(line, "ip=") {
			return strings.TrimPrefix(line, "ip="), nil
		}
	}
	return "", fmt.Errorf("IP address not found")
}

func readHistory() ([]string, error) {
	file, err := os.Open(ipHistoryPath)
	if os.IsNotExist(err) {
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var history []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		history = append(history, scanner.Text())
	}
	return history, scanner.Err()
}

func writeHistory(history []string) error {
	file, err := os.Create(ipHistoryPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, entry := range history {
		_, err := writer.WriteString(entry + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func trackIP(ipChan <-chan string) {
	var lastLoggedIP string
	for ip := range ipChan {
		history, err := readHistory()
		if err != nil {
			fmt.Println("Error reading IP history:", err)
			continue
		}

		entry := fmt.Sprintf("%s - %s", time.Now().Format("2006-01-02 15:04:05 MST"), ip)
		if len(history) == 0 || !strings.Contains(history[0], ip) {
			history = append([]string{entry}, history...)
			if err := writeHistory(history); err != nil {
				fmt.Println("Error writing IP history:", err)
			}
			fmt.Printf("📄 New IP logged: %s\n", entry)
			lastLoggedIP = ip
		} else {
			if lastLoggedIP != ip {
				fmt.Println("🤷 IP is already up to date")
				lastLoggedIP = ip
			}
		}
	}
}

func serveWeb() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

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
	ipChan := make(chan string)

	go func() {
		for {
			ip, err := getPublicIP()
			if err != nil {
				fmt.Println("Error fetching public IP:", err)
			} else {
				ipChan <- ip
			}
			time.Sleep(checkInterval)
		}
	}()

	go trackIP(ipChan)
	serveWeb()
}
