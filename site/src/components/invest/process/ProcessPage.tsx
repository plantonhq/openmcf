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
  CardText,
  Badge,
  Grid,
  List,
  FAQItemData,
} from '../explainer/shared';
import { Hero, Footer } from '../explainer/layout';
import { FAQ } from '../explainer/content';

// ============================================================================
// EXTERNAL LINKS
// ============================================================================

const SAFE_LEARN_URL = 'https://carta.com/sg/en/learn/startups/fundraising/convertible-securities/safes/';
const CARTA_URL = 'https://carta.com/';
const MERCURY_URL = 'https://mercury.com/';
const STRIPE_URL = 'https://stripe.com/';
const STRIPE_ATLAS_URL = 'https://stripe.com/atlas';

// ============================================================================
// REQUIRED INPUT FIELDS - Information needed from investor
// ============================================================================

interface RequiredInput {
  label: string;
  purpose: string;
  example: string;
  icon: string;
}

const REQUIRED_INPUTS: RequiredInput[] = [
  {
    label: 'Legal Name',
    purpose: 'Appears on the SAFE agreement',
    example: 'John M. Smith',
    icon: '📝',
  },
  {
    label: 'Email Address',
    purpose: 'Carta sends the signing link here',
    example: 'john@example.com',
    icon: '📧',
  },
  {
    label: 'Investment Amount',
    purpose: 'Amount for the SAFE document',
    example: '$10,000',
    icon: '💵',
  },
];

// ============================================================================
// REQUIRED INPUT CARD COMPONENT
// ============================================================================

const RequiredInputCard: FC<{ className?: string }> = ({ className = '' }) => (
  <div className={`mt-4 rounded-xl border border-[#10b981]/30 bg-emerald-500/5 p-4 ${className}`}>
    <div className="flex items-center gap-2 mb-3">
      <span className="text-[#10b981] font-semibold text-sm">What We Need From You</span>
      <Badge variant="success">Required</Badge>
    </div>
    <div className="space-y-3">
      {REQUIRED_INPUTS.map((input) => (
        <div
          key={input.label}
          className="flex items-start gap-3 p-3 rounded-lg bg-white/5 border border-white/10"
        >
          <span className="text-xl flex-shrink-0">{input.icon}</span>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <span className="font-medium text-white text-sm">{input.label}</span>
              <span className="text-xs text-white/40">—</span>
              <span className="text-xs text-[#10b981]">{input.purpose}</span>
            </div>
            <div className="mt-1 text-xs text-white/40">
              Example: <span className="text-white/60 font-mono">{input.example}</span>
            </div>
          </div>
        </div>
      ))}
    </div>
    <SmallText className="mt-3 text-white/50">
      Email <span className="text-white">swarup@planton.ai</span> or reply to the message where you received this link.
    </SmallText>
  </div>
);

// ============================================================================
// THANK YOU SECTION - Personal note to investors
// ============================================================================

const ThankYouSection: FC = () => (
  <Section>
    <Container>
      <div className="border-l-4 border-white/50 pl-6 py-2">
        <BodyText className="text-white/80 text-base sm:text-lg leading-relaxed">
          <span className="font-semibold text-white">Thank you for considering investing in Planton.</span>
        </BodyText>
        <BodyText className="mt-3 text-white/60 leading-relaxed">
          Your interest means a lot. Whether you first heard about us yesterday or have been following
          the journey for months, the fact that you&apos;re here—ready to take this step—is something
          I don&apos;t take lightly. Let&apos;s walk through the process together.
        </BodyText>
        <SmallText className="mt-4 text-white/50">— Swarup</SmallText>
      </div>
    </Container>
  </Section>
);

// ============================================================================
// PROCESS STEPS
// ============================================================================

interface ProcessStep {
  number: number;
  title: string;
  description: string;
  details?: string[];
  who: 'you' | 'swarup' | 'carta';
  /** If true, renders RequiredInputCard instead of details list */
  hasRequiredInputs?: boolean;
}

