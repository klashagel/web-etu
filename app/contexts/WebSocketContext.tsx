'use client';

import React, { createContext, useContext, useEffect, useState, useRef } from 'react';
import { useConfig } from './ConfigContext';
import { EVENT_ID_CONTROLLERS_UPDATED, WEBSOCKET_RECONNECT_INTERVAL } from '../constants';

interface WebSocketContextType {
  socket: WebSocket | null;
  messages: any[];
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined);

export const WebSocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { config } = useConfig();
  const [messages, setMessages] = useState<any[]>([]);
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const messageIdSet = useRef(new Set<string>());
  const queue = useRef<any[]>([]);

  useEffect(() => {
    if (!config.websocketUrl) {
      console.error("WebSocket URL is undefined");
      return;
    }

    const ws = new WebSocket(config.websocketUrl);
    setSocket(ws);

    ws.onmessage = (event) => {
      const rawData = event.data;
      console.log('Raw WebSocket data received:', rawData);

      try {
        const messagesArray = rawData.split('\n').filter(Boolean);
        messagesArray.forEach(message => {
          try {
            const parsedData = JSON.parse(message);
            console.log('Parsed WebSocket data:', parsedData);
            if (!messageIdSet.current.has(parsedData.messageid)) {
              messageIdSet.current.add(parsedData.messageid);
              queue.current.push(parsedData);
              console.log(`Queue length after push: ${queue.current.length}`);
              setMessages(prevMessages => [...prevMessages, parsedData]);
            }
            if (parsedData.eventid === EVENT_ID_CONTROLLERS_UPDATED) {
              // Handle the controllers updated event
            }
          } catch (parseError) {
            console.error('Error parsing JSON from message:', parseError);
            console.error('Message content:', message);
          }
        });
      } catch (err) {
        console.error('Error handling WebSocket message:', err);
        console.error('Received data:', rawData);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket Error: ', error);
    };

    return () => {
      ws.close();
    };
  }, [config.websocketUrl]);

  useEffect(() => {
    const interval = setInterval(() => {
      if (queue.current.length > 0) {
        console.log(`Queue length before processing: ${queue.current.length}`);
        setMessages(prevMessages => [...prevMessages, ...queue.current]);
        queue.current = [];
        console.log('Queue cleared');
      }
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  return (
    <WebSocketContext.Provider value={{ socket, messages }}>
      {children}
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (context === undefined) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  return context;
};