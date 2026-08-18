'use client';

import React from 'react';
import { Button } from '@/components/tour/ui/button';
import { motion } from 'framer-motion';

type Callout = { 
  number: number; 
  x: number; 
  y: number; 
  title: string; 
  description: string; 
};

interface CalloutCardProps {
  callout: Callout;
  index: number;
  isVisible: boolean;
  isCurrent: boolean;
  currentCalloutIndex: number;
  totalCallouts: number;
  onNumberClick: (index: number) => void;
  onNext: () => void;
  onPrevious: () => void;
}

export default function CalloutCard({
  callout,
  index,
  isVisible,
  isCurrent,
  currentCalloutIndex,
  totalCallouts,
  onNumberClick,
  onNext,
  onPrevious,
}: CalloutCardProps) {
  return (
    <motion.div
      key={index}
      initial={{ opacity: 0, scale: 0.8 }}
      animate={{ opacity: isVisible ? 1 : 0, scale: isVisible ? 1 : 0.8 }}
      transition={{ duration: 0.3, delay: callout.number * 0.1 }}
      className={`absolute ${isCurrent ? 'z-50' : 'z-20'}`}
      style={{
        left: `${callout.x}%`,
        top: `${callout.y}%`,
        transform: 'translate(-50%, -50%)',
      }}
    >
      <div className="relative">
        <button
          onClick={() => onNumberClick(index)}
          className={`relative w-8 h-8 bg-white rounded-full shadow-lg border-2 border-[#2a2a2a] flex items-center justify-center hover:bg-gray-50 transition-colors cursor-pointer ${
            isCurrent ? 'z-40' : 'z-10'
          }`}
        >
          <span className="text-[#111] font-bold text-sm">{callout.number}</span>
        </button>
        <div className={`absolute inset-0 w-8 h-8 bg-[#666] rounded-full animate-ping opacity-20 ${
          isCurrent ? 'z-10' : 'z-5'
        }`}></div>
        
        {isVisible && isCurrent && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{
              opacity: 1,
              y: 0,
            }}
            transition={{ delay: callout.number * 0.1 + 0.2 }}
            className={`absolute w-72 bg-white rounded-xl shadow-xl border border-gray-100 p-4 z-[100] ${
              // Vertical positioning - very conservative to prevent clipping
              callout.y > 60 ? 'bottom-full mb-2' : 'top-full mt-2'
            } ${
              // Horizontal positioning - very conservative to prevent clipping
              callout.x > 70 ? 'right-4' : 
              callout.x < 30 ? 'left-4' : 
              'left-1/2 transform -translate-x-1/2'
            }`}
          >
          <div className={`absolute w-4 h-4 bg-white border-l border-t border-gray-100 rotate-45 ${
            callout.y > 60 ? 'bottom-[-8px]' : 'top-[-8px]'
          } ${
            callout.x > 70 ? 'right-4' : 
            callout.x < 30 ? 'left-4' : 
            'left-1/2 transform -translate-x-1/2'
          }`}></div>
          
          <h4 className="font-semibold text-gray-900 mb-2">{callout.title}</h4>
          <p className="text-sm text-gray-600 leading-relaxed mb-4">
            {callout.description}
          </p>
          
          {/* Only show navigation controls on the current card */}
          {isCurrent && (
            <div className="flex justify-between items-center">
              <span className="text-xs text-gray-500">
                {currentCalloutIndex + 1} of {totalCallouts}
              </span>
              <div className="flex gap-2">
                {currentCalloutIndex > 0 && (
                  <Button
                    onClick={onPrevious}
                    size="sm"
                    variant="outline"
                    className="text-gray-600 border-gray-300 hover:bg-gray-50"
                  >
                    Previous
                  </Button>
                )}
                <Button
                  onClick={onNext}
                  size="sm"
                  className="bg-[#111] hover:bg-[#222] text-white"
                >
                  {currentCalloutIndex < totalCallouts - 1 ? 'Next' : 'Continue'}
                </Button>
              </div>
            </div>
          )}
          </motion.div>
        )}
      </div>
    </motion.div>
  );
}
