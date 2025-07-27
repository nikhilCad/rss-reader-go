import Sidebar from "./Sidebar";
import Article from "./Article";
import { useEffect, useState } from "preact/hooks";
import { Post } from "App";
import { useStore } from "../utils/store";


interface HomeProps {
  grouped: Record<string, Post[]>;
}

export default function Home({ grouped }: HomeProps) {

  const {
    readLinks,
    selected,
    refreshing,
    markAsRead,
    markAsUnread,
    setSelected,
    refreshFeeds,
  } = useStore();


  const [showUnread, setShowUnread] = useState(true);
  useEffect(() => {

    const handleGlobalSpacePress = (event: KeyboardEvent) => {
      if (event.key === ' ') {
        event.preventDefault();

        const targetTagName = (event.target as HTMLElement).tagName;
        if (targetTagName === 'INPUT' || targetTagName === 'TEXTAREA' || targetTagName === 'SELECT') {
          return;
        }

        window.open(selected?.link, "_blank");
      }
    };

    window.addEventListener('keydown', handleGlobalSpacePress);

    return () => {
      window.removeEventListener('keydown', handleGlobalSpacePress);
    };

  }, [selected]);

  return (

    <>
      <div id="toolbar">
        <button onClick={refreshFeeds} disabled={refreshing}>
          {refreshing ? "Refreshing..." : "Refresh Feeds"}
        </button>
        <button onClick={ () => setShowUnread(!showUnread)}>
          {showUnread ? "Showing: All" : "Showing: Unread"}
        </button>
      </div>
      <div id="main-container">
        <div id="left-pane">
          <Sidebar
            grouped={grouped}
            readLinks={readLinks}
            selected={selected}
            setSelected={setSelected}
            showUnread={showUnread}
            markAsRead={markAsRead}
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

  );

}