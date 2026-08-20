'use client';

import { FC } from 'react';
import { Stack, Typography } from '@mui/material';
import { Badge } from '@/components/landing-page/v3-2026-01-02-1000/shared';
import { SELF_HOSTED_LICENSE_SIZES, SELF_SERVE_SEAT_CEILING } from '@/data/pricing';

/**
 * The enterprise page opens with the self-serve-first posture: most teams
 * never need this page, and saying so builds the trust the rate card
 * spends. The numbers come from the pricing-truth module.
 */
export const EnterpriseHero: FC = () => {
  const sizeLow = SELF_HOSTED_LICENSE_SIZES[0].usdPerYear;
  const sizeHigh = SELF_HOSTED_LICENSE_SIZES[SELF_HOSTED_LICENSE_SIZES.length - 1].usdPerYear;
  return (
    <Stack className="items-center gap-5 pt-14 pb-4 px-4 bg-[#0a0a0a] text-center">
      <Badge>Enterprise</Badge>
      <Typography
        variant="h1"
        className="text-3xl md:text-4xl lg:text-5xl font-semibold text-white leading-tight tracking-tight max-w-[900px]"
      >
        Enterprise at Planton
      </Typography>
      <Typography className="text-sm md:text-base text-[#a0a0a0] max-w-[720px]">
        {`Under ${SELF_SERVE_SEAT_CEILING} seats, you don't need to talk to us at all — the self-serve license is $${(sizeLow / 1000).toFixed(0)}K–$${(sizeHigh / 1000).toFixed(0)}K a year, card and email, running today. Enterprise adds the things procurement actually needs: your identity provider, air-gap, compliance reporting, and a real SLA — at a published price.`}
      </Typography>
    </Stack>
  );
};
