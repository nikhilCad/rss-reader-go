import { Post } from "../App";

interface SidebarProps {
  grouped: Record<string, Post[]>;
  readLinks: Set<string>;
  selected: Post | null;
  setSelected: (post: Post) => void;
  showUnread: boolean;
  markAsRead: (link: string) => Promise<void>;
}

export default function Sidebar({
  grouped,
  readLinks,
  selected,
  setSelected,
  showUnread,
  markAsRead,
}: SidebarProps) {
  return (
    <div id="sidebar">
      {Object.keys(grouped).map((source) => (
        <div key={source}>
          <div className="feed-source">{source}</div>
          {grouped[source]
          .filter((post) => showUnread ? true : !readLinks.has(post.link))
          .map((post) => (
            <a
              className={`article-title${
                readLinks.has(post.link) ? " read" : ""
              }${selected && selected.link === post.link ? " active" : ""}`}
              href="#"
              onClick={(e) => {
                e.preventDefault();
                setSelected(post);
                markAsRead(post.link);
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
