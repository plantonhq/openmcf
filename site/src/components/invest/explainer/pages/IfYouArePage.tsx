'use client';

import { FC, Suspense } from 'react';
import {
  Section,
  Container,
  SectionTitle,
  BodyText,
  SmallText,
  Card,
  CardTitle,
  CardText,
  Callout,
  Badge,
  Grid,
  List,
  Table,
  PathIfYouAre,
} from '../shared';
import { useExplainerState } from '../hooks';
import { Hero, Footer, WhatsNext } from '../layout';
import { PathSelector, PathOption } from '../controls';

// ============================================================================
// PATH OPTIONS
// ============================================================================

const PATH_OPTIONS: PathOption<PathIfYouAre>[] = [
  {
    id: 'vc-angel',
    icon: '💼',
    title: "I've invested in DevTools/SaaS companies before",
    description: 'VC fund or angel with relevant portfolio',
  },
  {
    id: 'technical',
    icon: '⚙️',
    title: "I'm technical (engineer/founder)",
    description: 'You understand the tech, want to see the depth',
  },
  {
    id: 'friend',
    icon: '💜',
    title: "I'm a friend or supporter",
    description: 'Personal relationship, want to help',
  },
  {
    id: 'customer',
    icon: '🏢',
    title: 'I work at a company that might need this',
    description: 'Potential customer-investor',
  },
  {
    id: 'general',
    icon: '👀',
    title: 'Just browsing/learning',
    description: 'Curious about the opportunity',
  },
];

const VALID_PATHS: readonly PathIfYouAre[] = ['vc-angel', 'technical', 'friend', 'customer', 'general'];

// ============================================================================
// MAIN COMPONENT
// ============================================================================

