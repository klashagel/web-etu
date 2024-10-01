import React from 'react';

interface Sir4Props {
  controller: any; // Replace 'any' with a more specific type if possible
  onClick: () => void;
  onDoubleClick: () => void;
  isSelected: boolean;
}

const Sir4Controller: React.FC<Sir4Props> = ({ controller, onClick, onDoubleClick, isSelected }) => {
  return (
    <div 
      className={`w-full h-full p-4 rounded-lg shadow hover:shadow-md transition-shadow ${
        isSelected ? 'bg-green-200' : 'bg-green-100'
      }`}
      onClick={onClick}
      onDoubleClick={onDoubleClick}
    >
      <h3 className="text-lg font-bold">SIR4</h3>
      <p>IP: {controller.IP}</p>
      <p>Last Active: {new Date(controller.LASTACTIVE).toLocaleString()}</p>
      {/* Add more SIR4-specific details here */}
    </div>
  );
};

export default Sir4Controller;