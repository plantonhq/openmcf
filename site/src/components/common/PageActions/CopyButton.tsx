import React from 'react';
import { Button } from '@mui/material';
import {
  Check,
  ContentCopy,
  KeyboardArrowDown as ChevronDownIcon,
} from '@mui/icons-material';

interface CopyButtonProps {
  onClick: (event: React.MouseEvent<HTMLElement>) => void;
  copied?: boolean;
}

export const CopyButton: React.FC<CopyButtonProps> = ({ onClick, copied = false }) => {
  return (
    <>
      {/* Compact icon on mobile — subtle, sits next to the page title */}
      <button
        onClick={copied ? undefined : onClick}
        className="md:hidden p-1 text-gray-500 hover:text-white transition-colors flex-shrink-0"
        aria-label={copied ? 'Copied' : 'Page actions'}
      >
        {copied
          ? <Check sx={{ fontSize: 14 }} className="text-green-500" />
          : <ContentCopy sx={{ fontSize: 14 }} />}
      </button>

      {/* Full button on desktop */}
      <Button
        onClick={copied ? undefined : onClick}
        variant="outlined"
        size="small"
        className={`hidden md:inline-flex rounded-lg px-3 py-2 normal-case font-medium shadow-sm transition-all duration-200 ${
          copied
            ? '!border-green-700 !text-green-500'
            : 'hover:!bg-gray-700 text-gray-600 hover:text-white border-gray-300 hover:border-gray-700'
        }`}
        startIcon={
          copied
            ? <Check className="text-green-500" />
            : <ContentCopy className="text-gray-400 hover:text-gray-700 transition-colors duration-200" />
        }
        endIcon={
          copied
            ? undefined
            : <ChevronDownIcon className="text-gray-400 hover:text-gray-700 transition-colors duration-200" />
        }
      >
        {copied ? 'Copied!' : 'Copy page'}
      </Button>
    </>
  );
};
