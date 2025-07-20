package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

var db DB // Global database interface

// fetchAndParseRSS fetches the feed and parses it using ParseFeed from parser.go
func fetchAndParseRSS(url string) ([]Post, error) {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err := http.Get(url)
		if err != nil {
			lastErr = err
			log.Printf("Fetch attempt %d failed for %s: %v", attempt, url, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			log.Printf("Read attempt %d failed for %s: %v", attempt, url, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		posts, err := ParseFeed(body)
		if err != nil {
			lastErr = err
			log.Printf("Parse attempt %d failed for %s: %v", attempt, url, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		return posts, nil // Success!
	}
	return nil, lastErr
}

// postsHandler returns all posts from all feeds in the database as JSON
func postsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := db.ListFeeds()
	if err != nil {
		http.Error(w, "Failed to fetch feeds", http.StatusInternalServerError)
		return
	}
	var allPosts []Post
	for _, feed := range feeds {
		log.Println("Fetching post for:", feed.URL)
		posts, err := fetchAndParseRSS(feed.URL)
		if err != nil {
			log.Println("Error fetching/parsing feed:", feed.URL, err)
			continue // skip feeds that fail
		}
		// Attach source info for grouping in frontend
		for _, p := range posts {
			p := p
			p.Source = feed.URL
			allPosts = append(allPosts, p)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allPosts)
}

// feedsHandler handles GET, POST, DELETE for RSS feed URLs
func feedsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		feeds, err := db.ListFeeds()
		if err != nil {
			http.Error(w, "Failed to list feeds", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(feeds)
	case http.MethodPost:
		var req struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if err := db.AddFeed(req.URL); err != nil {
			http.Error(w, "Failed to add feed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		var req struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if err := db.RemoveFeed(req.URL); err != nil {
			http.Error(w, "Failed to remove feed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		links, err := db.ListRead()
		if err != nil {
			http.Error(w, "Failed to list read articles", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(links)
	case http.MethodPost:
		var req struct {
			Link string `json:"link"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Link == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if r.URL.Path == "/read" {
			db.MarkRead(req.Link)
		} else {
			db.MarkUnread(req.Link)
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	var err error
	db, err = NewSQLiteDB("./posts.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Remove hardcoded RSS fetch logic

	http.HandleFunc("/posts", postsHandler) // API endpoint
	http.HandleFunc("/feeds", feedsHandler) // Feed management API
	http.HandleFunc("/read", readHandler)
	http.HandleFunc("/unread", readHandler)

	// Sample RSS feed
	StartSampleFeeds()
	// Sample RSS end

	println("Server running at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
