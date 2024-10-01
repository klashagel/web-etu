'use client';

import React, { useEffect, useState } from 'react';
import Epic4Controller from '../components/controllers/Epic4Controller';
import Sir4Controller from '../components/controllers/Sir4Controller';
import UnknownController from '../components/controllers/UnknownController';
import { useRouter } from 'next/navigation';

interface Controller {
  LASTACTIVE: string;
  STATUS: number;
  IP: string;
  CTRLTYPE: string;
  Controller: any; // Replace 'any' with a more specific type if possible
}

const Dashboard = () => {
  const [controllers, setControllers] = useState<Controller[]>([]);
  const [selectedController, setSelectedController] = useState<string | null>(null);
  const router = useRouter();

  useEffect(() => {
    const fetchControllers = async () => {
      try {
        const response = await fetch('http://192.168.1.38:12345/etu/getallcontrollers');
        const data = await response.json();
        setControllers(data);
      } catch (error) {
        console.error('Error fetching controllers:', error);
      }
    };

    fetchControllers();
  }, []);

  const handleControllerClick = (controller: Controller) => {
    setSelectedController(controller.IP === selectedController ? null : controller.IP);
  };

  const handleControllerDoubleClick = (controller: Controller) => {
    router.push(`/controller/${controller.IP}`);
  };

  const renderController = (controller: Controller) => {
    const props = {
      controller,
      onClick: () => handleControllerClick(controller),
      onDoubleClick: () => handleControllerDoubleClick(controller),
      isSelected: controller.IP === selectedController,
    };

    switch (controller.CTRLTYPE) {
      case 'EPIC4':
        return <Epic4Controller {...props} />;
      case 'SIR4':
        return <Sir4Controller {...props} />;
      default:
        return <UnknownController {...props} />;
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-2xl md:text-3xl lg:text-4xl font-bold mb-6">Controller Dashboard</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {controllers.map((controller, index) => (
          <div key={index} className="cursor-pointer">
            {renderController(controller)}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Dashboard;