import { useEffect, useState, useRef } from "react";
import { ProjectCard } from "../../components/project-card/project-card";
import WidgetCard from "../../components/widget-card/widget-card";
import {
  listRepositories,
  searchRepositories,
} from "../../services/repositoryService";
import type { Repository } from "../../GQL/models/repository";
import "./projects.css";
import { useDebounce } from "../../utils/useDebounce";
import { FilterBar } from "../../components/filter-bar/filter-bar";
import { useDelayedLoading } from "../../utils/useDelayedLoading";

const PAGE_SIZE = 20;

export default function Projects() {
  const [projects, setProjects] = useState<Repository[]>([]);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const showLoadingMessage = useDelayedLoading(loading);
  const loaderRef = useRef<HTMLDivElement | null>(null);

  /* ---------------- Fetch repositories ---------------- */

  const fetchRepositories = async (reset = false) => {
    if (loading || (!hasMore && !reset)) return;

    setLoading(true);

    try {
      const res =
        debouncedSearch.trim().length > 0
          ? await searchRepositories(debouncedSearch, {
              limit: PAGE_SIZE,
              offset: reset ? 0 : offset,
            })
          : await listRepositories({
              limit: PAGE_SIZE,
              offset: reset ? 0 : offset,
            });

      setProjects(prev =>
        reset ? res.items : [...prev, ...res.items]
      );
      setHasMore(res.hasMore);
      setOffset(prev => (reset ? PAGE_SIZE : prev + PAGE_SIZE));
    } finally {
      setLoading(false);
    }
  };

  /* Reset on search */
  useEffect(() => {
    setOffset(0);
    setHasMore(true);
    fetchRepositories(true);
  }, [debouncedSearch]);

  /* Infinite scroll observer */
  useEffect(() => {
    if (!loaderRef.current) return;

    const observer = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting) {
        fetchRepositories();
      }
    });

    observer.observe(loaderRef.current);
    return () => observer.disconnect();
  }, [loaderRef.current, debouncedSearch, hasMore]);

  return (
    <div className="projects-page">
      <div className="projects-topbar">
        <div />
        <div className="projects-topbar-right">
          <FilterBar
            filters={[
              {
                key: "search",
                type: "text",
                placeholder: "Search repositories…",
                value: search,
                onChange: setSearch,
              },
            ]}
                  actions={
              <button className="projects-add-btn">
                Add New…
              </button>
            }
          />
        </div>
      </div>

      <div className="projects-grid">
        <div className="projects-column">
          <WidgetCard title="Last 30 days" actionLabel="Upgrade">
            <ul className="usage-list">
              <li><span>Edge Requests</span><span>15 / 1M</span></li>
              <li><span>Edge Request CPU Duration</span><span>0s / 1h</span></li>
              <li><span>Fast Data Transfer</span><span>166.7 KB / 100 GB</span></li>
            </ul>
          </WidgetCard>

          <WidgetCard title="Alerts" actionLabel="Upgrade to Observability Plus">
            <p className="widget-text">
              Automatically monitor your projects for anomalies.
            </p>
          </WidgetCard>

          <WidgetCard title="Recent Previews">
            <p className="widget-text">
              Your recent deployments will appear here.
            </p>
          </WidgetCard>
        
        </div>

        <div className="projects-column">
          {projects.map(repo => (
            <ProjectCard
              key={repo.id}
              icon="📦"
              title={repo.name}
              repoName={repo.fullName}
              repoLink={repo.htmlUrl}
              dateRange={new Date(repo.updatedAt).toLocaleDateString()}
              branch={repo.defaultBranch}
              status={repo.private ? "🔒" : "◉"}
              stars={repo.stars}
              forks={repo.forks}
            />
          ))}

          {showLoadingMessage && <p className="widget-text">Loading…</p>}

          {/* Sentinel for infinite scroll */}
          <div ref={loaderRef} style={{ height: 1 }} />
        </div>
      </div>
    </div>
  );
}
