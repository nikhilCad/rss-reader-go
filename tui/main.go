package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/rivo/tview"
)

type Post struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Source      string `json:"source"` // feed_name
	PubDate     string `json:"pubdate"`
}

type PostsResponse struct {
	FromCache bool   `json:"fromCache"`
	Articles  []Post `json:"articles"`
}

func fetchArticles() ([]Post, error) {
	resp, err := http.Get("http://localhost:8080/posts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var pr PostsResponse
	err = json.NewDecoder(resp.Body).Decode(&pr)
	return pr.Articles, err
}

func pad(s string, padLines int) string {
	padStr := strings.Repeat("\n", padLines)
	return padStr + s + padStr
}

// cleanHTML removes HTML tags and cleans up the text
func cleanHTML(input string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(input, "")

	// Clean up common HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	text = strings.ReplaceAll(text, "&apos;", "'")

	// Remove excessive whitespace
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	// Clean up line breaks
	text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	text = strings.TrimSpace(text)

	return text
}

// wrapText wraps text to fit within specified width
func wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	var lines []string
	words := strings.Fields(text)
	var currentLine []string
	currentLength := 0

	for _, word := range words {
		if currentLength+len(word)+1 > width && len(currentLine) > 0 {
			lines = append(lines, strings.Join(currentLine, " "))
			currentLine = []string{word}
			currentLength = len(word)
		} else {
			currentLine = append(currentLine, word)
			currentLength += len(word) + 1
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, strings.Join(currentLine, " "))
	}

	return strings.Join(lines, "\n")
}

func main() {
	app := tview.NewApplication()

	// Use Table instead of List for better multi-line support
	articleTable := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetSeparator(tview.Borders.Vertical)

	textView := tview.NewTextView()
	textView.SetDynamicColors(true).
		SetWrap(true).
		SetScrollable(true).
		SetTextAlign(tview.AlignLeft).
		SetBorder(true).
		SetTitle(" Article Content ")
	contentView := textView

	articles, err := fetchArticles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch articles: %v\n", err)
		os.Exit(1)
	}

	lastFeed := ""
	row := 0
	articleMap := make(map[int]Post) // Map row to article

	for _, article := range articles {
		// Add a header for each new feed
		if article.Source != lastFeed {
			articleTable.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("─ %s ─", article.Source)).
				SetTextColor(tview.Styles.SecondaryTextColor).
				SetSelectable(false))
			row++
			lastFeed = article.Source
		}

		// Wrap article title for multi-line display
		wrappedTitle := wrapText(article.Title, 50) // Adjust width as needed

		articleTable.SetCell(row, 0, tview.NewTableCell(wrappedTitle).
			SetMaxWidth(50).
			SetExpansion(1))

		articleMap[row] = article
		row++
	}

	// Set selection handler
	articleTable.SetSelectedFunc(func(row, column int) {
		if article, exists := articleMap[row]; exists {
			content := article.Content
			if content == "" {
				content = article.Description
			}

			// Clean HTML from content
			cleanContent := cleanHTML(content)

			contentView.Clear()
			fmt.Fprint(contentView, pad(fmt.Sprintf(
				"[yellow::b]%s[white]\n[gray]%s[white]\n\n%s\n\n[blue]Link: %s",
				article.Title, article.PubDate, cleanContent, article.Link,
			), 1))
			contentView.ScrollToBeginning()
		}
	})

	// Load first article by default
	if len(articles) > 0 {
		content := articles[0].Content
		if content == "" {
			content = articles[0].Description
		}
		cleanContent := cleanHTML(content)

		contentView.Clear()
		fmt.Fprint(contentView, pad(fmt.Sprintf(
			"[yellow::b]%s[white]\n[gray]%s[white]\n\n%s\n\n[blue]Link: %s",
			articles[0].Title, articles[0].PubDate, cleanContent, articles[0].Link,
		), 1))
		contentView.ScrollToBeginning()
	}

	flex := tview.NewFlex().
		AddItem(articleTable, 0, 1, true).
		AddItem(contentView, 0, 2, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
