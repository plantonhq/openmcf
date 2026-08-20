"use client";

import { Inter } from "next/font/google";
import { createTheme } from "@mui/material/styles";
import { colors } from "./colors";

const { primaryLight } = colors;

const inter = Inter({
  weight: ["300", "400", "500", "600", "700"],
  subsets: ["latin"],
  display: "swap",
});

const theme = createTheme({
  cssVariables: {
    colorSchemeSelector: "class",
  },
  colorSchemes: {
    light: true,
    dark: true,
  },
  palette: {
    mode: "dark",
    primary: {
      main: primaryLight[50],
    },
  },
  typography: {
    fontFamily: inter.style.fontFamily,
    allVariants: {
      letterSpacing: "-0.011em",
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
        },
      },
    },
    MuiIcon: {
      defaultProps: {
        baseClassName: "material-symbols-rounded",
      },
      styleOverrides: {
        root: {
          fontVariationSettings: "'FILL' 0, 'wght' 300, 'GRAD' 200, 'opsz' 48",
        },
      },
    },
  },
});

export { theme };
