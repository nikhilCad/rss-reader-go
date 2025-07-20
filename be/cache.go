package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

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

// FetchAndCacheFeed fetches a feed, parses it, and stores articles in the DB.
func FetchAndCacheFeed(feedURL string) error {
	posts, err := fetchAndParseRSS(feedURL)
	if err != nil {
		return err
	}
	for _, post := range posts {
		err := upsertArticle(post, feedURL)
		if err != nil {
			log.Printf("Failed to upsert article %s: %v", post.Link, err)
		}
	}
	return nil
}

// RefreshAllFeeds fetches and caches all feeds in the DB.
func RefreshAllFeeds() error {
	feeds, err := db.ListFeeds()
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		log.Println("Refreshing feed:", feed.URL)
		if err := FetchAndCacheFeed(feed.URL); err != nil {
			log.Println("Failed to refresh feed:", feed.URL, err)
		}
	}
	return nil
}

// GetCachedArticles returns all cached articles from the DB.
func GetCachedArticles() ([]Post, error) {
	rows, err := db.(*sqliteDB).db.Query(`
		SELECT title, link, description, content, source, pubdate FROM articles
		ORDER BY pubdate DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.Title, &p.Link, &p.Description, &p.Content, &p.Source, &p.PubDate)
		if err != nil {
			continue
		}
		articles = append(articles, p)
	}
	return articles, nil
}

// upsertArticle inserts or updates an article in the DB.
func upsertArticle(p Post, source string) error {
	_, err := db.(*sqliteDB).db.Exec(`
		INSERT INTO articles (title, link, description, content, source, pubdate, fetched_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(link) DO UPDATE SET
			title=excluded.title,
			description=excluded.description,
			content=excluded.content,
			source=excluded.source,
			pubdate=excluded.pubdate,
			fetched_at=excluded.fetched_at
	`, p.Title, p.Link, p.Description, p.Content, source, p.PubDate, time.Now().Format(time.RFC3339))
	return err
}
