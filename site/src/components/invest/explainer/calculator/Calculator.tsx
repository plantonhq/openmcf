'use client';

import { FC, useId } from 'react';
import { motion } from 'framer-motion';
import {
  Currency,
  formatCurrency,
  formatUSD,
  formatARR,
  EXIT_SCENARIOS,
  VALUATION_CAP,
  fadeInUp,
  defaultTransition,
  Badge,
} from '../shared';
import { useCalculator, ScenarioResult } from '../hooks';

interface CalculatorProps {
  /** Current currency */
  currency: Currency;
  /** Currency change handler */
  onCurrencyChange?: (currency: Currency) => void;
  className?: string;
}

/**
 * Interactive investment calculator with inline currency toggle.
 * 
 * The currency toggle is contextually placed here rather than in the global header
 * because this is where currency actually matters - when doing calculations.
 */
export const Calculator: FC<CalculatorProps> = ({ 
  currency, 
  onCurrencyChange,
  className = '' 
}) => {
  const sliderId = useId();
  
  const {
    state,
    results,
    setInvestmentAmount,
    setExitValuation,
    investmentPresets,
    exitValuationFormatted,
  } = useCalculator({ currency });

  return (
    <motion.div
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true }}
      variants={fadeInUp}
      transition={defaultTransition}
      className={`bg-white/5 border border-white/10 rounded-2xl p-4 sm:p-6 md:p-8 ${className}`}
    >
      {/* Header with title and currency toggle */}
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg sm:text-xl font-bold text-white">
          Investment Calculator
        </h3>
        
        {/* Inline Currency Toggle */}
        {onCurrencyChange && (
          <CurrencyToggle currency={currency} onChange={onCurrencyChange} />
        )}
      </div>

      {/* Investment Amount Section */}
      <div className="mb-6">
        <label className="block text-sm font-medium text-white/70 mb-3">
          Your Investment Amount
        </label>
        <div className="flex flex-wrap gap-2">
          {investmentPresets.map((amount) => {
            const usdAmount = currency === 'INR' ? amount / 83 : amount;
            const isActive = Math.abs(state.investmentAmount - usdAmount) < 100;
            
            return (
              <button
                key={amount}
                type="button"
                onClick={() => setInvestmentAmount(usdAmount)}
                className={`
                  px-4 py-2 rounded-lg text-sm font-medium transition-all
                  ${isActive
                    ? 'bg-white/30 text-white/70 border border-white/50'
                    : 'bg-white/5 text-white/70 border border-white/10 hover:bg-white/10'
                  }
                `}
              >
                {formatCurrency(usdAmount, currency)}
              </button>
            );
          })}
        </div>
      </div>

      {/* Exit Valuation Section */}
      <div className="mb-6">
        <label htmlFor={sliderId} className="block text-sm font-medium text-white/70 mb-3">
          Exit Valuation Scenario
        </label>
        
        {/* Slider */}
        <div className="mb-4">
          <input
            type="range"
            id={sliderId}
            min={VALUATION_CAP}
            max={500_000_000}
            step={1_000_000}
            value={state.exitValuation}
            onChange={(e) => setExitValuation(Number(e.target.value))}
            className="w-full h-2 bg-white/10 rounded-lg appearance-none cursor-pointer accent-white"
            aria-valuemin={VALUATION_CAP}
            aria-valuemax={500_000_000}
            aria-valuenow={state.exitValuation}
            aria-valuetext={exitValuationFormatted}
          />
          <div className="flex justify-between text-xs text-white/40 mt-1">
            <span>{formatUSD(VALUATION_CAP)} (cap)</span>
            <span>{formatUSD(500_000_000)}</span>
          </div>
          <div className="text-center mt-2">
            <span className="text-lg font-bold text-white">{exitValuationFormatted}</span>
            <span className="text-xs text-white/40 ml-1.5">
              ({formatARR(state.exitValuation, currency)})
            </span>
          </div>
        </div>

        {/* Scenario buttons */}
        <div className="flex flex-wrap gap-2">
          {Object.entries(EXIT_SCENARIOS).map(([key, scenario]) => {
            const isActive = state.exitValuation === scenario.value;
            
            return (
              <button
                key={key}
                type="button"
                onClick={() => setExitValuation(scenario.value)}
                className={`
                  px-3 py-1.5 rounded-lg text-xs font-medium transition-all
                  ${isActive
                    ? 'bg-white/30 text-white/70 border border-white/50'
                    : 'bg-white/5 text-white/60 border border-white/10 hover:bg-white/10'
                  }
                `}
              >
                {formatCurrency(scenario.value, currency)}
              </button>
            );
          })}
        </div>
      </div>

      {/* Fixed Values */}
      <div className="mb-6">
        <label className="block text-sm font-medium text-white/70 mb-3">
          Fixed Values
        </label>
        <div className="flex flex-wrap gap-2">
          <Badge variant="purple">Cap: {formatCurrency(VALUATION_CAP, currency)}</Badge>
          <Badge variant="default">Type: Post-money SAFE</Badge>
        </div>
      </div>

      {/* Results */}
      <div className="grid grid-cols-3 gap-3 sm:gap-4 p-4 bg-white/5 rounded-xl">
        <div className="text-center">
          <div className="text-xl sm:text-2xl md:text-3xl font-bold text-white">
            {results.ownershipFormatted}
          </div>
          <div className="text-xs sm:text-sm text-white/50">Ownership</div>
        </div>
        <div className="text-center">
          <div className="text-xl sm:text-2xl md:text-3xl font-bold text-[#10b981]">
            {results.returnValueFormatted}
          </div>
          <div className="text-xs sm:text-sm text-white/50">Value at Exit</div>
        </div>
        <div className="text-center">
          <div className={`text-xl sm:text-2xl md:text-3xl font-bold ${results.multipleColorClass}`}>
            {results.multipleFormatted}
          </div>
          <div className="text-xs sm:text-sm text-white/50">Return</div>
        </div>
      </div>
    </motion.div>
  );
};

