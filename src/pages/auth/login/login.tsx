import { useState } from 'react';
import "./login.css"
import { useNavigate } from 'react-router-dom';
import { login } from '../../../services/graphqlService';

const Login = () => {
  const navigate = useNavigate();

  const [uid, setUid] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const result = await login({ uid, password });

      // Store the token
      localStorage.setItem('authToken', result.login.token);

      // Store user info
      localStorage.setItem('user', JSON.stringify(result.login.user));

      console.log('Login successful:', result.login.user);

      // Navigate to dashboard
      navigate("/dashboard");
    } catch (err: any) {
      console.error('Login failed:', err);
      setError(err.message || 'Login failed. Please check your credentials.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center ">
      <div className="bg-white p-8 rounded-lg shadow-lg w-full max-w-sm card">
        <h2 className="text-2xl font-bold mb-6 text-center text-black">Login</h2>

        {error && (
          <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <label className="block text-gray-700">Username (UID)</label>
            <input
              type="text"
              value={uid}
              onChange={(e) => setUid(e.target.value)}
              className="w-full px-4 py-2 border focus:outline-none focus:ring-2 focus:ring-blue-400"
              placeholder="Enter your username (e.g., john.doe)"
              required
              disabled={loading}
            />
          </div>
          <div>
            <label className="block text-gray-700">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 border focus:outline-none focus:ring-2 focus:ring-blue-400"
              placeholder="Enter your password"
              required
              disabled={loading}
            />
          </div>
          <button
            type="submit"
            className="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600 transition-colors submit-button disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={loading}
          >
            {loading ? 'Logging in...' : 'Login'}
          </button>
        </form>
        </div>
      </div>
  );
};

export default Login;
