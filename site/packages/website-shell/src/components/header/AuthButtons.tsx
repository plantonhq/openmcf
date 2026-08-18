'use client';

import type { FC } from 'react';
import Link from 'next/link';
import { Button } from '@mui/material';
import { useLoggedIn } from '../../hooks/useLoggedIn';

const ctaSx = {
  height: 32,
  px: { xs: 1, sm: 1.5 },
  py: 0.5,
  fontSize: { xs: '0.75rem', sm: '0.875rem' },
  fontWeight: 500,
  whiteSpace: 'nowrap',
  borderRadius: '10px',
  textTransform: 'none',
} as const;

const whiteButtonStyle: React.CSSProperties = {
  backgroundColor: '#ffffff',
  color: '#000000',
};

export const DesktopAuthButtons: FC = () => {
  const loggedIn = useLoggedIn();

  if (loggedIn) {
    return (
      <Button LinkComponent={Link} href="/dashboard" style={whiteButtonStyle} sx={ctaSx}>
        Dashboard
      </Button>
    );
  }

  return (
    <>
      <Button
        LinkComponent={Link}
        href="/login"
        style={{ color: '#a0a0a0' }}
        sx={{
          display: { xs: 'none', sm: 'inline-flex' },
          fontSize: '0.875rem',
          fontWeight: 500,
          textTransform: 'none',
          borderRadius: '10px',
        }}
      >
        Sign in
      </Button>
      <Button LinkComponent={Link} href="/signup" style={whiteButtonStyle} sx={ctaSx}>
        Sign up
      </Button>
    </>
  );
};

const ctaFullWidthSx = {
  ...ctaSx,
  height: 40,
  width: '100%',
  justifyContent: 'center',
} as const;

export const MobileAuthButtons: FC = () => {
  const loggedIn = useLoggedIn();

  if (loggedIn) {
    return (
      <Button LinkComponent={Link} href="/dashboard" style={whiteButtonStyle} sx={ctaFullWidthSx}>
        Dashboard
      </Button>
    );
  }

  return (
    <>
      <Button
        LinkComponent={Link}
        href="/login"
        style={{ color: '#a0a0a0' }}
        sx={{
          width: '100%',
          justifyContent: 'center',
          textTransform: 'none',
          borderRadius: '10px',
          fontWeight: 500,
        }}
      >
        Sign in
      </Button>
      <Button LinkComponent={Link} href="/signup" style={whiteButtonStyle} sx={ctaFullWidthSx}>
        Sign up
      </Button>
    </>
  );
};
