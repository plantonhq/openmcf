'use client';

import { useEffect, useCallback } from 'react';

// Global state to survive component remounting (React Strict Mode)
// Using actual global variables, not useRef, to persist across remounts
let globalAutoClickTimer: NodeJS.Timeout | null = null;
let globalAutoClickExecuted = false;

interface AutoModalAndFillOptions {
  /** Whether to enable auto modal open functionality */
  enabled?: boolean;
  /** Delay before auto-clicking the first item (in milliseconds) */
  autoClickDelay?: number;
  /** Function to call when auto-clicking (should open the modal) */
  onAutoClick: () => void;
  /** Debug prefix for console logs */
  debugPrefix?: string;
}

interface AutoModalAndFillReturn {
  /** Manual trigger function to test the auto-click */
  triggerAutoClick: () => void;
}

/**
 * Reusable hook for auto modal open functionality
 * This is a simplified version that only handles the auto-click part
 * to ensure we don't break the existing working auto-fill logic
 *
 * @param options Configuration options for the hook
 * @returns Object with manual trigger function for testing
 *
 * @example
 * ```tsx
 * const { triggerAutoClick } = useAutoModalAndFill({
 *   enabled: AUTO_CLICK_FIRST_CARD,
 *   autoClickDelay: 1500,
 *   onAutoClick: () => handleCardClick(providers[0]),
 *   debugPrefix: 'Cloud Connections'
 * });
 * ```
 */
export const useAutoModalAndFill = (options: AutoModalAndFillOptions): AutoModalAndFillReturn => {
  const {
    enabled = true,
    autoClickDelay = 1500,
    onAutoClick,
    debugPrefix = 'Auto Modal & Fill',
  } = options;

  // Auto-click functionality
  useEffect(() => {
    const isIOS =
      /iPad|iPhone|iPod/.test(navigator.userAgent) ||
      /iPad|iPhone|iPod/.test(navigator.platform) ||
      (navigator.maxTouchPoints &&
        navigator.maxTouchPoints > 2 &&
        /MacIntel/.test(navigator.platform));

    if (enabled && !globalAutoClickExecuted && !globalAutoClickTimer) {
      // iOS-specific timing adjustments
      const adjustedDelay = isIOS ? autoClickDelay + 1000 : autoClickDelay; // Extra 1 second for iOS
      globalAutoClickExecuted = true; // Mark as triggered to prevent multiple timers

      globalAutoClickTimer = setTimeout(() => {
        try {
          onAutoClick();
        } catch (error) {
          console.error('Error calling onAutoClick:', error);
        }
        globalAutoClickTimer = null; // Clear the global timer
      }, adjustedDelay);
    }
  }, [enabled, autoClickDelay, onAutoClick, debugPrefix]);

  // Reset global state when component unmounts
  useEffect(() => {
    return () => {
      globalAutoClickExecuted = false;
      if (globalAutoClickTimer) {
        clearTimeout(globalAutoClickTimer);
        globalAutoClickTimer = null;
      }
    };
  }, [debugPrefix]);

  // Manual trigger function for testing
  const triggerAutoClick = useCallback(() => {
    try {
      onAutoClick();
    } catch (error) {
      console.error('Error in manual onAutoClick:', error);
    }
  }, [onAutoClick]);

  return {
    triggerAutoClick,
  };
};
