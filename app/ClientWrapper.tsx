'use client';

import React from 'react';
import { AppProvider } from './contexts/AppContext';

const ClientWrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return <AppProvider>{children}</AppProvider>;
};

export default ClientWrapper;