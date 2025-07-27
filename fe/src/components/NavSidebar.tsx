import "./NavSidebar.css";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faHouse, faRssSquare, faGear } from '@fortawesome/free-solid-svg-icons'

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
        <FontAwesomeIcon icon={faHouse} />
      </div>
      <div
        className={`nav-icon${active === "feeds" ? " active" : ""}`}
        title="Manage Feeds"
        onClick={() => setActive("feeds")}
      >
        <FontAwesomeIcon icon={faRssSquare} />
      </div>
      <div
        className={`nav-icon${active === "settings" ? " active" : ""}`}
        title="Settings"
        onClick={() => setActive("settings")}
      >
        <FontAwesomeIcon icon={faGear} />
      </div>
    </nav>
  );
}
