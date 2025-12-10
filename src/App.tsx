import {BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import './App.css'
import Login from './pages/auth/login/login';
import Dashboard from './pages/dashboard/dashboard';
import Layout from './layouts/layout';
import Projects from './pages/projects/projects';

function App() {

  return (
    <Router>
        <Routes>
          <Route path="/" element={<Navigate to="/login" replace />} />
          <Route path="/login" element={<Login />} />
          <Route element={<Layout />}>
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/projects" element={<Projects />} />
          </Route>
        </Routes>
    </Router>
  )
}

export default App
