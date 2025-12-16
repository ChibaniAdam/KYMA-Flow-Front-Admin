import { Tooltip } from "../tooltip/tooltip";
import "./project-card.css";
import GithubIcon from "../../assets/github-icon.svg";
import BranchIcon from "../../assets/branch-icon.svg";

interface ProjectCardProps {
  readonly icon?: string;
  readonly title: string;
  readonly repoName?: string;
  readonly repoLink?: string;
  readonly dateRange?: string;
  readonly branch?: string;
  readonly status?: string;
}

export function ProjectCard({
  icon = "⚛️",
  title,
  repoName,
  repoLink,
  dateRange,
  branch = "master",
  status = "◉",
}: ProjectCardProps) {
  return (
    <div className="project-card">
      <div className="project-left">
        <div className="project-icon">{icon}</div>
        <div className="project-meta">
          <h3 className="project-title">{title}</h3>

          {repoName && (
            <Tooltip text={repoName} color="var(--light-blue-600)" position="top">
    <a
  href={repoLink}
  target="_blank"
  className="project-subtext flex min-w-0 flex-row items-center gap-0.5 rounded-full p-0.5 pr-1.5 max-w-48"
>
<img src={GithubIcon} alt="Github Icon"/>

  <span className="min-w-0 overflow-hidden text-ellipsis whitespace-nowrap">
    {repoName}
  </span>
</a>

            </Tooltip>
          )}

          {dateRange && (
            <span className="project-date">
              {dateRange} — 
              <img src={BranchIcon} alt="branch icon"/>
              &nbsp;
              {branch}
            </span>
          )}
        </div>
      </div>

      <div className="project-right">
        <div className="project-status">{status}</div>
        <div className="project-menu">⋮</div>
      </div>
    </div>
  );
}
