/**
 * THE single source of pricing truth for the website.
 *
 * Every number a visitor can see -- plan cards, the landing page, the value
 * matrix, FAQ copy -- must read from here. The pricing page and the landing
 * page once disagreed publicly ($19 vs $20 per seat, and a 10x drift on the
 * automation-minute rate) precisely because these constants lived inside one
 * UI component while other components re-authored them. The law applies to
 * prose too: a sentence that names a price references a constant.
 *
 * Prices are market-aware: each market carries its own SET prices -- an
 * India price is a real India number, never an FX conversion of the US one.
 * Adding a market is one entry in MARKETS; the display layer follows.
 */

export type MarketId = 'us' | 'in';

export interface EnterpriseTierDisplay {
  name: string;
  /** Market-native display string (e.g. "$60K", "₹30 lakh"). */
  perYear: string;
  seatCeiling: number;
  support: string;
}

export interface Market {
  id: MarketId;
  /** Currency code shown on the market control. */
  currency: string;
  /** Region name shown beside the control. */
  region: string;
  symbol: string;
  teamSeatMonthly: number;
  /** Per seat per YEAR (two months free vs monthly). */
  teamSeatAnnual: number;
  /** Smallest purchasable AI credit pack. */
  creditPackStart: number;
  enterprise: EnterpriseTierDisplay[];
}

export const MARKETS: Record<MarketId, Market> = {
  us: {
    id: 'us',
    currency: 'USD',
    region: 'United States',
    symbol: '$',
    teamSeatMonthly: 20,
    teamSeatAnnual: 192,
    creditPackStart: 10,
    enterprise: [
      {
        name: 'Enterprise Standard',
        perYear: '$60K',
        seatCeiling: 100,
        support: 'Named business-hours support',
      },
      {
        name: 'Enterprise Plus',
        perYear: '$120K',
        seatCeiling: 250,
        support: '24×7 support with SLA',
      },
    ],
  },
  in: {
    id: 'in',
    currency: 'INR',
    region: 'India',
    symbol: '₹',
    teamSeatMonthly: 999,
    teamSeatAnnual: 9990,
    creditPackStart: 499,
    enterprise: [
      {
        name: 'Enterprise Standard',
        perYear: '₹30 lakh',
        seatCeiling: 100,
        support: 'Named business-hours support',
      },
      {
        name: 'Enterprise Plus',
        perYear: '₹60 lakh',
        seatCeiling: 250,
        support: '24×7 support with SLA',
      },
    ],
  },
};

export const DEFAULT_MARKET_ID: MarketId = 'us';

/**
 * Best-effort market detection from a browser locale (e.g. "en-IN").
 * The site is a static export -- there is no server to geolocate a request,
 * so the browser locale picks the DEFAULT and an explicit control lets the
 * visitor choose (the same posture as the buy page's country prefill).
 */
export const detectMarket = (locale?: string): MarketId => {
  if (!locale) return DEFAULT_MARKET_ID;
  try {
    const region = new Intl.Locale(locale).region;
    return region === 'IN' ? 'in' : DEFAULT_MARKET_ID;
  } catch {
    return DEFAULT_MARKET_ID;
  }
};

// The free tier's ONE bound (founder-decided 2026-08-13): seats. Everything
// else -- environments, cloud account connections, resources, services,
// automation minutes -- is deliberately unlimited on every plan (the meters
// competitors cap; environments cost the customer's cloud, not ours). The
// free tier never asks for a card and never bills -- at its cap it pauses
// new invites; nothing running is touched. A future tightening for NEW
// organizations would be a catalog edit that grandfathers existing ones.
export const FREE_TIER_SEATS = 3;

/**
 * Self-hosted community edition seat cap. Decided 2026-08-17 (it briefly
 * shipped as "unlimited", which made the paid rungs unsellable — a guest
 * read it live and asked "why would I even pay you?"). The sixth member
 * is the structural walk into the licensed rungs.
 */
export const COMMUNITY_SEAT_LIMIT = 5;

// Self-hosted licenses: yearly, billed in USD in every market (a global
// SKU; market-native license prices become one data edit if ever decided).
// The community edition stays free forever -- full core, unlimited seats.
export const SELF_HOSTED_LICENSE_RUNGS = [
  { usdPerYear: 1999, seatCeiling: 10 },
  { usdPerYear: 4999, seatCeiling: 25 },
] as const;

// Free full-experience evaluation on your own cluster: every capability
// unlocked for this many days; expiry steps down gently, never bricks.
export const EVALUATION_DAYS = 30;

// Where a self-hosted license purchase starts (the console's public
// buy page; no account required).
export const BUY_LICENSE_URL = 'https://planton.ai/license/buy';

// The self-serve ceiling: below this seat count nobody needs to talk to
// sales -- the license rungs are card-and-email purchases.
export const SELF_SERVE_SEAT_CEILING = 25;

