'use client';

import { Box, Stack, Typography } from '@mui/material';
import { FC } from 'react';
import Link from 'next/link';
import { Section, SectionTitle, SectionSubtitle, Card, CheckIcon, Badge, PrimaryButton, SecondaryButton } from './shared';
import { useMarket } from '@/components/market';

export const PricingSimplified: FC = () => {
  // Locale-detected market (no selector here; /pricing carries the control)
  // so this section and the pricing page can never tell different stories.
  const { market } = useMarket();
  return (
    <Section id="pricing">
      <Stack className="items-center text-center mb-16">
        <Badge className="mb-6">PRICING</Badge>
        <SectionTitle>
          Simple, Predictable Pricing
        </SectionTitle>
        <SectionSubtitle className="mx-auto">
          Start free. Pay per seat. AI on prepaid credits. No hidden costs.
        </SectionSubtitle>
      </Stack>

      <Box className="grid md:grid-cols-2 gap-6 mb-16 max-w-3xl mx-auto">
        <Card className="text-center">
          <Typography className="text-sm text-[#a0a0a0] uppercase tracking-wide mb-2">
            Seat-Based Subscription
          </Typography>
          <Typography className="text-4xl md:text-5xl font-bold text-white mb-2">
            {market.symbol}
            {market.teamSeatMonthly.toLocaleString()}
          </Typography>
          <Typography className="text-[#a0a0a0]">
            per developer / month
          </Typography>
        </Card>
        
        <Card className="text-center">
          <Typography className="text-sm text-[#a0a0a0] uppercase tracking-wide mb-2">
            Prepaid AI Credits
          </Typography>
          <Typography className="text-4xl md:text-5xl font-bold text-white mb-2">
            {market.symbol}
            {market.creditPackStart.toLocaleString()}
          </Typography>
          <Typography className="text-[#a0a0a0]">
            smallest pack — spend protection on by default
          </Typography>
        </Card>
      </Box>

      <Typography className="text-center text-sm text-[#666] mb-16">
        The free tier never asks for a card and never surprise-bills — at its limits it simply pauses.
      </Typography>

      <Card className="mb-16 max-w-4xl mx-auto">
        <Box className="grid lg:grid-cols-2 gap-8 items-center">
          <Box>
            <Badge className="mb-4">CASE STUDY</Badge>
            <Typography className="text-xl font-semibold text-white mb-4">
              iorta TechNext
            </Typography>
            <Box className="space-y-2 text-sm">
              <Box className="flex justify-between">
                <Typography className="text-[#a0a0a0]">Team Size:</Typography>
                <Typography className="text-white">7 developers</Typography>
              </Box>
              <Box className="flex justify-between">
                <Typography className="text-[#a0a0a0]">Monthly Spend:</Typography>
                <Typography className="text-white">~$450 (platform + usage)</Typography>
              </Box>
            </Box>
            
            <Box className="mt-6">
              <Typography className="text-sm text-[#666] uppercase tracking-wide mb-3">
                What They Get
              </Typography>
              <Stack className="gap-2">
                {[
                  'Complete infrastructure management (AWS ECS environment for SalesVerse)',
                  'Automated CI/CD for all services',
                  'Multi-environment deployments (dev, QA, staging, prod)',
                  '24/7 support access',
                  'Full audit trail for compliance',
                ].map((item, index) => (
                  <Box key={index} className="flex items-start gap-2">
                    <Box className="mt-0.5 flex-shrink-0">
                      <CheckIcon />
                    </Box>
                    <Typography className="text-xs text-[#b0b0b0]">{item}</Typography>
                  </Box>
                ))}
              </Stack>
            </Box>
          </Box>
          
          <Box className="text-center">
            <Box className="p-6 rounded-xl bg-white/5 border border-white/10">
              <Typography className="text-sm text-[#a0a0a0] mb-2">The Outcome</Typography>
              <Typography className="text-2xl font-bold text-white">
                Production infra, self-served
              </Typography>
              <Typography className="text-xs text-[#a0a0a0]">
                7 developers deploy independently — without growing the ops team
              </Typography>
            </Box>
          </Box>
        </Box>
      </Card>

      <Box className="text-center">
        <Link href="/pricing">
          <PrimaryButton className="text-lg px-8 py-4">
            View Full Pricing Details →
          </PrimaryButton>
        </Link>
        <Stack direction={{ xs: 'column', sm: 'row' }} className="justify-center gap-4 mt-4">
          <Link href="/book-demo">
            <SecondaryButton>
              Request a Demo
            </SecondaryButton>
          </Link>
        </Stack>
      </Box>
    </Section>
  );
};
