'use client';

import { FC, useRef, useState, useCallback } from 'react';
import { Box, Typography, IconButton } from '@mui/material';
import { ContentCopy, Check } from '@mui/icons-material';
import { motion, useInView, AnimatePresence } from 'framer-motion';

export interface TerminalLine {
  text: string;
  delay?: number;
  className?: string;
}

interface AnimatedTerminalProps {
  lines: TerminalLine[];
  title?: string;
  trigger?: 'onView' | 'immediate';
  lineDelay?: number;
  className?: string;
}

export const AnimatedTerminal: FC<AnimatedTerminalProps> = ({
  lines,
  title = 'Terminal',
  trigger = 'onView',
  lineDelay = 400,
  className = '',
}) => {
  const ref = useRef<HTMLDivElement>(null);
  const isInView = useInView(ref, { once: true, margin: '-80px' });
  const [copied, setCopied] = useState(false);

  const shouldAnimate = trigger === 'immediate' || isInView;

  const handleCopy = useCallback(() => {
    const text = lines.map((l) => l.text).join('\n');
    navigator.clipboard.writeText(text).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  }, [lines]);

  return (
    <Box
      ref={ref}
      className={`rounded-xl bg-[#0d0d0d] border border-[#2a2a2a] overflow-hidden ${className}`}
    >
      <Box className="flex items-center justify-between px-4 py-2.5 bg-[#1a1a1a] border-b border-[#2a2a2a]">
        <Box className="flex items-center gap-2">
          <Box className="w-2.5 h-2.5 rounded-full bg-[#3a3a3a]" />
          <Box className="w-2.5 h-2.5 rounded-full bg-[#3a3a3a]" />
          <Box className="w-2.5 h-2.5 rounded-full bg-[#3a3a3a]" />
          <Typography className="ml-2 text-xs text-[#555]">{title}</Typography>
        </Box>
        <IconButton onClick={handleCopy} size="small" className="!text-[#555] hover:!text-[#999]">
          {copied ? <Check sx={{ fontSize: 14 }} /> : <ContentCopy sx={{ fontSize: 14 }} />}
        </IconButton>
      </Box>

      <Box className="p-4 font-mono text-[13px] leading-relaxed min-h-[120px]">
        <AnimatePresence>
          {lines.map((line, i) => {
            const delay = line.delay ?? i * lineDelay;
            return (
              <motion.div
                key={i}
                initial={shouldAnimate ? { opacity: 0, x: -8 } : false}
                animate={shouldAnimate ? { opacity: 1, x: 0 } : undefined}
                transition={{
                  duration: 0.3,
                  delay: delay / 1000,
                  ease: 'easeOut',
                }}
                className={line.className ?? 'text-[#b0b0b0]'}
              >
                {line.text || '\u00A0'}
              </motion.div>
            );
          })}
        </AnimatePresence>
      </Box>
    </Box>
  );
};
