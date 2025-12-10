import type { ReactNode } from "react";
import "./widget-card.css";

interface WidgetCardProps {
  title: string;
  children: ReactNode;
  actionLabel?: string;
  onAction?: () => void;
}

export default function WidgetCard({
  title,
  children,
  actionLabel,
  onAction,
}: WidgetCardProps) {
  return (
    <div className="widget-card">
      <h3 className="widget-title">{title}</h3>

      <div className="widget-content">{children}</div>

      {actionLabel && (
        <button className="upgrade-btn" onClick={onAction}>
          {actionLabel}
        </button>
      )}
    </div>
  );
}
