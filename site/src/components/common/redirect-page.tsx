'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Box, Typography } from '@mui/material';

export const RedirectPage = ({ to }: { to: string }) => {
  const router = useRouter();

  useEffect(() => {
    router.replace(to);
  }, [router, to]);

  return (
    <Box className="min-h-screen flex items-center justify-center bg-[#0a0a0a]">
      <Typography className="text-[#666] text-sm">Redirecting...</Typography>
    </Box>
  );
};
