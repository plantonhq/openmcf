import { notFound } from 'next/navigation';
import MeetsDeckClient from './[date]/MeetsDeckClient';
import { listGuests, getLatestGuestConfig } from '@/components/meets/guests';

interface GuestPageProps {
  params: Promise<{
    guest: string;
  }>;
}

/**
 * Generate static params for all guest index pages
 * Required for static export
 */
export function generateStaticParams() {
  const guests = listGuests();
  const uniqueGuests = [...new Set(guests.map((key) => key.split('/')[0]))];
  return uniqueGuests.map((guest) => ({ guest }));
}

/**
 * Guest index page - renders the latest (upcoming) presentation directly
 *
 * Route: /meets/[guest]
 * Example: /meets/sep -> renders the latest SEP presentation
 *
 * This is the primary URL for sharing with guests. The most recent
 * meeting is always displayed at this URL, making it easy to remember
 * and type: planton.ai/meets/guest-name
 */
export default async function GuestPage({ params }: GuestPageProps) {
  const { guest } = await params;

  // Get the latest presentation for this guest
  const config = getLatestGuestConfig(guest);

  if (!config) {
    return notFound();
  }

  // Render the presentation directly
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
