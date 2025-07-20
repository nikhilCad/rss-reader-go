import { useRef, useState } from "preact/hooks";

interface ArticleTopBarProps {
  isRead: boolean;
  onMarkRead: () => void;
  onMarkUnread: () => void;
  postLink: string;
  onParsed: (
    parsed: { title: string; content: string; byline: string } | null,
    error: string | null
  ) => void;
  content: string;
}

export default function ArticleTopBar({
  isRead,
  onMarkRead,
  onMarkUnread,
  postLink,
  onParsed,
  content,
}: ArticleTopBarProps) {
  const [loading, setLoading] = useState(false);

  const [ttsActive, setTtsActive] = useState(false); // Whether TTS controls are shown
  const [ttsPlaying, setTtsPlaying] = useState(false); // Whether TTS is currently playing
  const utteranceRef = useRef<SpeechSynthesisUtterance | null>(null);

  // Dummy text for TTS (replace with actual article content as needed)
  const ttsText = content;

  // Start TTS
  const handleTtsStart = () => {
    if (!ttsActive) setTtsActive(true);
    if (utteranceRef.current) {
      window.speechSynthesis.cancel();
    }
    const utter = new window.SpeechSynthesisUtterance(ttsText);
    utter.onend = () => {
      setTtsPlaying(false);
      setTtsActive(false);
    };
    utter.onpause = () => setTtsPlaying(false);
    utter.onresume = () => setTtsPlaying(true);
    utter.onerror = () => {
      setTtsPlaying(false);
      setTtsActive(false);
    };
    utteranceRef.current = utter;
    window.speechSynthesis.speak(utter);
    setTtsPlaying(true);
  };

  // Pause/Resume TTS
  const handleTtsToggle = () => {
    if (!ttsPlaying) {
      window.speechSynthesis.resume();
      setTtsPlaying(true);
    } else {
      window.speechSynthesis.pause();
      setTtsPlaying(false);
    }
  };

  // Stop TTS
  const handleTtsStop = () => {
    window.speechSynthesis.cancel();
    setTtsPlaying(false);
    setTtsActive(false);
  };

  const handleParse = async () => {
    setLoading(true);
    onParsed(null, null);
    try {
      const res = await fetch(
        `/parse-article?url=${encodeURIComponent(postLink)}`
      );
      if (!res.ok) throw new Error("Failed to parse article");
      const data = await res.json();
      onParsed(data, null);
    } catch (e: any) {
      onParsed(null, e.message || "Unknown error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div id="article-toolbar">
      <button type="button" onClick={isRead ? onMarkUnread : onMarkRead}>
        {isRead ? "Mark as Unread" : "Mark as Read"}
      </button>
      <button
        type="button"
        onClick={handleParse}
        disabled={loading}
        style={{ marginLeft: 8 }}
      >
        {loading ? "Parsing..." : "Parse Article (experimental)"}
      </button>
      <button
        type="button"
        onClick={handleTtsStart}
        style={{ marginLeft: 8 }}
        disabled={ttsPlaying}
      >
        {ttsActive ? "Restart TTS" : "TTS"}
      </button>
      {ttsActive && (
        <>
          <button
            type="button"
            onClick={handleTtsToggle}
            style={{ marginLeft: 8 }}
          >
            {ttsPlaying ? "Pause" : "Play"}
          </button>
          <button
            type="button"
            onClick={handleTtsStop}
            style={{ marginLeft: 8 }}
          >
            Stop
          </button>
        </>
      )}
    </div>
  );
}
