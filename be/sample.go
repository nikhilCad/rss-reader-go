package main

import (
	"log"
	"net/http"
)

// Generate sample RSS XML 1
func sampleXML1() string {
	return `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
  <channel>
    <title>Sample Feed</title>
    <link>http://localhost:8081/sample.xml</link>
    <description>This is a test RSS feed</description>
    <item>
      <title>First Post</title>
      <link>http://localhost:8081/posts/1</link>
      <description>Hello from the sample feed!</description>
    </item>
    <item>
      <title>Second Post</title>
      <link>http://localhost:8081/posts/2</link>
      <description>Another sample post.</description>
    </item>
  </channel>
</rss>`
}

// Generate sample RSS XML 2
func sampleXML2() string {
	return `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
  <channel>
    <title>Sample Feed 2</title>
    <link>http://localhost:8082/sample2.xml</link>
    <description>This is another test RSS feed</description>
    <item>
      <title>Third Post</title>
      <link>http://localhost:8082/posts/3</link>
      <description>Hello from the second sample feed!</description>
    </item>
  </channel>
</rss>`
}

// StartSampleFeeds launches sample feed servers on 8081 and 8082
func StartSampleFeeds() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/sample.xml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/rss+xml")
			w.Write([]byte(sampleXML1()))
		})
		log.Println("Sample RSS feed available at http://localhost:8081/sample.xml")
		log.Fatal(http.ListenAndServe(":8081", mux))
	}()

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/sample2.xml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/rss+xml")
			w.Write([]byte(sampleXML2()))
		})
		log.Println("Sample RSS feed 2 available at http://localhost:8082/sample2.xml")
		log.Fatal(http.ListenAndServe(":8082", mux))
	}()
}
