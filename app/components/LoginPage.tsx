'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Moon, Sun } from 'lucide-react'; // Import icons

const LoginPage: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isDarkTheme, setIsDarkTheme] = useState(true); // State for theme
  const router = useRouter();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    // Temporary login logic for testing
    if (username === 'admin' && password === 'password123') {
      router.push('/dashboard');
    } else {
      setError('Invalid username or password');
    }

    // When implementing the actual API call, replace the above logic with:
    /*
    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });

      if (response.ok) {
        router.push('/dashboard');
      } else {
        const data = await response.json();
        setError(data.message || 'Login failed');
      }
    } catch (err) {
      setError('An error occurred. Please try again.');
    }
    */
  };

  const toggleTheme = () => {
    setIsDarkTheme(!isDarkTheme);
  };

  return (
    <div className={`flex h-screen ${isDarkTheme ? 'bg-gray-900 text-white' : 'bg-gray-100 text-black'}`}>
      {/* Theme toggle button */}
      <button
        onClick={toggleTheme}
        className="absolute top-4 right-4 p-2 rounded-full bg-opacity-20 bg-gray-600"
      >
        {isDarkTheme ? <Sun className="h-6 w-6" /> : <Moon className="h-6 w-6" />}
      </button>

      {/* Left side */}
      <div className="flex-1 flex flex-col justify-between p-8">
        <div>
          <h1 className="text-2xl font-bold">⌘ Andritz AB</h1>
        </div>
        <div>
          <blockquote className="text-xl">
            "This library has saved me countless hours of work and helped me
            deliver stunning designs to my clients faster than ever before."
          </blockquote>
          <p className="mt-4">Sofia Davis</p>
        </div>
      </div>

      {/* Right side */}
      <div className="flex-1 flex items-center justify-center">
        <div className="w-96">
          <h2 className="text-3xl font-bold mb-4">Login</h2>
          <p className={`mb-6 ${isDarkTheme ? 'text-gray-400' : 'text-gray-600'}`}>
            Enter your credentials to access your account
          </p>
          {error && <p className="text-red-500 mb-4">{error}</p>}
          <form onSubmit={handleLogin}>
            <input
              type="text"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className={`w-full p-3 mb-4 rounded ${isDarkTheme ? 'bg-gray-800' : 'bg-gray-200'}`}
              required
            />
            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={`w-full p-3 mb-4 rounded ${isDarkTheme ? 'bg-gray-800' : 'bg-gray-200'}`}
              required
            />
            <button
              type="submit"
              className={`w-full p-3 mb-4 rounded font-semibold ${
                isDarkTheme ? 'bg-white text-black' : 'bg-black text-white'
              }`}
            >
              Login
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;