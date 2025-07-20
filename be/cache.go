package main

import (
	"database/sql"
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
// If the DB is nil, it logs a warning and returns an empty slice.
func GetCachedArticles(db interface{}) ([]Post, error) {
	// Defensive: Check if db is nil or not the expected type before querying
	sqliteDB, ok := db.(*sqliteDB)
	if !ok || sqliteDB.db == nil {
		log.Println("WARNING: Database handle is nil in GetCachedArticles, returning empty result")
		return []Post{}, nil // Proceed with empty result, no error
	}

	// Query the articles table, including enclosure fields
	rows, err := sqliteDB.db.Query(`
		SELECT title, link, description, content, source, pubdate, enclosure_url, enclosure_type, enclosure_length
		FROM articles
	`)
	if err != nil {
		log.Printf("DB query error in GetCachedArticles: %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var enclosureURL, enclosureType, enclosureLength sql.NullString

		err := rows.Scan(
			&post.Title,
			&post.Link,
			&post.Description,
			&post.Content,
			&post.Source,
			&post.PubDate,
			&enclosureURL,
			&enclosureType,
			&enclosureLength,
		)
		if err != nil {
			log.Printf("Row scan error in GetCachedArticles: %v", err)
			continue // Skip this row, but keep going
		}

		// Only set Enclosure if URL is present
		if enclosureURL.Valid && enclosureURL.String != "" {
			post.Enclosure = &Enclosure{
				URL:    enclosureURL.String,
				Type:   enclosureType.String,
				Length: enclosureLength.String,
			}
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// upsertArticle inserts or updates an article in the DB.
func upsertArticle(p Post, source string) error {
	_, err := db.(*sqliteDB).db.Exec(`
	INSERT INTO articles (
		title, link, description, content, source, pubdate, fetched_at,
		enclosure_url, enclosure_type, enclosure_length
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(link) DO UPDATE SET
		title=excluded.title,
		description=excluded.description,
		content=excluded.content,
		source=excluded.source,
		pubdate=excluded.pubdate,
		fetched_at=excluded.fetched_at,
		enclosure_url=excluded.enclosure_url,
		enclosure_type=excluded.enclosure_type,
		enclosure_length=excluded.enclosure_length
`,
		p.Title,
		p.Link,
		p.Description,
		p.Content,
		source,
		p.PubDate,
		time.Now().Format(time.RFC3339),
		// Enclosure fields: use empty string if nil
		func() string {
			if p.Enclosure != nil {
				return p.Enclosure.URL
			} else {
				return ""
			}
		}(),
		func() string {
			if p.Enclosure != nil {
				return p.Enclosure.Type
			} else {
				return ""
			}
		}(),
		func() string {
			if p.Enclosure != nil {
				return p.Enclosure.Length
			} else {
				return ""
			}
		}(),
	)
	return err
}
