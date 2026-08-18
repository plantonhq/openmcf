import { Metadata } from 'next';
import { Box } from '@mui/material';
import {
  EnterpriseHero,
  EnterpriseHowItWorks,
  EnterpriseRateCards,
} from '@/components/enterprise';
import { PricingCta } from '@/components/pricing';

export const metadata: Metadata = {
  title: 'Enterprise | Planton',
  description:
    'Enterprise at Planton: published rate cards for the US and India, enterprise identity, air-gap, compliance reporting, and real SLAs. Self-serve below 25 seats — no sales call required.',
};

export default function EnterprisePage() {
  return (
    <Box>
      <EnterpriseHero />
      <EnterpriseRateCards />
      <EnterpriseHowItWorks />
      <PricingCta />
    </Box>
  );
}
