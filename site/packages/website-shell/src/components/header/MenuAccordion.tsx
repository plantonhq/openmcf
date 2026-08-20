'use client';

import type { ReactNode, SyntheticEvent } from 'react';
import { Typography } from '@mui/material';
import { ExpandLess, ExpandMore } from '@mui/icons-material';
import {
  ShellAccordion,
  ShellAccordionDetails,
  ShellAccordionSummary,
} from './styled';

interface MenuAccordionProps {
  expanded: boolean;
  title: string;
  children: ReactNode;
  onChange: (event: SyntheticEvent, isExpanded: boolean) => void;
}

export function MenuAccordion({ expanded, title, children, onChange }: MenuAccordionProps) {
  return (
    <ShellAccordion expanded={expanded} onChange={onChange}>
      <ShellAccordionSummary>
        <Typography className="shell-accordion-title">{title}</Typography>
        {expanded ? <ExpandLess fontSize="small" /> : <ExpandMore fontSize="small" />}
      </ShellAccordionSummary>
      <ShellAccordionDetails>{children}</ShellAccordionDetails>
    </ShellAccordion>
  );
}
