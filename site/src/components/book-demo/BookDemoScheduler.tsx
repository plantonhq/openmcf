'use client';

import { Box, Typography } from '@mui/material';
import { useEffect } from 'react';
import dynamic from 'next/dynamic';
import type { DemoFormData } from './types';

const Cal = dynamic(
  () => import('@calcom/embed-react').then((mod) => mod.default),
  {
    ssr: false,
    loading: () => (
      <Box className="w-full min-h-[500px] md:min-h-[600px] rounded-xl bg-[#111] border border-[#2a2a2a] flex items-center justify-center">
        <Box className="flex flex-col items-center gap-3">
          <Box className="w-8 h-8 border-2 border-[#2a2a2a] border-t-white/60 rounded-full animate-spin" />
          <Typography className="text-sm text-[#555]">
            Loading scheduler...
          </Typography>
        </Box>
      </Box>
    ),
  },
);

interface BookDemoSchedulerProps {
  formData: DemoFormData;
}

export function BookDemoScheduler({ formData }: BookDemoSchedulerProps) {
  useEffect(() => {
    (async function () {
      const { getCalApi: api } = await import('@calcom/embed-react');
      const cal = await api({ namespace: '60min' });
      cal('ui', { hideEventTypeDetails: false, layout: 'month_view' });
    })();
  }, []);

  return (
    <Box className="rounded-xl bg-[#111] border border-[#2a2a2a]">
      <Cal
        namespace="60min"
        calLink="swarup-donepudi/60min"
        style={{ width: '100%', minHeight: '600px' }}
        config={{
          layout: 'month_view',
          useSlotsViewOnSmallScreen: 'true',
          name: `${formData.firstName} ${formData.lastName}`.trim(),
          email: formData.workEmail,
        }}
      />
    </Box>
  );
}
