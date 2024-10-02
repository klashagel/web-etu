import React, { useState, useRef, useEffect } from 'react';

interface ResizableControllerProps {
  children: React.ReactNode;
  minWidth: number;
  minHeight: number;
  isSelected: boolean;
}

const ResizableController: React.FC<ResizableControllerProps> = ({ children, minWidth, minHeight, isSelected }) => {
  const [size, setSize] = useState({ width: 300, height: 200 });
  const [isResizing, setIsResizing] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const resizeHandleRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!isResizing || !ref.current || !resizeHandleRef.current) return;

      const newWidth = e.clientX - ref.current.getBoundingClientRect().left;
      const newHeight = e.clientY - ref.current.getBoundingClientRect().top;

      setSize({
        width: Math.max(newWidth, minWidth),
        height: Math.max(newHeight, minHeight)
      });
    };

    const handleMouseUp = () => {
      setIsResizing(false);
    };

    if (isResizing) {
      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
    }

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [isResizing, minWidth, minHeight]);

  const startResize = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsResizing(true);
  };

  return (
    <div
      ref={ref}
      className="relative"
      style={{ width: `${size.width}px`, height: `${size.height}px` }}
    >
      {children}
      {isSelected && (
        <div
          ref={resizeHandleRef}
          className="absolute bottom-0 right-0 w-4 h-4 bg-gray-400 cursor-se-resize"
          onMouseDown={startResize}
          style={{ touchAction: 'none' }}
        />
      )}
    </div>
  );
};

export default ResizableController;