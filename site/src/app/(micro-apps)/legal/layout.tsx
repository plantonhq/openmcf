import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import InvestHeader from '@/components/invest/InvestHeader';

export const metadata: Metadata = {
  title: 'Legal - Planton',
  description: 'Legal information and investor updates for Planton',
};

const inter = Inter({
  weight: ['400', '500', '600', '700', '800'],
  subsets: ['latin'],
  display: 'swap',
});

/**
 * Layout for legal pages including investor updates.
 * Uses the parent HeaderLogo from micro-apps layout for consistent branding.
 * Adds InvestHeader with Home button for navigation back to /invest.
 */
export default function LegalLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className={`isolate ${inter.className}`}>
      <InvestHeader alwaysShow />
      {children}
    </div>
  );
}
