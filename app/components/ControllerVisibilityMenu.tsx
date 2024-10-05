'use client';

import React, { useState } from 'react';
import { Settings } from 'lucide-react';

interface ControllerVisibilityMenuProps {
  showUnknown: boolean;
  toggleUnknown: () => void;
}

const ControllerVisibilityMenu: React.FC<ControllerVisibilityMenuProps> = ({ showUnknown, toggleUnknown }) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center space-x-2 text-white bg-gray-800 hover:bg-gray-700 px-4 py-2 rounded-md"
      >
        <Settings size={20} />
        <span>Settings</span>
      </button>
      {isOpen && (
        <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1">
          <button
            onClick={() => {
              toggleUnknown();
              setIsOpen(false);
            }}
            className="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
          >
            <input
              type="checkbox"
              checked={showUnknown}
              onChange={toggleUnknown}
              className="mr-2"
            />
            Show Unknown Controllers
          </button>
        </div>
      )}
    </div>
  );
};

export default ControllerVisibilityMenu;