'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Moon, Sun } from 'lucide-react'; // Import icons
import Image from 'next/image'; // Import the Next.js Image component

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
    <div className={`flex flex-col lg:flex-row min-h-screen ${isDarkTheme ? 'bg-gray-900 text-white' : 'bg-gray-100 text-black'}`}>
      {/* Theme toggle button */}
      <button
        onClick={toggleTheme}
        className="absolute top-4 right-4 p-2 rounded-full bg-opacity-20 bg-gray-600 z-20"
      >
        {isDarkTheme ? <Sun className="h-6 w-6" /> : <Moon className="h-6 w-6" />}
      </button>

      {/* Left side */}
      <div className="lg:flex-1 relative overflow-hidden">
        <div className="absolute inset-0 z-0">
          <Image
            src="/images/digital-brain-technology-wallpaper-0410-5697354-1.jpg"
            alt="Digital Brain Technology"
            layout="fill"
            objectFit="cover"
            className="opacity-50"
          />
        </div>
        <div className="relative z-10 p-8 lg:p-16 flex flex-col justify-between h-full">
          <h1 className="text-2xl lg:text-4xl font-bold mb-4">⌘ Andritz AB</h1>
          <div>
            <blockquote className="text-lg lg:text-xl italic">
              &ldquo;This library has saved me countless hours of work and helped me
              deliver stunning designs to my clients faster than ever before.&rdquo;
            </blockquote>
            <p className="mt-4 font-semibold">Sofia Davis</p>
          </div>
        </div>
      </div>

      {/* Right side */}
      <div className="lg:flex-1 flex items-center justify-center p-8 lg:p-16">
        <div className="w-full max-w-md">
          <h2 className="text-2xl lg:text-3xl font-bold mb-6">Login</h2>
          <p className={`mb-8 ${isDarkTheme ? 'text-gray-400' : 'text-gray-600'}`}>
            Enter your credentials to access your account
          </p>
          {error && <p className="text-red-500 mb-4">{error}</p>}
          <form onSubmit={handleLogin} className="space-y-6">
            <div>
              <label htmlFor="username" className="block mb-2 font-medium">Username</label>
              <input
                id="username"
                type="text"
                placeholder="Enter your username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className={`w-full p-3 rounded ${isDarkTheme ? 'bg-gray-800' : 'bg-white border border-gray-300'}`}
                required
              />
            </div>
            <div>
              <label htmlFor="password" className="block mb-2 font-medium">Password</label>
              <input
                id="password"
                type="password"
                placeholder="Enter your password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className={`w-full p-3 rounded ${isDarkTheme ? 'bg-gray-800' : 'bg-white border border-gray-300'}`}
                required
              />
            </div>
            <button
              type="submit"
              className={`w-full p-3 rounded font-semibold transition-colors ${
                isDarkTheme ? 'bg-blue-600 hover:bg-blue-700 text-white' : 'bg-blue-500 hover:bg-blue-600 text-white'
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