'use client';

import { FC, Suspense } from 'react';
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
  Callout,
  Metric,
  Badge,
  Grid,
  List,
  Table,
  Step,
  Steps,
  PathAndYouGet,
  formatCurrency,
  VALUATION_CAP,
  FAQItemData,
} from '../shared';
import { useExplainerState, useCalculator } from '../hooks';
import { Hero, Footer, WhatsNext } from '../layout';
import { PathSelector, PathOption } from '../controls';
import { Calculator, ScenarioTable } from '../calculator';
import { FAQ } from '../content';

// ============================================================================
// PATH OPTIONS
// ============================================================================

const PATH_OPTIONS: PathOption<PathAndYouGet>[] = [
  {
    id: 'beginner',
    icon: '🌱',
    title: "I've never invested in a startup before",
    description: "Start with the basics—we'll explain everything from scratch",
  },
  {
    id: 'intermediate',
    icon: '📊',
    title: 'I understand equity, but not SAFEs',
    description: 'You know investments—let us explain our specific structure',
  },
  {
    id: 'advanced',
    icon: '🎯',
    title: "I'm familiar with SAFEs and venture investing",
    description: 'Skip to the terms and numbers',
  },
  {
    id: 'friend',
    icon: '💜',
    title: "I'm a friend or supporter",
    description: 'The personal version—with extra risk warnings',
  },
];

const VALID_PATHS: readonly PathAndYouGet[] = ['beginner', 'intermediate', 'advanced', 'friend'];

// ============================================================================
// FAQ DATA
// ============================================================================

const FAQ_ITEMS: FAQItemData[] = [
  {
    question: "What's a SAFE?",
    answer: "Simple Agreement for Future Equity. You invest now, get ownership later at a discounted price. It's a 3-page legal document created by Y Combinator, now standard for seed-stage investments.",
  },
  {
    question: 'When do I get actual shares?',
    answer: 'At Series A (our target: 18-24 months). Until then, you hold a SAFE—a legal right to shares in the future.',
  },
  {
    question: "What's a valuation cap?",
    answer: <>The maximum price at which your investment converts to shares. Our $7M cap <span className="text-white/50 text-sm">(~$583K ARR / ~$49K MRR)</span> means you convert as if the company is worth $7M maximum, even if we raise Series A at $20M <span className="text-white/50 text-sm">(~$1.7M ARR / ~$139K MRR)</span> or higher.</>,
  },
  {
    question: 'Can I lose everything?',
    answer: <><strong>Yes.</strong> Startup investing is high-risk. You could lose your entire investment. Only invest money you can afford to lose.</>,
  },
  {
    question: "What's the minimum investment?",
    answer: "$15,000-$20,000 for direct investors. For smaller amounts, we work with coordinators who pool multiple investors into a single check of $15K+. This lets us focus on building while still making the opportunity accessible.",
  },
  {
    question: 'Do I get voting rights?',
    answer: "Voting rights and board seats come with actual equity at Series A, when your SAFE converts to shares.",
  },
  {
    question: 'Why SAFE instead of a convertible note?',
    answer: "SAFE has no interest, no maturity date, and simpler terms. A convertible note is debt (loan) with 5-8% annual interest and a deadline. SAFE means no debt hanging over our heads, no artificial deadlines, and full focus on building.",
  },
  {
    question: 'What updates will I receive?',
    answer: 'Monthly investor email (MRR, wins, challenges, asks), quarterly metrics (detailed financials, runway, KPIs), and ad-hoc calls when you can help.',
  },
];

// ============================================================================
// MAIN COMPONENT
// ============================================================================

