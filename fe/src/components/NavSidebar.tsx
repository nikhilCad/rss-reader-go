import "./NavSidebar.css";

interface NavSidebarProps {
  active: string;
  setActive: (page: string) => void;
}

export default function NavSidebar({ active, setActive }: NavSidebarProps) {
  return (
    <nav className="nav-sidebar">
      <div
        className={`nav-icon${active === "home" ? " active" : ""}`}
        title="Home"
        onClick={() => setActive("home")}
      >
        🏠
      </div>
      <div
        className={`nav-icon${active === "feeds" ? " active" : ""}`}
        title="Manage Feeds"
        onClick={() => setActive("feeds")}
      >
        📚
      </div>
      <div
        className={`nav-icon${active === "settings" ? " active" : ""}`}
        title="Settings"
        onClick={() => setActive("settings")}
      >
        ⚙️
      </div>
    </nav>
  );
}
