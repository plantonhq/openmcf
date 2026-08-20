'use client';

import { Box, Stack, Typography } from '@mui/material';
import Image from 'next/image';
import { useState, useCallback } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import confetti from 'canvas-confetti';
import {
  Speed as SpeedIcon,
  CloudQueue as CloudIcon,
  Code as CodeIcon,
  LockOpen as LockIcon,
} from '@mui/icons-material';
import { BookDemoForm } from './BookDemoForm';
import { BookDemoScheduler } from './BookDemoScheduler';
import type { Phase, DemoFormData } from './types';
import { PLATFORM_STATS } from '@/data/platform-stats';

const cloudProviders = [
  { src: '/_site/images/providers/aws.svg', alt: 'AWS' },
  { src: '/_site/images/providers/gcp.svg', alt: 'GCP' },
  { src: '/_site/images/providers/azure.svg', alt: 'Azure' },
  { src: '/_site/images/providers/kubernetes.svg', alt: 'Kubernetes' },
  { src: '/_site/images/providers/cloudflare.svg', alt: 'Cloudflare' },
];

const metrics = [
  { value: '<1 hr', label: 'Average setup time' },
  { value: PLATFORM_STATS.DEPLOYMENT_MODULE_COUNT, label: 'Cloud resource kinds' },
  { value: '100%', label: 'Customer retention' },
];

const valueBullets = [
  {
    icon: <SpeedIcon sx={{ fontSize: 20 }} />,
    title: 'Production-ready in under an hour',
    description:
      'Deploy complete cloud environments with one command. VPC, load balancer, DNS, certificates — all orchestrated.',
  },
  {
    icon: <CodeIcon sx={{ fontSize: 20 }} />,
    title: 'Push code, get deployments',
    description:
      'Built-in CI/CD from git push to production URL. No Dockerfiles, no pipeline YAML.',
  },
  {
    icon: <CloudIcon sx={{ fontSize: 20 }} />,
    title: 'Your cloud, your credentials',
    description:
      'Everything runs in your AWS, GCP, or Azure account. Secrets never leave your infrastructure.',
  },
  {
    icon: <LockIcon sx={{ fontSize: 20 }} />,
    title: 'Open source foundation',
    description:
      'Every infrastructure module is open source. Export manifests and continue independently.',
  },
];

function fireConfetti() {
  const colors = ['#ffffff', '#ededed', '#a0a0a0', '#10b981'];
  confetti({
    particleCount: 80,
    spread: 70,
    origin: { y: 0.3 },
    colors,
    disableForReducedMotion: true,
  });
}

const fadeVariants = {
  initial: { opacity: 0, y: 12 },
  animate: { opacity: 1, y: 0, transition: { duration: 0.4, ease: 'easeOut' } },
  exit: { opacity: 0, y: -12, transition: { duration: 0.25, ease: 'easeIn' } },
};

