import { useEffect, useState } from "preact/hooks";
import "./app.css";
import Sidebar from "./components/Sidebar";
import Article from "./components/Article";
import { useStore } from "./utils/store";
import NavSidebar from "./components/NavSidebar";
import FeedControls from "./components/FeedControls";

export interface Feed {
  id: number;
  url: string;
  feed_name: string;
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
  content: string;
}

export default function App() {
  const {
    feeds,
    posts,
    readLinks,
    selected,
    refreshing,
    reload,
    addFeedUrl,
    removeFeedUrl,
    markAsRead,
    markAsUnread,
    setSelected,
    refreshFeeds,
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
          <>
            <div id="toolbar">
              <button onClick={refreshFeeds} disabled={refreshing}>
                {refreshing ? "Refreshing..." : "Refresh Feeds"}
              </button>
            </div>
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
                  markRead={markAsRead}
                  markUnread={markAsUnread}
                />
              </div>
            </div>
          </>
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
              <li>Podcast</li>
              <li>Better UI</li>
              <li>Better Feed Management</li>
              <li>AI summary of all feed titles</li>
              <li>AI summary of fetched feed content</li>
              <li>Better full article fetching</li>
              <li>Actual Setting Page</li>
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}
