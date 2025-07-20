package main

import (
	"github.com/mmcdole/gofeed"
)

// Post represents a blog post or article.
type Post struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Source      string `json:"source"`
	PubDate     string `json:"pubdate"`
}

func ParseFeed(data []byte) ([]Post, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(string(data))
	if err != nil {
		return nil, err
	}
	posts := make([]Post, 0, len(feed.Items))
	for _, item := range feed.Items {
		posts = append(posts, Post{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Source:      feed.Title, // or set as needed
		})
	}
	return posts, nil
}
