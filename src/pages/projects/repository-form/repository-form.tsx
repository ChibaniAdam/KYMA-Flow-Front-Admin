import type { Repository } from "../../../GQL/models/repository";

interface Props {
  formData: Partial<Repository>;
  setFormData: (data: Partial<Repository>) => void;
}

export const RepositoryForm = ({ formData, setFormData }: Props) => {
  return (
    <form className="repository-form" onSubmit={e => e.preventDefault()}>
      <div className="form-group">
        <label htmlFor="name">Name</label>
        <input
          id="name"
          type="text"
          value={formData.name || ""}
          onChange={e => setFormData({ ...formData, name: e.target.value })}
          required
        />
      </div>

      <div className="form-group">
        <label htmlFor="description">Description</label>
        <textarea
          id="description"
          value={formData.description || ""}
          onChange={e => setFormData({ ...formData, description: e.target.value })}
        />
      </div>

      <div className="form-group">
        <label>
          <input type="checkbox" checked={formData.private || false} onChange={e => setFormData({ ...formData, private: e.target.checked })}/> Private
        </label>
      </div>
    </form>
  );
};
