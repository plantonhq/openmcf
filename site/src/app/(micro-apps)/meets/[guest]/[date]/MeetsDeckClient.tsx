'use client';

import MeetsDeck, { MeetsDeckProps } from '@/components/meets/MeetsDeck';

/**
 * Client component wrapper for MeetsDeck
 * This allows the page to be a server component with generateStaticParams
 * while the actual deck functionality runs on the client
 */
export default function MeetsDeckClient(props: MeetsDeckProps) {
  return <MeetsDeck {...props} />;
}
