import { Outlet } from "react-router-dom";
import { Header } from "../components/header/header";
import { useState } from "react";
import "./layout.css";

export default function Layout() {
  const [isSidebarOpen, setSidebarOpen] = useState(false);

  const toggleSidebar = () => setSidebarOpen(!isSidebarOpen);

  return (
    <div className="layout">
      <Header/>

      <aside className={`layout-sidebar ${isSidebarOpen ? "open" : "closed"}`}>
        <div className="sidebar-toggle" onClick={toggleSidebar}>
          <span className={`arrow ${isSidebarOpen ? "open" : ""}`}>➤</span>
        </div>
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
