'use client';

import React from 'react';
import LoginPage from './components/LoginPage';
import Dashboard from './dashboard/page';
import { useAppContext } from './contexts/AppContext';

export default function Home() {
  const { user, isLoading } = useAppContext();

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return user ? <Dashboard /> : <LoginPage />;
}
