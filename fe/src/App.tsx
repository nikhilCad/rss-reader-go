import { useEffect, useState } from "preact/hooks";
import "./app.css";
import Sidebar from "./Sidebar";
import Article from "./Article";

export interface Feed {
  url: string;
}

export interface PostResponse {
  fromCache: boolean;
  articles: Post[];
}

export interface Post {
  title: string;
  link: string;
  description: string;
  source?: string;
}

export default function App() {
  const [feeds, setFeeds] = useState<Feed[]>([]);
  const [readLinks, setReadLinks] = useState<Set<string>>(new Set());
  const [posts, setPosts] = useState<Post[]>([]);
  const [selected, setSelected] = useState<Post | null>(null);
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    reloadFeedsAndPosts();
  }, []);

  async function fetchFeeds(): Promise<Feed[]> {
    const res = await fetch("/feeds");
    return await res.json();
  }

  async function fetchReadLinks(): Promise<string[]> {
    const res = await fetch("/read");
    return await res.json();
  }

  async function fetchPosts(): Promise<PostResponse> {
    const res = await fetch("/posts");
    return await res.json();
  }

  async function reloadFeedsAndPosts() {
    const [feedsFetched, read, postsFetched] = await Promise.all([
      fetchFeeds(),
      fetchReadLinks(),
      fetchPosts(),
    ]);
    setFeeds(feedsFetched);
    setReadLinks(new Set(read));
    setPosts(postsFetched?.articles);
    if (postsFetched?.articles.length > 0)
      setSelected(postsFetched?.articles[0]);
  }

  async function addFeed(url: string) {
    await fetch("/feeds", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url }),
    });
    await reloadFeedsAndPosts();
  }

  async function removeFeed(url: string) {
    await fetch("/feeds", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url }),
    });
    await reloadFeedsAndPosts();
  }

  async function markRead(link: string) {
    await fetch("/read", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ link }),
    });
    setReadLinks(new Set([...readLinks, link]));
  }

  async function markUnread(link: string) {
    await fetch("/unread", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ link }),
    });
    const newSet = new Set(readLinks);
    newSet.delete(link);
    setReadLinks(newSet);
  }

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await fetch("/refresh", { method: "POST" });
      await reloadFeedsAndPosts();
    } catch (e) {
      alert("Failed to refresh feeds");
    }
    setRefreshing(false);
  };

  const grouped: Record<string, Post[]> = {};
  posts.forEach((post) => {
    const source = post.source || "Feed";
    if (!grouped[source]) grouped[source] = [];
    grouped[source].push(post);
  });

  return (
    <div>
      <h1>RSS Reader</h1>
      <FeedControls feeds={feeds} addFeed={addFeed} removeFeed={removeFeed} />
      <button onClick={handleRefresh} disabled={refreshing}>
        {refreshing ? "Refreshing..." : "Refresh Feeds"}
      </button>
      <div id="main-container">
        <div id="left-pane">
          <Sidebar
            grouped={grouped}
            readLinks={readLinks}
            selected={selected}
            setSelected={setSelected}
          />
        </div>
        <div id="right-pane">
          <Article
            post={selected}
            isRead={!!(selected && readLinks.has(selected.link))}
            markRead={markRead}
            markUnread={markUnread}
          />
        </div>
      </div>
    </div>
  );
}

interface FeedControlsProps {
  feeds: Feed[];
  addFeed: (url: string) => void;
  removeFeed: (url: string) => void;
}

function FeedControls({ feeds, addFeed, removeFeed }: FeedControlsProps) {
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
          onInput={(e) => setUrl((e.target as HTMLInputElement).value)}
        />
        <button type="submit">Add Feed</button>
      </form>
      <div id="feed-list-controls">
        {feeds.map((feed) => (
          <div className="feed-url-row" key={feed.url}>
            <span className="feed-url-text">{feed.url}</span>
            <button
              className="remove-feed-btn"
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
