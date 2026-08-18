'use client';

import { ThemeProvider } from '@mui/material/styles';
import type { Theme } from '@mui/material/styles';
import type { PropsWithChildren } from 'react';
import { websiteTheme } from '../theme/websiteTheme';

interface WebsiteThemeProviderProps extends PropsWithChildren {
  /** Supply a custom theme created via `createWebsiteTheme(overrides)`. */
  theme?: Theme;
}

/**
 * Wraps children in the planton.ai website MUI theme.
 *
 * CssBaseline is intentionally omitted — the host application (planton.ai
 * website or the console) owns global resets. Nesting a second CssBaseline
 * inside a consumer's provider tree would fight the host's baseline styles.
 */
export function WebsiteThemeProvider({
  theme = websiteTheme,
  children,
}: WebsiteThemeProviderProps) {
  return <ThemeProvider theme={theme}>{children}</ThemeProvider>;
}
