'use client';

import React, { createContext, useContext, useState, useEffect } from 'react';
import { ConfigProvider } from './ConfigContext';
import { WebSocketProvider } from './WebSocketContext';
import { ControllersDataProvider } from './ControllersDataContext';

interface User {
  username: string;
  // Add other user properties as needed
}

interface AppContextType {
  user: User | null;
  setUser: (user: User | null) => void;
  isLoading: boolean;
}

const AppContext = createContext<AppContextType | undefined>(undefined);

export const AppProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check for user session here
    // For now, we'll just set isLoading to false after a short delay
    const timer = setTimeout(() => setIsLoading(false), 1000);
    return () => clearTimeout(timer);
  }, []);

  return (
    <ConfigProvider>
      <WebSocketProvider>
        <ControllersDataProvider>
          <AppContext.Provider value={{ user, setUser, isLoading }}>
            {children}
          </AppContext.Provider>
        </ControllersDataProvider>
      </WebSocketProvider>
    </ConfigProvider>
  );
};

export const useAppContext = () => {
  const context = useContext(AppContext);
  if (context === undefined) {
    throw new Error('useAppContext must be used within an AppProvider');
  }
  return context;
};