function AndYouGetPageContent() {
  const { state, selectPath, clearPath, setCurrency, hasPathSelected } = useExplainerState<PathAndYouGet>({
    validPaths: VALID_PATHS,
  });

  const calculator = useCalculator({ currency: state.currency });

  const isPath = (path: PathAndYouGet) => state.selectedPath === path;
  const showSection = (paths: PathAndYouGet[] | 'all') =>
    hasPathSelected && (paths === 'all' || paths.includes(state.selectedPath!));

  return (
    <div className="min-h-screen bg-[#0a0a0a] pt-16">
      <main>
        <Hero
          title="What You Get When You Invest"
          subtitle="You're considering putting money into Planton. Let us explain exactly what that means and what you can expect in return."
        />

        <PathSelector<PathAndYouGet>
          options={PATH_OPTIONS}
          selectedPath={state.selectedPath}
          onSelect={selectPath}
          onClear={clearPath}
          title="Help us explain this in a way that makes sense to you."
          subtitle="How familiar are you with startup investing?"
        />

        <div id="content-start" />

        {/* ================================================================
            BEGINNER PATH CONTENT
            ================================================================ */}
        <Section visible={isPath('beginner')}>
          <Container>
            <SectionTitle>
              What You&apos;re Investing In
            </SectionTitle>

            <Card variant="highlight" className="mb-6">
              <CardTitle color="#ededed">
                Think of it like buying a concert ticket before they announce the price.
              </CardTitle>
            </Card>

            <BodyText className="mb-6">
              <p className="mb-4">
                You give us money now, and when we do our &quot;big funding round&quot; later (called a Series A),
                your money automatically converts into ownership in the company.
              </p>
              <p className="mb-4">
                The cool part? You get a special deal because you believed in us early.
                There&apos;s a &quot;cap&quot; of <strong className="text-white">$7M</strong> <span className="text-white/40 text-sm">(~$583K ARR / ~$49K MRR)</span>—which means
                no matter how valuable we become at Series A, you get to buy in as if we&apos;re only worth $7M.
                If we raise at $20M <span className="text-white/40 text-sm">(~$1.7M ARR / ~$139K MRR)</span>? You still get the $7M price. It&apos;s like locking in early-bird pricing
                before the concert sells out.
              </p>
            </BodyText>

            <Callout variant="highlight" icon="📋" title={<strong>SAFE = Simple Agreement for Future Equity</strong>} className="mb-6">
              It&apos;s a promise: &quot;Give us money now, get ownership later, at a great price because you took the risk early.&quot;
            </Callout>

            <Grid cols={2} className="mb-6">
              <Card>
                <CardTitle className="mb-4">Understanding the SAFE Timeline</CardTitle>
                <BodyText className="mb-4">
                  A SAFE is a right to future ownership, not immediate shares. Here&apos;s what happens:
                </BodyText>
                <List
                  items={[
                    { text: 'Your investment converts to actual shares at Series A (target: 18-24 months)' },
                    { text: 'Voting rights come with those shares at conversion' },
                    { text: 'Returns come at exit (5-8 years), not through dividends' },
                  ]}
                />
              </Card>

              <Card variant="success">
                <CardTitle color="#10b981" className="mb-3">What you get with a SAFE</CardTitle>
                <List
                  items={[
                    { icon: '✓', iconColor: 'emerald', text: 'A legal right to shares in the future at a locked-in price' },
                    { icon: '✓', iconColor: 'emerald', text: 'Early-bird pricing ($7M cap) regardless of how much we grow' },
                    { icon: '✓', iconColor: 'emerald', text: 'Priority over founders if things go wrong (SAFE holders get paid first)' },
                  ]}
                />
                <BodyText className="mt-4">
                  This structure lets us move fast and focus on building rather than governance overhead.
                </BodyText>
              </Card>
            </Grid>

            {/* Valuation Explanation */}
            <Card className="mb-6">
              <CardTitle className="mb-4">&quot;How do you even value a company at $7M?&quot;</CardTitle>
              <BodyText className="mb-4">
                This is a great question. Let us explain how company valuations actually work.
              </BodyText>
              <BodyText className="mb-4">
                Companies are valued based on their revenue. In the DevTools industry, valuations are typically calculated as:
              </BodyText>

              <Callout variant="highlight" className="mb-4">
                <p className="text-center text-lg font-semibold text-white">
                  Annual Revenue × 10 to 15 = Company Valuation
                </p>
              </Callout>

              <BodyText className="mb-4">
                This &quot;10 to 15x&quot; is called a revenue multiple. It&apos;s not something we made up—it&apos;s how
                the market values DevTools companies based on historical data.
              </BodyText>

              <Table
                columns={[
                  { header: 'Annual Revenue', accessor: 'revenue' },
                  { header: 'At 10x Multiple', accessor: 'at10x' },
                  { header: 'At 15x Multiple', accessor: 'at15x' },
                ]}
                data={[
                  { revenue: '$500,000/year', at10x: '$5M valuation', at15x: '$7.5M valuation' },
                  { revenue: '$700,000/year', at10x: <span className="text-white">$7M valuation</span>, at15x: '$10.5M valuation' },
                  { revenue: '$1,000,000/year', at10x: '$10M valuation', at15x: '$15M valuation' },
                ]}
                className="mb-4"
              />

              <Callout variant="warning" icon="📍" title="Where we are now" className="mb-4">
                <p className="mb-2">Current revenue: ~$300-350/month (~$4,200/year)</p>
                <p>Just started getting paying customers 6 months ago</p>
              </Callout>

              <BodyText className="mt-4 mb-4"><strong>The $7M cap assumes:</strong></BodyText>
              <List
                items={[
                  { icon: '→', iconColor: 'pink', text: "We'll reach about $500K-$700K in annual revenue by Series A" },
                  { icon: '→', iconColor: 'pink', text: "That's roughly 40-50 customers paying $15K/year each" },
                  { icon: '→', iconColor: 'pink', text: 'Or 10-15 enterprise customers at $50K/year' },
                ]}
              />

              <BodyText className="mt-4">
                <strong className="text-[#10b981]">The cap protects you.</strong> The harder we work, the more the cap benefits you.
              </BodyText>
            </Card>
          </Container>
        </Section>

        {/* ================================================================
            INTERMEDIATE PATH CONTENT
            ================================================================ */}
        <Section visible={isPath('intermediate')}>
          <Container>
            <SectionTitle>
              What You&apos;re Investing In
            </SectionTitle>

            <Card variant="highlight" className="mb-6">
              <CardTitle color="#ededed">It&apos;s like a convertible note, but simpler.</CardTitle>
            </Card>

            <BodyText className="mb-6">
              If you&apos;ve invested in stocks or equity before, here&apos;s how a SAFE differs from direct equity:
            </BodyText>

            <Table
              columns={[
                { header: 'Aspect', accessor: 'aspect' },
                { header: 'Direct Equity', accessor: 'equity' },
                { header: 'SAFE', accessor: 'safe' },
              ]}
              data={[
                { aspect: <strong>Valuation</strong>, equity: 'Agreed NOW, shares issued immediately', safe: <span className="text-white">Deferred to Series A, shares issued LATER</span> },
                { aspect: <strong>Ownership</strong>, equity: 'You know exactly what % you own today', safe: 'You know the cap, but exact % depends on Series A price' },
                { aspect: <strong>Complexity</strong>, equity: 'Requires full valuation negotiation', safe: <span className="text-[#10b981]">3-page document, straightforward terms</span> },
              ]}
              className="mb-6"
            />

            <SectionTitle className="text-xl">Why we chose SAFE instead of a priced round</SectionTitle>
            <BodyText className="mb-4">
              We believe in staying lean, retaining autonomy, and moving fast. SAFE enables this by design.
            </BodyText>

            <BodyText className="mb-4"><strong>SAFE lets us:</strong></BodyText>
            <List
              items={[
                { icon: '✓', iconColor: 'emerald', text: 'Get capital quickly without valuation debates' },
                { icon: '✓', iconColor: 'emerald', text: 'Focus on building product, not negotiating terms' },
                { icon: '✓', iconColor: 'emerald', text: 'Keep decision-making simple and fast' },
              ]}
              className="mb-6"
            />

            <Grid cols={2}>
              <Card variant="warning">
                <CardTitle color="#a0a0a0">Understanding SAFE Structure</CardTitle>
                <List
                  items={[
                    { icon: '→', iconColor: 'amber', text: 'Shares come at Series A (target: 18-24 months)' },
                    { icon: '→', iconColor: 'amber', text: 'Voting rights arrive with those shares' },
                    { icon: '→', iconColor: 'amber', text: 'Exact ownership % determined at Series A valuation' },
                  ]}
                />
              </Card>

              <Card variant="success">
                <CardTitle color="#10b981">What you get in exchange</CardTitle>
                <List
                  items={[
                    { icon: '✓', iconColor: 'emerald', text: 'Significant pricing discount ($7M cap vs likely $15-25M Series A)' },
                    { icon: '✓', iconColor: 'emerald', text: 'Early access to a company with product already built' },
                    { icon: '✓', iconColor: 'emerald', text: 'Simpler terms, aligned incentives' },
                  ]}
                />
              </Card>
            </Grid>
          </Container>
        </Section>

        {/* ================================================================
            ADVANCED PATH CONTENT
            ================================================================ */}
        <Section visible={isPath('advanced')}>
          <Container>
            <SectionTitle>
              The Terms
            </SectionTitle>

            <Card variant="highlight" className="mb-6">
              <CardTitle color="#ededed">Standard YC post-money SAFE, $7M cap.</CardTitle>
            </Card>

            <Table
              columns={[
                { header: 'Term', accessor: 'term' },
                { header: 'Value', accessor: 'value' },
              ]}
              data={[
                { term: <strong>Valuation cap</strong>, value: formatCurrency(VALUATION_CAP, state.currency) },
                { term: <strong>Discount rate</strong>, value: 'None (cap-only structure)' },
                { term: <strong>SAFE type</strong>, value: 'Post-money (clear ownership calculation)' },
                { term: <strong>MFN clause</strong>, value: 'None (same terms for all investors)' },
                { term: <strong>Pro-rata rights</strong>, value: 'Available for checks $50K+ (follow-on in future rounds)' },
                { term: <strong>Terms negotiable</strong>, value: 'For strategic investors with significant value-add' },
              ]}
              className="mb-6"
            />

            <SectionTitle className="text-xl">Why the $7M cap is attractive</SectionTitle>
            <List
              items={[
                { icon: '→', iconColor: 'pink', text: 'DevTools Series A typically $15-25M valuations' },
                { icon: '→', iconColor: 'pink', text: '$7M cap = seed-stage pricing with significant discount to likely Series A' },
                { icon: '→', iconColor: 'pink', text: 'At $20M Series A: $7M cap gives ~3x more shares than investing at Series A pricing' },
              ]}
              className="mb-6"
            />

            <SectionTitle className="text-xl">Current traction</SectionTitle>
            <Grid cols={2} className="mb-6">
              <Card>
                <CardTitle>Revenue</CardTitle>
                <CardText>~$300-350/month (early subscription customers)</CardText>
              </Card>
              <Card>
                <CardTitle>Timeline</CardTitle>
                <CardText>Just started getting paying customers 6 months ago</CardText>
              </Card>
              <Card>
                <CardTitle>Pipeline</CardTitle>
                <CardText>Active conversations with larger US companies</CardText>
              </Card>
              <Card>
                <CardTitle>Enterprise</CardTitle>
                <CardText>Self-hosted option: targeting $20-30K/year subscriptions</CardText>
              </Card>
            </Grid>

            <Callout variant="success" icon="⚡" title="Technical moat">
              Product is exceptionally well-engineered (Bazel monorepo, protobuf/gRPC). 3-person team doing the work of 20+ engineers with AI tooling.
            </Callout>
          </Container>
        </Section>

        {/* ================================================================
            FRIEND PATH CONTENT
            ================================================================ */}
        <Section visible={isPath('friend')}>
          <Container>
            <Callout variant="danger" icon="⚠️" title="First, the important stuff" className="mb-6">
              <p className="mb-4">
                <strong>Investing in startups is high-risk. Most startups fail. You could lose every dollar you put in.</strong>
              </p>
              <p className="mb-4">
                Only invest money you can genuinely afford to lose—money that, if it disappeared tomorrow, wouldn&apos;t change your life.
              </p>
              <p className="text-white">
                <strong>Our relationship matters more than your investment. If losing this money would strain our friendship, please don&apos;t invest. We mean that.</strong>
              </p>
            </Callout>

            <SectionTitle>
              Now, if you&apos;re still interested...
            </SectionTitle>

            <BodyText className="mb-6">
              When you invest, you&apos;re buying a &quot;SAFE&quot;—a Simple Agreement for Future Equity. It&apos;s basically a promise:
              your money today converts into company ownership later, when we raise a bigger funding round (called Series A).
            </BodyText>

            <Card variant="highlight" className="mb-6">
              <BodyText>
                You get a special deal for believing in us early. There&apos;s a <strong className="text-white">$7M &quot;cap&quot;</strong>—meaning
                no matter how valuable Planton becomes, your money converts as if we were worth only $7M. That&apos;s the reward for taking the risk now.
              </BodyText>
            </Card>

            <Callout icon="⏱️" title="Timeline" className="mb-6">
              Your money sits as a SAFE for 18-24 months until we hit our milestones and raise Series A.
              Then it converts to actual shares. Full exit (when you&apos;d actually see cash) is realistically <strong>5-8 years away</strong>.
            </Callout>
          </Container>
        </Section>

        {/* ================================================================
            THE NUMBERS (All paths)
            ================================================================ */}
        <Section visible={showSection('all')}>
          <Container>
            <SectionTitle gradient className="text-center">The Numbers</SectionTitle>

            <Metric
              value={formatCurrency(500_000, state.currency)}
              label="Total Raise via SAFE Notes"
              sublabel="~18 months runway"
              highlight
              className="mb-8 max-w-sm mx-auto"
            />

            <Grid cols={3} className="mb-8">
              <Card className="text-center">
                <div className="text-2xl font-bold text-white">{formatCurrency(350_000, state.currency)}</div>
                <SmallText>US Institutional (IEP visa)</SmallText>
              </Card>
              <Card className="text-center">
                <div className="text-2xl font-bold text-white">{formatCurrency(150_000, state.currency)}</div>
                <SmallText>Angels & Supporters</SmallText>
              </Card>
              <Card className="text-center">
                <div className="text-2xl font-bold text-[#a0a0a0]">{formatCurrency(15_000, state.currency)}-{formatCurrency(20_000, state.currency)}</div>
                <SmallText>Minimum Check</SmallText>
              </Card>
            </Grid>

            <Callout icon="💡" title="Why this structure?" className="mb-8">
              For smaller amounts, we work with coordinators who pool multiple investors into a single check of $15K+.
              This lets us focus time on building while still making the opportunity accessible. <strong>Sweet spot: $50K-$100K</strong>
            </Callout>

            <Calculator 
              currency={state.currency} 
              onCurrencyChange={setCurrency}
              className="mb-8" 
            />

            <ScenarioTable
              currency={state.currency}
              investmentAmount={calculator.state.investmentAmount}
              scenarios={calculator.scenarioResults}
              className="mb-8"
            />
          </Container>
        </Section>

        {/* ================================================================
            HOW CONVERSION WORKS (All paths)
            ================================================================ */}
        <Section visible={showSection('all')} background="gradient-subtle">
          <Container>
            <SectionTitle gradient className="text-center">How Conversion Works</SectionTitle>
            <SectionSubtitle className="text-center mx-auto mb-8">
              Step-by-step example: Planton raises Series A at $20M valuation <span className="text-white/40 text-sm">(~$1.7M ARR / ~$139K MRR)</span>
            </SectionSubtitle>

            <Steps className="mb-8">
              <Step number={1} title="What's the Series A price per share?">
                <p>
                  $20M ÷ 1,000,000 shares = <strong className="text-white">$20 per share</strong>
                </p>
                <SmallText>This is what new Series A investors pay.</SmallText>
              </Step>

              <Step number={2} title="How does the $7M cap work?">
                <p className="mb-2">As a SAFE holder, you get the LOWER of:</p>
                <ul className="list-disc list-inside mb-2 text-white/60">
                  <li>Series A price ($20/share), OR</li>
                  <li>Cap price ($7M ÷ 1,000,000 shares = $7/share)</li>
                </ul>
                <p>$7 is less than $20, so you convert at <strong className="text-[#10b981]">$7/share</strong>.</p>
              </Step>

              <Step number={3} title="How many shares do you get?">
                <p className="mb-2"><strong>If you invested $1,000:</strong></p>
                <ul className="list-disc list-inside mb-2 text-white/60">
                  <li>Without cap: $1,000 ÷ $20 = <span className="text-white/40">50 shares</span></li>
                  <li>With $7M cap: $1,000 ÷ $7 = <strong className="text-[#10b981]">143 shares</strong></li>
                </ul>
                <Badge variant="success">The cap gives you ~3x more shares for the same investment.</Badge>
              </Step>
            </Steps>

            <Card>
              <CardTitle className="mb-4">How the cap protects you as we grow</CardTitle>
              <Table
                columns={[
                  { header: 'Series A Valuation', accessor: 'valuation' },
                  { header: 'Without Cap', accessor: 'without' },
                  { header: 'With $7M Cap', accessor: 'with' },
                  { header: 'Your Benefit', accessor: 'benefit' },
                ]}
                data={[
                  { valuation: '$7M (at cap)', without: 'Same', with: 'Same', benefit: 'No difference' },
                  { valuation: '$10M', without: '7% ownership', with: '10% ownership', benefit: <span className="text-[#10b981]">+43% more</span> },
                  { valuation: '$20M', without: '3.5% ownership', with: '10% ownership', benefit: <span className="text-[#10b981]">+186% more</span> },
                  { valuation: '$30M', without: '2.3% ownership', with: '10% ownership', benefit: <span className="text-white">+328% more</span> },
                ]}
                highlightRow={3}
              />
              <SmallText className="mt-4">Percentages shown for $700K investment example</SmallText>
            </Card>
          </Container>
        </Section>

        {/* ================================================================
            TIMELINE (All paths)
            ================================================================ */}
        <Section visible={showSection('all')}>
          <Container>
            <SectionTitle gradient className="text-center">Timeline & Milestones</SectionTitle>

            <Grid cols={2} className="mb-8">
              <Card variant="highlight">
                <div className="text-4xl font-bold text-white">18-24</div>
                <CardTitle>Months to Series A</CardTitle>
                <CardText>Your SAFE converts automatically when we raise Series A</CardText>
              </Card>
              <Card>
                <div className="text-4xl font-bold text-white">5-8</div>
                <CardTitle>Years to Exit</CardTitle>
                <CardText>Realistic timeline to see actual cash (acquisition or IPO)</CardText>
              </Card>
            </Grid>

            <Callout variant="success" icon="🎯" title="18-Month Milestones" className="mb-8">
              <Grid cols={2} className="mt-4">
                <div className="flex items-start gap-2">
                  <span className="text-[#10b981]">→</span>
                  <span>50 Enterprise Clients</span>
                </div>
                <div className="flex items-start gap-2">
                  <span className="text-[#10b981]">→</span>
                  <span>$100K Monthly Recurring Revenue</span>
                </div>
                <div className="flex items-start gap-2">
                  <span className="text-[#10b981]">→</span>
                  <span>Planton DevOps AI Agents in Production</span>
                </div>
                <div className="flex items-start gap-2">
                  <span className="text-[#10b981]">→</span>
                  <span>Demonstrated Product-Market Fit</span>
                </div>
              </Grid>
            </Callout>

            <SectionTitle className="text-xl">What If We Don&apos;t Raise Series A?</SectionTitle>
            <Grid cols={2} className="mb-6">
              <Card>
                <Badge variant="success" className="mb-2">Acquisition</Badge>
                <CardText>SAFE converts based on acquisition price. Cap still protects your conversion price.</CardText>
              </Card>
              <Card>
                <Badge variant="purple" className="mb-2">Profitable Bootstrap</Badge>
                <CardText>We&apos;d negotiate a conversion mechanism. This is actually a good problem—means company is successful.</CardText>
              </Card>
              <Card>
                <Badge variant="danger" className="mb-2">Dissolution</Badge>
                <CardText>SAFE holders have preference over founders. Realistically: if no assets, investors likely lose most/all.</CardText>
              </Card>
              <Card>
                <Badge variant="warning" className="mb-2">Zombie (Operating but no raise)</Badge>
                <CardText>SAFE sits unconverted—no payout, no conversion. This is startup risk.</CardText>
              </Card>
            </Grid>
          </Container>
        </Section>

        {/* ================================================================
            RISKS (All paths)
            ================================================================ */}
        <Section visible={showSection('all')} background="danger-subtle">
          <Container>
            <SectionTitle color="#ef4444" className="text-center">What Could Go Wrong</SectionTitle>
            <SectionSubtitle className="text-center mx-auto mb-8">
              Let us be direct: most startups fail.
            </SectionSubtitle>

            <Callout variant="danger" icon="📊" title="Industry Statistics" className="mb-8">
              <List
                items={[
                  { icon: '•', iconColor: 'red', text: <><strong>~90%</strong> of startups fail overall</> },
                  { icon: '•', iconColor: 'red', text: <><strong>~70-80%</strong> of seed-stage companies never reach Series A</> },
                  { icon: '•', iconColor: 'amber', text: <>DevTools/B2B SaaS has <em>better</em> odds than consumer, but still risky</> },
                ]}
              />
            </Callout>

            <SectionTitle className="text-xl">Specific Risks for Planton</SectionTitle>
            <Grid cols={2} className="mb-8">
              <Card variant="danger">
                <CardTitle>Can&apos;t get traction</CardTitle>
                <CardText>Great engineering but no product-market fit</CardText>
              </Card>
              <Card variant="danger">
                <CardTitle>Run out of runway</CardTitle>
                <CardText>18 months passes, milestones not hit, can&apos;t raise Series A</CardText>
              </Card>
              <Card variant="danger">
                <CardTitle>Competition</CardTitle>
                <CardText>HashiCorp, GitLab, AWS add similar features</CardText>
              </Card>
              <Card variant="warning">
                <CardTitle>Founder burnout</CardTitle>
                <CardText>Building a startup is grueling—16-hour days, 3am starts</CardText>
              </Card>
            </Grid>

            <Callout variant="warning" icon="⚠️" title="What happens if Planton fails" className="mb-8">
              SAFE holders have preference in dissolution (paid before founders). Realistically: if company has no assets left,
              investors lose entire investment. This is startup risk—venture capital, not savings bonds.
            </Callout>

            <SectionTitle gradient className="text-xl">Why Invest Despite the Risks?</SectionTitle>
            <Grid cols={2} className="mb-6">
              <Card variant="success">
                <CardTitle color="#10b981">Product exists and works</CardTitle>
                <CardText>Most seed companies are pre-product. We have paying customers.</CardText>
              </Card>
              <Card variant="success">
                <CardTitle color="#10b981">Revenue has started</CardTitle>
                <CardText>Transactions are flowing. ~$300-350/month and growing.</CardText>
              </Card>
              <Card variant="success">
                <CardTitle color="#10b981">Lean operation</CardTitle>
                <CardText>$10-20K/month burn rate. We&apos;re not burning cash recklessly.</CardText>
              </Card>
              <Card variant="success">
                <CardTitle color="#10b981">Technical moat</CardTitle>
                <CardText>Bazel monorepo, protobuf/gRPC architecture. Hard to replicate quickly.</CardText>
              </Card>
            </Grid>

            <Callout variant="highlight" icon="🚀" title="Founder commitment">
              <strong>$500K personal savings invested. 3.5 years building.</strong> This isn&apos;t a hobby.
            </Callout>
          </Container>
        </Section>

        {/* ================================================================
            COMPARABLE COMPANIES (All paths)
            ================================================================ */}
        <Section visible={showSection('all')}>
          <Container>
            <SectionTitle gradient className="text-center">DevTools Success Stories</SectionTitle>
            <SectionSubtitle className="text-center mx-auto mb-8">
              For context on what&apos;s possible in this space
            </SectionSubtitle>

            <Table
              columns={[
                { header: 'Company', accessor: 'company' },
                { header: 'What They Do', accessor: 'description' },
                { header: 'Outcome', accessor: 'outcome' },
              ]}
              data={[
                { company: <strong>GitHub</strong>, description: 'Code hosting & collaboration', outcome: <span className="text-[#10b981]">Acquired by Microsoft for <strong>$7.5B</strong></span> },
                { company: <strong>GitLab</strong>, description: 'DevOps platform', outcome: <span className="text-[#10b981]">IPO at <strong>$11B</strong> valuation</span> },
                { company: <strong>Postman</strong>, description: 'API development tools', outcome: <span className="text-white">Valued at <strong>$5.6B</strong> (from India!)</span> },
                { company: <strong>HashiCorp</strong>, description: 'Infrastructure automation', outcome: <span className="text-[#10b981]">IPO at <strong>~$14B</strong></span> },
                { company: <strong>Datadog</strong>, description: 'Monitoring & observability', outcome: <span className="text-[#10b981]">IPO at <strong>$7.8B</strong>, now ~$40B</span> },
              ]}
              highlightRow={2}
              className="mb-8"
            />

            <Callout variant="highlight" icon="🇮🇳" title="Postman is especially relevant">
              An Indian DevTools company that went global. They proved it can be done from here.
              We&apos;re not saying we&apos;ll become the next Postman—but they proved Indian DevTools companies can build products that compete globally.
            </Callout>
          </Container>
        </Section>

        {/* ================================================================
            FAQs (All paths)
            ================================================================ */}
        <Section visible={showSection('all')} background="gradient-subtle">
          <Container>
            <SectionTitle gradient className="text-center">Frequently Asked Questions</SectionTitle>
            <FAQ items={FAQ_ITEMS} className="mt-8" />
          </Container>
        </Section>

        {/* ================================================================
            INDIAN INVESTOR SECTION (INR only)
            ================================================================ */}
        <Section visible={hasPathSelected && state.currency === 'INR'}>
          <Container>
            <SectionTitle color="#a0a0a0" className="text-center">For Indian Investors</SectionTitle>

            <Callout variant="warning" icon="🇮🇳" title="iSAFE / CCPS Structure" className="mb-6">
              For Indian investors, we use an iSAFE (Indian SAFE) or CCPS (Compulsorily Convertible Preference Shares) structure,
              depending on regulatory requirements. Same economic terms as the US SAFE.
            </Callout>

            <Grid cols={2} className="mb-6">
              <Card>
                <div className="text-2xl font-bold text-[#a0a0a0]">₹58 Cr</div>
                <CardTitle>Valuation Cap</CardTitle>
                <CardText>Equivalent to $7M USD</CardText>
              </Card>
              <Card>
                <div className="text-2xl font-bold text-white">₹25 Lakh+</div>
                <CardTitle>Minimum Investment</CardTitle>
                <CardText>~$30K USD equivalent</CardText>
              </Card>
            </Grid>

            <Callout>
              All examples and calculations on this page work in INR when you toggle the currency.
              Same math, same percentages, same potential returns—just denominated in rupees.
            </Callout>
          </Container>
        </Section>

        {/* ================================================================
            WHAT'S NEXT (All paths)
            ================================================================ */}
        {showSection('all') && <WhatsNext />}
      </main>

      <Footer />
    </div>
  );
}

export const AndYouGetPage: FC = () => {
  return (
    <Suspense fallback={<div className="min-h-screen bg-[#0a0a0a]" />}>
      <AndYouGetPageContent />
    </Suspense>
  );
};

export default AndYouGetPage;
