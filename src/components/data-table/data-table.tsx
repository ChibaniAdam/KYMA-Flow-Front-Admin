import "./data-table.css";

interface Column<T> {
  key: string;
  header: string;
  render?: (row: T) => React.ReactNode;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  loading?: boolean;
  emptyMessage?: string;
  onEdit?: (row: T) => void;
  onDelete?: (row: T) => void;
}

export function DataTable<T extends { uid: string }>({
  columns,
  data,
  loading,
  emptyMessage = "No data found.",
  onEdit,
  onDelete,
}: Readonly<DataTableProps<T>>) {
  if (loading) {
    return (
      <div className="skeleton-table">
        {Array.from({ length: 8 }).map((_, i) => (
          <div key={i} className="skeleton-row" />
        ))}
      </div>
    );
  }

  return (
    <table className="data-table">
      <thead>
        <tr>
          {columns.map((c) => (
            <th key={c.key}>{c.header}</th>
          ))}
          {(onEdit || onDelete) && <th />}
        </tr>
      </thead>

      <tbody>
        {data.map((row) => (
          <tr key={row.uid}>
            {columns.map((c) => (
              <td key={c.key}>
                {c.render ? c.render(row) : (row as any)[c.key]}
              </td>
            ))}

            {(onEdit || onDelete) && (
              <td className="actions">
                {onEdit && (
                  <button
                    className="update-btn"
                    onClick={() => onEdit(row)}
                  >
                    Edit
                  </button>
                )}
                {onDelete && (
                  <button
                    className="delete-btn"
                    onClick={() => onDelete(row)}
                  >
                    Delete
                  </button>
                )}
              </td>
            )}
          </tr>
        ))}

        {data.length === 0 && (
          <tr>
            <td colSpan={columns.length + 1} className="empty-msg">
              {emptyMessage}
            </td>
          </tr>
        )}
      </tbody>
    </table>
  );
}
