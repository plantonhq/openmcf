'use client';

import React, { useEffect, useState } from 'react';
import { CopyButton } from '@/components/common/PageActions/CopyButton';
import { ActionsMenu } from '@/components/common/PageActions/ActionsMenu';
import { MarkdownViewDialog } from '@/components/common/PageActions/MarkdownViewDialog';

interface PageActionsProps {
  markdownContent: string;
  title?: string;
  path: string;
  /** Resolved path to the actual .md file (e.g. "/docs/platform/index.md"). Falls back to `${path}.md`. */
  rawPath?: string;
  hideCopyMarkdown?: boolean;
  hideViewMarkdown?: boolean;
}

export const PageActions: React.FC<PageActionsProps> = ({
  markdownContent,
  title = 'Documentation',
  path,
  rawPath,
  hideCopyMarkdown = false,
  hideViewMarkdown = false,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [showCopied, setShowCopied] = useState(false);
  const [showViewDialog, setShowViewDialog] = useState(false);

  // Auto-reset copied state after 2 seconds
  useEffect(() => {
    if (showCopied) {
      const timer = setTimeout(() => setShowCopied(false), 2000);
      return () => clearTimeout(timer);
    }
  }, [showCopied]);

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleCopyAsMarkdown = async () => {
    try {
      await navigator.clipboard.writeText(markdownContent);
      setShowCopied(true);
    } catch {
      // Fallback for older browsers
      const textArea = document.createElement('textarea');
      textArea.value = markdownContent;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
      setShowCopied(true);
    }
    handleMenuClose();
  };

  const handleOpenRaw = () => {
    window.open(rawPath || `${path}.md`, '_blank');
    handleMenuClose();
  };

  const handleViewAsMarkdown = () => {
    setShowViewDialog(true);
    handleMenuClose();
  };

  const handleCloseViewDialog = () => {
    setShowViewDialog(false);
  };

  if (hideCopyMarkdown && hideViewMarkdown) {
    return <></>;
  }

  return (
    <>
      <CopyButton onClick={handleMenuOpen} copied={showCopied} />

      <ActionsMenu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
        onCopyAsMarkdown={handleCopyAsMarkdown}
        onViewAsMarkdown={handleViewAsMarkdown}
        onOpenRaw={handleOpenRaw}
        hideCopyMarkdown={hideCopyMarkdown}
        hideViewMarkdown={hideViewMarkdown}
      />

      <MarkdownViewDialog
        open={showViewDialog}
        onClose={handleCloseViewDialog}
        markdownContent={markdownContent}
        title={title}
        onCopyAsMarkdown={handleCopyAsMarkdown}
        onOpenSource={handleOpenRaw}
      />
    </>
  );
};
