'use client';

import { FC } from 'react';
import Link from 'next/link';
import { motion } from 'framer-motion';
import {
  fadeInUp,
  staggerContainer,
  defaultTransition,
  Container,
  Badge,
} from '../explainer/shared';
import { Footer } from '../explainer/layout';

// ============================================================================
// NAVIGATION CARD DATA
// ============================================================================

interface NavCard {
  title: string;
  description: string;
  href: string;
  icon: string;
  badge?: string;
  primary?: boolean;
}

const PRIMARY_CARDS: NavCard[] = [
  {
    title: 'The Pitch Deck',
    description: '13-slide presentation covering problem, solution, traction, and ask',
    href: '/invest/deck',
    icon: '📊',
    badge: '13 slides',
    primary: true,
  },
  {
    title: 'The Opportunity',
    description: 'Market context, alternatives, and why we believe this is worth building',
    href: '/invest/opportunity',
    icon: '🎯',
  },
  {
    title: 'How to Invest',
    description: 'Step-by-step process using YC SAFE via Carta',
    href: '/invest/process',
    icon: '📝',
  },
  {
    title: 'Investor Updates',
    description: 'What we shipped, where we are, and where we are headed',
    href: '/legal/investor-updates',
    icon: '📰',
    badge: 'For investors',
  },
];

interface SecondaryCard {
  title: string;
  description: string;
  href: string;
  icon: string;
}

const SECONDARY_CARDS: SecondaryCard[] = [
  {
    title: 'What You Get',
    description: 'Investment terms explained for your background',
    href: '/invest/and-you-get',
    icon: '💎',
  },
  {
    title: 'What We Look For',
    description: 'The investors we work best with',
    href: '/invest/if-you-are',
    icon: '🤝',
  },
  {
    title: 'Carta Walkthrough',
    description: 'See the actual SAFE signing process',
    href: '/invest/steps/carta',
    icon: '📋',
  },
];

// ============================================================================
// CREDIBILITY METRICS
// ============================================================================

interface CredibilityMetric {
  value: string;
  label: string;
}

const CREDIBILITY_METRICS: CredibilityMetric[] = [
  { value: '3+', label: 'Years Building' },
  { value: '$500K+', label: 'Self-Funded' },
  { value: '100%', label: 'Customer Retention' },
  { value: 'DE C-Corp', label: 'Structure' },
];

// ============================================================================
// HERO SECTION
// ============================================================================

const HeroSection: FC = () => (
  <section className="relative py-16 sm:py-20 md:py-28">
    <Container className="relative text-center">
      <motion.div
        initial="hidden"
        animate="visible"
        variants={staggerContainer}
      >
        {/* Key terms badge */}
        <motion.div variants={fadeInUp} transition={defaultTransition}>
          <Badge className="mb-6 text-sm px-4 py-1.5">
            Seed Stage — Raising $500K SAFE at $7M Cap
          </Badge>
        </motion.div>

        {/* Headline */}
        <motion.h1
          variants={fadeInUp}
          transition={{ ...defaultTransition, delay: 0.1 }}
          className="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-semibold text-white tracking-tight mb-6"
        >
          Invest in Planton
        </motion.h1>

        {/* Tagline */}
        <motion.p
          variants={fadeInUp}
          transition={{ ...defaultTransition, delay: 0.2 }}
          className="text-base sm:text-lg md:text-xl text-[#a0a0a0] max-w-2xl mx-auto mb-10"
        >
          The Self-Service Cloud Platform — AI designs infrastructure, teams
          deploy it as templates, in their own cloud accounts. Accessible to
          every company, not just well-funded startups.
        </motion.p>

        {/* Primary CTA */}
        <motion.div
          variants={fadeInUp}
          transition={{ ...defaultTransition, delay: 0.3 }}
        >
          <Link
            href="/invest/deck"
            className="inline-flex items-center gap-2 px-5 py-2.5 rounded-lg bg-[#fff] text-black font-medium text-sm hover:bg-gray-200 transition-all duration-300 hover:-translate-y-0.5"
          >
            View Pitch Deck
            <span>→</span>
          </Link>
        </motion.div>
      </motion.div>
    </Container>
  </section>
);

// ============================================================================
// NAVIGATION CARDS SECTION
// ============================================================================

const NavigationCard: FC<NavCard> = ({ title, description, href, icon, badge }) => (
  <Link href={href} className="block group">
    <motion.div
      variants={fadeInUp}
      className="relative h-full rounded-xl border p-5 sm:p-6 transition-all duration-300 bg-[#151515] border-[#2a2a2a] hover:border-[#3a3a3a] hover:bg-[#1a1a1a]"
    >
      <div className="flex items-start justify-between mb-3">
        <span className="text-2xl">{icon}</span>
        {badge && (
          <Badge className="text-xs">
            {badge}
          </Badge>
        )}
      </div>
      <h3 className="text-base sm:text-lg font-semibold text-white mb-2">
        {title}
      </h3>
      <p className="text-sm text-[#a0a0a0] leading-relaxed">
        {description}
      </p>
      <div className="mt-4 text-sm font-medium flex items-center gap-1 text-[#666] group-hover:text-white group-hover:gap-2 transition-all duration-300">
        Explore
        <span className="group-hover:translate-x-1 transition-transform duration-300">→</span>
      </div>
    </motion.div>
  </Link>
);

