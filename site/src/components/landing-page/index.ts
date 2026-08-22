// Active landing page version (v4 - infra-first positioning refresh).
// v3 remains on disk for rollback: point this line back at it to revert.
export * from './v4-2026-08-17-1700';

// Legacy v1 layout primitives still used by feature, solution, agent, and hackathon pages.
// Named re-exports (not `export *`) to avoid Turbopack's ESM proxy conflict when
// re-exporting a 'use client' barrel across a server module boundary.
export {
  BlurCard,
  PageMain,
  PageSection,
  PageSectionBackgroundContainer,
  Pill,
  SectionContainer,
  SectionTitle,
} from './v1-2025-02-15-0800/components';

