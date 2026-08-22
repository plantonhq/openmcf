'use client';

import { FC } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home } from 'lucide-react';

interface InvestHeaderProps {
  /** Always show the home button regardless of current path */
  alwaysShow?: boolean;
}

/**
 * Header component for invest pages that adds a Home button on the top right.
 * Only shows on non-landing pages by default (i.e., not on /invest itself).
 * Use alwaysShow=true for pages outside /invest/* (like /legal/investor-updates).
 */
export const InvestHeader: FC<InvestHeaderProps> = ({ alwaysShow = false }) => {
  const pathname = usePathname();
  
  // Don't show home button on the landing page itself (unless alwaysShow is true)
  if (!alwaysShow && (pathname === '/invest' || pathname === '/invest/')) {
    return null;
  }

  return (
    <Link
      href="/invest"
      className="fixed top-[23px] right-8 z-[9999] p-2 rounded-lg bg-white/5 border border-white/10 hover:bg-white/10 hover:border-white/20 transition-all group"
      title="Back to Investor Home"
    >
      <Home className="w-5 h-5 text-white/60 group-hover:text-white transition-colors" />
    </Link>
  );
};

export default InvestHeader;
