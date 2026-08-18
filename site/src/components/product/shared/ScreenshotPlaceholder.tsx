'use client';

import { FC } from 'react';
import { Box, Typography } from '@mui/material';

interface ScreenshotPlaceholderProps {
  caption?: string;
  aspectRatio?: string;
  className?: string;
}

export const ScreenshotPlaceholder: FC<ScreenshotPlaceholderProps> = ({
  caption,
  aspectRatio = '16 / 9',
  className = '',
}) => (
  <Box className={`max-w-4xl mx-auto ${className}`}>
    <Box className="rounded-xl border border-[#2a2a2a] overflow-hidden bg-[#111]">
      {/* Browser chrome */}
      <Box className="flex items-center gap-2 px-4 py-2.5 bg-[#1a1a1a] border-b border-[#2a2a2a]">
        <Box className="w-2.5 h-2.5 rounded-full bg-[#3a3a3a]" />
        <Box className="w-2.5 h-2.5 rounded-full bg-[#3a3a3a]" />
        <Box className="w-2.5 h-2.5 rounded-full bg-[#3a3a3a]" />
        <Box className="flex-1 mx-4">
          <Box className="max-w-xs mx-auto h-5 rounded-md bg-[#2a2a2a] flex items-center justify-center">
            <Typography className="text-[10px] text-[#555]">planton.ai</Typography>
          </Box>
        </Box>
      </Box>

      {/* Content area with dot grid */}
      <Box
        sx={{ aspectRatio }}
        className="relative"
      >
        <Box
          className="absolute inset-0"
          sx={{
            backgroundImage: 'radial-gradient(circle, #2a2a2a 1px, transparent 1px)',
            backgroundSize: '24px 24px',
          }}
        />
      </Box>
    </Box>
    {caption && (
      <Typography className="text-xs text-[#555] text-center mt-3">{caption}</Typography>
    )}
  </Box>
);
