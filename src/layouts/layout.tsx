import React, { useState } from "react";
import type {ReactNode} from "react";
import "./Layout.css";

interface LayoutProps {
  children: ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [isSidebarOpen, setSidebarOpen] = useState(true);

  const toggleSidebar = () => setSidebarOpen(!isSidebarOpen);

  return (
    <div className="layout">
      <header className="layout-topbar">
         <button className={`hamburger ${isSidebarOpen ? "open" : ""}`} onClick={toggleSidebar}>
          <span></span>
          <span></span>
          <span></span>
        </button>
        <h1>Header</h1>
      </header>

      <div className="layout-body">
        <aside className={`layout-sidebar ${isSidebarOpen ? "open" : "closed"}`}>
          <ul>
            <li>Tool 1</li>
            <li>Tool 2</li>
            <li>Tool 3</li>
          </ul>
        </aside>

        <main
          className={`layout-content ${isSidebarOpen ? "sidebar-open" : "sidebar-closed"}`}
        >
          {children}
        </main>
      </div>
    </div>
  );
};

export default Layout;
