'use client';

import React, { useEffect, useState } from 'react';
import Epic4Controller from '../components/controllers/Epic4Controller';
import Sir4Controller from '../components/controllers/Sir4Controller';
import UnknownController from '../components/controllers/UnknownController';
import DraggableController from '../components/DraggableController';
import UserMenu from '../components/UserMenu';
import { useRouter } from 'next/navigation';

interface Controller {
  LASTACTIVE: string;
  STATUS: number;
  IP: string;
  CTRLTYPE: string;
  Controller: any;
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

  const renderController = (controller: Controller, index: number) => {
    const isSelected = controller.IP === selectedController;
    const props = {
      controller,
      onClick: () => handleControllerClick(controller),
      onDoubleClick: () => handleControllerDoubleClick(controller),
      isSelected,
    };

    let ControllerComponent;
    switch (controller.CTRLTYPE) {
      case 'EPIC4':
        ControllerComponent = Epic4Controller;
        break;
      case 'SIR4':
        ControllerComponent = Sir4Controller;
        break;
      default:
        ControllerComponent = UnknownController;
    }

    // Calculate initial grid position with increased horizontal spacing
    const columns = 3; // Reduced number of columns for more horizontal space
    const itemWidth = 300; // Increased width
    const itemHeight = 200; // Keep the same height
    const horizontalSpacing = 50; // Increased horizontal spacing
    const verticalSpacing = 20; // Keep the same vertical spacing

    const initialX = (index % columns) * (itemWidth + horizontalSpacing);
    const initialY = Math.floor(index / columns) * (itemHeight + verticalSpacing);

    return (
      <DraggableController key={controller.IP} initialX={initialX} initialY={initialY} isSelected={isSelected}>
        <ControllerComponent {...props} />
      </DraggableController>
    );
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-gray-900 text-white p-4">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-xl sm:text-2xl font-bold">Controller Dashboard</h1>
          <UserMenu />
        </div>
      </header>
      <main className="container mx-auto px-4 py-8">
        <div className="relative h-[calc(100vh-100px)]">
          {controllers.map((controller, index) => renderController(controller, index))}
        </div>
      </main>
    </div>
  );
};

export default Dashboard;