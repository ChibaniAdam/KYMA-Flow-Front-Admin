import { useEffect, useState, useRef } from "react";
import { ProjectCard } from "../../components/project-card/project-card";
import WidgetCard from "../../components/widget-card/widget-card";
import {
  listRepositories,
  searchRepositories,
  updateRepository,
  deleteRepository,
} from "../../services/repositoryService";
import type { Repository, RepositoryPage } from "../../GQL/models/repository";
import "./projects.css";
import { useDebounce } from "../../utils/useDebounce";
import { FilterBar } from "../../components/filter-bar/filter-bar";
import { useDelayedLoading } from "../../utils/useDelayedLoading";
import { Modal } from "../../components/modal/modal";
import { RepositoryForm } from "./repository-form/repository-form";
import { ConfirmationModal } from "../../components/confirmation-modal/confirmation-modal"; 

const PAGE_SIZE = 5;

export default function Projects() {
  const [repositories, setRepositories] = useState<Repository[]>([]);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const showLoadingMessage = useDelayedLoading(loading);
  const loaderRef = useRef<HTMLDivElement | null>(null);

  const [showModal, setShowModal] = useState(false);
  const [editingRepo, setEditingRepo] = useState<Repository | null>(null);
  const [formData, setFormData] = useState<Partial<Repository>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const [deleteRepo, setDeleteRepo] = useState<Repository | null>(null);
  const [deleting, setDeleting] = useState(false);

  /* ---------------- Fetch repositories ---------------- */
const fetchRepositories = async (reset = false) => {
  if (loading) return;

  setLoading(true);

  try {
    const currentOffset = reset ? 0 : offset;

    const res: RepositoryPage =
      debouncedSearch.trim().length > 0
        ? await searchRepositories(debouncedSearch, {
            limit: PAGE_SIZE,
            offset: currentOffset,
          })
        : await listRepositories({
            limit: PAGE_SIZE,
            offset: currentOffset,
          });

    setRepositories(prev => {
      const merged = reset ? res.items : [...prev, ...res.items];

      const uniqueMap = new Map<string | number, Repository>();
      merged.forEach(repo => uniqueMap.set(repo.id, repo));
      const uniqueRepos = Array.from(uniqueMap.values());

      return uniqueRepos;
    });

    setOffset(currentOffset + res.items.length);
    setHasMore(res.hasMore);
  } finally {
    setLoading(false);
  }
};



  /* ---------------- Reset on search ---------------- */
  useEffect(() => {
    setOffset(0);
    setHasMore(true);
    fetchRepositories(true);
  }, [debouncedSearch]);

  /* ---------------- Infinite scroll ---------------- */
useEffect(() => {
  const onScroll = () => {
    if (!hasMore || loading) return;
    if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 10) {
      fetchRepositories();
    }
  };

  window.addEventListener("scroll", onScroll);
  return () => window.removeEventListener("scroll", onScroll);
}, [hasMore, loading, debouncedSearch]);


  /* ---------------- Form Handlers ---------------- */
  const handleEditClick = (repo: Repository) => {
    setEditingRepo(repo);
    setFormData({
      name: repo.name,
      description: repo.description,
      private: repo.private,
    });
    setSubmitError(null);
    setShowModal(true);
  };

  const handleCreateClick = () => {
    setEditingRepo(null);
    setFormData({ name: "", description: "", private: false });
    setSubmitError(null);
    setShowModal(true);
  };

  const handleSubmit = async () => {
    if (!formData.name) {
      setSubmitError("Repository name is required");
      return;
    }

    setSubmitting(true);
    setSubmitError(null);

    try {
      if (editingRepo) {
        await updateRepository(editingRepo.owner.login, editingRepo.name);
      } else {
        // Create repo mutation/service here
      }
      setShowModal(false);
      fetchRepositories(true);
    } catch (err: any) {
      setSubmitError(err.message || "Failed to save repository");
    } finally {
      setSubmitting(false);
    }
  };

  /* ---------------- Delete Handlers ---------------- */
  const handleDeleteClick = (repo: Repository) => {
    setDeleteRepo(repo); // open confirmation modal
  };

  const confirmDelete = async () => {
    if (!deleteRepo) return;
    setDeleting(true);
    try {
      await deleteRepository(deleteRepo.owner.login, deleteRepo.name);
      setDeleteRepo(null);
      fetchRepositories(true);
    } finally {
      setDeleting(false);
    }
  };

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
            actions={<button className="projects-add-btn" onClick={handleCreateClick}>Add New…</button>}
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
            <p className="widget-text">Automatically monitor your projects for anomalies.</p>
          </WidgetCard>
          <WidgetCard title="Recent Previews">
            <p className="widget-text">Your recent deployments will appear here.</p>
          </WidgetCard>
        </div>

        <div className="projects-column">
          {repositories.map(repo => (
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
              onEdit={() => handleEditClick(repo)}
              onDelete={() => handleDeleteClick(repo)} 
            />
          ))}

          {showLoadingMessage && <p className="widget-text">Loading…</p>}
          {!showLoadingMessage && repositories.length===0 && <p className="widget-text">No repositories found.</p>}

          <div ref={loaderRef} style={{ height: 1 }} />
        </div>
      </div>

      {showModal && (
        <Modal
          title={editingRepo ? "Update Repository" : "Create Repository"}
          subtitle={editingRepo ? "Modify repository details below." : "Fill in the information below to add a new repository."}
          onClose={() => setShowModal(false)}
          footer={
            <>
              <button className="cancel-btn" onClick={() => setShowModal(false)} disabled={submitting}>Cancel</button>
              <button className="submit-btn" onClick={handleSubmit} disabled={submitting}>
                {submitting ? "Saving..." : editingRepo ? "Update" : "Create"}
              </button>
            </>
          }
        >
          {submitError && <div className="form-error">{submitError}</div>}
          <RepositoryForm formData={formData} setFormData={setFormData} />
        </Modal>
      )}

      {deleteRepo && (
        <ConfirmationModal
          message={`Are you sure you want to delete repository "${deleteRepo.fullName}"? This action cannot be undone.`}
          onCancel={() => setDeleteRepo(null)}
          onConfirm={confirmDelete}
          loading={deleting}
        />
      )}
    </div>
  );
}
