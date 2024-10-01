import React from 'react';
import Draggable from 'react-draggable';
import ResizableController from './ResizableController';

interface DraggableControllerProps {
  children: React.ReactNode;
  initialX: number;
  initialY: number;
}

const DraggableController: React.FC<DraggableControllerProps> = ({ children, initialX, initialY }) => {
  return (
    <Draggable
      defaultPosition={{ x: initialX, y: initialY }}
      bounds="parent"
    >
      <div className="absolute">
        <ResizableController minWidth={200} minHeight={150}>
          {children}
        </ResizableController>
      </div>
    </Draggable>
  );
};

export default DraggableController;