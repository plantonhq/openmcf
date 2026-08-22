/**
 * Landing page v4 (2026-08-17): the infra-first positioning refresh.
 *
 * What changed from v3, and why:
 * - The hero states the umbrella positioning from src/data/positioning.ts
 *   (never a hub analogy — that is the three-level rule documented there).
 * - The ROI calculator and both comparison tables are GONE, not moved:
 *   they carried unvalidated savings math ($150K DevOps hire) and stale
 *   prices that the claims discipline bans.
 * - Prices and platform stats read from src/data/pricing.ts and
 *   src/data/platform-stats.ts; testimonials are verbatim from their
 *   source of record.
 *
 * Section order mirrors the argument: what it is → how it works → the two
 * halves → why the data is trustworthy → why the security is trustworthy →
 * why there is no lock-in → how to run it → who already does → start.
 */
export { HeroSection } from './HeroSection';
export { HowItWorks } from './HowItWorks';
export { TwoHubs } from './TwoHubs';
export { VerifiedBeforeDeploy } from './VerifiedBeforeDeploy';
export { YourCloudYourControl } from './YourCloudYourControl';
export { OpenSource } from './OpenSource';
export { ThreeWaysToRun } from './ThreeWaysToRun';
export { Proof } from './Proof';
export { FinalCTA } from './FinalCTA';
