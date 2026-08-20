'use client';

import {
  createContext,
  FC,
  ReactNode,
  useContext,
  useState,
  useSyncExternalStore,
} from 'react';
import { Box, Typography } from '@mui/material';
import {
  DEFAULT_MARKET_ID,
  detectMarket,
  Market,
  MarketId,
  MARKETS,
} from '@/data/pricing';

/**
 * Market awareness for price display on a static site.
 *
 * The site is a static export -- no server sees the request, so detection
 * is client-side from the browser's two locality signals (locale region +
 * IANA timezone; see detectMarket). The detected market is a GATE, not
 * just a default (founder direction 2026-08-20):
 *
 * - Detected outside India: prices pin to USD and NO market control
 *   renders -- INR is unreachable.
 * - Detected in India: INR is the default, and the control offers ₹ INR
 *   (left, pre-selected) and $ USD (right).
 *
 * Display truth vs charge truth: this is presentation only. The actual
 * charge is always fixed server-side from the catalog's regional prices at
 * checkout.
 */

// The browser's locale and timezone never change within a page's lifetime;
// reading them through useSyncExternalStore keeps hydration honest (server
// snapshot = the default market) without effect-driven state.
const emptySubscribe = () => () => {};
const clientSnapshot = () =>
  detectMarket(navigator.language, Intl.DateTimeFormat().resolvedOptions().timeZone);
const serverSnapshot = () => DEFAULT_MARKET_ID;

/**
 * The gated market choice: outside India the market is pinned to USD (an
 * override cannot reach INR because the computed id ignores it); in India
 * the default is INR and the visitor may switch.
 */
const useMarketChoice = (): [MarketId, (id: MarketId) => void, boolean] => {
  const detected = useSyncExternalStore(emptySubscribe, clientSnapshot, serverSnapshot);
  const [chosen, setChosen] = useState<MarketId | null>(null);
  const switchable = detected === 'in';
  const marketId = switchable ? (chosen ?? 'in') : 'us';
  return [marketId, setChosen, switchable];
};

interface MarketContextValue {
  market: Market;
  marketId: MarketId;
  setMarketId: (id: MarketId) => void;
  /** True only for visitors detected in India — the market control's render condition. */
  switchable: boolean;
}

const MarketContext = createContext<MarketContextValue | null>(null);

export const MarketProvider: FC<{ children: ReactNode }> = ({ children }) => {
  const [marketId, setMarketId, switchable] = useMarketChoice();
  return (
    <MarketContext.Provider
      value={{ market: MARKETS[marketId], marketId, setMarketId, switchable }}
    >
      {children}
    </MarketContext.Provider>
  );
};

/**
 * The shared market fact. Outside a provider (e.g. the landing page's
 * pricing section, which carries no selector) it falls back to a local
 * detected value — the India gate applies identically on both paths.
 */
export const useMarket = (): MarketContextValue => {
  const ctx = useContext(MarketContext);
  const [fallbackId, setFallbackId, fallbackSwitchable] = useMarketChoice();
  if (ctx) return ctx;
  return {
    market: MARKETS[fallbackId],
    marketId: fallbackId,
    setMarketId: setFallbackId,
    switchable: fallbackSwitchable,
  };
};

// India-first order (founder direction 2026-08-20): the control renders only
// for visitors detected in India, so INR sits left and pre-selected, USD
// right as the comparison view.
const SELECTOR_ORDER: Market[] = [MARKETS.in, MARKETS.us];

/**
 * Compact segmented control for choosing the displayed market. Renders
 * NOTHING outside India — non-India visitors see USD only, with no path to
 * INR (the gate in useMarketChoice makes that structural, not cosmetic).
 */
export const MarketSelector: FC<{ className?: string }> = ({ className = '' }) => {
  const { marketId, setMarketId, switchable } = useMarket();
  if (!switchable) return null;
  return (
    <Box className={`inline-flex items-center gap-2 ${className}`}>
      <Box className="inline-flex rounded-full border border-[#2a2a2a] bg-[#111] p-0.5">
        {SELECTOR_ORDER.map((market) => (
          <button
            key={market.id}
            type="button"
            onClick={() => setMarketId(market.id)}
            aria-pressed={marketId === market.id}
            className={`px-3 py-1 rounded-full text-xs font-medium transition-all duration-200 cursor-pointer border-0 ${
              marketId === market.id
                ? 'bg-white text-black'
                : 'bg-transparent text-[#a0a0a0] hover:text-white'
            }`}
          >
            {market.symbol} {market.currency}
          </button>
        ))}
      </Box>
      <Typography className="text-xs text-[#666]">
        Prices shown for {MARKETS[marketId].region}
      </Typography>
    </Box>
  );
};