function IfYouArePageContent() {
  const { state, selectPath, clearPath, hasPathSelected } = useExplainerState<PathIfYouAre>({
    validPaths: VALID_PATHS,
  });

  const isPath = (path: PathIfYouAre) => state.selectedPath === path;
  const showSection = (paths: PathIfYouAre[] | 'all') =>
    hasPathSelected && (paths === 'all' || paths.includes(state.selectedPath!));

  return (
    <div className="min-h-screen bg-[#0a0a0a] pt-16">
      <main>
        <Hero
          title="What We're Looking For in an Investor"
          subtitle="Who we partner with matters as much as how much we raise. This page explains what makes someone a great fit and how we work with investors."
        />

        <PathSelector<PathIfYouAre>
          options={PATH_OPTIONS}
          selectedPath={state.selectedPath}
          onSelect={selectPath}
          onClear={clearPath}
          title="Tell us about yourself so we can make this relevant."
        />

        <div id="content-start" />

        {/* ================================================================
            OUR IDEAL PARTNER (All paths)
            ================================================================ */}
        <Section visible={showSection('all')}>
          <Container>
            <SectionTitle className="text-center">
              Our Ideal Partner
            </SectionTitle>

            <Card variant="highlight" className="mb-6">
              <BodyText className="text-center">
                Our ideal investor understands the B2B DevTools market and can genuinely appreciate what we&apos;ve built.
                We&apos;re a SaaS/subscription-based developer tools company—think GitHub/GitLab category.
              </BodyText>
            </Card>

            <SectionTitle className="text-xl">What We Value Most</SectionTitle>

            <Grid cols={2} className="mb-8">
              <Card variant="success">
                <CardTitle color="#10b981">Technical understanding</CardTitle>
                <CardText>
                  Not necessarily that you can code, but that you appreciate what it means when a 3-person team builds
                  a Bazel monorepo with protobuf/gRPC architecture. That&apos;s not typical.
                </CardText>
              </Card>
              <Card variant="success">
                <CardTitle color="#10b981">Relevant experience</CardTitle>
                <CardText>
                  B2B DevTools, SaaS, developer platforms. You&apos;ve seen how these companies grow, what works, what doesn&apos;t.
                </CardText>
              </Card>
              <Card variant="success">
                <CardTitle color="#10b981">Network to mid-market companies</CardTitle>
                <CardText>
                  Our sweet spot is companies with 100-1000 employees. Introductions to CTOs and VPs of Engineering are gold.
                </CardText>
              </Card>
              <Card variant="success">
                <CardTitle color="#10b981">Patience</CardTitle>
                <CardText>
                  We&apos;re playing the long game. We&apos;re working fast, but not cutting corners.
                  We&apos;re building for the long term, which aligns best with patient capital.
                </CardText>
              </Card>
            </Grid>
          </Container>
        </Section>

        {/* ================================================================
            GEOGRAPHY SECTION (All paths)
            ================================================================ */}
        <Section visible={showSection('all')} background="gradient-subtle">
          <Container>
            <SectionTitle color="#ededed" className="text-center">Geography Matters</SectionTitle>

            <Callout variant="highlight" icon="🇺🇸" title={<strong>US investors are strongly preferred.</strong>} className="mb-6">
              <p className="mt-2">
                Here&apos;s the honest reason: We&apos;re prioritizing <strong className="text-white">$350K from a US institutional investor</strong> to
                qualify for the International Entrepreneurial Parole (IEP) visa. This enables direct US market access, which is critical for growth.
              </p>
            </Callout>

            <Card className="mb-6">
              <CardTitle className="mb-3">Context</CardTitle>
              <List
                items={[
                  { icon: '→', iconColor: 'pink', text: 'The founding team has Bay Area experience' },
                  { icon: '→', iconColor: 'pink', text: 'Currently based in Hyderabad for focused execution (16-hour days, 3am starts)' },
                  { icon: '→', iconColor: 'pink', text: 'US market access is essential for scaling' },
                ]}
              />
            </Card>

            <SectionTitle className="text-lg">Geographic Priority</SectionTitle>
            <div className="flex gap-2 flex-wrap justify-center">
              <Badge variant="cyan">1. US (priority — visa + market access)</Badge>
              <Badge variant="purple">2. India</Badge>
              <Badge variant="default">3. Europe</Badge>
              <Badge variant="default">4. Australia</Badge>
            </div>
          </Container>
        </Section>

        {/* ================================================================
            VC/ANGEL SPECIFIC SECTION
            ================================================================ */}
        <Section visible={isPath('vc-angel')}>
          <Container>
            <SectionTitle>
              For VC Funds & Angels
            </SectionTitle>

            <Card variant="highlight" className="mb-6">
              <CardTitle color="#ededed">The Dream Investor Profile</CardTitle>
              <CardText className="mt-2">
                Someone from investment firms like Irregular Expressions (SF-based, technical founders investing in DevTools)—a
                combination of capital to accelerate, hunger to build global tech products, and experience/network to make it happen.
              </CardText>
            </Card>

            <SectionTitle className="text-lg">What gets them excited</SectionTitle>
            <List
              items={[
                { icon: '✓', iconColor: 'emerald', text: 'Recognizing that a 3-person team built a product that looks like it came from a 50-person Silicon Valley company' },
                { icon: '✓', iconColor: 'emerald', text: 'Understanding the potential of the underserved mid-market DevOps segment' },
                { icon: '✓', iconColor: 'emerald', text: 'Seeing the engineering depth' },
              ]}
              className="mb-6"
            />

            <SectionTitle className="text-lg">Fund Profile We&apos;re Looking For</SectionTitle>
            <Grid cols={3} className="mb-6">
              <Card className="text-center">
                <div className="text-xl font-bold text-white">$50M-$200M</div>
                <SmallText>Fund Size</SmallText>
              </Card>
              <Card className="text-center">
                <div className="text-xl font-bold text-white">DevTools</div>
                <SmallText>Focus Area</SmallText>
              </Card>
              <Card className="text-center">
                <div className="text-xl font-bold text-[#10b981]">Global</div>
                <SmallText>Indian founders going global</SmallText>
              </Card>
            </Grid>

            <Callout variant="warning" icon="📋" title="Board seat">
              Ideally no board seat at seed stage (SAFE = no governance). Observer seat only if required.
            </Callout>
          </Container>
        </Section>

        {/* ================================================================
            TECHNICAL INVESTOR SECTION
            ================================================================ */}
        <Section visible={isPath('technical')}>
          <Container>
            <SectionTitle>
              For Technical Investors
            </SectionTitle>

            <Card variant="cyan" className="mb-6">
              <BodyText>
                While we&apos;re not seeking DevOps coaching, if you&apos;re a high-quality software engineer passionate about
                cloud-native software—that earns respect.
              </BodyText>
            </Card>

            <SectionTitle className="text-lg">How technical involvement works</SectionTitle>
            <Grid cols={2} className="mb-6">
              <Card variant="success">
                <CardTitle color="#10b981">High involvement welcome</CardTitle>
                <CardText>
                  If you have a proven background operating DevTools companies, we&apos;d love you involved in decision-making
                </CardText>
              </Card>
              <Card>
                <CardTitle>Advisory welcome</CardTitle>
                <CardText>If you&apos;re technical but not in our domain, advisory is welcome</CardText>
              </Card>
            </Grid>

            <Callout>
              <strong>Technical credibility is earned, not assumed.</strong>
            </Callout>
          </Container>
        </Section>

        {/* ================================================================
            FRIEND INVESTOR SECTION
            ================================================================ */}
        <Section visible={isPath('friend')}>
          <Container>
            <SectionTitle>
              For Friends & Supporters
            </SectionTitle>

            <Callout variant="danger" icon="⚠️" title="Important: Read This First" className="mb-6">
              <p className="mb-4">Friends are already supporting us (borrowing money to stay bootstrapped). If you&apos;re considering investing:</p>
              <List
                items={[
                  { icon: '!', iconColor: 'red', text: <><strong>Only money you can afford to lose</strong> (startup risk is real)</> },
                  { icon: '!', iconColor: 'red', text: 'No liquidity for years (5-10 year horizon)' },
                  { icon: '!', iconColor: 'red', text: "We'll treat you like any other investor" },
                  { icon: '!', iconColor: 'red', text: "Don't invest more than you're okay losing" },
                ]}
              />
            </Callout>

            <SectionTitle className="text-lg">Guidelines for Friends</SectionTitle>
            <Table
              columns={[
                { header: '', accessor: 'label' },
                { header: '', accessor: 'value' },
              ]}
              data={[
                { label: <strong>Typical check size</strong>, value: '$10K-$25K' },
                { label: <strong>Terms</strong>, value: 'Same as everyone else ($7M cap SAFE)' },
                { label: <strong>Involvement expected</strong>, value: 'None (unless you have relevant expertise)' },
                { label: <strong>Priority</strong>, value: 'Relationship comes first' },
              ]}
              className="mb-6"
            />

            <Callout variant="highlight">
              <p className="text-center">
                <strong className="text-white">Our relationship matters more to us than your investment.</strong>
              </p>
            </Callout>
          </Container>
        </Section>

        {/* ================================================================
            CUSTOMER INVESTOR SECTION
            ================================================================ */}
        <Section visible={isPath('customer')}>
          <Container>
            <SectionTitle>
              For Potential Customer-Investors
            </SectionTitle>

            <Card variant="highlight" className="mb-6">
              <BodyText>
                If your company might need Planton, you&apos;re exactly who we want to talk to.
                Customer-investors have the best signal on product-market fit.
              </BodyText>
            </Card>

            <SectionTitle className="text-lg">Why this is valuable</SectionTitle>
            <List
              items={[
                { icon: '✓', iconColor: 'emerald', text: 'You can evaluate the product firsthand' },
                { icon: '✓', iconColor: 'emerald', text: 'Your feedback directly shapes the roadmap' },
                { icon: '✓', iconColor: 'emerald', text: 'Reference customer + investor = powerful signal' },
              ]}
              className="mb-6"
            />

            <SectionTitle className="text-lg">Strategic/Corporate Investment Considerations</SectionTitle>
            <Grid cols={2}>
              <Card variant="success">
                <CardTitle color="#10b981">Potentially interesting</CardTitle>
                <List
                  items={[
                    { icon: '→', iconColor: 'emerald', text: 'Cloud providers (AWS, GCP, Azure) — marketplace, credits' },
                    { icon: '→', iconColor: 'emerald', text: 'DevTools companies (non-competitive) — integration' },
                  ]}
                />
              </Card>
              <Card variant="warning">
                <CardTitle color="#a0a0a0">Concerns</CardTitle>
                <List
                  items={[
                    { icon: '!', iconColor: 'amber', text: 'Strong value of autonomy—careful about limiting future options' },
                    { icon: '!', iconColor: 'amber', text: 'No exclusivity clauses' },
                  ]}
                />
              </Card>
            </Grid>
          </Container>
        </Section>

        {/* ================================================================
            WHAT WE NEED BEYOND CAPITAL (All paths)
            ================================================================ */}
        <Section visible={showSection('all')}>
          <Container>
            <SectionTitle gradient className="text-center">What We Need Beyond Capital</SectionTitle>

            <Grid cols={2} className="mb-8">
              <Card>
                <Badge variant="pink" className="mb-3">Priority #1</Badge>
                <CardTitle>Market Access</CardTitle>
                <CardText className="mb-3">
                  Mid-market is our sweet spot (100-1000 employees). They want control over their architecture, no real solution serving them.
                </CardText>
                <SmallText>How you can help: Introductions to CTOs/VPs of Engineering at mid-market companies</SmallText>
              </Card>

              <Card>
                <Badge variant="cyan" className="mb-3">Priority #2</Badge>
                <CardTitle>US Institutional Investment</CardTitle>
                <CardText className="mb-3">
                  $350K+ from a US investor qualifies for IEP visa—critical path to US market.
                </CardText>
                <SmallText>How you can help: Connections to seed-stage DevTools funds</SmallText>
              </Card>

              <Card>
                <Badge variant="purple" className="mb-3">Valuable</Badge>
                <CardTitle>Operational Guidance</CardTitle>
                <CardText>GTM strategy for DevTools, scaling/hiring, Series A prep, operations/finance, legal/compliance navigation</CardText>
              </Card>

              <Card>
                <Badge variant="default" className="mb-3">Helpful</Badge>
                <CardTitle>Network Effects</CardTitle>
                <CardText>DevTools GTM/sales experts, senior platform engineers (future hiring), CNCF/DevOps community leaders</CardText>
              </Card>
            </Grid>
          </Container>
        </Section>

        {/* ================================================================
            DEALBREAKERS (All paths)
            ================================================================ */}
        <Section visible={showSection('all')} background="danger-subtle">
          <Container>
            <SectionTitle color="#ef4444" className="text-center">What We&apos;re Looking For</SectionTitle>

            <SectionTitle className="text-lg">Non-Negotiable Requirements</SectionTitle>
            <Grid cols={2} className="mb-8">
              <Card variant="success">
                <CardTitle color="#10b981">Believes in the vision</CardTitle>
                <CardText>Genuinely excited about building a DevTools company, not just looking for returns</CardText>
              </Card>
              <Card variant="success">
                <CardTitle color="#10b981">Proven relevant background</CardTitle>
                <CardText>Credible, verifiable experience in B2B tech, SaaS, or DevTools. Depth and relevance matter more than portfolio size.</CardText>
              </Card>
              <Card variant="success">
                <CardTitle color="#10b981">Tech-savvy</CardTitle>
                <CardText>Can appreciate the engineering depth. Doesn&apos;t need to code, but understands what quality looks like.</CardText>
              </Card>
              <Card variant="success">
                <CardTitle color="#10b981">Respects founder autonomy</CardTitle>
                <CardText>Comfortable with founder-led decision-making. Advisory input welcome; we value collaborative partnership.</CardText>
              </Card>
            </Grid>
          </Container>
        </Section>

        {/* ================================================================
            CHECK SIZES & STRUCTURE (All paths)
            ================================================================ */}
        <Section visible={showSection('all')}>
          <Container>
            <SectionTitle gradient className="text-center">Check Sizes & Structure</SectionTitle>

            <Table
              columns={[
                { header: '', accessor: 'label' },
                { header: '', accessor: 'value' },
              ]}
              data={[
                { label: <strong>Minimum</strong>, value: <span className="text-[#a0a0a0]">$15K-$20K</span> },
                { label: <strong>Sweet spot</strong>, value: <span className="text-white">$50K-$100K</span> },
                { label: <strong>Maximum (single investor)</strong>, value: <span className="text-white">$350K (institutional)</span> },
                { label: <strong>Total raise</strong>, value: <strong>$500K</strong> },
                { label: <strong>Close type</strong>, value: 'Rolling (no artificial deadline)' },
              ]}
              highlightRow={1}
              className="mb-8"
            />

            <SectionTitle className="text-lg">Why These Numbers</SectionTitle>
            <Grid cols={3} className="mb-8">
              <Card>
                <Badge variant="warning" className="mb-2">Minimum</Badge>
                <CardText>
                  Our minimum of $15K-$20K helps us maintain efficient operations while serving all investors well.
                  Exception: pooled investors coordinated into one check.
                </CardText>
              </Card>
              <Card variant="highlight">
                <Badge variant="pink" className="mb-2">Sweet Spot</Badge>
                <CardText>Meaningful capital, reasonable investor count, manageable communication.</CardText>
              </Card>
              <Card>
                <Badge variant="cyan" className="mb-2">Maximum</Badge>
                <CardText>One US institutional investor at $350K qualifies for IEP visa. That&apos;s the priority.</CardText>
              </Card>
            </Grid>

            <SectionTitle className="text-lg">Ideal Mix</SectionTitle>
            <Callout className="mb-6">
              <List
                items={[
                  { icon: '→', iconColor: 'cyan', text: <><strong>1-2 institutional investors</strong> at $150K-$350K (anchors, IEP qualification)</> },
                  { icon: '→', iconColor: 'pink', text: <><strong>3-5 angels</strong> at $25K-$100K (network, expertise)</> },
                  { icon: '→', text: <><strong>Optional: friends/supporters</strong> at $15K-$50K</> },
                ]}
              />
            </Callout>

            <Callout variant="highlight" icon="📋" title="Rolling Close">
              <p className="mt-2">
                We&apos;re cash-constrained, so we&apos;re taking investments as they come. Each investor signs their own SAFE when ready.
                Same terms for everyone ($7M cap). No &quot;lead investor&quot; requirement. <strong>No artificial urgency.</strong> Better to find right investors than rush.
              </p>
            </Callout>
          </Container>
        </Section>

        {/* ================================================================
            WHAT ACTUALLY MATTERS (All paths)
            ================================================================ */}
        <Section visible={showSection('all')} background="gradient-subtle">
          <Container>
            <SectionTitle gradient className="text-center">What Actually Matters</SectionTitle>

            <Card className="mb-6">
              <CardTitle className="mb-4">What Makes a Strong Partner</CardTitle>
              <BodyText className="mb-4">The backgrounds that help us most:</BodyText>
              <List
                items={[
                  { icon: '✓', iconColor: 'emerald', text: 'Relevant operational experience in B2B/DevTools/SaaS' },
                  { icon: '✓', iconColor: 'emerald', text: 'Ability to open doors at mid-market companies' },
                  { icon: '✓', iconColor: 'emerald', text: "Technical understanding of what we're building" },
                  { icon: '✓', iconColor: 'emerald', text: 'Shared vision and patience for the journey' },
                ]}
              />
              <BodyText className="mt-4 italic">
                We care more about what you can contribute than credentials or titles.
              </BodyText>
            </Card>

            <SectionTitle className="text-lg">Common Misconceptions</SectionTitle>
            <Grid cols={2}>
              <Card>
                <CardText><strong>&quot;They just want money&quot;</strong><br />No. We want partners who understand and can help.</CardText>
              </Card>
              <Card>
                <CardText><strong>&quot;They need millions&quot;</strong><br />$100K gives 6 months runway. We&apos;re capital-efficient.</CardText>
              </Card>
              <Card>
                <CardText><strong>&quot;They&apos;re desperate&quot;</strong><br />We have clear criteria and prioritize finding partners who align with our values.</CardText>
              </Card>
              <Card>
                <CardText><strong>&quot;They want hands-off investors&quot;</strong><br />We want engaged investors IF they have relevant expertise.</CardText>
              </Card>
            </Grid>
          </Container>
        </Section>

        {/* ================================================================
            POST-INVESTMENT EXPECTATIONS (All paths)
            ================================================================ */}
        <Section visible={showSection('all')}>
          <Container>
            <SectionTitle gradient className="text-center">Post-Investment Expectations</SectionTitle>

            <Grid cols={2} className="mb-8">
              <Card>
                <CardTitle color="#ededed">What You Can Expect From Us</CardTitle>
                <BodyText className="mb-3"><strong>Communication:</strong></BodyText>
                <List
                  items={[
                    { icon: '✓', iconColor: 'emerald', text: 'Monthly written update (MRR, wins, challenges, asks)' },
                    { icon: '✓', iconColor: 'emerald', text: 'Quarterly metrics email (detailed financials, runway, KPIs)' },
                    { icon: '✓', iconColor: 'emerald', text: 'Available for calls when you have something valuable to offer' },
                    { icon: '✕', iconColor: 'red', text: 'No real-time dashboard (overhead to maintain)' },
                  ]}
                />
                <SmallText className="text-[#10b981] mt-3">Honest about wins AND challenges.</SmallText>
              </Card>

              <Card>
                <CardTitle color="#ededed">What We Value From Investors</CardTitle>
                <BodyText className="mb-3"><strong>Helpful:</strong></BodyText>
                <List
                  items={[
                    { icon: '✓', iconColor: 'emerald', text: 'Customer introductions when you know someone relevant' },
                    { icon: '✓', iconColor: 'emerald', text: 'Strategic advice when you have domain expertise' },
                    { icon: '✓', iconColor: 'emerald', text: 'Network introductions proactively' },
                  ]}
                />
                <BodyText className="mt-4 italic">
                  The most helpful investors know when to engage and when to let us execute.
                </BodyText>
              </Card>
            </Grid>

            <SectionTitle className="text-lg">Governance</SectionTitle>
            <Table
              columns={[
                { header: '', accessor: 'label' },
                { header: '', accessor: 'value' },
              ]}
              data={[
                { label: <strong>Board seat</strong>, value: 'No (SAFEs = no governance)' },
                { label: <strong>Observer rights</strong>, value: 'Optional for larger investors' },
                { label: <strong>Monthly calls</strong>, value: 'Only if investor has expertise to share' },
                { label: <strong>Quarterly reviews</strong>, value: 'Optional—written updates preferred' },
                { label: <strong>Ad-hoc availability</strong>, value: 'Yes—when you can genuinely help' },
              ]}
              className="mb-6"
            />

            <Callout>
              <p className="text-center"><strong>Involvement is earned, not expected.</strong></p>
            </Callout>

            <SectionTitle className="text-lg mt-8">Boundaries</SectionTitle>
            <Grid cols={2}>
              <Card>
                <CardTitle>Response time</CardTitle>
                <CardText>24-48 hours for non-urgent. Same day if genuinely urgent.</CardText>
              </Card>
              <Card>
                <CardTitle>Communication</CardTitle>
                <CardText>Email preferred (async). Scheduled calls for complex discussions.</CardText>
              </Card>
            </Grid>

            <Callout variant="highlight" className="mt-6">
              <p className="text-center italic">&quot;At every step we value mindful health and avoiding stress at all costs.&quot;</p>
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

export const IfYouArePage: FC = () => {
  return (
    <Suspense fallback={<div className="min-h-screen bg-[#0a0a0a]" />}>
      <IfYouArePageContent />
    </Suspense>
  );
};

export default IfYouArePage;
