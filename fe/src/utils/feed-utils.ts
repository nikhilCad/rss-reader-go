import {
  fetchFeeds,
  fetchReadLinks,
  fetchPosts,
  addFeed,
  removeFeed,
  markRead,
  markUnread,
} from "./api";

export async function reloadAll() {
  const [feedsFetched, read, postsFetched] = await Promise.all([
    fetchFeeds(),
    fetchReadLinks(),
    fetchPosts(),
  ]);
  return {
    feeds: feedsFetched,
    readLinks: new Set(read),
    posts: postsFetched?.articles || [],
  };
}

export { addFeed, removeFeed, markRead, markUnread };
