'use client';

import { FC } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import {
  Section,
  Container,
  SectionTitle,
  SectionSubtitle,
  BodyText,
  SmallText,
  Card,
  CardTitle,
  Badge,
} from '../explainer/shared';
import { Hero, Footer } from '../explainer/layout';

// ============================================================================
// WALKTHROUGH STEPS
// ============================================================================

interface WalkthroughStep {
  image: string;
  title: string;
  description: string;
  highlight?: string;
}

const WALKTHROUGH_STEPS: WalkthroughStep[] = [
  {
    image: '/_site/images/carta-walkthrough/01-cap-table.png',
    title: 'Cap Table Overview',
    description:
      'Your ownership is tracked professionally on Carta. This shows the current cap table with fully diluted shares, ownership percentages, and amount raised — all visible in real-time.',
    highlight: 'Real-time ownership tracking',
  },
  {
    image: '/_site/images/carta-walkthrough/02-raise-funds.png',
    title: 'Raise Funds Dashboard',
    description:
      'Carta provides an industry-standard SAFE workflow. Identity verification is complete, and SAFEs can be created, signed, and funded in one simple workflow.',
    highlight: 'Industry-standard SAFE workflow',
  },
  {
    image: '/_site/images/carta-walkthrough/03-safe-terms.png',
    title: 'SAFE Terms Configuration',
    description:
      'The exact terms we discussed: $7M valuation cap, Delaware governing law, Pre-money dilution type. Carta generates the standard YC SAFE document automatically.',
    highlight: '$7M cap, Delaware, Pre-money',
  },
  {
    image: '/_site/images/carta-walkthrough/04-funding-method.png',
    title: 'Funding Method',
    description:
      'Mercury bank integration for secure, direct transfers. Funds arrive in 4-5 business days with automatic SAFE issuance once funds are confirmed. No fees, reduced wire fraud risk.',
    highlight: 'Mercury bank integration',
  },
  {
    image: '/_site/images/carta-walkthrough/05-review-safe.png',
    title: 'Final Review',
    description:
      'Before sending to the investor, everything is reviewed: financing information, investor details, payment account, and legal agreements. One click sends the SAFE for signature.',
    highlight: 'Send to investor',
  },
];

// ============================================================================
// MAIN COMPONENT
// ============================================================================

export const CartaWalkthroughPage: FC = () => {
  return (
    <div className="min-h-screen bg-[#0a0a0a] pt-16">
      <main>
        {/* Hero Section */}
        <Hero
          title="Carta Walkthrough"
          subtitle="See the actual SAFE creation process. These screenshots show exactly how your investment is handled on Carta — the industry standard for cap table management."
        />

        {/* Back Navigation */}
        <Section>
          <Container>
            <div className="flex items-center justify-between mb-8">
              <Link
                href="/invest/steps"
                className="inline-flex items-center gap-2 text-white/60 hover:text-white transition-colors"
              >
                ← Back to Investment Process
              </Link>
              <a
                href="https://carta.com/"
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-white hover:underline"
              >
                Visit Carta.com ↗
              </a>
            </div>

            {/* Context Card */}
            <Card variant="highlight" className="mb-12">
              <div className="flex flex-col md:flex-row items-start md:items-center gap-4">
                <div className="flex-1">
                  <CardTitle color="#ededed" className="mb-2">
                    What You&apos;re Looking At
                  </CardTitle>
                  <BodyText>
                    These are real screenshots from Planton&apos;s Carta account showing the SAFE
                    creation workflow. This is the same interface used to manage investments for
                    companies like Stripe, Coinbase, and Robinhood.
                  </BodyText>
                </div>
                <Badge variant="cyan">40,000+ Companies</Badge>
              </div>
            </Card>
          </Container>
        </Section>

        {/* Walkthrough Steps */}
        <Section background="gradient-subtle">
          <Container>
            <SectionTitle gradient className="text-center">
              Step-by-Step Process
            </SectionTitle>
            <SectionSubtitle className="text-center mx-auto mb-12">
              From cap table to sending the SAFE — here&apos;s exactly what happens.
            </SectionSubtitle>

            <div className="space-y-12">
              {WALKTHROUGH_STEPS.map((step, index) => (
                <Card key={index} className="overflow-hidden">
                  {/* Step Header */}
                  <div className="flex items-center gap-4 mb-4">
                    <div className="w-10 h-10 rounded-full bg-[#2a2a2a] flex items-center justify-center text-white font-bold">
                      {index + 1}
                    </div>
                    <div className="flex-1">
                      <CardTitle>{step.title}</CardTitle>
                    </div>
                    {step.highlight && (
                      <Badge variant="purple">{step.highlight}</Badge>
                    )}
                  </div>

                  {/* Description */}
                  <BodyText className="mb-6">{step.description}</BodyText>

                  {/* Screenshot */}
                  <div className="rounded-lg overflow-hidden border border-white/10">
                    <Image
                      src={step.image}
                      alt={`${step.title} - Carta screenshot`}
                      width={1200}
                      height={675}
                      className="w-full"
                    />
                  </div>
                </Card>
              ))}
            </div>
          </Container>
        </Section>

        {/* CTA Section */}
        <Section>
          <Container>
            <div className="text-center py-8">
              <SectionTitle gradient className="text-center mb-4">
                Ready to Invest?
              </SectionTitle>
              <SectionSubtitle className="text-center mx-auto mb-8">
                Your SAFE will be processed through this exact workflow.
              </SectionSubtitle>

              <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
                <a
                  href="mailto:swarup@planton.ai?subject=Investment Interest - Planton&body=Hi Swarup,%0A%0AI'm interested in investing in Planton.%0A%0AName: %0AEmail: %0AIntended Amount: $%0A%0A"
                  className="inline-flex items-center gap-2 px-5 py-2.5 rounded-lg bg-[#fff] text-black font-medium text-sm hover:bg-gray-200 transition-all duration-300 hover:-translate-y-0.5"
                >
                  Express Interest →
                </a>
                <Link
                  href="/invest/steps"
                  className="inline-flex items-center gap-2 px-6 py-3 rounded-lg border border-white/20 text-white hover:bg-white/5 transition-colors"
                >
                  Investment Process
                </Link>
                <Link
                  href="/invest/and-you-get"
                  className="inline-flex items-center gap-2 px-6 py-3 rounded-lg border border-white/10 text-white/60 hover:bg-white/5 transition-colors"
                >
                  What You Get
                </Link>
              </div>

              <SmallText className="mt-8">
                Questions? Email{' '}
                <a
                  href="mailto:swarup@planton.ai"
                  className="text-white hover:underline"
                >
                  swarup@planton.ai
                </a>
              </SmallText>
            </div>
          </Container>
        </Section>
      </main>

      <Footer />
    </div>
  );
};

export default CartaWalkthroughPage;
