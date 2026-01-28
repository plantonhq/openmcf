/**
 * Design Tokens for OpenMCF
 * 
 * A modern, semantic design system inspired by Linear, Vercel, and Stripe.
 * Tokens are organized by purpose rather than arbitrary number scales.
 */

// =============================================================================
// COLORS - Dark Theme
// =============================================================================

export const darkColors = {
  // Surface hierarchy (layered depth system)
  background: {
    base: '#0A0A0B',        // Page background
    raised: '#111113',      // Cards, panels
    elevated: '#18181B',    // Popovers, modals, dropdowns
    overlay: '#27272A',     // Hover states, overlays
    subtle: '#0F0F10',      // Subtle sections
  },

  // Text hierarchy
  text: {
    primary: '#FAFAFA',     // High emphasis - headings, important text
    secondary: '#A1A1AA',   // Medium emphasis - body text
    muted: '#71717A',       // Low emphasis - help text, captions
    disabled: '#52525B',    // Disabled state
    inverse: '#09090B',     // Text on light backgrounds
  },

  // Border system
  border: {
    default: '#27272A',     // Subtle borders
    strong: '#3F3F46',      // Emphasized borders
    focus: '#3B82F6',       // Focus rings
    subtle: '#1F1F23',      // Very subtle borders
  },

  // Brand accent
  accent: {
    primary: '#3B82F6',     // Primary accent (blue)
    hover: '#60A5FA',       // Hover state
    pressed: '#2563EB',     // Active/pressed state
    subtle: 'rgba(59, 130, 246, 0.15)',  // Subtle backgrounds
    muted: 'rgba(59, 130, 246, 0.08)',   // Very subtle tint
  },

  // Semantic colors
  semantic: {
    success: '#22C55E',
    successSubtle: 'rgba(34, 197, 94, 0.15)',
    successMuted: 'rgba(34, 197, 94, 0.08)',
    
    warning: '#F59E0B',
    warningSubtle: 'rgba(245, 158, 11, 0.15)',
    warningMuted: 'rgba(245, 158, 11, 0.08)',
    
    error: '#EF4444',
    errorSubtle: 'rgba(239, 68, 68, 0.15)',
    errorMuted: 'rgba(239, 68, 68, 0.08)',
    
    info: '#06B6D4',
    infoSubtle: 'rgba(6, 182, 212, 0.15)',
    infoMuted: 'rgba(6, 182, 212, 0.08)',
  },

  // Interactive states
  interactive: {
    hover: 'rgba(255, 255, 255, 0.05)',
    active: 'rgba(255, 255, 255, 0.08)',
    selected: 'rgba(59, 130, 246, 0.15)',
  },
};

// =============================================================================
// COLORS - Light Theme
// =============================================================================

export const lightColors = {
  // Surface hierarchy
  background: {
    base: '#FFFFFF',        // Page background
    raised: '#FAFAFA',      // Cards, panels
    elevated: '#FFFFFF',    // Popovers, modals, dropdowns
    overlay: '#F4F4F5',     // Hover states, overlays
    subtle: '#F9FAFB',      // Subtle sections
  },

  // Text hierarchy
  text: {
    primary: '#09090B',     // High emphasis
    secondary: '#52525B',   // Medium emphasis
    muted: '#71717A',       // Low emphasis
    disabled: '#A1A1AA',    // Disabled state
    inverse: '#FAFAFA',     // Text on dark backgrounds
  },

  // Border system
  border: {
    default: '#E4E4E7',     // Subtle borders
    strong: '#D4D4D8',      // Emphasized borders
    focus: '#3B82F6',       // Focus rings
    subtle: '#F4F4F5',      // Very subtle borders
  },

  // Brand accent
  accent: {
    primary: '#2563EB',     // Primary accent (blue - slightly darker for light mode)
    hover: '#3B82F6',       // Hover state
    pressed: '#1D4ED8',     // Active/pressed state
    subtle: 'rgba(37, 99, 235, 0.1)',   // Subtle backgrounds
    muted: 'rgba(37, 99, 235, 0.05)',   // Very subtle tint
  },

  // Semantic colors
  semantic: {
    success: '#16A34A',
    successSubtle: 'rgba(22, 163, 74, 0.1)',
    successMuted: 'rgba(22, 163, 74, 0.05)',
    
    warning: '#D97706',
    warningSubtle: 'rgba(217, 119, 6, 0.1)',
    warningMuted: 'rgba(217, 119, 6, 0.05)',
    
    error: '#DC2626',
    errorSubtle: 'rgba(220, 38, 38, 0.1)',
    errorMuted: 'rgba(220, 38, 38, 0.05)',
    
    info: '#0891B2',
    infoSubtle: 'rgba(8, 145, 178, 0.1)',
    infoMuted: 'rgba(8, 145, 178, 0.05)',
  },

  // Interactive states
  interactive: {
    hover: 'rgba(0, 0, 0, 0.04)',
    active: 'rgba(0, 0, 0, 0.06)',
    selected: 'rgba(37, 99, 235, 0.1)',
  },
};

