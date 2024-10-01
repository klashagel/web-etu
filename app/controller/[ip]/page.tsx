'use client';

import React, { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';

const ControllerDetail = () => {
  const params = useParams();
  const [controller, setController] = useState(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchControllerDetails = async () => {
      if (!params || !params.ip) {
        setError('No IP address provided');
        return;
      }

      try {
        const response = await fetch(`http://192.168.1.38:12345/etu/getcontroller/${params.ip}`);
        if (!response.ok) {
          throw new Error('Failed to fetch controller details');
        }
        const data = await response.json();
        setController(data);
      } catch (error) {
        console.error('Error fetching controller details:', error);
        setError('Failed to fetch controller details');
      }
    };

    fetchControllerDetails();
  }, [params]);

  if (error) {
    return <div className="container mx-auto p-4 text-red-500">{error}</div>;
  }

  if (!controller) {
    return <div className="container mx-auto p-4">Loading...</div>;
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Controller Detail</h1>
      <p>IP: {controller.IP}</p>
      <p>Type: {controller.CTRLTYPE}</p>
      {/* Add more details here */}
    </div>
  );
};

export default ControllerDetail;