export function BookDemoPage() {
  const [phase, setPhase] = useState<Phase>('form');
  const [formData, setFormData] = useState<DemoFormData | null>(null);

  const handleFormSuccess = useCallback((data: DemoFormData) => {
    setFormData(data);
    setPhase('scheduler');
    fireConfetti();
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, []);

  const firstName = formData?.firstName ?? '';

  return (
    <Box className="relative min-h-screen bg-[#0a0a0a] flex items-center">
      {/* Grid background */}
      <Box
        className="absolute inset-0 pointer-events-none opacity-[0.03]"
        sx={{
          backgroundImage: `
            linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px),
            linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px)
          `,
          backgroundSize: '60px 60px',
        }}
      />

      <Box className="relative z-10 w-full max-w-7xl mx-auto px-4 md:px-8 py-12 md:py-16">
        <Box className="grid grid-cols-1 lg:grid-cols-12 gap-10 lg:gap-16 items-start">
          {/* ── Left Column ── */}
          <Box className="lg:col-span-7">
            <AnimatePresence mode="wait">
              {phase === 'form' ? (
                <motion.div key="phase1-left" {...fadeVariants}>
                  <Phase1Left />
                </motion.div>
              ) : (
                <motion.div key="phase2-left" {...fadeVariants}>
                  <Phase2Left firstName={firstName} />
                </motion.div>
              )}
            </AnimatePresence>
          </Box>

          {/* ── Right Column ── */}
          <Box className="lg:col-span-5">
            <AnimatePresence mode="wait">
              {phase === 'form' ? (
                <motion.div key="phase1-right" {...fadeVariants}>
                  <BookDemoForm onSuccess={handleFormSuccess} />
                </motion.div>
              ) : (
                <motion.div key="phase2-right" {...fadeVariants}>
                  <BookDemoScheduler formData={formData!} />
                </motion.div>
              )}
            </AnimatePresence>
          </Box>
        </Box>
      </Box>
    </Box>
  );
}

/* ═══════════════════════════════════════════════════════════════════
   Phase 1 — Left Column
   ═══════════════════════════════════════════════════════════════════ */

function Phase1Left() {
  return (
    <Stack className="gap-8">
      {/* Headline */}
      <Stack className="gap-4">
        <Typography
          variant="h1"
          className="text-3xl sm:text-4xl md:text-[44px] font-semibold text-white leading-[1.15] tracking-tight"
        >
          See Planton in action
        </Typography>
        <Typography className="text-base md:text-lg text-[#a0a0a0] leading-relaxed max-w-xl">
          Get a walkthrough of how teams deploy production infrastructure
          and ship services across AWS, GCP, and Azure &mdash; without a
          dedicated DevOps team.
        </Typography>
      </Stack>

      {/* Value bullets */}
      <Stack className="gap-5">
        {valueBullets.map((bullet) => (
          <Stack key={bullet.title} direction="row" className="gap-4 items-start">
            <Box className="w-10 h-10 rounded-lg bg-white/5 border border-[#2a2a2a] flex items-center justify-center flex-shrink-0 text-[#a0a0a0]">
              {bullet.icon}
            </Box>
            <Stack className="gap-0.5">
              <Typography className="text-sm font-semibold text-white">
                {bullet.title}
              </Typography>
              <Typography className="text-sm text-[#666] leading-relaxed">
                {bullet.description}
              </Typography>
            </Stack>
          </Stack>
        ))}
      </Stack>

      {/* Provider logos */}
      <Box className="flex flex-wrap items-center gap-5 md:gap-6">
        {cloudProviders.map((provider) => (
          <Box
            key={provider.alt}
            className="opacity-40 hover:opacity-70 transition-opacity"
            sx={{ height: { xs: 22, sm: 28 }, width: 'auto' }}
          >
            <Image
              src={provider.src}
              alt={provider.alt}
              width={72}
              height={28}
              style={{ height: '100%', width: 'auto', objectFit: 'contain' }}
            />
          </Box>
        ))}
      </Box>

      {/* Metrics strip */}
      <Box className="flex flex-wrap items-center gap-6 md:gap-8 pt-2 border-t border-[#2a2a2a]/50">
        {metrics.map((metric) => (
          <Box key={metric.label} className="text-left">
            <Typography className="text-xl md:text-2xl font-bold text-white">
              {metric.value}
            </Typography>
            <Typography className="text-xs text-[#666] mt-0.5">
              {metric.label}
            </Typography>
          </Box>
        ))}
      </Box>
    </Stack>
  );
}

/* ═══════════════════════════════════════════════════════════════════
   Phase 2 — Left Column
   ═══════════════════════════════════════════════════════════════════ */

const prepTips = [
  'Discovery call — no obligations',
  'Come with your current stack and infrastructure pain points in mind',
  "We'll walk through your specific use case — not a generic pitch",
  'Technical conversation tailored to your team\u2019s needs',
];

function Phase2Left({ firstName }: { firstName: string }) {
  return (
    <Stack className="gap-8">
      {/* Confirmation */}
      <Stack className="gap-3">
        <Typography
          variant="h1"
          className="text-3xl sm:text-4xl md:text-[44px] font-semibold text-white leading-[1.15] tracking-tight"
        >
          Thanks, {firstName}!
        </Typography>
        <Typography className="text-base md:text-lg text-[#a0a0a0] leading-relaxed max-w-xl">
          Your request is in. Someone from the Planton team will reach out
          via email within 1&ndash;2 business days.
        </Typography>
      </Stack>

      {/* Soft bridge */}
      <Box className="flex items-center gap-4 max-w-xl">
        <Box className="flex-1 h-px bg-[#2a2a2a]" />
        <Typography className="text-xs text-[#555] uppercase tracking-widest flex-shrink-0">
          or
        </Typography>
        <Box className="flex-1 h-px bg-[#2a2a2a]" />
      </Box>

      <Stack className="gap-3">
        <Typography className="text-base md:text-lg text-white font-medium">
          Prefer to skip the wait?
        </Typography>
        <Typography className="text-sm text-[#a0a0a0] leading-relaxed max-w-xl">
          Pick a time that works for you. We&rsquo;ll have a technical
          conversation tailored to your infrastructure and team.
        </Typography>
      </Stack>

      {/* What to expect */}
      <Stack className="gap-3 pt-4 border-t border-[#2a2a2a]/50">
        <Typography className="text-xs uppercase tracking-widest text-[#555] font-medium">
          What to expect
        </Typography>
        <Stack className="gap-2.5">
          {prepTips.map((tip) => (
            <Stack key={tip} direction="row" className="gap-3 items-start">
              <Box className="w-1.5 h-1.5 rounded-full bg-[#10b981]/60 mt-1.5 flex-shrink-0" />
              <Typography className="text-sm text-[#a0a0a0] leading-relaxed">
                {tip}
              </Typography>
            </Stack>
          ))}
        </Stack>
      </Stack>
    </Stack>
  );
}
