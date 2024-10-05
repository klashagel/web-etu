'use client';

import React from 'react';
import { useControllersData } from '../contexts/ControllersDataContext';

const SomeComponent: React.FC = () => {
  const { data, loading, error, refresh, getControllerField } = useControllersData();

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <h2>Controllers Data:</h2>
      <button onClick={refresh}>Refresh Data</button>
      <ul>
        {data?.map((controller, index) => (
          <li key={index}>
            IP: {controller.Controller.ip}
            {/* Access other fields as needed */}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default SomeComponent;