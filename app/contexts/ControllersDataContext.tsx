'use client';

import React, { createContext, useState, useEffect, useCallback, useMemo, useRef, useContext } from 'react';
import axios from 'axios';
import { useConfig } from './ConfigContext';
import { useWebSocket } from './WebSocketContext';
import { EVENT_ID_CONTROLLERS_UPDATED } from '../constants';

interface Controller {
  Controller: {
    ip: string;
    fields: any;
  };
}

interface ControllersDataContextType {
  data: Controller[] | null;
  loading: boolean;
  error: Error | null;
  refresh: () => Promise<void>;
  getControllerField: (ip: string | null, registerPath: string[]) => any;
}

export const ControllersDataContext = createContext<ControllersDataContextType | undefined>(undefined);

export const ControllersDataProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { config } = useConfig();
  const { messages } = useWebSocket();
  const [data, setData] = useState<Controller[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const initialMount = useRef(true);
  const processedMessageIdsRef = useRef(new Set<string>());

  const fetchData = useCallback(async () => {
    if (!config.apiUrl) {
      console.error("Error fetching data: REST URL is undefined");
      setError(new Error("REST URL is undefined"));
      setLoading(false);
      return;
    }

    setLoading(true);
    try {
      const response = await axios.get(`${config.apiUrl}/etu/getallcontrollers`);
      setData(response.data);
      setError(null);
    } catch (err) {
      console.error("Error fetching data: ", err);
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, [config.apiUrl]);

  useEffect(() => {
    if (initialMount.current) {
      initialMount.current = false;
      fetchData();
    }
  }, [fetchData]);

  useEffect(() => {
    messages.forEach((message: any) => {
      if (message.eventid === EVENT_ID_CONTROLLERS_UPDATED) {
        if (!processedMessageIdsRef.current.has(message.messageid)) {
          processedMessageIdsRef.current.add(message.messageid);
          fetchData();
        }
      }
    });
  }, [messages, fetchData]);

  const contextValue = useMemo(() => {
    const getControllerField = (ip: string | null, registerPath: string[]) => {
      if (!data || !Array.isArray(data)) return null;

      const filteredData = ip
        ? data.filter((item) => item.Controller && item.Controller.ip === ip)
        : data;

      if (filteredData.length > 0) {
        const controller = filteredData[0].Controller;
        if (controller && controller.fields) {
          return registerPath.reduce(
            (acc, key) => (acc ? acc[key] : null),
            controller.fields
          );
        }
      }
      return null;
    };

    return {
      data,
      loading,
      error,
      refresh: fetchData,
      getControllerField,
    };
  }, [data, loading, error, fetchData]);

  return (
    <ControllersDataContext.Provider value={contextValue}>
      {children}
    </ControllersDataContext.Provider>
  );
};

export const useControllersData = () => {
  const context = useContext(ControllersDataContext);
  if (context === undefined) {
    throw new Error('useControllersData must be used within a ControllersDataProvider');
  }
  return context;
};