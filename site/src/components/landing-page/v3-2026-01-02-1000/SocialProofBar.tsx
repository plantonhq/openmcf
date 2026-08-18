'use client';

import { Box, Stack, Typography } from '@mui/material';
import Image from 'next/image';
import { FC } from 'react';
import { Section, Badge } from './shared';

const customers = [
  {
    name: 'Jai.CX',
    url: 'https://jai.cx/',
    logo: '/_site/images/customers/logos/jai-cx.svg',
    logoClass: 'max-h-8',
  },
  {
    name: 'TYNYBAY',
    url: 'https://tynybay.com',
    logo: '/_site/images/customers/logos/tynybay.png',
    logoClass: 'max-h-20', // Larger to compensate for logo padding
  },
  {
    name: 'iorta TechNext',
    url: 'https://iorta.in',
    logo: '/_site/images/customers/logos/iorta.svg',
    logoClass: 'max-h-8',
  },
  {
    name: 'Odwen',
    url: 'https://odwen.co.in',
    logo: '/_site/images/customers/logos/odwen.svg',
    logoClass: 'max-h-8',
  },
];


export const SocialProofBar: FC = () => {
  return (
    <Section className="py-12 md:py-16">
      {/* Section header */}
      <Stack className="items-center text-center mb-10">
        <Typography className="text-sm text-[#a0a0a0] uppercase tracking-widest mb-3">
          Trusted by Fast-Moving Teams
        </Typography>
      </Stack>

      {/* Customer logos */}
      <Box className="flex flex-wrap items-center justify-center gap-8 md:gap-12">
        {customers.map((customer, index) => (
          <a
            key={index}
            href={customer.url}
            target="_blank"
            rel="noopener noreferrer"
            className="group h-10 md:h-12 flex items-center justify-center"
          >
            <Image
              src={customer.logo}
              alt={customer.name}
              width={120}
              height={48}
              className={`${customer.logoClass} w-auto object-contain brightness-0 invert opacity-70 group-hover:opacity-100 transition-opacity`}
            />
          </a>
        ))}
        <Badge variant="default" className="opacity-70">
          + others in beta
        </Badge>
      </Box>
    </Section>
  );
};

