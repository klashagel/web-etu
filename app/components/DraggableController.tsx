import React from 'react';
import Draggable from 'react-draggable';
import ResizableController from './ResizableController';

interface DraggableControllerProps {
  children: React.ReactNode;
  initialX: number;
  initialY: number;
  isSelected: boolean;
}

const DraggableController: React.FC<DraggableControllerProps> = ({ children, initialX, initialY, isSelected }) => {
  return (
    <Draggable
      defaultPosition={{ x: initialX, y: initialY }}
      bounds="parent"
      handle=".drag-handle"
    >
      <div className="absolute">
        <ResizableController minWidth={200} minHeight={150} isSelected={isSelected}>
          {children}
          {isSelected && <div className="drag-handle absolute top-0 left-0 w-full h-6 bg-gray-200 opacity-50 cursor-move" />}
        </ResizableController>
      </div>
    </Draggable>
  );
};

export default DraggableController;