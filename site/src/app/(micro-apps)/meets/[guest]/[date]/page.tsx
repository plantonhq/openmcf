import { notFound } from 'next/navigation';
import MeetsDeckClient from './MeetsDeckClient';
import { getGuestConfig, listGuests } from '@/components/meets/guests';

interface MeetsPageProps {
  params: Promise<{
    guest: string;
    date: string;
  }>;
}

/**
 * Generate static params for all guest presentations
 * Required for static export
 */
export function generateStaticParams() {
  const guests = listGuests();
  return guests.map((key) => {
    const [guest, date] = key.split('/');
    return { guest, date };
  });
}

/**
 * Dynamic route for guest meeting presentations
 *
 * Route: /meets/[guest]/[date]
 * Example: /meets/sep/2026-01-23-1400
 *
 * The date format is yyyy-mm-dd-hhmm where:
 * - yyyy: 4-digit year
 * - mm: 2-digit month
 * - dd: 2-digit day
 * - hhmm: hours and minutes in 24-hour format
 *
 * The guest and date params are used to look up the appropriate
 * slide configuration from the guests registry.
 *
 * Historical meetings can be accessed at their specific dated URLs,
 * while the latest meeting is always available at /meets/[guest]
 */
export default async function MeetsPage({ params }: MeetsPageProps) {
  const { guest, date } = await params;
  const config = getGuestConfig(guest, date);

  if (!config) {
    return notFound();
  }

  return (
    <MeetsDeckClient
      slides={config.slides}
      guest={config.guest}
      meetingDate={config.meetingDate}
      presenter={config.presenter}
      company={config.company}
      location={config.location}
    />
  );
}
