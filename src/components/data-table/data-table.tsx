import { useMemo, useState } from "react";
import "./data-table.css";

interface Column<T> {
  key: keyof T;
  header: string;
  sortable?: boolean;
  render?: (row: T) => React.ReactNode;
}

interface DataTableProps<T extends { uid: string }> {
  columns: Column<T>[];
  data: T[];
  loading?: boolean;
  emptyMessage?: string;

  onEdit?: (row: T) => void;
  onDelete?: (row: T) => void;

  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;

  selectable?: boolean;
  selectedRows?: string[];
  onSelectionChange?: (ids: string[]) => void;
}

export function DataTable<T extends { uid: string }>({
  columns,
  data,
  loading,
  emptyMessage = "No data found.",
  onEdit,
  onDelete,
  page,
  pageSize = 10,
  total,
  onPageChange,
  selectable = false,
  selectedRows = [],
  onSelectionChange,
}: Readonly<DataTableProps<T>>) {
  /* ───────────── Sorting ───────────── */
  const [sortKey, setSortKey] = useState<keyof T | null>(null);
  const [sortDir, setSortDir] = useState<"asc" | "desc">("asc");

  const sortedData = useMemo(() => {
    if (!sortKey) return data;

    return [...data].sort((a, b) => {
      const va = a[sortKey];
      const vb = b[sortKey];

      if (va == null || vb == null) return 0;
      if (va < vb) return sortDir === "asc" ? -1 : 1;
      if (va > vb) return sortDir === "asc" ? 1 : -1;
      return 0;
    });
  }, [data, sortKey, sortDir]);

  const handleSort = (key: keyof T) => {
    if (sortKey === key) {
      setSortDir((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setSortKey(key);
      setSortDir("asc");
    }
  };

  /* ───────────── Pagination ───────────── */
  const totalPages = Math.ceil(total / pageSize);
  const paginatedData = data;



  /* ───────────── Selection ───────────── */
  const toggleRow = (id: string) => {
    if (!onSelectionChange) return;

    onSelectionChange(
      selectedRows.includes(id)
        ? selectedRows.filter((r) => r !== id)
        : [...selectedRows, id]
    );
  };

  const toggleAll = () => {
    if (!onSelectionChange) return;

    if (selectedRows.length === paginatedData.length) {
      onSelectionChange([]);
    } else {
      onSelectionChange(paginatedData.map((r) => r.uid));
    }
  };

  /* ───────────── Loading ───────────── */
  if (loading) {
    return (
      <div className="skeleton-table">
        {Array.from({ length: pageSize }).map((_, i) => (
          <div key={i} className="skeleton-row" />
        ))}
      </div>
    );
  }

  return (
    <>
      <table className="data-table">
        <thead>
          <tr>
            {selectable && (
              <th>
                <input
                  type="checkbox"
                  checked={
                    paginatedData.length > 0 &&
                    selectedRows.length === paginatedData.length
                  }
                  onChange={toggleAll}
                />
              </th>
            )}

            {columns.map((c) => (
              <th
                key={String(c.key)}
                className={c.sortable ? "sortable" : ""}
                onClick={() => c.sortable && handleSort(c.key)}
              >
                {c.header}
                {sortKey === c.key && (
                  <span className="sort-indicator">
                    {sortDir === "asc" ? "▲" : "▼"}
                  </span>
                )}
              </th>
            ))}

            {(onEdit || onDelete) && <th />}
          </tr>
        </thead>

        <tbody>
          {paginatedData.map((row) => (
            <tr key={row.uid}>
              {selectable && (
                <td>
                  <input
                    type="checkbox"
                    checked={selectedRows.includes(row.uid)}
                    onChange={() => toggleRow(row.uid)}
                  />
                </td>
              )}

              {columns.map((c) => (
                <td key={String(c.key)}>
                  {c.render ? c.render(row) : String(row[c.key])}
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

          {paginatedData.length === 0 && (
            <tr>
              <td
                colSpan={
                  columns.length +
                  (selectable ? 1 : 0) +
                  (onEdit || onDelete ? 1 : 0)
                }
                className="empty-msg"
              >
                {emptyMessage}
              </td>
            </tr>
          )}
        </tbody>
      </table>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="table-pagination">
          <button disabled={page === 1} onClick={() => onPageChange(1)}>⏮</button>
          <button disabled={page === 1} onClick={() => onPageChange(page - 1)}>◀</button>

          <span>Page {page} / {totalPages}</span>

          <button
            disabled={page === totalPages}
            onClick={() => onPageChange(page + 1)}
          >
            ▶
          </button>
          <button
            disabled={page === totalPages}
            onClick={() => onPageChange(totalPages)}
          >
            ⏭
          </button>
        </div>
      )}

    </>
  );
}
