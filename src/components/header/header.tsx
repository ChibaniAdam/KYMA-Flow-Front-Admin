import { NavLink, useLocation } from "react-router-dom";
import "./header.css";
import { useEffect, useRef, useState } from "react";

export function Header({
  isSidebarOpen,
  toggleSidebar
}: {
  isSidebarOpen: boolean;
  toggleSidebar: () => void;
}) {
  const navRef = useRef<HTMLDivElement>(null);
  const location = useLocation();
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const [underline, setUnderline] = useState({ left: 0, width: 0, transition: "" });
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setMenuOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);
  useEffect(() => {
    if (!navRef.current) return;

    const activeLink = navRef.current.querySelector<HTMLAnchorElement>("a.active");
    if (!activeLink) return;

    const { offsetLeft, offsetWidth } = activeLink;

    setUnderline((prev: any) => {
      if (prev.left < offsetLeft) {
        return { left: prev.left, width: offsetWidth, transition: "all 0.3s ease" };
      } else {
        return { left: offsetLeft, width: offsetWidth, transition: "all 0.3s ease" };
      }
    });

    const timeout = setTimeout(() => {
      setUnderline({ left: offsetLeft-5, width: offsetWidth+10, transition: "all 0.3s ease" });
    }, 20);

    return () => clearTimeout(timeout);
  }, [location]);
  return (
    <header className="topbar-container">
      <div className="topbar-row">
        
        <div className="topbar-left">
          <button 
            className={`hamburger ${isSidebarOpen ? "open" : ""}`}
            onClick={toggleSidebar}
          >
            <span></span>
            <span></span>
            <span></span>
          </button>

          <div className="project">
            <div className="logo">▲</div>
            <span className="project-name">KYMA Flow</span>
          </div>
        </div>

        <div className="topbar-right" ref={menuRef}>
        <div 
          className="topbar-avatar" 
          onClick={() => setMenuOpen(!menuOpen)}
        />
          {menuOpen && (
          <div className="avatar-menu">
            <button className="menu-item">Account Settings</button>
            <button className="menu-item">Logout</button>
          </div>
        )}
      </div>
      </div>
       <nav className="topbar-menu" ref={navRef}>
        <NavLink to="/projects" className={({ isActive }) => (isActive ? "active" : "")}>Projects</NavLink>
        <NavLink to="/dashboard" className={({ isActive }) => (isActive ? "active" : "")}>Dashboard</NavLink>
        <NavLink to="/deployments" className={({ isActive }) => (isActive ? "active" : "")}>Deployments</NavLink>
        <NavLink to="/activity" className={({ isActive }) => (isActive ? "active" : "")}>Activity</NavLink>
        <NavLink to="/domains" className={({ isActive }) => (isActive ? "active" : "")}>Domains</NavLink>
        <NavLink to="/usage" className={({ isActive }) => (isActive ? "active" : "")}>Usage</NavLink>
        <NavLink to="/support" className={({ isActive }) => (isActive ? "active" : "")}>Support</NavLink>
        <NavLink to="/settings" className={({ isActive }) => (isActive ? "active" : "")}>Settings</NavLink>

        <span
          className="magic-underline"
          style={{
            left: underline.left,
            width: underline.width,
            transition: underline.transition,
          }}
        />
      </nav>
    </header>
  );
}
