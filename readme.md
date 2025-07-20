# RSS Reader Go

A modern RSS reader application built with Go backend, Preact frontend, and a terminal user interface. This project provides multiple ways to consume RSS feeds with caching, article parsing, and various user interfaces.

## 🚀 Features

### Backend (Go)

- **RSS Feed Management**: Add, remove, and manage RSS feed subscriptions
- **Smart Caching**: SQLite-based caching system to avoid rate limits and ensure fast loading. [Detailed implementation explanation](be/readme.md)
- **Article Parsing**: Extract full article content from web pages using go-readability
- **Podcast Support**: Handle media enclosures for podcast episodes with audio player support
- **Read Status Tracking**: Mark articles as read/unread with persistent storage
- **Sample Feeds**: Built-in sample RSS feeds for testing and demonstration
- **RESTful API**: Clean API endpoints for all operations
- **Background Refresh**: Manual feed refresh without blocking normal operations

### Frontend (Preact + Vite)

- **Modern UI**: Lightweight React-like interface using Preact for smaller bundle size
- **Feed Controls**: Easy feed management with URL and name input
- **Article Reader**: Clean article reading interface with sidebar navigation
- **Text-to-Speech**: Built-in TTS functionality for audio article consumption
- **Article Parsing**: On-demand full article content extraction
- **Read Status**: Visual indicators for read/unread articles
- **Responsive Design**: Multi-pane layout with navigation sidebar
- **Hot Reload**: Fast development with Vite's HMR

### TUI (Terminal User Interface)

- **Terminal-based**: Full-featured RSS reader that runs in the terminal
- **Article Browser**: Navigate articles using keyboard shortcuts
- **Content Display**: Clean text rendering with HTML tag removal
- **Feed Grouping**: Articles organized by RSS feed source
- **Text Wrapping**: Proper text formatting for terminal display
- **Mouse Support**: Click navigation in supported terminals

## 🛠️ Setup & Installation

### Prerequisites

- Go 1.19+ installed
- Node.js 16+ installed
- Git

### Backend Setup

1. **Install Go dependencies:**

   ```sh
   go mod tidy
   ```

2. **Run the backend:**
   ```sh
   cd be
   go run *.go
   ```
   The backend will be available at http://localhost:8080

### Frontend Setup

1. **Navigate to frontend directory:**

   ```sh
   cd fe
   ```

2. **Install dependencies:**

   ```sh
   npm install
   ```

3. **Start development server:**

   ```sh
   npm run dev
   ```

   The frontend will be available at http://localhost:5173

4. **Build for production:**
   ```sh
   npm run build
   ```
   Static files will be generated in `fe/dist/`

### TUI Setup

1. **Run the terminal interface:**
   ```sh
   cd tui
   go run main.go
   ```
   Note: The backend must be running first at http://localhost:8080

## 🏗️ Build Process

### Backend Build

```sh
cd be
go build -o rss-reader-backend *.go
./rss-reader-backend
```

### Frontend Build

```sh
cd fe
npm run build
# Serve the dist/ folder with any static server
```

### TUI Build

```sh
cd tui
go build -o rss-reader-tui main.go
./rss-reader-tui
```

## 📡 API Endpoints

- `GET /posts` - Retrieve all cached articles
- `POST /refresh` - Refresh all RSS feeds
- `GET /feeds` - List all subscribed feeds
- `POST /feeds` - Add a new RSS feed
- `DELETE /feeds` - Remove a feed
- `GET /read` - List read article links
- `POST /read` - Mark article as read
- `POST /unread` - Mark article as unread
- `GET /parse-article?url=<url>` - Parse full article content

## 💡 Usage Tips

- Use the frontend for the best visual experience
- Use the TUI for terminal-only environments or lightweight usage
- The backend caches all articles locally for fast loading
- Use the refresh button to fetch new articles from RSS feeds
- The system supports both regular RSS feeds and podcast feeds with audio

## 🔧 Development

- Frontend uses Vite proxy to communicate with the Go backend
- Backend includes sample RSS feeds for testing
- All data is stored in SQLite database (`be/posts.db`)
- Hot reload available for both frontend and backend development

---

_Built with ❤️ using Go, Preact, and modern web technologies_