interface ScenarioTableProps {
  /** Current currency */
  currency: Currency;
  /** Investment amount for header */
  investmentAmount: number;
  /** Scenario results */
  scenarios: ScenarioResult[];
  className?: string;
}

export const ScenarioTable: FC<ScenarioTableProps> = ({
  currency,
  investmentAmount,
  scenarios,
  className = '',
}) => {
  return (
    <div className={`bg-white/5 border border-white/10 rounded-xl p-4 sm:p-6 ${className}`}>
      <h4 className="text-base sm:text-lg font-bold text-white mb-4">
        Pre-Calculated Scenarios ({formatCurrency(investmentAmount, currency)} Investment)
      </h4>
      
      <div className="overflow-x-auto -mx-4 px-4">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-white/10">
              <th className="text-left py-3 px-2 font-semibold text-white/80">Scenario</th>
              <th className="text-left py-3 px-2 font-semibold text-white/80">Exit Valuation</th>
              <th className="text-left py-3 px-2 font-semibold text-white/80">Your Return</th>
              <th className="text-left py-3 px-2 font-semibold text-white/80">Multiple</th>
            </tr>
          </thead>
          <tbody>
            {scenarios.map((scenario) => (
              <tr
                key={scenario.key}
                className={`border-b border-white/5 ${scenario.key === 'optimistic' ? 'bg-white/5' : ''}`}
              >
                <td className="py-3 px-2 text-white/60">{scenario.label}</td>
                <td className="py-3 px-2 text-white/60">
                  {scenario.exitValueFormatted}
                  <span className="text-[10px] text-white/30 ml-1">
                    ({formatARR(scenario.exitValuation, currency)})
                  </span>
                </td>
                <td className={`py-3 px-2 font-medium ${scenario.key === 'optimistic' ? 'text-white' : 'text-white/60'}`}>
                  {scenario.returnValueFormatted}
                </td>
                <td className={`py-3 px-2 font-medium ${scenario.multipleColorClass}`}>
                  {scenario.multipleFormatted}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      
      <p className="text-xs text-white/40 mt-4">
        Note: Actual returns depend on dilution from future funding rounds. These are illustrative scenarios, not guarantees.
      </p>
    </div>
  );
};

// ============================================================================
// CURRENCY TOGGLE
// ============================================================================

interface CurrencyToggleProps {
  currency: Currency;
  onChange: (currency: Currency) => void;
}

/**
 * Compact currency toggle for inline use within sections.
 * 
 * Appears contextually where currency matters (calculator, scenario table)
 * rather than as a global header control.
 */
const CurrencyToggle: FC<CurrencyToggleProps> = ({ currency, onChange }) => {
  return (
    <div
      className="flex rounded-lg border border-white/10 overflow-hidden bg-white/5"
      role="radiogroup"
      aria-label="Currency selection"
    >
      <button
        type="button"
        role="radio"
        aria-checked={currency === 'USD'}
        onClick={() => onChange('USD')}
        className={`
          px-2.5 py-1 text-xs font-medium transition-all duration-200
          ${currency === 'USD'
            ? 'bg-white/10 text-white/70'
            : 'text-white/40 hover:text-white/60'
          }
        `}
      >
        USD
      </button>
      <button
        type="button"
        role="radio"
        aria-checked={currency === 'INR'}
        onClick={() => onChange('INR')}
        className={`
          px-2.5 py-1 text-xs font-medium transition-all duration-200
          ${currency === 'INR'
            ? 'bg-white/10 text-white/70'
            : 'text-white/40 hover:text-white/60'
          }
        `}
      >
        INR
      </button>
    </div>
  );
};

export default Calculator;
