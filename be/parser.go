package main

import (
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
)

// Enclosure represents a media enclosure (e.g., podcast audio).
type Enclosure struct {
	URL    string `json:"url"`
	Type   string `json:"type,omitempty"`
	Length string `json:"length,omitempty"`
}

// Post represents a blog post or article.
type Post struct {
	Title       string     `json:"title"`
	Link        string     `json:"link"`
	Content     string     `json:"content"`
	Description string     `json:"description"`
	Source      string     `json:"source"`
	PubDate     string     `json:"pubdate"`
	Enclosure   *Enclosure `json:"enclosure,omitempty"`
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
		var enclosure *Enclosure
		if len(item.Enclosures) > 0 {
			enc := item.Enclosures[0]
			enclosure = &Enclosure{
				URL:    enc.URL,
				Type:   enc.Type,
				Length: enc.Length,
			}
		}
		posts = append(posts, Post{
			Title:       item.Title,
			Link:        item.Link,
			Content:     item.Content,
			Description: item.Description,
			Source:      feed.Title, // or set as needed
			Enclosure:   enclosure,
		})
	}
	return posts, nil
}
