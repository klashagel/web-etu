import React from 'react';

interface UnknownProps {
  controller: any; // Replace 'any' with a more specific type if possible
  onClick: () => void;
  onDoubleClick: () => void;
  isSelected: boolean;
}

const UnknownController: React.FC<UnknownProps> = ({ controller, onClick, onDoubleClick, isSelected }) => {
  return (
    <div 
      className={`p-4 rounded-lg shadow hover:shadow-md transition-shadow ${
        isSelected ? 'bg-gray-300' : 'bg-gray-100'
      }`}
      onClick={onClick}
      onDoubleClick={onDoubleClick}
    >
      <h3 className="text-lg font-bold">Unknown</h3>
      <p>IP: {controller.IP}</p>
      <p>Last Active: {new Date(controller.LASTACTIVE).toLocaleString()}</p>
      {/* Add more generic details here */}
    </div>
  );
};

export default UnknownController;