import { createTheme, alpha, type ThemeOptions } from '@mui/material/styles';
import { tabClasses, buttonClasses } from '@mui/material';
import { websiteColors, websiteGrey } from './colors';

/**
 * Default font stack for the website theme.
 *
 * Uses the CSS variable set by next/font/google in the host app's root
 * layout. Falls back through the standard system font stack so the theme
 * works correctly even outside a Next.js context (e.g. Storybook, tests).
 */
const WEBSITE_FONT_FAMILY =
  'var(--font-inter), Inter, system-ui, -apple-system, sans-serif';

/**
 * Base ThemeOptions that define the planton.ai marketing visual identity.
 *
 * This is intentionally a plain object (not a Theme instance) so that
 * `createWebsiteTheme` can deep-merge caller overrides before calling
 * `createTheme` once — avoiding the double-creation penalty.
 */
const baseThemeOptions: ThemeOptions = {
  breakpoints: {
    values: { xs: 0, sm: 640, md: 768, lg: 1024, xl: 1280 },
  },

  palette: {
    mode: 'dark',
    background: {
      default: websiteColors.background.default,
      paper: websiteColors.background.paper,
    },
    text: {
      primary: websiteColors.text.primary,
      secondary: websiteColors.text.secondary,
      disabled: websiteColors.text.disabled,
    },
    divider: websiteColors.border.default,
    primary: {
      main: websiteColors.cta.background,
      contrastText: websiteColors.cta.text,
    },
    secondary: {
      main: websiteColors.background.tertiary,
      light: websiteGrey[300],
      dark: websiteColors.background.paper,
      contrastText: websiteColors.text.primary,
    },
    grey: websiteGrey,
  },

  typography: {
    fontFamily: WEBSITE_FONT_FAMILY,
    allVariants: {
      letterSpacing: '-0.011em',
    },
  },

  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none' as const,
        },
        containedSecondary: {
          color: websiteColors.text.primary,
          backgroundColor: websiteColors.background.tertiary,
          border: `1px solid ${websiteColors.border.default}`,
          '&:hover': {
            backgroundColor: websiteColors.border.hover,
            boxShadow: 'none',
          },
        },
      },
    },

    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
        },
      },
    },

    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
        },
      },
    },

    MuiCard: {
      styleOverrides: {
        root: ({ theme }) => ({
          border: `1px solid ${theme.palette.divider}`,
        }),
      },
    },

    MuiIcon: {
      defaultProps: {
        baseClassName: 'material-symbols-rounded',
      },
      styleOverrides: {
        root: {
          fontVariationSettings:
            "'FILL' 0, 'wght' 300, 'GRAD' 200, 'opsz' 48",
        },
      },
    },

    MuiTabs: {
      styleOverrides: {
        root: ({ theme }) => ({
          minHeight: 40,
          width: 'fit-content',
          borderRadius: 0,
          backgroundColor: theme.palette.background.paper,
          padding: theme.spacing(1.5),
          alignItems: 'center',
          '& .MuiTabs-scroller': {
            overflow: 'visible !important',
          },
        }),
        flexContainer: ({ theme }) => ({
          gap: theme.spacing(1.5),
        }),
        indicator: ({ theme }) => ({
          bottom: theme.spacing(-1.5),
        }),
      },
    },

    MuiTab: {
      styleOverrides: {
        root: ({ theme }) => ({
          minHeight: 0,
          backgroundColor: 'transparent',
          color: theme.palette.text.primary,
          cursor: 'pointer',
          fontSize: theme.spacing(1.5),
          fontWeight: 500,
          padding: theme.spacing(1, 1.5),
          borderRadius: theme.shape.borderRadius + 2,
          textTransform: 'none' as const,
          [`& .${tabClasses.icon}`]: {
            marginBottom: 0,
            color: theme.palette.text.secondary,
          },
          [`&.${tabClasses.selected}, &:hover`]: {
            backgroundColor: theme.palette.grey[900],
            color: theme.palette.text.primary,
            [`& .${tabClasses.icon}`]: {
              color: theme.palette.text.primary,
            },
          },
          [`&.${buttonClasses.disabled}`]: {
            opacity: 0.5,
            cursor: 'not-allowed',
          },
        }),
      },
    },

    MuiInputBase: {
      styleOverrides: {
        root: ({ theme }) => ({
          minHeight: 34,
          backgroundColor: theme.palette.background.paper,
          borderRadius: `${theme.shape.borderRadius * 2}px !important`,
          margin: '4px',
          '& fieldset': {
            borderColor: theme.palette.divider,
          },
          '&:hover:not(.Mui-disabled):not(.Mui-focused):not(.Mui-error)': {
            '& fieldset': {
              borderColor: theme.palette.text.primary,
              transition: 'border-color 0.5s ease-in-out',
            },
          },
          '&.Mui-focused': {
            '&.Mui-error': {
              '& fieldset': {
                boxShadow: `0 0 0 3px ${alpha(theme.palette.error.main, 0.2)}`,
                transition: 'box-shadow 0.5s ease-in-out',
              },
            },
            '&:not(.Mui-error)': {
              '& fieldset': {
                boxShadow: `0 0 0 3px ${alpha(theme.palette.text.secondary, 0.25)}`,
                borderColor: `${theme.palette.text.primary} !important`,
                transition: 'box-shadow 0.5s ease-in-out',
              },
            },
          },
          '& .MuiInputBase-input': {
            fontSize: 12,
            fontWeight: 400,
            padding: '8.5px 8px',
            '&.Mui-disabled': {
              color: theme.palette.text.secondary,
            },
          },
          '& .MuiSelect-iconOutlined': {
            color: theme.palette.text.secondary,
          },
          '&.MuiAutocomplete-inputRoot': {
            padding: 0,
          },
          '& .MuiInputAdornment-root': {
            marginRight: 0,
          },
        }),
      },
    },

    MuiCheckbox: {
      styleOverrides: {
        root: ({ theme }) => ({
          '& .MuiSvgIcon-root': {
            height: '1.3rem',
          },
          '&.Mui-checked .MuiSvgIcon-root path': {
            fill: theme.palette.text.primary,
          },
        }),
      },
    },

    MuiTooltip: {
      styleOverrides: {
        tooltip: ({ theme }) => ({
          backgroundColor: theme.palette.grey[900],
          color: theme.palette.text.primary,
          fontSize: 12,
          fontWeight: 400,
          border: `1px solid ${theme.palette.divider}`,
        }),
        arrow: ({ theme }) => ({
          color: theme.palette.grey[900],
          '&::before': {
            border: `1px solid ${theme.palette.divider}`,
          },
        }),
      },
    },

    MuiDivider: {
      styleOverrides: {
        root: ({ theme }) => ({
          '&:before': {
            borderColor: theme.palette.divider,
          },
          '&:after': {
            borderColor: theme.palette.divider,
          },
        }),
        vertical: ({ theme }) => ({
          borderWidth: '1px',
          margin: '0 10px',
          borderColor: theme.palette.divider,
        }),
      },
    },

    MuiSvgIcon: {
      styleOverrides: {
        fontSizeSmall: {
          fontSize: '16px',
        },
      },
    },

    MuiLink: {
      styleOverrides: {
        root: ({ theme }) => ({
          color: theme.palette.text.primary,
          textDecorationColor: theme.palette.text.primary,
        }),
      },
    },

    MuiFormControlLabel: {
      styleOverrides: {
        label: {
          lineHeight: 1,
        },
      },
    },
  },
};

/**
 * Create a website theme with optional overrides.
 *
 * The factory accepts a partial `ThemeOptions` that is shallow-merged at
 * the top level with the base options. Use this when the consumer needs a
 * minor tweak (e.g. a different font family) without forking the palette
 * or component overrides.
 *
 * For deep customisation, pass a full `ThemeOptions` — `createTheme` will
 * deep-merge palettes and component overrides internally.
 */
export function createWebsiteTheme(overrides?: ThemeOptions) {
  if (!overrides) {
    return createTheme(baseThemeOptions);
  }
  return createTheme(baseThemeOptions, overrides);
}

/** Default website theme instance — ready to use with no overrides. */
export const websiteTheme = createWebsiteTheme();