// =============================================================================
// SPACING
// =============================================================================

export const spacing = {
  px: '1px',
  0: '0',
  0.5: '2px',
  1: '4px',
  1.5: '6px',
  2: '8px',
  2.5: '10px',
  3: '12px',
  3.5: '14px',
  4: '16px',
  5: '20px',
  6: '24px',
  7: '28px',
  8: '32px',
  9: '36px',
  10: '40px',
  11: '44px',
  12: '48px',
  14: '56px',
  16: '64px',
  20: '80px',
  24: '96px',
  28: '112px',
  32: '128px',
};

// =============================================================================
// TYPOGRAPHY
// =============================================================================

export const typography = {
  fontFamily: {
    sans: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    mono: '"JetBrains Mono", "Fira Code", "SF Mono", Consolas, monospace',
  },

  fontSize: {
    xs: '0.75rem',      // 12px
    sm: '0.8125rem',    // 13px
    base: '0.875rem',   // 14px
    md: '0.9375rem',    // 15px
    lg: '1rem',         // 16px
    xl: '1.125rem',     // 18px
    '2xl': '1.25rem',   // 20px
    '3xl': '1.5rem',    // 24px
    '4xl': '2rem',      // 32px
    '5xl': '2.5rem',    // 40px
  },

  fontWeight: {
    normal: 400,
    medium: 500,
    semibold: 600,
    bold: 700,
  },

  lineHeight: {
    none: 1,
    tight: 1.2,
    snug: 1.375,
    normal: 1.5,
    relaxed: 1.625,
    loose: 2,
  },

  letterSpacing: {
    tighter: '-0.02em',
    tight: '-0.01em',
    normal: '0',
    wide: '0.01em',
    wider: '0.02em',
  },
};

// =============================================================================
// SHADOWS
// =============================================================================

export const shadows = {
  none: 'none',
  
  // Subtle shadows for cards and raised elements
  xs: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
  sm: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)',
  md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)',
  lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)',
  xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)',
  '2xl': '0 25px 50px -12px rgba(0, 0, 0, 0.25)',

  // Dark mode specific shadows (with darker opacity)
  dark: {
    xs: '0 1px 2px 0 rgba(0, 0, 0, 0.2)',
    sm: '0 1px 3px 0 rgba(0, 0, 0, 0.3), 0 1px 2px -1px rgba(0, 0, 0, 0.3)',
    md: '0 4px 6px -1px rgba(0, 0, 0, 0.3), 0 2px 4px -2px rgba(0, 0, 0, 0.3)',
    lg: '0 10px 15px -3px rgba(0, 0, 0, 0.3), 0 4px 6px -4px rgba(0, 0, 0, 0.3)',
    xl: '0 20px 25px -5px rgba(0, 0, 0, 0.4), 0 8px 10px -6px rgba(0, 0, 0, 0.4)',
    '2xl': '0 25px 50px -12px rgba(0, 0, 0, 0.5)',
  },

  // Glow effects for buttons and interactive elements
  glow: {
    accent: '0 0 20px rgba(59, 130, 246, 0.3)',
    success: '0 0 20px rgba(34, 197, 94, 0.3)',
    error: '0 0 20px rgba(239, 68, 68, 0.3)',
  },

  // Focus ring shadow
  focus: '0 0 0 2px rgba(59, 130, 246, 0.5)',
  focusError: '0 0 0 2px rgba(239, 68, 68, 0.5)',
};