const INVESTMENT_STEPS: ProcessStep[] = [
  {
    number: 1,
    title: 'Express Interest',
    description: 'Provide the information needed to create your SAFE agreement.',
    hasRequiredInputs: true,
    who: 'you',
  },
  {
    number: 2,
    title: 'SAFE Created on Carta',
    description: 'Swarup creates the investment agreement using your information.',
    details: [
      'YC Post-Money SAFE with $7M valuation cap',
      'No discount, pre-money valuation',
      'Your legal name and email are entered exactly as provided',
    ],
    who: 'swarup',
  },
  {
    number: 3,
    title: 'Document Sent via Carta',
    description: 'You receive the SAFE document for review.',
    details: [
      'Email from Carta with the SAFE agreement',
      'Standard YC document — 3 pages, straightforward terms',
      'Review at your own pace, ask questions if needed',
    ],
    who: 'carta',
  },
  {
    number: 4,
    title: 'Electronic Signature',
    description: 'Sign the document digitally through Carta.',
    details: [
      'Click the link in the email to review and sign',
      'Industry-standard e-signature process',
      'Takes about 2 minutes',
    ],
    who: 'you',
  },
  {
    number: 5,
    title: 'Wire Funds',
    description: 'Transfer your investment to Planton\'s bank account.',
    details: [
      'Carta provides wire instructions in the document',
      'Funds go to Planton\'s Mercury bank account',
      'Same account that receives customer payments',
    ],
    who: 'you',
  },
  {
    number: 6,
    title: 'Investor Portal Access',
    description: 'You\'re officially an investor with full visibility.',
    details: [
      'Access to Carta investor portal',
      'Real-time cap table visibility',
      'Track your investment as the company grows',
    ],
    who: 'carta',
  },
];

// ============================================================================
// WHAT YOU RECEIVE
// ============================================================================

interface Deliverable {
  icon: string;
  title: string;
  description: string;
}

const DELIVERABLES: Deliverable[] = [
  {
    icon: '📄',
    title: 'Signed SAFE Agreement',
    description: 'Your legally binding investment contract, stored securely on Carta.',
  },
  {
    icon: '🔐',
    title: 'Carta Investor Portal',
    description: 'Real-time access to the cap table, your ownership, and company updates.',
  },
  {
    icon: '📧',
    title: 'Monthly Investor Updates',
    description: 'Regular email updates on MRR, wins, challenges, and how you can help.',
  },
  {
    icon: '📊',
    title: 'Quarterly Metrics',
    description: 'Detailed financials, runway, KPIs, and strategic outlook.',
  },
  {
    icon: '📞',
    title: 'Direct Access',
    description: 'Swarup\'s personal contact for questions, introductions, or advice.',
  },
];

// ============================================================================
// FAQ DATA
// ============================================================================

const PROCESS_FAQ: FAQItemData[] = [
  {
    question: 'How long does the whole process take?',
    answer: 'From expressing interest to having your funds in Planton\'s account: typically 3-5 business days. Signing takes 2 minutes; wire transfer timing depends on your bank.',
  },
  {
    question: 'What is Carta?',
    answer: 'Carta is the industry-standard platform for cap table management, used by 40,000+ companies including Stripe, Coinbase, and Robinhood. It handles equity, SAFEs, and investor relations professionally.',
  },
  {
    question: 'What is a YC SAFE?',
    answer: 'Simple Agreement for Future Equity, created by Y Combinator. It\'s a 3-page legal document used by thousands of startups. Your investment converts to shares at Series A, at a price capped at $7M valuation.',
  },
  {
    question: 'What does "pre-money" mean?',
    answer: 'The $7M valuation cap is calculated before the Series A investment goes in. This is the most investor-friendly SAFE structure.',
  },
  {
    question: 'Can I invest from outside the US?',
    answer: 'Yes. Carta supports international investors. Wire transfer instructions will be provided for your bank. For Indian investors, we can also discuss iSAFE/CCPS structures if needed.',
  },
  {
    question: 'Is there a minimum investment?',
    answer: 'Minimum check size is $5K. This allows us to keep administrative overhead manageable while welcoming a broader range of supporters.',
  },
  {
    question: 'What if I have questions about the SAFE terms?',
    answer: 'We encourage you to have a lawyer review if you\'d like. The YC SAFE is standardized, so any startup lawyer will be familiar with it. We\'re also happy to walk through every line with you.',
  },
  {
    question: 'When will I see a return on my investment?',
    answer: 'SAFE converts to shares at Series A (target: 18-24 months). Actual cash returns come at exit—acquisition or IPO—realistically 5-8 years away. This is long-term venture investing.',
  },
];

