'use client';

import { FC } from 'react';
import Link from 'next/link';

interface FooterProps {
  className?: string;
}

const INVEST_PAGES = [
  { href: '/invest', label: 'Home' },
  { href: '/invest/deck', label: 'Deck' },
  { href: '/invest/and-you-get', label: 'What You Get' },
  { href: '/invest/if-you-are', label: 'What We Look For' },
  { href: '/invest/opportunity', label: 'Opportunity' },
  { href: '/invest/process', label: 'Process' },
  { href: '/legal/investor-updates', label: 'Updates' },
];

export const Footer: FC<FooterProps> = ({ className = '' }) => {
  return (
    <footer className={`py-8 border-t border-[#2a2a2a] ${className}`}>
      <div className="w-full max-w-3xl mx-auto px-4 sm:px-6">
        <nav className="flex flex-wrap justify-center gap-4 mb-6">
          {INVEST_PAGES.map((page) => (
            <Link
              key={page.href}
              href={page.href}
              className="text-xs sm:text-sm text-[#666] hover:text-white transition-colors duration-300"
            >
              {page.label}
            </Link>
          ))}
        </nav>

        <p className="text-xs sm:text-sm text-[#666] text-center">
          Last updated: February 2, 2026
        </p>
      </div>
    </footer>
  );
};

export default Footer;