const NavigationCardsSection: FC = () => (
  <section className="py-12 sm:py-16">
    <Container maxWidth="4xl">
      <motion.div
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: '-50px' }}
        variants={staggerContainer}
        className="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6"
      >
        {PRIMARY_CARDS.map((card) => (
          <NavigationCard key={card.href} {...card} />
        ))}
      </motion.div>
    </Container>
  </section>
);

// ============================================================================
// SECONDARY CARDS SECTION
// ============================================================================

const SecondaryCardItem: FC<SecondaryCard> = ({ title, description, href, icon }) => (
  <Link href={href} className="block group">
    <motion.div
      variants={fadeInUp}
      className="h-full rounded-xl border border-[#2a2a2a] bg-[#111] p-4 sm:p-5 hover:border-[#3a3a3a] hover:bg-[#151515] transition-all duration-300"
    >
      <div className="flex items-center gap-3 mb-2">
        <span className="text-lg">{icon}</span>
        <h4 className="text-sm sm:text-base font-semibold text-white">
          {title}
        </h4>
        <span className="ml-auto text-[#3a3a3a] group-hover:text-white group-hover:translate-x-1 transition-all duration-300">
          →
        </span>
      </div>
      <p className="text-xs sm:text-sm text-[#666] leading-relaxed pl-8">
        {description}
      </p>
    </motion.div>
  </Link>
);

const SecondaryCardsSection: FC = () => (
  <section className="pb-12 sm:pb-16">
    <Container maxWidth="4xl">
      <motion.div
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: '-50px' }}
        variants={staggerContainer}
        className="grid grid-cols-1 sm:grid-cols-3 gap-3 sm:gap-4"
      >
        {SECONDARY_CARDS.map((card) => (
          <SecondaryCardItem key={card.href} {...card} />
        ))}
      </motion.div>
    </Container>
  </section>
);

// ============================================================================
// CREDIBILITY STRIP SECTION
// ============================================================================

const CredibilityStripSection: FC = () => (
  <section className="py-12 sm:py-16 border-t border-b border-[#2a2a2a]">
    <Container maxWidth="4xl">
      <motion.div
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={staggerContainer}
        className="grid grid-cols-2 md:grid-cols-4 gap-6 sm:gap-8"
      >
        {CREDIBILITY_METRICS.map((metric, index) => (
          <motion.div
            key={metric.label}
            variants={fadeInUp}
            transition={{ ...defaultTransition, delay: index * 0.1 }}
            className="text-center"
          >
            <div className="text-2xl sm:text-3xl md:text-4xl font-bold tracking-tight text-white mb-1">
              {metric.value}
            </div>
            <div className="text-xs sm:text-sm text-[#666]">
              {metric.label}
            </div>
          </motion.div>
        ))}
      </motion.div>
    </Container>
  </section>
);

// ============================================================================
// READY TO TALK CTA SECTION
// ============================================================================

const ReadyToTalkSection: FC = () => (
  <section className="py-16 sm:py-20">
    <Container>
      <motion.div
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={fadeInUp}
        transition={defaultTransition}
        className="text-center"
      >
        <h2 className="text-xl sm:text-2xl md:text-3xl font-semibold text-white tracking-tight mb-4">
          Ready to Talk?
        </h2>
        <p className="text-sm sm:text-base text-[#a0a0a0] mb-8 max-w-lg mx-auto">
          Have questions or ready to express interest? Reach out directly.
        </p>
        <a
          href="mailto:swarup@planton.ai?subject=Investment Interest - Planton&body=Hi Swarup,%0A%0AI'm interested in learning more about investing in Planton.%0A%0A"
          className="inline-flex items-center gap-2 px-5 py-2.5 rounded-lg border border-[#3a3a3a] text-white hover:border-white hover:bg-white/5 transition-all duration-300"
        >
          swarup@planton.ai
          <span>→</span>
        </a>
      </motion.div>
    </Container>
  </section>
);

// ============================================================================
// MAIN COMPONENT
// ============================================================================

export const InvestLandingPage: FC = () => {
  return (
    <div className="min-h-screen bg-[#0a0a0a] pt-16">
      <main>
        <HeroSection />
        <NavigationCardsSection />
        <SecondaryCardsSection />
        <CredibilityStripSection />
        <ReadyToTalkSection />
      </main>
      <Footer />
    </div>
  );
};

export default InvestLandingPage;
