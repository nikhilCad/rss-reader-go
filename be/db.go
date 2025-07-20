package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Add Feed struct
// Feed represents an RSS feed source
type Feed struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	FeedName string `json:"feed_name"`
}

// DB is an interface for database operations.
type DB interface {
	InsertPost(post Post) error
	GetAllPosts() ([]Post, error)
	Close() error
	// Feed management
	AddFeed(url string, name string) error
	RemoveFeed(url string) error
	ListFeeds() ([]Feed, error)
	// Add to DB interface
	MarkRead(link string) error
	MarkUnread(link string) error
	ListRead() ([]string, error)
}

// sqliteDB implements DB using SQLite.
type sqliteDB struct {
	db *sql.DB
}

// NewSQLiteDB creates a new SQLite database and returns a DB interface.
func NewSQLiteDB(path string) (DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	createPosts := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		link TEXT UNIQUE,
		description TEXT
	);`
	_, err = db.Exec(createPosts)
	if err != nil {
		return nil, err
	}
	createFeeds := `
	CREATE TABLE IF NOT EXISTS feeds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT UNIQUE,
		feed_name TEXT
	);`
	_, err = db.Exec(createFeeds)
	if err != nil {
		return nil, err
	}
	createRead := `
	CREATE TABLE IF NOT EXISTS read_articles (
		link TEXT PRIMARY KEY
	);`
	_, err = db.Exec(createRead)
	if err != nil {
		return nil, err
	}

	const createArticlesTableSQL = `
		CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			link TEXT UNIQUE,
			description TEXT,
			content TEXT,
			source TEXT,
			pubdate TEXT,
			fetched_at TEXT,
			enclosure_url TEXT,     
			enclosure_type TEXT,    
			enclosure_length TEXT    
		);
`
	_, err = db.Exec(createArticlesTableSQL)
	if err != nil {
		return nil, err
	}

	return &sqliteDB{db: db}, nil
}

func (s *sqliteDB) InsertPost(post Post) error {
	_, err := s.db.Exec("INSERT OR IGNORE INTO posts (title, link, description) VALUES (?, ?, ?)", post.Title, post.Link, post.Description)
	return err
}

func (s *sqliteDB) GetAllPosts() ([]Post, error) {
	rows, err := s.db.Query("SELECT title, link, description FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Title, &post.Link, &post.Description); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// Feed management methods
func (s *sqliteDB) AddFeed(url, name string) error {
	_, err := s.db.Exec("INSERT OR IGNORE INTO feeds (url, feed_name) VALUES (?, ?)", url, name)
	return err
}

func (s *sqliteDB) RemoveFeed(url string) error {
	_, err := s.db.Exec("DELETE FROM feeds WHERE url = ?", url)
	if err != nil {
		return err
	}
	// Also delete articles from this feed
	_, err = s.db.Exec("DELETE FROM articles WHERE source = ?", url)
	return err
}

func (s *sqliteDB) ListFeeds() ([]Feed, error) {
	rows, err := s.db.Query("SELECT id, url, feed_name FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var feeds []Feed
	for rows.Next() {
		var f Feed
		if err := rows.Scan(&f.ID, &f.URL, &f.FeedName); err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}
	return feeds, nil
}

func (s *sqliteDB) MarkRead(link string) error {
	_, err := s.db.Exec("INSERT OR IGNORE INTO read_articles (link) VALUES (?)", link)
	return err
}
func (s *sqliteDB) MarkUnread(link string) error {
	_, err := s.db.Exec("DELETE FROM read_articles WHERE link = ?", link)
	return err
}
func (s *sqliteDB) ListRead() ([]string, error) {
	rows, err := s.db.Query("SELECT link FROM read_articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var links []string
	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (s *sqliteDB) Close() error {
	return s.db.Close()
}
