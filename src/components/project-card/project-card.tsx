import { Tooltip } from "../tooltip/tooltip";
import "./project-card.css";

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
              <a href={repoLink} target="_blank" className="project-subtext flex min-w-0 flex-none flex-row items-center gap-0.5 rounded-full p-0.5 pr-1.5 w-fit max-w-48">
                <svg
                  aria-label="github"
                  height="14"
                  viewBox="0 0 14 14"
                  width="14"
                  className="m-0.5 shrink-0"
                >
                  <path
                    d="M7 .175c-3.872 0-7 3.128-7 7 0 3.084 2.013 5.71 4.79 6.65.35.066.482-.153.482-.328v-1.181c-1.947.415-2.363-.941-2.363-.941-.328-.81-.787-1.028-.787-1.028-.634-.438.044-.416.044-.416.7.044 1.071.722 1.071.722.635 1.072 1.641.766 2.035.59.066-.459.24-.765.437-.94-1.553-.175-3.193-.787-3.193-3.456 0-.766.262-1.378.721-1.881-.065-.175-.306-.897.066-1.86 0 0 .59-.197 1.925.722a6.754 6.754 0 0 1 1.75-.24c.59 0 1.203.087 1.75.24 1.335-.897 1.925-.722 1.925-.722.372.963.131 1.685.066 1.86.46.48.722 1.115.722 1.88 0 2.691-1.641 3.282-3.194 3.457.24.219.481.634.481 1.29v1.926c0 .197.131.415.481.328C11.988 12.884 14 10.259 14 7.175c0-3.872-3.128-7-7-7z"
                    fill="white"
                    fillRule="nonzero"
                  />
                </svg>
                {repoName}
              </a>
            </Tooltip>
          )}

          {dateRange && (
            <span className="project-date">
              {dateRange} — 
              <svg
                className="flex-none"
                data-testid="geist-icon"
                height="16"
                strokeLinejoin="round"
                viewBox="0 0 16 16"
                width="16"
              >
                <path
                  d="M4 6.25V14.25"
                  stroke="currentColor"
                  strokeWidth="1.5"
                  strokeLinecap="square"
                  strokeLinejoin="round"
                />
                <path
                  fillRule="evenodd"
                  clipRule="evenodd"
                  d="M10.5 12C10.5 12.8284 11.1716 13.5 12 13.5C12.8284 13.5 13.5 12.8284 13.5 12C13.5 11.1716 12.8284 10.5 12 10.5C11.1716 10.5 10.5 11.1716 10.5 12ZM9.079 12.6869C9.38957 14.0127 10.5795 15 12 15C13.6569 15 15 13.6569 15 12C15 10.3431 13.6569 9 12 9C10.6293 9 9.47333 9.91924 9.1149 11.1749C8.05096 10.9929 7.0611 10.4857 6.28769 9.71231C5.51428 8.9389 5.0071 7.94904 4.82513 6.8851C6.08076 6.52667 7 5.37069 7 4C7 2.34315 5.65685 1 4 1C2.34315 0.999999 1 2.34315 1 4C1 5.42051 1.98728 6.61042 3.3131 6.921C3.51279 8.37102 4.18025 9.72619 5.22703 10.773C6.2738 11.8197 7.62898 12.4872 9.079 12.6869ZM2.5 4C2.5 4.82843 3.17157 5.5 4 5.5C4.82843 5.5 5.5 4.82843 5.5 4C5.5 3.17157 4.82843 2.5 4 2.5C3.17157 2.5 2.5 3.17157 2.5 4Z"
                  fill="currentColor"
                />
              </svg>
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
