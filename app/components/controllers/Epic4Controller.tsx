import React from 'react';

interface Epic4Props {
  controller: any; // Replace 'any' with a more specific type if possible
  onClick: () => void;
  onDoubleClick: () => void;
  isSelected: boolean;
}

const Epic4Controller: React.FC<Epic4Props> = ({ controller, onClick, onDoubleClick, isSelected }) => {
  return (
    <div 
      className={`w-full h-full p-4 rounded-lg shadow hover:shadow-md transition-shadow ${
        isSelected ? 'bg-blue-200' : 'bg-blue-100'
      }`}
      onClick={onClick}
      onDoubleClick={onDoubleClick}
    >
      <h3 className="text-lg font-bold">EPIC4</h3>
      <p>IP: {controller.IP}</p>
      <p>Last Active: {new Date(controller.LASTACTIVE).toLocaleString()}</p>
      {/* Add more EPIC4-specific details here */}
    </div>
  );
};

export default Epic4Controller;