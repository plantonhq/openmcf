import React from 'react';
import {
  Menu,
  MenuItem,
  ListItemIcon,
  ListItemText,
} from '@mui/material';
import {
  ContentCopy,
  OpenInNew,
  Visibility as ViewIcon,
} from '@mui/icons-material';

interface ActionsMenuProps {
  anchorEl: HTMLElement | null;
  open: boolean;
  onClose: () => void;
  onCopyAsMarkdown: () => void;
  onViewAsMarkdown: () => void;
  onOpenRaw: () => void;
  hideCopyMarkdown?: boolean;
  hideViewMarkdown?: boolean;
}

export const ActionsMenu: React.FC<ActionsMenuProps> = ({
  anchorEl,
  open,
  onClose,
  onCopyAsMarkdown,
  onViewAsMarkdown,
  onOpenRaw,
  hideCopyMarkdown = false,
  hideViewMarkdown = false,
}) => {
  // If both options are hidden, don't render the menu
  if (hideCopyMarkdown && hideViewMarkdown) {
    return null;
  }

  return (
    <Menu
      anchorEl={anchorEl}
      open={open}
      onClose={onClose}
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'right',
      }}
      transformOrigin={{
        vertical: 'top',
        horizontal: 'right',
      }}
      slotProps={{
        paper: {
          className: 'bg-[#111] border border-[#2a2a2a]',
        },
      }}
    >
      {!hideCopyMarkdown && (
        <MenuItem onClick={onCopyAsMarkdown} className="text-[#a0a0a0] hover:!bg-white/10">
          <ListItemIcon>
            <ContentCopy className="text-[#a0a0a0]" />
          </ListItemIcon>
          <ListItemText primary="Copy as Markdown" />
        </MenuItem>
      )}
      {!hideViewMarkdown && (
        <MenuItem onClick={onViewAsMarkdown} className="text-[#a0a0a0] hover:!bg-white/10">
          <ListItemIcon>
            <ViewIcon className="text-[#a0a0a0]" />
          </ListItemIcon>
          <ListItemText primary="View as Markdown" />
        </MenuItem>
      )}
      <MenuItem onClick={onOpenRaw} className="text-[#a0a0a0] hover:!bg-white/10">
        <ListItemIcon>
          <OpenInNew className="text-[#a0a0a0]" />
        </ListItemIcon>
        <ListItemText primary="Open Raw" />
      </MenuItem>
    </Menu>
  );
};
