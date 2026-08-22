'use client';

import { Box, Typography } from '@mui/material';
import { FC } from 'react';
import {
  Section,
  SectionTitle,
  SectionSubtitle,
  Grid,
  Metric,
  TestimonialCard,
} from './shared';

/**
 * Social proof with a hard rule: testimonials are quoted VERBATIM from
 * their source of record (company/sales/customers/... in the platform
 * repo), attributed to the person who said them — never paraphrased,
 * never attributed to a company. No dollar figures anywhere in this
 * section: the approved value story is "run production infra without
 * growing your ops team", not a savings number.
 */

const stats = [
  { value: '<1 hr', label: 'Infrastructure Setup' },
  { value: '2023', label: 'In Production Since' },
  { value: '0', label: 'Security Incidents' },
  { value: '100%', label: 'Customer Retention' },
];

const testimonials = [
  {
    name: 'Sai Saketh',
    role: 'Junior DevOps Engineer',
    company: 'iorta TechNext',
    location: 'India',
    quote:
      'As a junior DevOps engineer with almost no AWS experience, Planton enabled me to provide a very mature developer experience to our entire 7-person dev team. They can quickly deploy services to multiple environments without me having to deal with learning AWS from scratch or rewriting complex infrastructure code.',
  },
  {
    name: 'Rohit Reddy Gopu',
    role: 'CEO',
    company: 'TynyBay',
    location: 'India',
    quote:
      'For one client project where the client mandated GCP but our DevOps engineer had no GCP experience, Planton allowed us to successfully deliver the entire infrastructure. We essentially got full DevOps capabilities for GCP without needing GCP expertise on our team.',
  },
  {
    name: 'Balaji Borra',
    role: 'DevOps Engineer',
    company: 'TynyBay',
    location: 'India',
    quote:
      'Planton has dramatically improved my efficiency. I no longer have to deal with the grunt work of rewriting Terraform configurations between client projects. I can now manage multiple client environments simultaneously and provide a much better experience for all the developers I support.',
  },
];

export const Proof: FC = () => (
  <Section id="proof">
    <Box className="text-center mb-10">
      <SectionTitle>Teams Ship With Planton</SectionTitle>
      <SectionSubtitle className="mx-auto">
        Consulting firms and product teams, running in production since 2023.
      </SectionSubtitle>
    </Box>

    <Box className="grid grid-cols-2 md:grid-cols-4 gap-6 max-w-3xl mx-auto mb-12">
      {stats.map((stat) => (
        <Metric key={stat.label} value={stat.value} label={stat.label} />
      ))}
    </Box>

    <Grid cols={3} className="max-w-5xl mx-auto mb-8">
      {testimonials.map((testimonial) => (
        <TestimonialCard
          key={testimonial.name}
          name={testimonial.name}
          role={testimonial.role}
          company={testimonial.company}
          location={testimonial.location}
          quote={testimonial.quote}
        />
      ))}
    </Grid>

    <Typography className="text-center text-sm text-[#a0a0a0]">
      And the platform is its own hardest customer —{' '}
      <strong className="text-white">Planton runs on Planton</strong>, from
      its infrastructure to the pipelines that ship it.
    </Typography>
  </Section>
);
