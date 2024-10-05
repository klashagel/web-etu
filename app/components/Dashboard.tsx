'use client';

import React from 'react';
import { useAppContext } from '../contexts/AppContext';
import { useConfig } from '../contexts/ConfigContext';

const Dashboard: React.FC = () => {
  const { user } = useAppContext();
  const { config, updateConfig } = useConfig();

  const toggleTheme = () => {
    updateConfig({ theme: config.theme === 'light' ? 'dark' : 'light' });
  };

  return (
    <div className={`p-4 ${config.theme === 'dark' ? 'bg-gray-800 text-white' : 'bg-white text-black'}`}>
      <h1 className="text-2xl font-bold mb-4">Welcome to the Dashboard</h1>
      <p className="mb-4">Hello, {user?.username}!</p>
      <p className="mb-4">API URL: {config.apiUrl}</p>
      <p className="mb-4">Max Controllers: {config.maxControllers}</p>
      <button
        onClick={toggleTheme}
        className={`px-4 py-2 rounded ${config.theme === 'dark' ? 'bg-white text-black' : 'bg-gray-800 text-white'}`}
      >
        Toggle Theme
      </button>
    </div>
  );
};

export default Dashboard;