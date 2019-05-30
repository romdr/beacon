package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Send to URL
func sendToURL(hostMetrics *HostMetrics, url string) {
	// To json
	buffer, err := json.Marshal(hostMetrics)
	if err != nil {
		log.Printf("ERROR: Could not send to URL, could not encode to json: %s\n", err)
		return
	}

	// Set a timeout for the request
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	// Send it
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(buffer))
	if err != nil {
		log.Printf("ERROR: Could not send to URL: %s\n", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("ERROR: Heartbeat URL %s responded with: %d\n", url, resp.StatusCode)
		return
	}
}
