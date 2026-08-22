import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Planton Tour',
  description: 'Interactive tour of the Planton platform',
};

export default function TourLayout({ children }: { children: React.ReactNode }) {
  return <div className="tour-layout">{children}</div>;
}
