import React, { useState } from "react";

interface TooltipProps {
  text: string;
  position?:
    | "top"
    | "bottom"
    | "left"
    | "right"
    | "top-left"
    | "top-right"
    | "bottom-left"
    | "bottom-right";
  color?: string;
  children: React.ReactNode;
}

const positionClasses: Record<string, string> = {
  top: "bottom-full left-1/2 -translate-x-1/2 mb-2",
  bottom: "top-full left-1/2 -translate-x-1/2 mt-2",
  left: "right-full top-1/2 -translate-y-1/2 mr-2",
  right: "left-full top-1/2 -translate-y-1/2 ml-2",
  "top-left": "bottom-full left-0 mb-2",
  "top-right": "bottom-full right-0 mb-2",
  "bottom-left": "top-full left-0 mt-2",
  "bottom-right": "top-full right-0 mt-2",
};

export const Tooltip: React.FC<TooltipProps> = ({ text, position = "top", color = "var(--dark-purple-500)", children }) => {
  const [visible, setVisible] = useState(false);

  return (
    <div
      className="relative inline-block"
      onMouseEnter={() => setVisible(true)}
      onMouseLeave={() => setVisible(false)}
    >
      {children}

      {visible && (
        <div
          style={{ backgroundColor: color, borderColor: color }}
          className={`absolute px-2 py-1 text-sm rounded-lg shadow-lg bg-gray-800 text-white whitespace-nowrap z-50 ${
            positionClasses[position]
          }`}
        >
          {text}
        </div>
      )}
    </div>
  );
};