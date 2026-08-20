import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import InvestHeader from '@/components/invest/InvestHeader';
import './invest.css';

export const metadata: Metadata = {
  title: 'Invest in Planton - The Self-Service Cloud Platform',
  description:
    'Join us in building the platform that makes cloud infrastructure accessible to every company. Seed stage investment opportunity.',
};

const inter = Inter({
  weight: ['400', '500', '600', '700', '800'],
  subsets: ['latin'],
  display: 'swap',
});

/**
 * Layout for all investor-related pages.
 * 
 * Uses the parent HeaderLogo from micro-apps layout for consistent branding.
 * Adds InvestHeader with Home button for navigation back to /invest.
 * Provides Inter font and isolate context for all /invest/* routes.
 */
export default function InvestLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className={`isolate ${inter.className}`}>
      <InvestHeader />
      {children}
    </div>
  );
}
