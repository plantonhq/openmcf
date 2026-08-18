'use client';

import { Box, Stack, Typography } from '@mui/material';
import { FC } from 'react';
import { Section, SectionTitle, SectionSubtitle, Card, Badge } from './shared';

const securityModels = [
  {
    label: 'OAuth-Based Access',
    description: 'Short-lived tokens, no stored credentials',
    tag: 'Quick Start',
  },
  {
    label: 'Trust Relationship',
    description: 'Planton assumes IAM role in your account',
    tag: 'Recommended',
  },
  {
    label: 'Customer-Hosted Runner',
    description: 'All operations execute inside your network',
    tag: 'Maximum Security',
  },
];

const securityPillars = [
  { label: 'Zero-Trust Architecture' },
  { label: 'Open Source Audit' },
  { label: 'Scoped IAM Permissions' },
  { label: 'Encrypted Tunnels' },
];

export const Security: FC = () => {
  return (
    <Section id="security">
      <Stack className="items-center text-center mb-12">
        <Badge className="mb-6">SECURITY</Badge>
        <SectionTitle>
          Your Cloud, Your Control
        </SectionTitle>
        <SectionSubtitle className="mx-auto">
          Choose your security model. All operations happen in your cloud account.
        </SectionSubtitle>
      </Stack>

      <Box className="flex flex-wrap justify-center gap-3 mb-12">
        {securityPillars.map((pillar) => (
          <Box
            key={pillar.label}
            className="flex items-center gap-2 px-4 py-2 rounded-full bg-white/5 border border-white/20 text-white"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
            <Typography className="text-xs md:text-sm font-medium">{pillar.label}</Typography>
          </Box>
        ))}
      </Box>

      <Box className="grid md:grid-cols-3 gap-4 mb-10">
        {securityModels.map((model) => (
          <Card key={model.label} className="text-center">
            <Badge className="mb-3">{model.tag}</Badge>
            <Typography className="text-white font-medium mb-2">{model.label}</Typography>
            <Typography className="text-sm text-[#a0a0a0]">{model.description}</Typography>
          </Card>
        ))}
      </Box>

      <Card className="p-0 overflow-hidden max-w-3xl mx-auto">
        <Box className="p-4 border-b border-[#2a2a2a] bg-[#0f0f0f]">
          <Typography className="text-sm font-medium text-white">
            Scoped IAM Permissions (Example)
          </Typography>
        </Box>
        <Box className="p-6 font-mono text-sm bg-[#0a0a0a] overflow-x-auto">
          <pre className="text-[#a0a0a0]">
{`{
  "Effect": "Allow",
  "Action": [
    "ecs:CreateService",
    "ecs:UpdateService",
    "ecs:DeleteService"
  ],
  "Resource": "arn:aws:ecs:*:*:service/planton-*"
}`}
          </pre>
        </Box>
        <Box className="p-4 border-t border-[#2a2a2a] bg-[#0f0f0f]">
          <Typography className="text-xs text-[#666]">
            Only resources with <code className="text-white">planton-</code> prefix.
            Never full account access.
          </Typography>
        </Box>
      </Card>
    </Section>
  );
};
