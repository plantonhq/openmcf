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
 * The site is a static export -- no server sees the request, so the browser
 * locale picks the DEFAULT market (an India visitor lands on ₹
 * automatically) and an explicit control lets anyone view any market. This
 * mirrors the product's own buy-page posture: locale prefill, user control.
 *
 * Display truth vs charge truth: this is presentation only. The actual
 * charge is always fixed server-side from the catalog's regional prices at
 * checkout.
 */

// The browser locale never changes within a page's lifetime; reading it
// through useSyncExternalStore keeps hydration honest (server snapshot =
// the default market) without effect-driven state.
const emptySubscribe = () => () => {};
const clientSnapshot = () => detectMarket(navigator.language);
const serverSnapshot = () => DEFAULT_MARKET_ID;

/** Locale-detected market id + an explicit override the visitor controls. */
const useMarketChoice = (): [MarketId, (id: MarketId) => void] => {
  const detected = useSyncExternalStore(emptySubscribe, clientSnapshot, serverSnapshot);
  const [chosen, setChosen] = useState<MarketId | null>(null);
  return [chosen ?? detected, setChosen];
};

interface MarketContextValue {
  market: Market;
  marketId: MarketId;
  setMarketId: (id: MarketId) => void;
}

const MarketContext = createContext<MarketContextValue | null>(null);

export const MarketProvider: FC<{ children: ReactNode }> = ({ children }) => {
  const [marketId, setMarketId] = useMarketChoice();
  return (
    <MarketContext.Provider
      value={{ market: MARKETS[marketId], marketId, setMarketId }}
    >
      {children}
    </MarketContext.Provider>
  );
};

/**
 * The shared market fact. Outside a provider (e.g. the landing page's
 * pricing section, which carries no selector) it falls back to a local
 * locale-detected value.
 */
export const useMarket = (): MarketContextValue => {
  const ctx = useContext(MarketContext);
  const [fallbackId, setFallbackId] = useMarketChoice();
  if (ctx) return ctx;
  return {
    market: MARKETS[fallbackId],
    marketId: fallbackId,
    setMarketId: setFallbackId,
  };
};

/** Compact segmented control for choosing the displayed market. */
export const MarketSelector: FC<{ className?: string }> = ({ className = '' }) => {
  const { marketId, setMarketId } = useMarket();
  return (
    <Box className={`inline-flex items-center gap-2 ${className}`}>
      <Box className="inline-flex rounded-full border border-[#2a2a2a] bg-[#111] p-0.5">
        {Object.values(MARKETS).map((market) => (
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
