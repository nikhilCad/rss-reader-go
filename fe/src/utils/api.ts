import { Feed, PostResponse } from "../App";

export async function fetchFeeds(): Promise<Feed[]> {
  const res = await fetch("/feeds");
  return await res.json();
}

export async function fetchReadLinks(): Promise<string[]> {
  const res = await fetch("/read");
  return await res.json();
}

export async function fetchPosts(): Promise<PostResponse> {
  const res = await fetch("/posts");
  return await res.json();
}

export async function addFeed(url: string, name: string) {
  await fetch("/feeds", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(name ? { url, name } : { url }),
  });
}

export async function removeFeed(url: string) {
  await fetch("/feeds", {
    method: "DELETE",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ url }),
  });
}

export async function markRead(link: string) {
  await fetch("/read", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ link }),
  });
}

export async function markUnread(link: string) {
  await fetch("/unread", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ link }),
  });
}
