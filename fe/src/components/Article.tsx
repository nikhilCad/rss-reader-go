import { Post } from "../App";

interface ArticleProps {
  post: Post | null;
  isRead: boolean;
  markRead: (link: string) => void;
  markUnread: (link: string) => void;
}

export default function Article({
  post,
  isRead,
  markRead,
  markUnread,
}: ArticleProps) {
  if (!post) return null;
  return (
    <div id="article-detail">
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
        style={{ marginTop: "1em" }}
        dangerouslySetInnerHTML={{ __html: post.description }}
      />
      <div
        style={{ marginTop: "1em" }}
        dangerouslySetInnerHTML={{ __html: post.content }}
      />
    </div>
  );
}
