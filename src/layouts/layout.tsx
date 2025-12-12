import { Outlet } from "react-router-dom";
import { Header } from "../components/header/header";
import { useState } from "react";
import "./layout.css";

export default function Layout() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const toggleSidebar = () => setIsSidebarOpen(prev => !prev);

  return (
    <div className="layout">
      <Header/>

      <aside className={`layout-sidebar ${isSidebarOpen ? "open" : "closed"}`}>
        <button
          type="button"
          className="sidebar-toggle"
          onClick={toggleSidebar}
          aria-pressed={isSidebarOpen}
          aria-label="Toggle sidebar"
        >
          <span className={`arrow ${isSidebarOpen ? "open" : ""}`}>➤</span>
        </button>
        <ul>
          <li>Tool 1</li>
          <li>Tool 2</li>
          <li>Tool 3</li>
        </ul>
      </aside>

      <main className="layout-content">
        <Outlet />
      </main>
    </div>
  );
}
