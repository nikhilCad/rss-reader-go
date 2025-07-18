package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// DB is an interface for database operations.
type DB interface {
	InsertPost(post Post) error
	GetAllPosts() ([]Post, error)
	Close() error
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
	createTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		link TEXT UNIQUE,
		description TEXT
	);`
	_, err = db.Exec(createTable)
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

func (s *sqliteDB) Close() error {
	return s.db.Close()
}
