import { Outlet } from "react-router-dom";
import { Header } from "../components/header/header";
import { useState } from "react";
import "./layout.css";

export default function Layout() {
  const [isSidebarOpen, setSidebarOpen] = useState(true);

  const toggleSidebar = () => setSidebarOpen(!isSidebarOpen);

  return (
    <div className="layout">
      <Header
        isSidebarOpen={isSidebarOpen}
        toggleSidebar={toggleSidebar}
      />

      <div className="layout-body">
        <aside className={`layout-sidebar ${isSidebarOpen ? "open" : "closed"}`}>
          <ul>
            <li>Tool 1</li>
            <li>Tool 2</li>
            <li>Tool 3</li>
          </ul>
        </aside>

        <main className={`layout-content ${isSidebarOpen ? "sidebar-open" : "sidebar-closed"}`}>
          <Outlet /> 
        </main>
      </div>
    </div>
  );
}
