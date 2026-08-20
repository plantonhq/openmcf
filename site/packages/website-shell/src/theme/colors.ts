/**
 * Website color palette — the single source of truth for the planton.ai
 * marketing visual identity.
 *
 * Authoritative sources:
 *   - public/branding/design-system.md  (Website Palette table)
 *   - src/theme/colors.ts               (primary ramp)
 *   - src/theme/docs.ts                 (Tailwind class tokens)
 *
 * The grey ramp maps the website's primary scale (src/theme/colors.ts)
 * into MUI's standard 50–900 range. The mapping inverts the source
 * numbering so that lower MUI keys = lighter shades, matching MUI
 * convention (grey[50] is lightest, grey[900] is darkest).
 */

export const websiteGrey = {
  // --- Standard MUI keys (50–900) ---
  50: '#f5f5f5',
  100: '#ededed',
  200: '#d4d4d4',
  300: '#a0a0a0',
  400: '#888888',
  500: '#666666',
  600: '#555555',
  700: '#3a3a3a',
  800: '#2a2a2a',
  900: '#1a1a1a',

  // --- Console-compatible keys ---
  // The console theme uses a non-standard 0–100 grey ramp (steps of 10).
  // Shared styled components (StyledFilterRow, etc.) reference these keys
  // directly via theme.palette.grey[N]. Adding the keys that don't collide
  // with the standard MUI ramp (50, 100, … are already defined above)
  // ensures those components resolve to visible colors in the website theme.
  20: '#111111',  // surface-secondary  (≈ background.paper)
  70: '#1a1a1a',  // input / pill bg    (≈ background.tertiary)
  80: '#1a1a1a',  // interactive hover   (≈ background.tertiary)
} as const;

export const websiteColors = {
  background: {
    default: '#0a0a0a',
    paper: '#111111',
    tertiary: '#1a1a1a',
  },
  text: {
    primary: '#a0a0a0',
    secondary: '#a0a0a0',
    disabled: '#666666',
  },
  border: {
    default: '#2a2a2a',
    hover: '#3a3a3a',
  },
  cta: {
    background: '#ffffff',
    text: '#000000',
  },
} as const;
