package main

import (
	"time"

	"github.com/go-shiori/go-readability"
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

// ArticleParseResult represents the parsed article data.
type ArticleParseResult struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Byline  string `json:"byline"`
}

// ParseArticleFromURL fetches and parses an article using go-readability.
func ParseArticleFromURL(urlStr string) (ArticleParseResult, error) {
	article, err := readability.FromURL(urlStr, 30*time.Second)
	if err != nil {
		return ArticleParseResult{}, err
	}

	return ArticleParseResult{
		Title:   article.Title,
		Content: article.Content,
		Byline:  article.Byline,
	}, nil
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
			Content:     item.Content,
			Description: item.Description,
			Source:      feed.Title, // or set as needed
		})
	}
	return posts, nil
}
