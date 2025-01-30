package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	CreatedAt string `json:"created_At"`
}

func fetch(username string) ([]Event, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close() // Close AFTER reading

	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	
	body, err := io.ReadAll(resp.Body) // ✅ Read before closing
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// // Debugging: Print raw response (REMOVE this after testing)
	// fmt.Println("Raw Response Body:", string(body))

	var events []Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return events, nil
}


func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: github-activity <username>")
		return
	}
	username := os.Args[1]
	events, err := fetch(username)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	if len(events) == 0 {
		fmt.Println("no recent activity found")
		return
	}
	fmt.Println("recent activity for", username, ":")
	for _, event := range events {
		fmt.Println("-", event.Type, "in repository", event.Repo.Name, "at", event.CreatedAt)
	}
}
