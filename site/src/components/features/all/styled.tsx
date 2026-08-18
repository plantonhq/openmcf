'use client';

import { Box } from '@mui/material';
import { styled } from '@mui/material/styles';

export const StyledGradientBorderBox = styled(Box)(() => ({
  position: 'relative',
  '&::before': {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    background: 'linear-gradient(180deg, #8A3391 0%, #03A587 100%)',
    borderRadius: 'inherit',
    padding: '3px',
    WebkitMask: 'linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0)',
    WebkitMaskComposite: 'destination-out',
    maskComposite: 'exclude',
    pointerEvents: 'none',
  },
}));

export const StyledGradientBorderBox2 = styled(Box)(() => ({
  position: 'relative',
  '&::before': {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    background: 'linear-gradient(180deg, #816FF2 0%, #104EAE 100%)',
    borderRadius: 'inherit',
    padding: '6px',
    WebkitMask: 'linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0)',
    WebkitMaskComposite: 'destination-out',
    maskComposite: 'exclude',
    pointerEvents: 'none',
  },
}));
