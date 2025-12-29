import "../project-card/project-card.css";
interface DepartmentCardProps {
  icon?: string;
  title: string;
  isOpen: boolean;
  onToggle: () => void;
  checked: boolean;
  onCheck: () => void;
}

export function DepartmentCard({
  icon = "🏢",
  title,
  isOpen,
  onToggle,
  checked,
  onCheck,
}: DepartmentCardProps) {
  return (
    <div className="project-card cursor-pointer" onClick={onToggle}>
      <div className="project-left">
        <input
          type="checkbox"
          checked={checked}
          onChange={onCheck}
          onClick={(e) => e.stopPropagation()}
        />

        <div className="project-icon">{icon}</div>

        <div className="project-meta">
          <h3 className="project-title">{title}</h3>
        </div>
      </div>

      <div className="project-right">
        <div className="project-menu">
          {isOpen ? "▲" : "▼"}
        </div>
      </div>
    </div>
  );
}
