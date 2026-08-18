'use client';

import { useState, useMemo, useCallback } from 'react';
import {
  Currency,
  VALUATION_CAP,
  EXIT_SCENARIOS,
  INVESTMENT_PRESETS_USD,
  INVESTMENT_PRESETS_INR,
  formatCurrency,
  formatPercent,
  formatMultiple,
} from '../shared';

/**
 * Calculator state
 */
export interface CalculatorState {
  /** Investment amount in USD */
  investmentAmount: number;
  /** Exit valuation in USD */
  exitValuation: number;
}

/**
 * Calculated results
 */
export interface CalculatorResults {
  /** Ownership percentage (decimal) */
  ownership: number;
  /** Formatted ownership percentage */
  ownershipFormatted: string;
  /** Return value in USD */
  returnValue: number;
  /** Formatted return value */
  returnValueFormatted: string;
  /** Return multiple */
  multiple: number;
  /** Formatted return multiple */
  multipleFormatted: string;
  /** Multiple color class for styling */
  multipleColorClass: string;
}

/**
 * Scenario calculation result
 */
export interface ScenarioResult {
  key: keyof typeof EXIT_SCENARIOS;
  label: string;
  exitValuation: number;
  exitValueFormatted: string;
  returnValue: number;
  returnValueFormatted: string;
  multiple: number;
  multipleFormatted: string;
  multipleColorClass: string;
}

interface UseCalculatorOptions {
  /** Current currency for formatting */
  currency: Currency;
  /** Initial investment amount in USD */
  initialInvestment?: number;
  /** Initial exit valuation in USD */
  initialExitValuation?: number;
}

interface UseCalculatorReturn {
  /** Current calculator state */
  state: CalculatorState;
  /** Calculated results */
  results: CalculatorResults;
  /** Set investment amount (in USD) */
  setInvestmentAmount: (amount: number) => void;
  /** Set exit valuation (in USD) */
  setExitValuation: (valuation: number) => void;
  /** Investment preset amounts based on currency */
  investmentPresets: readonly number[];
  /** Calculate results for all exit scenarios */
  scenarioResults: ScenarioResult[];
  /** Formatted investment amount */
  investmentFormatted: string;
  /** Formatted exit valuation */
  exitValuationFormatted: string;
}

/**
 * Get color class based on return multiple
 */
function getMultipleColorClass(multiple: number): string {
  if (multiple >= 10) return 'text-emerald-400';
  if (multiple >= 4) return 'text-pink-400';
  return 'text-amber-400';
}

/**
 * Hook for investment calculator logic
 * 
 * Calculates ownership percentage, return value, and multiples based on:
 * - Investment amount
 * - Exit valuation
 * - Valuation cap ($7M)
 */
export function useCalculator(options: UseCalculatorOptions): UseCalculatorReturn {
  const {
    currency,
    initialInvestment = 25_000,
    initialExitValuation = 100_000_000,
  } = options;
  
  const [state, setState] = useState<CalculatorState>({
    investmentAmount: initialInvestment,
    exitValuation: initialExitValuation,
  });
  
  // Calculate ownership percentage based on investment and cap
  const ownership = useMemo(() => {
    return state.investmentAmount / VALUATION_CAP;
  }, [state.investmentAmount]);
  
  // Calculate return value at exit
  const returnValue = useMemo(() => {
    return ownership * state.exitValuation;
  }, [ownership, state.exitValuation]);
  
  // Calculate return multiple
  const multiple = useMemo(() => {
    return returnValue / state.investmentAmount;
  }, [returnValue, state.investmentAmount]);
  
  // Compiled results with formatting
  const results = useMemo<CalculatorResults>(() => ({
    ownership,
    ownershipFormatted: formatPercent(ownership),
    returnValue,
    returnValueFormatted: formatCurrency(returnValue, currency),
    multiple,
    multipleFormatted: formatMultiple(multiple),
    multipleColorClass: getMultipleColorClass(multiple),
  }), [ownership, returnValue, multiple, currency]);
  
  // Calculate results for all exit scenarios
  const scenarioResults = useMemo<ScenarioResult[]>(() => {
    return Object.entries(EXIT_SCENARIOS).map(([key, scenario]) => {
      const scenarioOwnership = state.investmentAmount / VALUATION_CAP;
      const scenarioReturnValue = scenarioOwnership * scenario.value;
      const scenarioMultiple = scenarioReturnValue / state.investmentAmount;
      
      return {
        key: key as keyof typeof EXIT_SCENARIOS,
        label: scenario.label,
        exitValuation: scenario.value,
        exitValueFormatted: formatCurrency(scenario.value, currency),
        returnValue: scenarioReturnValue,
        returnValueFormatted: formatCurrency(scenarioReturnValue, currency),
        multiple: scenarioMultiple,
        multipleFormatted: formatMultiple(scenarioMultiple),
        multipleColorClass: getMultipleColorClass(scenarioMultiple),
      };
    });
  }, [state.investmentAmount, currency]);
  
  // Get investment presets based on currency
  const investmentPresets = useMemo(() => {
    return currency === 'INR' ? INVESTMENT_PRESETS_INR : INVESTMENT_PRESETS_USD;
  }, [currency]);
  
  const setInvestmentAmount = useCallback((amount: number) => {
    setState(prev => ({ ...prev, investmentAmount: amount }));
  }, []);
  
  const setExitValuation = useCallback((valuation: number) => {
    setState(prev => ({ ...prev, exitValuation: valuation }));
  }, []);
  
  return {
    state,
    results,
    setInvestmentAmount,
    setExitValuation,
    investmentPresets,
    scenarioResults,
    investmentFormatted: formatCurrency(state.investmentAmount, currency),
    exitValuationFormatted: formatCurrency(state.exitValuation, currency),
  };
}
