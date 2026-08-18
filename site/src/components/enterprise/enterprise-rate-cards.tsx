'use client';

import { FC } from 'react';
import Link from 'next/link';
import { Box, Grid2, Stack, Typography } from '@mui/material';
import {
  CheckIcon,
  FeatureTitle,
  PrimaryButton,
} from '@/components/landing-page/v3-2026-01-02-1000/shared';
import { MARKETS } from '@/data/pricing';

/**
 * Published rate anchors, BOTH markets printed side by side — deliberately,
 * unlike the pricing page's market toggle: against quote-only competitors,
 * a printed number is the trust move, and enterprise buyers compare
 * regions. India prices are set India numbers, never an FX conversion.
 */

const tierBullets: Record<string, string[]> = {
  'Enterprise Standard': [
    'Enterprise SSO (SAML/OIDC) and SCIM directory sync — Coming Soon',
    'Supported air-gap posture',
    'Compliance reporting',
    'Named business-hours support',
  ],
  'Enterprise Plus': [
    'Everything in Enterprise Standard',
    '24×7 support with SLA and a dedicated channel',
    'Procurement and security-review assistance',
    'Contracting on your paper',
  ],
};

export const EnterpriseRateCards: FC = () => {
  return (
    <Stack className="items-center px-4 md:px-8 py-8 gap-4 bg-[#0a0a0a]">
      <Grid2 container spacing={4} className="w-full max-w-5xl items-stretch">
        {MARKETS.us.enterprise.map((tier, tierIndex) => {
          const indiaTier = MARKETS.in.enterprise[tierIndex];
          return (
            <Grid2 size={{ xs: 12, md: 6 }} key={tier.name} className="flex">
              <Box className="w-full rounded-xl border border-[#2a2a2a] bg-[#151515] hover:border-[#3a3a3a] transition-all duration-300 p-6 md:p-8">
                <Stack className="gap-5 h-full justify-between">
                  <Stack className="gap-4">
                    <FeatureTitle>{tier.name}</FeatureTitle>
                    <Box>
                      <Box className="flex items-baseline gap-1.5">
                        <Typography className="text-3xl font-bold text-white">
                          {tier.perYear}
                        </Typography>
                        <Typography className="text-sm text-[#a0a0a0]">/year</Typography>
                      </Box>
                      <Typography className="text-sm text-[#a0a0a0] mt-1">
                        {indiaTier.perYear}/year in India
                      </Typography>
                    </Box>
                    <Typography className="text-sm font-medium text-[#c0c0c0]">
                      Up to {tier.seatCeiling} seats
                    </Typography>
                    <Stack className="gap-2">
                      {tierBullets[tier.name].map((bullet) => (
                        <Box key={bullet} className="flex items-start gap-2">
                          <Box className="mt-1 flex-shrink-0">
                            <CheckIcon />
                          </Box>
                          <Typography className="text-sm text-[#c0c0c0] leading-relaxed">
                            {bullet}
                          </Typography>
                        </Box>
                      ))}
                    </Stack>
                  </Stack>
                  <Link href="/book-demo" className="self-start">
                    <PrimaryButton>Talk to Us</PrimaryButton>
                  </Link>
                </Stack>
              </Box>
            </Grid2>
          );
        })}
      </Grid2>
      <Typography className="text-sm text-[#8a8a8a] max-w-[760px] text-center">
        Features marked Coming Soon are decided packaging that is still
        shipping — they become live as they ship.
      </Typography>
    </Stack>
  );
};
