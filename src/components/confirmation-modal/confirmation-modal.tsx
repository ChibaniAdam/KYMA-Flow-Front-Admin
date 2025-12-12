import "./confirmation-modal.css";

interface ConfirmationModalProps {
  title?: string;
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
}

export const ConfirmationModal = ({
  title = "Confirm Deletion",
  message,
  onConfirm,
  onCancel,
  loading = false,
}: ConfirmationModalProps) => {
  return (
    <div className="modal-backdrop" onClick={onCancel}>
      <div
        className="modal-card modal-3parts"
        onClick={(e) => e.stopPropagation()} 
      >
        <div className="modal-header">
          <h2>{title}</h2>
        </div>
        <div className="modal-body">
          <p>{message}</p>
        </div>
        <div className="modal-footer">
          <button className="cancel-btn" onClick={onCancel} disabled={loading}>
            Cancel
          </button>
          <button
            className="submit-btn"
            onClick={onConfirm}
            disabled={loading}
          >
            {loading ? "Deleting..." : "Delete"}
          </button>
        </div>
      </div>
    </div>
  );
};
