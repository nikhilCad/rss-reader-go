import { useEffect, useState } from "preact/hooks";
import { Post } from "../App";
import ArticleTopBar from "./ArticleTopBar";

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
  const [parsed, setParsed] = useState<null | {
    title: string;
    content: string;
    byline: string;
  }>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setParsed(null);
    setError(null);
  }, [post && post?.link]);

  useEffect(() => {
    setParsed(null);
    setError(null);
  }, [post && post?.link]);

  if (!post) return null;

  const content = parsed ? parsed.content : post.description || post.content;

  return (
    <div id="article-detail">
      <ArticleTopBar
        isRead={isRead}
        onMarkRead={() => markRead(post.link)}
        onMarkUnread={() => markUnread(post.link)}
        postLink={post.link}
        onParsed={(parsed, error) => {
          setParsed(parsed);
          setError(error);
        }}
        content={post.title + content}
      />
      <h2>{post.title}</h2>
      <a href={post.link} target="_blank" rel="noopener noreferrer">
        Read original
      </a>
      {error && <div style={{ color: "red" }}>Error: {error}</div>}
      <div
        style={{ marginTop: "1em" }}
        dangerouslySetInnerHTML={{
          __html: content,
        }}
      />
      {parsed && parsed.byline && (
        <div style={{ marginTop: "1em", fontStyle: "italic" }}>
          By: {parsed.byline}
        </div>
      )}
    </div>
  );
}
