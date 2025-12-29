import React from "react";
import { ProjectCard } from "../project-card/project-card";

interface RepoItem {
  id: string;
  title: string;
  repoName?: string;
}

interface RepoListProps {
  disabled?: boolean;
  onSelect?: (id: string) => void;
  items?: RepoItem[];
  title?: string;
}

const RepoList: React.FC<RepoListProps> = ({ disabled = false, onSelect, items, title = "Repo List" }) => {
  const defaultItems: RepoItem[] = items ?? [
    { id: "1", title: "KYMA Flow", repoName: "ChibaniAdam/KYMA-Flow-Front-Admin" },
    { id: "2", title: "KYMA Flow", repoName: "ChibaniAdam/KYMA-Flow-Front-Admin" },
    { id: "3", title: "KYMA Flow", repoName: "ChibaniAdam/KYMA-Flow-Front-Admin" },
    { id: "4", title: "KYMA Flow", repoName: "ChibaniAdam/KYMA-Flow-Front-Admin" },
  ];

  return (
    <div className={`${disabled ? "opacity-40 pointer-events-none" : ""} w-full max-w-sm p-4 sm:p-6 bg-neutral-primary-soft border-l-indigo-200 rounded shadow-xs`}>
      <h5 className="mb-2 text-base md:text-xl font-semibold text-heading">
        {title}
      </h5>

      <p className="text-body">
        Connect with one of our available department list or create a new one.
      </p>

      <ul className="my-6 space-y-3">
        {defaultItems.map((it) => (
          <li key={it.id}>
            <button
              type="button"
              onClick={() => onSelect && onSelect(it.id)}
              className="w-full text-left"
              disabled={disabled}
            >
                  
              <div className="projects-column">
                 
                <ProjectCard
                  icon="⚛️"
                  title={it.title}
                  repoName={it.repoName}
                  repoLink="https://github.com/ChibaniAdam/KYMA-Flow-Front-Admin"
                  dateRange="Oct 28"
                  branch="master"
                  status="◉"
                />
              </div>
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default RepoList;
