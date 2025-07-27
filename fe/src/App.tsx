import { useEffect, useState } from "preact/hooks";
import "./app.css";
import { useStore } from "./utils/store";
import NavSidebar from "./components/NavSidebar";
import FeedControls from "./components/FeedControls";
import Home from "./components/Home";

export interface Feed {
  id: number;
  url: string;
  feed_name: string;
}

export interface PostResponse {
  fromCache: boolean;
  articles: Post[];
}

export interface Enclosure {
  url: string;
  type?: string;
  length?: string;
}

export interface Post {
  title: string;
  link: string;
  description: string;
  source?: string;
  content: string;
  pubdate: string;
  enclosure?: Enclosure;
}

export default function App() {
  const {
    feeds,
    posts,
    reload,
    addFeedUrl,
    removeFeedUrl,
  } = useStore();

  useEffect(() => {
    reload();
  }, []);

  const grouped: Record<string, Post[]> = {};
  posts?.forEach((post) => {
    const feed = feeds.find(
      (f) => f.url === post.source || f.feed_name === post.source
    );
    const key = feed?.feed_name || post.source || "Feed";
    if (!grouped[key]) grouped[key] = [];
    grouped[key].push(post);
  });

  const [activePage, setActivePage] = useState("home");

  return (
    <div>
      <NavSidebar active={activePage} setActive={setActivePage} />
      <div style={{ marginLeft: 60 }}>
        {activePage === "home" && (
          <Home
            grouped={grouped}
          />
        )}
        {activePage === "feeds" && (
          <FeedControls
            feeds={feeds}
            addFeed={addFeedUrl}
            removeFeed={removeFeedUrl}
          />
        )}
        {activePage === "settings" && (
          <div style={{ padding: 32 }}>
            Settings page (coming soon)
            <ul>
              <li>Parse Youtube Links</li>
              <li>Better UI</li>
              <li>Better Feed Management</li>
              <li>AI summary of all feed titles</li>
              <li>AI summary of fetched feed content</li>
              <li>Better full article fetching</li>
              <li>Actual Setting Page</li>
              <li>Podcast - Done</li>
              <li>TTS - Done</li>
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}
