package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

// Post represents a blog post or article.
type Post struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

// rssItem matches the structure of an RSS <item>
type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

// rssChannel matches the structure of an RSS <channel>
type rssChannel struct {
	Items []rssItem `xml:"item"`
}

// rssFeed matches the structure of the RSS root
type rssFeed struct {
	Channel rssChannel `xml:"channel"`
}

var db DB // Global database interface

// fetchAndParseRSS fetches the RSS feed and returns a slice of Post
func fetchAndParseRSS(url string) ([]Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed rssFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, err
	}

	posts := make([]Post, 0, len(feed.Channel.Items))
	for _, item := range feed.Channel.Items {
		posts = append(posts, Post{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
		})
	}
	return posts, nil
}

// postsHandler returns all posts from the database as JSON
func postsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := db.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to fetch posts from database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func main() {
	var err error
	db, err = NewSQLiteDB("./posts.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Fetch and store posts from RSS feed at startup
	const rssURL = "https://xkcd.com/rss.xml"
	posts, err := fetchAndParseRSS(rssURL)
	if err != nil {
		panic("Failed to fetch or parse RSS feed: " + err.Error())
	}
	for _, post := range posts {
		_ = db.InsertPost(post) // This will insert duplicates if run multiple times
	}

	http.HandleFunc("/", renderIndex)      // Serve the main page
	http.Handle("/static/", serveStatic()) // Serve static files

	http.HandleFunc("/posts", postsHandler) // API endpoint

	println("Server running at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
