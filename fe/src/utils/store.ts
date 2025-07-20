import { create } from "zustand";
import {
  reloadAll,
  addFeed,
  removeFeed,
  markRead,
  markUnread,
} from "./feed-utils";
import type { Feed, Post } from "../App";

interface StoreState {
  feeds: Feed[];
  posts: Post[];
  readLinks: Set<string>;
  selected: Post | null;
  refreshing: boolean;
  reload: () => Promise<void>;
  addFeedUrl: (url: string, feed_name: string) => Promise<void>;
  removeFeedUrl: (url: string) => Promise<void>;
  markAsRead: (link: string) => Promise<void>;
  markAsUnread: (link: string) => Promise<void>;
  setSelected: (post: Post | null) => void;
  refreshFeeds: () => Promise<void>;
}

export const useStore = create<StoreState>((set, get) => ({
  feeds: [],
  posts: [],
  readLinks: new Set(),
  selected: null,
  refreshing: false,

  setSelected: (post) => set({ selected: post }),

  reload: async () => {
    const { feeds, readLinks, posts } = await reloadAll();
    set({ feeds, readLinks, posts });

    // Update selected if needed
    const { selected } = get();
    if (posts.length > 0 && posts.some((a) => a.link === selected?.link)) {
      // keep current selection
    } else if (posts.length > 0) {
      set({ selected: posts[0] });
    } else {
      set({ selected: null });
    }
  },

  addFeedUrl: async (url, name) => {
    await addFeed(url, name);
    await get().reload();
  },

  removeFeedUrl: async (url) => {
    await removeFeed(url);
    await get().reload();
  },

  markAsRead: async (link) => {
    await markRead(link);
    set((state) => ({
      readLinks: new Set([...state.readLinks, link]),
    }));
  },

  markAsUnread: async (link) => {
    await markUnread(link);
    set((state) => {
      const newSet = new Set(state.readLinks);
      newSet.delete(link);
      return { readLinks: newSet };
    });
  },

  refreshFeeds: async () => {
    set({ refreshing: true });
    try {
      await fetch("/refresh", { method: "POST" });
      await get().reload();
    } catch (e) {
      alert("Failed to refresh feeds");
    }
    set({ refreshing: false });
  },
}));
