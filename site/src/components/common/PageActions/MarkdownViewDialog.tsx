import React, { useEffect, useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Box,
  Typography,
  IconButton,
  Button,
  Stack,
} from '@mui/material';
import { Check, ContentCopy, OpenInNew } from '@mui/icons-material';

interface MarkdownViewDialogProps {
  open: boolean;
  onClose: () => void;
  markdownContent: string;
  title: string;
  onCopyAsMarkdown: () => void;
  onOpenSource: () => void;
}

export const MarkdownViewDialog: React.FC<MarkdownViewDialogProps> = ({
  open,
  onClose,
  markdownContent,
  title,
  onCopyAsMarkdown,
  onOpenSource,
}) => {
  const [copied, setCopied] = useState(false);

  // Reset copied state when dialog closes; auto-reset after 2s on copy.
  useEffect(() => {
    if (!open) {
      const raf = requestAnimationFrame(() => setCopied(false));
      return () => cancelAnimationFrame(raf);
    }
  }, [open]);

  useEffect(() => {
    if (copied) {
      const timer = setTimeout(() => setCopied(false), 2000);
      return () => clearTimeout(timer);
    }
  }, [copied]);

  const handleCopy = () => {
    onCopyAsMarkdown();
    setCopied(true);
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="md"
      fullWidth
      slotProps={{
        paper: {
          className: 'bg-gray-900 border border-gray-700',
        },
      }}
    >
      <DialogTitle className="text-white border-b border-gray-700">
        <Box className="flex items-center justify-between">
          <Typography variant="h6">{title} - Markdown Source</Typography>
          <Stack direction="row" className="gap-2">
            <IconButton
              onClick={copied ? undefined : handleCopy}
              className={`transition-colors duration-200 ${
                copied ? 'text-green-500' : 'text-gray-400 hover:text-white'
              }`}
              size="small"
              aria-label={copied ? 'Copied' : 'Copy as Markdown'}
            >
              {copied ? <Check /> : <ContentCopy />}
            </IconButton>
            <IconButton
              onClick={onOpenSource}
              className="text-gray-400 hover:text-white transition-colors duration-200"
              size="small"
              aria-label="Open raw"
            >
              <OpenInNew />
            </IconButton>
          </Stack>
        </Box>
      </DialogTitle>
      <DialogContent className="p-0">
        <Box className="p-4">
          <pre className="text-sm text-gray-300 bg-gray-800 p-4 rounded-lg overflow-x-auto max-h-96">
            <code>{markdownContent}</code>
          </pre>
        </Box>
      </DialogContent>
      <DialogActions className="border-t border-gray-700 p-4">
        <Button
          onClick={onClose}
          className="text-gray-300 hover:text-white transition-colors duration-200"
        >
          Close
        </Button>
      </DialogActions>
    </Dialog>
  );
};
