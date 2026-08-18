'use client';

import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';

interface CloudProviderModalProps {
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
}

export const CloudProviderModal: React.FC<CloudProviderModalProps> = ({
  isOpen,
  onClose,
  children,
}) => {
  const handleModalClose = () => {
    onClose();
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
          onClick={handleModalClose}
        >
          <motion.div
            initial={{ scale: 0.95, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.95, opacity: 0 }}
            className="bg-white rounded-2xl shadow-2xl w-full h-full max-w-7xl max-h-[95vh] overflow-hidden"
            onClick={(e) => e.stopPropagation()}
          >
            {children}
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};
