'use client';

import { FC } from 'react';
import { Box, BoxProps } from '@mui/material';

export const FrameImage: FC<BoxProps & { src: string }> = ({ src, className, ...props }) => {
  return (
    <Box className={`relative rounded-xl md:rounded-[42px] p-4 md:p-16 ${className}`} {...props}>
      <img src={src} className="rounded-md md:rounded-xl relative z-10" />
      <img src="/_site/images/frame.png" className="absolute top-0 right-0 bottom-0 left-0 z-0" />
    </Box>
  );
};
