package main

import (
	"encoding/json"
	"net/http"
)

// postsHandler returns all posts from all feeds in the database as JSON
func postsHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := GetCachedArticles()
	if err != nil {
		http.Error(w, "Failed to fetch cached articles", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		FromCache bool   `json:"fromCache"`
		Articles  []Post `json:"articles"`
	}{
		FromCache: true, // Always true for now
		Articles:  articles,
	}
	json.NewEncoder(w).Encode(response)
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

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := RefreshAllFeeds()
	if err != nil {
		http.Error(w, "Failed to refresh feeds", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
