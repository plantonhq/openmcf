'use client';

import React from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { IconButton } from '@mui/material';
import { Menu as MenuIcon } from '@mui/icons-material';
import { SearchBar } from '@/app/docs/components/SearchBar';

interface DocsHeaderProps {
  onMenuToggle: () => void;
}

export const DocsHeader: React.FC<DocsHeaderProps> = ({ onMenuToggle }) => {
  return (
    <nav className="fixed top-0 w-full bg-slate-950/95 backdrop-blur-sm border-b border-slate-800 z-50">
      <div className="max-w-full mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Left side: Hamburger (mobile) + Logo */}
          <div className="flex items-center gap-3">
            {/* Mobile hamburger menu - only show on mobile */}
            <div className="md:hidden">
              <IconButton 
                onClick={onMenuToggle} 
                size="small"
                sx={{ color: 'white' }}
              >
                <MenuIcon />
              </IconButton>
            </div>
            
            {/* Logo - same as landing page */}
            <Link href="/" className="flex items-center gap-3">
              <Image 
                src="/icon.png" 
                alt="OpenMCF logo" 
                width={36} 
                height={36} 
                className="h-9 w-auto object-contain" 
                priority 
              />
              <span className="text-xl font-bold text-white hidden sm:block">OpenMCF</span>
            </Link>
          </div>

          {/* Right side: Search */}
          <SearchBar />
        </div>
      </div>
    </nav>
  );
};

