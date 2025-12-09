import {BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import './App.css'
import Login from './pages/auth/login/login';
import Dashboard from './pages/dashboard/dashboard';
import Layout from './layouts/layout';

function App() {

  return (
    <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/dashboard" element={<Layout><Dashboard /></Layout>} />
        </Routes>
    </Router>
  )
}

export default App
