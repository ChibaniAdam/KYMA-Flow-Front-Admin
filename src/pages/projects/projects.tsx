import { ProjectCard } from "../../components/project-card/project-card";
import WidgetCard from "../../components/widget-card/widget-card";
import "./projects.css";

export default function Projects() {
  return (
    <div className="projects-page">

      <div className="projects-topbar">
        <input
          className="projects-search"
          type="text"
          placeholder="Search Projects..."
        />
        <div className="projects-actions">
          <button className="projects-view-btn">☰</button>
          <button className="projects-view-btn">⬛</button>
          <button className="projects-add-btn">Add New…</button>
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

          <WidgetCard
            title="Alerts"
            actionLabel="Upgrade to Observability Plus"
          >
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
          <ProjectCard
            icon="⚛️"
            title="KYMA Flow"
            repoName="ChibaniAdam/KYMA-Flow-Front-Admin"
            repoLink="https://github.com/ChibaniAdam/KYMA-Flow-Front-Admin"
            dateRange="Oct 28"
            branch="master"
            status="◉"
          />
        </div>

      </div>
    </div>
  );
}
