'use client';

import { FC, ReactNode, useRef } from 'react';
import { Box, Typography } from '@mui/material';
import { motion, useInView } from 'framer-motion';

export interface FlowStep {
  icon: ReactNode;
  label: string;
  sublabel?: string;
}

interface FlowStepsProps {
  steps: FlowStep[];
  className?: string;
}

export const FlowSteps: FC<FlowStepsProps> = ({ steps, className = '' }) => {
  const ref = useRef<HTMLDivElement>(null);
  const isInView = useInView(ref, { once: true, margin: '-40px' });

  return (
    <Box ref={ref} className={`w-full ${className}`}>
      {/* Desktop: horizontal */}
      <Box className="hidden md:flex items-center justify-center">
        {steps.map((step, i) => (
          <Box key={step.label} className="flex items-center">
            <motion.div
              initial={{ opacity: 0, scale: 0.8 }}
              animate={isInView ? { opacity: 1, scale: 1 } : { opacity: 0, scale: 0.8 }}
              transition={{ duration: 0.4, delay: i * 0.15, ease: 'easeOut' }}
              className="flex flex-col items-center text-center"
            >
              <Box className="w-12 h-12 rounded-xl bg-white/10 border border-[#2a2a2a] flex items-center justify-center text-white mb-2">
                {step.icon}
              </Box>
              <Typography className="text-xs font-semibold text-white whitespace-nowrap">
                {step.label}
              </Typography>
              {step.sublabel && (
                <Typography className="text-[10px] text-[#666] mt-0.5 max-w-[100px]">
                  {step.sublabel}
                </Typography>
              )}
            </motion.div>
            {i < steps.length - 1 && (
              <motion.div
                initial={{ opacity: 0, scaleX: 0 }}
                animate={isInView ? { opacity: 1, scaleX: 1 } : { opacity: 0, scaleX: 0 }}
                transition={{ duration: 0.3, delay: i * 0.15 + 0.2 }}
                className="mx-3 flex items-center"
                style={{ originX: 0 }}
              >
                <Box className="w-8 lg:w-12 h-px bg-[#3a3a3a]" />
                <Box
                  className="w-0 h-0"
                  sx={{
                    borderTop: '4px solid transparent',
                    borderBottom: '4px solid transparent',
                    borderLeft: '6px solid #3a3a3a',
                  }}
                />
              </motion.div>
            )}
          </Box>
        ))}
      </Box>

      {/* Mobile: vertical */}
      <Box className="flex md:hidden flex-col items-center gap-1">
        {steps.map((step, i) => (
          <Box key={step.label} className="flex flex-col items-center">
            <motion.div
              initial={{ opacity: 0, y: 12 }}
              animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 12 }}
              transition={{ duration: 0.4, delay: i * 0.12 }}
              className="flex items-center gap-3"
            >
              <Box className="w-10 h-10 rounded-lg bg-white/10 border border-[#2a2a2a] flex items-center justify-center text-white flex-shrink-0">
                {step.icon}
              </Box>
              <Box>
                <Typography className="text-xs font-semibold text-white">
                  {step.label}
                </Typography>
                {step.sublabel && (
                  <Typography className="text-[10px] text-[#666]">{step.sublabel}</Typography>
                )}
              </Box>
            </motion.div>
            {i < steps.length - 1 && (
              <Box className="w-px h-5 bg-[#2a2a2a] my-1" />
            )}
          </Box>
        ))}
      </Box>
    </Box>
  );
};
