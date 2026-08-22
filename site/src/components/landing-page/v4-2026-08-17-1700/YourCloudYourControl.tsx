'use client';

import { Box, Typography } from '@mui/material';
import { FC } from 'react';
import {
  Section,
  SectionTitle,
  SectionSubtitle,
  Grid,
  Card,
  FeatureTitle,
  BodyText,
  Badge,
} from './shared';

/**
 * The trust section. Three credential models in increasing order of
 * security, with the keyless line as the closer. All claims here are
 * shipped behavior — nothing on this section may describe roadmap.
 */

const models = [
  {
    badge: 'Quick Start',
    title: 'Stored as References',
    description:
      'Credentials are saved as encrypted references, resolved only at execution time on the runner — never sitting in a web form or a plain database row.',
  },
  {
    badge: 'Recommended',
    title: 'Keyless Connections',
    description:
      'No stored credential at all. Your cloud trusts a short-lived, connection-scoped token, and you can revoke the trust from your own console — without asking us.',
  },
  {
    badge: 'Maximum Security',
    title: 'Runner In Your Network',
    description:
      'Deployments execute on a runner inside your own network, connecting outbound-only on port 443. No inbound firewall rules, no VPN.',
  },
];

export const YourCloudYourControl: FC = () => (
  <Section id="security">
    <Box className="text-center mb-10">
      <SectionTitle>Your Cloud, Your Control</SectionTitle>
      <SectionSubtitle className="mx-auto">
        Everything deploys into your own cloud account — your credentials,
        your state, your audit trail. Choose how much trust you extend.
      </SectionSubtitle>
    </Box>

    <Grid cols={3} className="max-w-5xl mx-auto mb-8">
      {models.map((model) => (
        <Card key={model.title} className="text-center">
          <Box className="mb-3">
            <Badge>{model.badge}</Badge>
          </Box>
          <FeatureTitle className="mb-2">{model.title}</FeatureTitle>
          <BodyText>{model.description}</BodyText>
        </Card>
      ))}
    </Grid>

    <Typography className="text-center text-sm text-[#a0a0a0] max-w-2xl mx-auto">
      Every change carries a full audit record: who, what, when, and the
      exact configuration that was deployed.
    </Typography>
  </Section>
);
