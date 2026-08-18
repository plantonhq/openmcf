'use client';

import { Box } from '@mui/material';
import { useRouter } from 'next/navigation';
import { useEffect, useCallback, PropsWithChildren } from 'react';

export default function FocusLayout({ children }: PropsWithChildren) {
  const router = useRouter();

  const handleClose = useCallback(() => {
    router.push('/');
  }, [router]);

  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') {
        handleClose();
      }
    }
    window.addEventListener('keydown', onKeyDown);
    return () => window.removeEventListener('keydown', onKeyDown);
  }, [handleClose]);

  return (
    <Box className="relative min-h-screen bg-[#0a0a0a]">
      <button
        onClick={handleClose}
        aria-label="Close"
        className="fixed top-5 right-5 z-50 w-10 h-10 rounded-full border border-[#2a2a2a] bg-[#111] flex items-center justify-center text-[#a0a0a0] hover:text-white hover:border-[#3a3a3a] transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/20"
      >
        <svg
          width="18"
          height="18"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M18 6L6 18M6 6l12 12" />
        </svg>
      </button>
      {children}
    </Box>
  );
}
