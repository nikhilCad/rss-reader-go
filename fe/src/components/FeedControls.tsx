import { useState } from "preact/hooks";
import { Feed } from "../App";

interface FeedControlsProps {
  feeds: Feed[];
  addFeed: (url: string, feed_name: string) => void;
  removeFeed: (url: string) => void;
}

export default function FeedControls({
  feeds,
  addFeed,
  removeFeed,
}: FeedControlsProps) {
  const [url, setUrl] = useState("");
  const [name, setName] = useState("");
  return (
    <div id="feed-controls">
      <form
        id="add-feed-form"
        onSubmit={(e) => {
          e.preventDefault();
          if (url.trim()) {
            addFeed(url.trim(), name.trim()); // pass both
            setUrl("");
            setName("");
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
        <input
          type="text"
          id="feed-name"
          placeholder="Optional feed name"
          value={name}
          onInput={(e) => setName((e.target as HTMLInputElement).value)}
        />
        <button type="submit">Add Feed</button>
      </form>
      <div id="feed-list-controls">
        {feeds?.map((feed) => (
          <div className="feed-url-row" key={feed.url}>
            <span className="feed-url-text">{feed.feed_name || feed.url}</span>
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
