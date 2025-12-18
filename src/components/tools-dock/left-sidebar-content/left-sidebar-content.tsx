import GiteaIcon from "../../icons/gitea-icon";

const ITEMS = [
  { id: "gitea", label: "Gitea", subtitle: "Source Code Hosting", icon: <GiteaIcon /> },
  { id: "sonarqube", label: "SonarQube", subtitle: "Code Quality & Security", icon: <GiteaIcon /> },
  { id: "settings", label: "Settings", subtitle: "Configure Tools", icon: <GiteaIcon /> },
];

export function LeftSidebarContent() {
  return (
    <ul className="sidebar__list">
      {ITEMS.map(item => (
        <li key={item.id} className="sidebar__item">
          <button className="sidebar__link" data-tooltip={`${item.label} - ${item.subtitle}`}>
            <span className="icon">{item.icon}</span>
            <span className="text">
              <span className="label">{item.label}</span>
              <span className="subtitle">{item.subtitle}</span>
            </span>
          </button>
        </li>
      ))}
    </ul>
  );
}
