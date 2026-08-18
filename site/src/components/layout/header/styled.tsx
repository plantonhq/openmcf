'use client';

import {
  AccordionProps,
  Drawer,
  AccordionSummary as MuiAccordionSummary,
  Accordion as MuiAccordion,
  AccordionDetails as MuiAccordionDetails,
  AccordionSummaryProps,
  SvgIcon,
  SvgIconProps,
} from '@mui/material';
import { styled } from '@mui/material/styles';

export const SvgIconSizeResponsive = styled(SvgIcon)<SvgIconProps>(({ theme }) => ({
  fontSize: 16,
  [theme.breakpoints.up(767)]: {
    fontSize: 24,
  },
}));

export const StyledSvgIcon = styled(SvgIconSizeResponsive)<SvgIconProps>(({ theme }) => ({
  '& *': {
    fill: 'white',
  },
  [theme.breakpoints.up(767)]: {
    '& *': {
      fill: 'currentColor',
    },
  },
}));

export const StyledDrawer = styled(Drawer)(({ theme }) => ({
  '& > .MuiPaper-root': {
    padding: theme.spacing(2.5, 3.5),
    width: '90vw',
    maxWidth: 350,
    backgroundColor: '#111111',
  },
}));

export const Accordion = styled((props: AccordionProps) => (
  <MuiAccordion disableGutters elevation={0} square {...props} />
))(({ theme }) => ({
  padding: 0,
  background: 'transparent',
  '&.Mui-expanded': {
    display: 'flex',
    flexDirection: 'column',
    gap: theme.spacing(3),
  },
  '&::before': {
    display: 'none',
  },
}));

export const AccordionSummary = styled((props: AccordionSummaryProps) => (
  <MuiAccordionSummary {...props} />
))(({ theme }) => ({
  padding: theme.spacing(0),
  minHeight: 'unset',
  '& .MuiAccordionSummary-content': {
    margin: 0,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: theme.spacing(1),
  },
  '& .summary-content': {
    fontSize: theme.spacing(2),
    fontWeight: 600,
    color: '#999999',
  },
  '&.Mui-expanded': {
    '& .summary-content': {
      color: 'white',
    },
  },
}));

export const AccordionDetails = styled(MuiAccordionDetails)(({ theme }) => ({
  padding: theme.spacing(0),
  borderBottom: 'none',
}));
