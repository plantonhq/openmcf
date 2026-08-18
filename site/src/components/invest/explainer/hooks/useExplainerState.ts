'use client';

import { useCallback, useEffect, useState } from 'react';
import { Currency } from '../shared';

/**
 * State for the explainer pages, synced with URL hash for shareability
 */
export interface ExplainerState<TPath extends string> {
  /** Currently selected content path (knowledge level or background) */
  selectedPath: TPath | null;
  /** Selected currency for display */
  currency: Currency;
}

interface UseExplainerStateOptions<TPath extends string> {
  /** Valid path values for this page */
  validPaths: readonly TPath[];
  /** Default currency (defaults to USD) */
  defaultCurrency?: Currency;
}

interface UseExplainerStateReturn<TPath extends string> {
  /** Current state */
  state: ExplainerState<TPath>;
  /** Select a content path */
  selectPath: (path: TPath) => void;
  /** Clear the selected path (reset selection) */
  clearPath: () => void;
  /** Set the currency */
  setCurrency: (currency: Currency) => void;
  /** Check if a path is selected */
  isPathSelected: (path: TPath) => boolean;
  /** Check if any path is selected */
  hasPathSelected: boolean;
}

/**
 * Get hash value from window.location (browser only)
 */
function getHash(): string {
  if (typeof window === 'undefined') return '';
  return window.location.hash.replace('#', '');
}

/**
 * Hook for managing explainer page state with URL sync
 * 
 * State is persisted to URL for shareability:
 * - #beginner - selected content path (clean hash-based)
 * - Currency stored in localStorage (doesn't need to be in URL)
 */
export function useExplainerState<TPath extends string>(
  options: UseExplainerStateOptions<TPath>
): UseExplainerStateReturn<TPath> {
  const { validPaths, defaultCurrency = 'USD' } = options;
  
  // Initialize state - will be synced with hash on mount
  const [state, setState] = useState<ExplainerState<TPath>>({
    selectedPath: null,
    currency: defaultCurrency,
  });
  
  // Sync state from URL hash + localStorage on mount / when deps change.
  useEffect(() => {
    const syncFromExternal = () => {
      const hash = getHash();
      const savedCurrency = localStorage.getItem('investor-currency');

      const newPath = hash && validPaths.includes(hash as TPath)
        ? (hash as TPath)
        : null;
      const newCurrency: Currency = savedCurrency === 'INR' ? 'INR' : defaultCurrency;

      setState({ selectedPath: newPath, currency: newCurrency });
    };

    const raf = requestAnimationFrame(syncFromExternal);
    return () => cancelAnimationFrame(raf);
  }, [validPaths, defaultCurrency]);
  
  // Listen for hash changes (browser back/forward)
  useEffect(() => {
    const handleHashChange = () => {
      const hash = getHash();
      const newPath = hash && validPaths.includes(hash as TPath)
        ? (hash as TPath)
        : null;
      
      setState(prev => {
        if (prev.selectedPath !== newPath) {
          return { ...prev, selectedPath: newPath };
        }
        return prev;
      });
    };
    
    window.addEventListener('hashchange', handleHashChange);
    return () => window.removeEventListener('hashchange', handleHashChange);
  }, [validPaths]);
  
  const selectPath = useCallback((path: TPath) => {
    setState(prev => ({ ...prev, selectedPath: path }));
    // Update hash without triggering navigation
    if (typeof window !== 'undefined') {
      window.history.replaceState(null, '', `#${path}`);
    }
  }, []);
  
  const clearPath = useCallback(() => {
    setState(prev => ({ ...prev, selectedPath: null }));
    // Clear hash
    if (typeof window !== 'undefined') {
      window.history.replaceState(null, '', window.location.pathname);
    }
  }, []);
  
  const setCurrency = useCallback((currency: Currency) => {
    setState(prev => ({ ...prev, currency }));
    // Persist currency preference to localStorage
    if (typeof window !== 'undefined') {
      if (currency === 'INR') {
        localStorage.setItem('investor-currency', 'INR');
      } else {
        localStorage.removeItem('investor-currency');
      }
    }
  }, []);
  
  const isPathSelected = useCallback((path: TPath) => {
    return state.selectedPath === path;
  }, [state.selectedPath]);
  
  return {
    state,
    selectPath,
    clearPath,
    setCurrency,
    isPathSelected,
    hasPathSelected: state.selectedPath !== null,
  };
}
