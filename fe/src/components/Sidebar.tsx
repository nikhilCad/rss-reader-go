import { Post } from "../App";

interface SidebarProps {
  grouped: Record<string, Post[]>;
  readLinks: Set<string>;
  selected: Post | null;
  setSelected: (post: Post) => void;
}

export default function Sidebar({
  grouped,
  readLinks,
  selected,
  setSelected,
}: SidebarProps) {
  return (
    <div id="sidebar">
      {Object.keys(grouped).map((source) => (
        <div key={source}>
          <div className="feed-source">{source}</div>
          {grouped[source].map((post) => (
            <a
              className={`article-title${
                readLinks.has(post.link) ? " read" : ""
              }${selected && selected.link === post.link ? " active" : ""}`}
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
