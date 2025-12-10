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

  const [showRightArrow, setShowRightArrow] = useState(false);
  const [showLeftArrow, setShowLeftArrow] = useState(false);

  const [underline, setUnderline] = useState({
    left: 0,
    width: 0,
    transition: ""
  });

  useEffect(() => {
    if (!navRef.current) return;

    const el = navRef.current;

    const checkScroll = () => {
      const isOverflowing = el.scrollWidth > el.clientWidth;
      
      if (!isOverflowing) {
        setShowRightArrow(false);
        setShowLeftArrow(false);
        return;
      }

      setShowLeftArrow(el.scrollLeft > 10);
      setShowRightArrow(el.scrollLeft < el.scrollWidth - el.clientWidth - 10);
    };

    checkScroll();
    el.addEventListener("scroll", checkScroll);
    window.addEventListener("resize", checkScroll);

    return () => {
      el.removeEventListener("scroll", checkScroll);
      window.removeEventListener("resize", checkScroll);
    };
  }, []);

  const scrollToEnd = () => {
    if (!navRef.current) return;
    navRef.current.scrollTo({
      left: navRef.current.scrollWidth,
      behavior: "smooth"
    });
  };

  const scrollToStart = () => {
    if (!navRef.current) return;
    navRef.current.scrollTo({
      left: 0,
      behavior: "smooth"
    });
  };

  /* Underline animation logic */
  useEffect(() => {
    if (!navRef.current) return;

    const activeLink =
      navRef.current.querySelector<HTMLAnchorElement>("a.active");
    if (!activeLink) return;

    const { offsetLeft, offsetWidth } = activeLink;

    setUnderline(() => ({
      left: offsetLeft,
      width: offsetWidth,
      transition: "all 0.3s ease"
    }));
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

      <div className="topbar-menu-wrapper">

        {showLeftArrow && (
          <button className="scroll-arrow left" onClick={scrollToStart}>
            ◀
          </button>
        )}

        <nav className="topbar-menu" ref={navRef}>
          <NavLink to="/projects">Projects</NavLink>
          <NavLink to="/dashboard">Dashboard</NavLink>
          <NavLink to="/deployments">Deployments</NavLink>
          <NavLink to="/activity">Activity</NavLink>
          <NavLink to="/domains">Domains</NavLink>
          <NavLink to="/usage">Usage</NavLink>
          <NavLink to="/support">Support</NavLink>
          <NavLink to="/settings">Settings</NavLink>

          <span
            className="magic-underline"
            style={{
              left: underline.left,
              width: underline.width,
              transition: underline.transition
            }}
          />
        </nav>

        {showRightArrow && (
          <button className="scroll-arrow right" onClick={scrollToEnd}>
            ▶
          </button>
        )}

      </div>
    </header>
  );
}