// =============================================================================
// BORDER RADIUS
// =============================================================================

export const borderRadius = {
  none: '0',
  sm: '4px',
  md: '6px',
  lg: '8px',
  xl: '12px',
  '2xl': '16px',
  '3xl': '24px',
  full: '9999px',
};

// =============================================================================
// TRANSITIONS
// =============================================================================

export const transitions = {
  // Durations
  duration: {
    fast: '100ms',
    normal: '150ms',
    slow: '200ms',
    slower: '300ms',
  },

  // Easing functions
  easing: {
    default: 'cubic-bezier(0.4, 0, 0.2, 1)',
    in: 'cubic-bezier(0.4, 0, 1, 1)',
    out: 'cubic-bezier(0, 0, 0.2, 1)',
    inOut: 'cubic-bezier(0.4, 0, 0.2, 1)',
    spring: 'cubic-bezier(0.34, 1.56, 0.64, 1)',
  },

  // Pre-composed transitions
  all: {
    fast: 'all 100ms cubic-bezier(0.4, 0, 0.2, 1)',
    normal: 'all 150ms cubic-bezier(0.4, 0, 0.2, 1)',
    slow: 'all 200ms cubic-bezier(0.4, 0, 0.2, 1)',
  },

  // Common property transitions
  colors: 'background-color 150ms, border-color 150ms, color 150ms, fill 150ms, stroke 150ms',
  opacity: 'opacity 150ms cubic-bezier(0.4, 0, 0.2, 1)',
  shadow: 'box-shadow 150ms cubic-bezier(0.4, 0, 0.2, 1)',
  transform: 'transform 150ms cubic-bezier(0.4, 0, 0.2, 1)',
};

// =============================================================================
// Z-INDEX
// =============================================================================

export const zIndex = {
  hide: -1,
  base: 0,
  docked: 10,
  dropdown: 1000,
  sticky: 1100,
  banner: 1200,
  overlay: 1300,
  modal: 1400,
  popover: 1500,
  skipLink: 1600,
  toast: 1700,
  tooltip: 1800,
};

// =============================================================================
// LAYOUT
// =============================================================================

export const layout = {
  sidebar: {
    collapsed: 64,
    expanded: 240,
  },
  header: {
    height: 56,
  },
  content: {
    maxWidth: 1400,
    padding: 24,
  },
};

// =============================================================================
// COMPONENT-SPECIFIC TOKENS
// =============================================================================

export const components = {
  button: {
    height: {
      sm: 32,
      md: 36,
      lg: 40,
    },
    padding: {
      sm: '0 12px',
      md: '0 16px',
      lg: '0 20px',
    },
    borderRadius: borderRadius.lg,
  },

  input: {
    height: {
      sm: 32,
      md: 36,
      lg: 40,
    },
    borderRadius: borderRadius.lg,
  },

  card: {
    borderRadius: borderRadius.xl,
    padding: spacing[6],
  },

  table: {
    headerHeight: 44,
    rowHeight: 52,
    compactRowHeight: 40,
  },

  modal: {
    borderRadius: borderRadius['2xl'],
    padding: spacing[6],
  },
};

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Get design tokens for a specific mode
 */
export const getTokens = (mode: 'dark' | 'light') => ({
  colors: mode === 'dark' ? darkColors : lightColors,
  spacing,
  typography,
  shadows,
  borderRadius,
  transitions,
  zIndex,
  layout,
  components,
});

export type DesignTokens = ReturnType<typeof getTokens>;
export type ColorTokens = typeof darkColors;
