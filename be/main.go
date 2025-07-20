package main

import (
	"net/http"
)

var db DB // Global database interface

func main() {
	var err error
	db, err = NewSQLiteDB("./posts.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/posts", postsHandler)
	http.HandleFunc("/feeds", feedsHandler)
	http.HandleFunc("/read", readHandler)
	http.HandleFunc("/unread", readHandler)
	http.HandleFunc("/refresh", refreshHandler)
	http.HandleFunc("/parse-article", ParseArticleHandler)

	// Sample RSS feed
	StartSampleFeeds()
	// Sample RSS end

	println("Server running at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
