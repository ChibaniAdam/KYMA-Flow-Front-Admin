import GiteaIcon from "../../icons/gitea-icon";

const ITEMS = [
  { id: "gitea", label: "Gitea", icon: <GiteaIcon /> },
  { id: "repos", label: "Repositories", icon: <GiteaIcon /> },
  { id: "settings", label: "Settings", icon: <GiteaIcon /> },
];

export function LeftSidebarContent() {
  return (
    <ul className="sidebar__list">
      {ITEMS.map(item => (
        <li key={item.id} className="sidebar__item">
          <button className="sidebar__link" data-tooltip={item.label}>
            <span className="icon">{item.icon}</span>
            <span className="text">{item.label}</span>
          </button>
        </li>
      ))}
    </ul>
  );
}
