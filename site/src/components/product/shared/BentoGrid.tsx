'use client';

import { FC, ReactNode } from 'react';
import { Box } from '@mui/material';

type BentoSpan = 'default' | 'wide' | 'tall';

interface BentoGridProps {
  children: ReactNode;
  className?: string;
}

export const BentoGrid: FC<BentoGridProps> = ({ children, className = '' }) => (
  <Box
    className={`grid grid-cols-1 md:grid-cols-2 gap-4 ${className}`}
    sx={{
      '& > *': { minHeight: 0 },
    }}
  >
    {children}
  </Box>
);

interface BentoItemProps {
  children: ReactNode;
  span?: BentoSpan;
  className?: string;
}

const spanClasses: Record<BentoSpan, string> = {
  default: '',
  wide: 'md:col-span-2',
  tall: 'md:row-span-2',
};

export const BentoItem: FC<BentoItemProps> = ({ children, span = 'default', className = '' }) => (
  <Box
    className={`
      rounded-xl bg-[#151515] border border-[#2a2a2a]
      p-5 md:p-6 overflow-hidden
      hover:border-[#3a3a3a] hover:bg-[#1a1a1a] transition-all duration-300
      ${spanClasses[span]}
      ${className}
    `}
  >
    {children}
  </Box>
);
