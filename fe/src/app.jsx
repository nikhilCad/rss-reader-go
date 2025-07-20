import { useEffect, useState } from "preact/hooks";
import "./App.css";

export default function App() {
  const [feeds, setFeeds] = useState([]);
  const [readLinks, setReadLinks] = useState(new Set());
  const [posts, setPosts] = useState([]);
  const [selected, setSelected] = useState(null);

  // Fetch feeds and read links on mount
  useEffect(() => {
    reloadFeedsAndPosts();
  }, []);

  async function fetchFeeds() {
    const res = await fetch("/feeds");
    return await res.json();
  }

  async function fetchReadLinks() {
    const res = await fetch("/read");
    return await res.json();
  }

  async function fetchPosts() {
    const res = await fetch("/posts");
    return await res.json();
  }

  async function reloadFeedsAndPosts() {
    const [feeds, read, posts] = await Promise.all([
      fetchFeeds(),
      fetchReadLinks(),
      fetchPosts(),
    ]);
    setFeeds(feeds);
    setReadLinks(new Set(read));
    setPosts(posts);
    if (posts.length > 0) setSelected(posts[0]);
  }

  async function addFeed(url) {
    await fetch("/feeds", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url }),
    });
    await reloadFeedsAndPosts();
  }

  async function removeFeed(url) {
    await fetch("/feeds", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url }),
    });
    await reloadFeedsAndPosts();
  }

  async function markRead(link) {
    await fetch("/read", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ link }),
    });
    setReadLinks(new Set([...readLinks, link]));
  }

  async function markUnread(link) {
    await fetch("/unread", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ link }),
    });
    const newSet = new Set(readLinks);
    newSet.delete(link);
    setReadLinks(newSet);
  }

  // Group posts by source
  const grouped = {};
  posts.forEach((post) => {
    const source = post.source || "Feed";
    if (!grouped[source]) grouped[source] = [];
    grouped[source].push(post);
  });

  // UI
  return (
    <div>
      <h1>RSS Reader</h1>
      <FeedControls feeds={feeds} addFeed={addFeed} removeFeed={removeFeed} />
      <div id="main-container">
        <div id="left-pane">
          <FeedList
            grouped={grouped}
            readLinks={readLinks}
            selected={selected}
            setSelected={setSelected}
          />
        </div>
        <div id="right-pane">
          {selected && (
            <ArticleDetail
              post={selected}
              isRead={readLinks.has(selected.link)}
              markRead={markRead}
              markUnread={markUnread}
            />
          )}
        </div>
      </div>
    </div>
  );
}

function FeedControls({ feeds, addFeed, removeFeed }) {
  const [url, setUrl] = useState("");
  return (
    <div id="feed-controls">
      <form
        id="add-feed-form"
        onSubmit={(e) => {
          e.preventDefault();
          if (url.trim()) {
            addFeed(url.trim());
            setUrl("");
          }
        }}
      >
        <input
          type="url"
          id="feed-url"
          placeholder="Add RSS feed URL..."
          required
          value={url}
          onInput={(e) => setUrl(e.target.value)}
        />
        <button type="submit">Add Feed</button>
      </form>
      <div id="feed-list-controls">
        {feeds.map((feed) => (
          <div class="feed-url-row" key={feed.url}>
            <span class="feed-url-text">{feed.url}</span>
            <button
              class="remove-feed-btn"
              onClick={() => removeFeed(feed.url)}
              type="button"
            >
              Remove
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}

function FeedList({ grouped, readLinks, selected, setSelected }) {
  return (
    <div id="feed-list">
      {Object.keys(grouped).map((source) => (
        <div key={source}>
          <div class="feed-source">{source}</div>
          {grouped[source].map((post) => (
            <a
              class={`article-title${readLinks.has(post.link) ? " read" : ""}${
                selected && selected.link === post.link ? " active" : ""
              }`}
              href="#"
              onClick={(e) => {
                e.preventDefault();
                setSelected(post);
              }}
              key={post.link}
            >
              {post.title}
            </a>
          ))}
        </div>
      ))}
    </div>
  );
}

function ArticleDetail({ post, isRead, markRead, markUnread }) {
  return (
    <div>
      <div id="article-toolbar">
        <button
          type="button"
          onClick={() => (isRead ? markUnread(post.link) : markRead(post.link))}
        >
          {isRead ? "Mark as Unread" : "Mark as Read"}
        </button>
      </div>
      <h2>{post.title}</h2>
      <a href={post.link} target="_blank" rel="noopener noreferrer">
        Read original
      </a>
      <div
        style="margin-top:1em;"
        dangerouslySetInnerHTML={{ __html: post.description }}
      />
    </div>
  );
}
