import { SlideConfig } from '../MeetsDeck';

// ============================================================================
// GUEST CONFIGURATION TYPES
// ============================================================================

export interface GuestConfig {
  slides: SlideConfig[];
  guest: string;
  meetingDate: string;
  presenter?: string;
  company?: string;
  location?: string;
}

// ============================================================================
// GUEST REGISTRY
// ============================================================================

// Import guest configurations
import { sepConfig } from './sep/config';
import { niravConfig } from './nirav/config';
import { clearRouteConfig } from './clear-route/config';
import { rahulGulatiConfig } from './rahul-gulati/config';

/**
 * Registry of all guest meeting configurations
 *
 * Key format: "guest/date" (e.g., "sep/2026-01-23-1400")
 *
 * The date format is yyyy-mm-dd-hhmm where:
 * - yyyy: 4-digit year
 * - mm: 2-digit month
 * - dd: 2-digit day
 * - hhmm: hours and minutes in 24-hour format
 *
 * Example: "sep/2026-01-23-1400" for January 23, 2026 at 2:00 PM
 */
const guestRegistry: Record<string, GuestConfig> = {
  'sep/2026-01-23-1400': sepConfig,
  'nirav/2026-05-08-2030': niravConfig,
  'clear-route/2026-08-12-1100': clearRouteConfig,
  'rahul-gulati/2026-08-17-1700': rahulGulatiConfig,
};

// ============================================================================
// LOOKUP FUNCTIONS
// ============================================================================

/**
 * Get guest configuration by guest ID and date
 *
 * @param guest - Guest identifier (e.g., "sep")
 * @param date - Meeting date (e.g., "2026-01-23-1400")
 * @returns GuestConfig if found, null otherwise
 */
export function getGuestConfig(
  guest: string,
  date: string
): GuestConfig | null {
  const key = `${guest}/${date}`;
  return guestRegistry[key] || null;
}

/**
 * Get the latest (upcoming) guest configuration for a given guest
 *
 * This function returns the most recent meeting by sorting dates in descending order.
 * The URL /meets/[guest] always shows the upcoming/latest meeting.
 *
 * @param guest - Guest identifier (e.g., "sep")
 * @returns GuestConfig if found, null otherwise
 */
export function getLatestGuestConfig(guest: string): GuestConfig | null {
  const guestDates = Object.keys(guestRegistry)
    .filter((key) => key.startsWith(`${guest}/`))
    .map((key) => key.split('/')[1])
    .sort()
    .reverse();

  if (guestDates.length === 0) {
    return null;
  }

  return guestRegistry[`${guest}/${guestDates[0]}`] || null;
}

/**
 * List all available guest presentations
 *
 * @returns Array of guest/date keys
 */
export function listGuests(): string[] {
  return Object.keys(guestRegistry);
}