// ============================================================================
// MAIN COMPONENT
// ============================================================================

export const ProcessPage: FC = () => {
  const getWhoLabel = (who: ProcessStep['who']) => {
    switch (who) {
      case 'you':
        return { text: 'Your Action', variant: 'success' as const };
      case 'swarup':
        return { text: 'Swarup', variant: 'pink' as const };
      case 'carta':
        return { text: 'Carta', variant: 'cyan' as const };
    }
  };

  return (
    <div className="min-h-screen bg-[#0a0a0a] pt-16">
      <main>
        {/* Hero Section */}
        <Hero
          title="How to Invest"
          subtitle="A clear, step-by-step process using industry-standard tools. No surprises, no hidden complexity."
        />

        {/* Thank You Section */}
        <ThankYouSection />

        {/* Overview */}
        <Section>
          <Container>
            <Card variant="highlight" className="mb-8">
              <div className="flex flex-col md:flex-row items-start md:items-center gap-4 md:gap-8">
                <div className="flex-1">
                  <CardTitle color="#ededed" className="mb-2">
                    YC <a href={SAFE_LEARN_URL} target="_blank" rel="noopener noreferrer" className="underline hover:text-white/70">SAFE</a> + Carta = Standardized, Safe, Professional
                  </CardTitle>
                  <BodyText>
                    We use the same investment infrastructure as companies funded by Y Combinator,
                    Andreessen Horowitz, and Sequoia. Your investment is handled with the same
                    professionalism as a $10M round.
                  </BodyText>
                </div>
                <div className="flex items-center gap-4">
                  <a href={SAFE_LEARN_URL} target="_blank" rel="noopener noreferrer" className="text-center hover:opacity-80 transition-opacity">
                    <div className="w-12 h-12 rounded-lg flex items-center justify-center">
                      <Image
                        src="/_site/images/logos/yc.svg"
                        alt="Y Combinator"
                        width={48}
                        height={48}
                        className="rounded"
                      />
                    </div>
                    <SmallText className="mt-1">SAFE</SmallText>
                  </a>
                  <div className="text-2xl text-white/20">+</div>
                  <a href={CARTA_URL} target="_blank" rel="noopener noreferrer" className="text-center hover:opacity-80 transition-opacity">
                    <div className="w-16 h-12 rounded-lg bg-white/10 flex items-center justify-center p-2">
                      <Image
                        src="/_site/images/logos/carta.svg"
                        alt="Carta"
                        width={60}
                        height={26}
                        className="invert brightness-200"
                      />
                    </div>
                    <SmallText className="mt-1">Platform</SmallText>
                  </a>
                </div>
              </div>
            </Card>

            <Grid cols={3} className="mb-8">
              <Card className="text-center">
                <div className="text-4xl font-bold text-[#10b981]">$7M</div>
                <CardTitle className="mt-2">Valuation Cap</CardTitle>
                <CardText>No discount structure</CardText>
              </Card>
              <Card className="text-center">
                <div className="text-4xl font-bold text-white">3</div>
                <CardTitle className="mt-2">Page Document</CardTitle>
                <CardText>Standard YC <a href={SAFE_LEARN_URL} target="_blank" rel="noopener noreferrer" className="text-white hover:underline">SAFE</a></CardText>
              </Card>
              <Card className="text-center">
                <div className="text-4xl font-bold text-white">~5</div>
                <CardTitle className="mt-2">Days to Complete</CardTitle>
                <CardText>From interest to funded</CardText>
              </Card>
            </Grid>
          </Container>
        </Section>

        {/* Step-by-Step Process */}
        <Section background="gradient-subtle">
          <Container>
            <SectionTitle gradient className="text-center">Step-by-Step Process</SectionTitle>
            <SectionSubtitle className="text-center mx-auto mb-8">
              From &quot;I&apos;m interested&quot; to &quot;I&apos;m an investor&quot; in six clear steps.
            </SectionSubtitle>

            <div className="space-y-6">
              {INVESTMENT_STEPS.map((step) => {
                const whoLabel = getWhoLabel(step.who);
                return (
                  <Card
                    key={step.number}
                    className="relative overflow-hidden"
                    variant={step.hasRequiredInputs ? 'highlight' : 'default'}
                  >
                    {/* Step number badge */}
                    <div className="absolute top-0 right-0 w-16 h-16 bg-[#1a1a1a] flex items-center justify-center">
                      <span className="text-2xl font-bold text-white">{step.number}</span>
                    </div>

                    <div className="pr-16">
                      <div className="flex items-center gap-3 mb-2">
                        <CardTitle>{step.title}</CardTitle>
                        <Badge variant={whoLabel.variant}>{whoLabel.text}</Badge>
                        {step.hasRequiredInputs && (
                          <Badge variant="warning">Action Required</Badge>
                        )}
                      </div>
                      <BodyText className="mb-4">{step.description}</BodyText>

                      {/* Render RequiredInputCard for Step 1, otherwise render details list */}
                      {step.hasRequiredInputs ? (
                        <RequiredInputCard />
                      ) : (
                        step.details && step.details.length > 0 && (
                          <List
                            items={step.details.map((detail) => ({
                              icon: step.who === 'you' ? '→' : '•',
                              iconColor: step.who === 'you' ? 'emerald' : 'default',
                              text: detail,
                            }))}
                          />
                        )
                      )}
                    </div>
                  </Card>
                );
              })}
            </div>
          </Container>
        </Section>

        {/* What You Receive */}
        <Section>
          <Container>
            <SectionTitle gradient className="text-center">What You Receive</SectionTitle>
            <SectionSubtitle className="text-center mx-auto mb-8">
              Beyond the legal documents, here&apos;s what being an investor means.
            </SectionSubtitle>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {DELIVERABLES.map((item) => (
                <Card key={item.title} className="flex items-start gap-4">
                  <div className="text-3xl">{item.icon}</div>
                  <div>
                    <CardTitle className="mb-1">{item.title}</CardTitle>
                    <CardText>{item.description}</CardText>
                  </div>
                </Card>
              ))}
            </div>
          </Container>
        </Section>

        {/* Banking & Infrastructure */}
        <Section background="gradient-subtle">
          <Container>
            <SectionTitle gradient className="text-center">Our Infrastructure</SectionTitle>
            <SectionSubtitle className="text-center mx-auto mb-8">
              Professional tools, transparent processes, industry standards.
            </SectionSubtitle>

            <Grid cols={2} className="mb-8">
              <Card>
                <div className="flex items-center gap-3 mb-4">
                  <a href={MERCURY_URL} target="_blank" rel="noopener noreferrer" className="w-12 h-12 rounded-lg bg-white/10 flex items-center justify-center p-2 hover:bg-white/20 transition-colors">
                    <Image
                      src="/_site/images/logos/mercury.svg"
                      alt="Mercury"
                      width={80}
                      height={18}
                      className="invert brightness-200"
                    />
                  </a>
                  <div>
                    <CardTitle>Mercury</CardTitle>
                    <SmallText>Business Banking</SmallText>
                  </div>
                </div>
                <BodyText>
                  Our primary business bank account. Customer payments via Stripe settle here.
                  Investment wires go here. Clear, auditable, startup-standard banking.
                </BodyText>
              </Card>

              <Card>
                <div className="flex items-center gap-3 mb-4">
                  <a href={CARTA_URL} target="_blank" rel="noopener noreferrer" className="w-12 h-12 rounded-lg bg-white/10 flex items-center justify-center p-2 hover:bg-white/20 transition-colors">
                    <Image
                      src="/_site/images/logos/carta.svg"
                      alt="Carta"
                      width={80}
                      height={35}
                      className="invert brightness-200"
                    />
                  </a>
                  <div>
                    <CardTitle>Carta</CardTitle>
                    <SmallText>Cap Table & Equity</SmallText>
                  </div>
                </div>
                <BodyText>
                  Industry-standard equity management. 40,000+ companies use Carta including
                  Stripe, Coinbase, and Robinhood. Your investment is tracked professionally.
                </BodyText>
              </Card>

              <Card>
                <div className="flex items-center gap-3 mb-4">
                  <a href={STRIPE_URL} target="_blank" rel="noopener noreferrer" className="w-12 h-12 rounded-lg bg-white/10 flex items-center justify-center p-1 hover:bg-white/20 transition-colors">
                    <Image
                      src="/_site/images/logos/stripe.png"
                      alt="Stripe"
                      width={40}
                      height={40}
                      className="object-contain"
                    />
                  </a>
                  <div>
                    <CardTitle>Stripe</CardTitle>
                    <SmallText>Payment Processing</SmallText>
                  </div>
                </div>
                <BodyText>
                  Customer subscriptions, licenses, and prepaid AI credits processed through Stripe.
                  Automatic payouts to Mercury. Clean, auditable revenue flow.
                </BodyText>
              </Card>

              <Card>
                <div className="flex items-center gap-3 mb-4">
                  <a href={STRIPE_ATLAS_URL} target="_blank" rel="noopener noreferrer" className="w-12 h-12 rounded-lg bg-white/10 flex items-center justify-center p-1 hover:bg-white/20 transition-colors">
                    <Image
                      src="/_site/images/logos/stripe-atlas.png"
                      alt="Stripe Atlas"
                      width={40}
                      height={40}
                      className="object-contain"
                    />
                  </a>
                  <div>
                    <CardTitle>Stripe Atlas</CardTitle>
                    <SmallText>Company Formation</SmallText>
                  </div>
                </div>
                <BodyText>
                  PlantonCloud Inc. is a Delaware C-Corp formed via Stripe Atlas.
                  Standard US startup structure. EIN: 92-2345321.
                </BodyText>
              </Card>
            </Grid>

            {/* Carta Walkthrough Preview */}
            <Card className="overflow-hidden">
              <CardTitle className="mb-4">See It In Action</CardTitle>
              <BodyText className="mb-4">
                Here&apos;s what the actual SAFE creation process looks like on Carta — the same interface
                you&apos;ll see when your investment is processed.
              </BodyText>
              <Link href="/invest/steps/carta" className="block">
                <Image
                  src="/_site/images/carta-walkthrough/05-review-safe.png"
                  alt="Carta SAFE review screen showing $7M valuation cap, Delaware governing law, and Mercury bank integration"
                  width={800}
                  height={450}
                  className="rounded-lg border border-white/10 hover:opacity-90 transition-opacity w-full"
                />
              </Link>
              <div className="mt-4 flex items-center justify-between">
                <SmallText>
                  Final review screen before sending SAFE to investor
                </SmallText>
                <Link
                  href="/invest/steps/carta"
                  className="text-sm text-white hover:underline"
                >
                  View full walkthrough →
                </Link>
              </div>
            </Card>
          </Container>
        </Section>

        {/* Terms Summary */}
        <Section>
          <Container>
            <SectionTitle gradient className="text-center">Terms Summary</SectionTitle>
            <SectionSubtitle className="text-center mx-auto mb-8">
              Quick reference to the key investment terms.
            </SectionSubtitle>

            <Card className="mb-8">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <SmallText className="text-white font-medium mb-1">Investment Type</SmallText>
                  <BodyText>YC Post-Money <a href={SAFE_LEARN_URL} target="_blank" rel="noopener noreferrer" className="text-white hover:underline">SAFE</a></BodyText>
                </div>
                <div>
                  <SmallText className="text-white font-medium mb-1">Valuation Cap</SmallText>
                  <BodyText>$7M</BodyText>
                </div>
                <div>
                  <SmallText className="text-white font-medium mb-1">Discount</SmallText>
                  <BodyText>None (cap-only structure)</BodyText>
                </div>
                <div>
                  <SmallText className="text-white font-medium mb-1">Minimum Investment</SmallText>
                  <BodyText>$5K</BodyText>
                </div>
                <div>
                  <SmallText className="text-white font-medium mb-1">Pro-Rata Rights</SmallText>
                  <BodyText>Available for $50K+ investments</BodyText>
                </div>
                <div>
                  <SmallText className="text-white font-medium mb-1">Conversion Trigger</SmallText>
                  <BodyText>Series A (target: 18-24 months)</BodyText>
                </div>
              </div>
            </Card>

            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link
                href="/invest/and-you-get"
                className="inline-flex items-center gap-2 px-6 py-3 rounded-lg border border-white/20 text-white hover:bg-white/5 transition-colors"
              >
                Detailed Terms Explanation →
              </Link>
              <Link
                href="/invest/opportunity"
                className="inline-flex items-center gap-2 px-6 py-3 rounded-lg border border-white/10 text-white/60 hover:bg-white/5 transition-colors"
              >
                Market Opportunity
              </Link>
            </div>
          </Container>
        </Section>

        {/* FAQ */}
        <Section background="gradient-subtle">
          <Container>
            <SectionTitle gradient className="text-center">Process FAQ</SectionTitle>
            <SectionSubtitle className="text-center mx-auto mb-8">
              Common questions about the investment process.
            </SectionSubtitle>

            <FAQ items={PROCESS_FAQ} className="mb-8" />
          </Container>
        </Section>

        {/* CTA */}
        <Section>
          <Container>
            <Card variant="highlight" className="py-8">
              <SectionTitle color="#ededed" className="mb-4 text-center">
                Ready to Begin?
              </SectionTitle>
              <BodyText className="max-w-xl mx-auto mb-6 text-center">
                We&apos;ll have your <a href={SAFE_LEARN_URL} target="_blank" rel="noopener noreferrer" className="text-white hover:underline">SAFE</a> document ready within 24 hours of receiving your information.
              </BodyText>

              {/* Visual Input Checklist */}
              <div className="max-w-md mx-auto mb-8">
                <div className="text-xs uppercase tracking-wider text-white/40 mb-3 text-center">
                  Include in your message
                </div>
                <div className="space-y-2">
                  {REQUIRED_INPUTS.map((input, index) => (
                    <div
                      key={input.label}
                      className="flex items-center gap-3 p-3 rounded-lg bg-white/5 border border-white/10"
                    >
                      <div className="w-6 h-6 rounded-full bg-[#2a2a2a] flex items-center justify-center text-xs font-bold text-white">
                        {index + 1}
                      </div>
                      <div className="flex-1">
                        <span className="text-white font-medium text-sm">{input.label}</span>
                        <span className="text-white/30 mx-2">—</span>
                        <span className="text-white/50 text-sm">{input.purpose}</span>
                      </div>
                      <span className="text-lg">{input.icon}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className="text-center">
                <a
                  href="mailto:swarup@planton.ai?subject=Investment Interest - Planton&body=Hi Swarup,%0A%0AI'm interested in investing in Planton.%0A%0ALegal Name (as it should appear on the SAFE): %0AEmail (where Carta should send the agreement): %0AInvestment Amount: $%0A%0A"
                  className="inline-flex items-center gap-2 px-5 py-2.5 rounded-lg bg-[#fff] text-black font-medium text-sm hover:bg-gray-200 transition-all duration-300 hover:-translate-y-0.5"
                >
                  Email swarup@planton.ai →
                </a>

                <SmallText className="mt-6">
                  Or reply to the message where you received this link with these three items.
                </SmallText>
              </div>
            </Card>

            {/* After Investment */}
            <div className="mt-8 text-center">
              <SectionSubtitle className="mb-4">After you invest...</SectionSubtitle>
              <Link
                href="/legal/investor-updates"
                className="inline-flex items-center gap-2 px-4 py-2 rounded-lg border border-white/20 text-white hover:bg-white/5 transition-colors"
              >
                View Investor Updates →
              </Link>
            </div>
          </Container>
        </Section>
      </main>

      <Footer />
    </div>
  );
};

export default ProcessPage;
