'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { Target, Users, Bot, Rocket, HelpCircle, X, ArrowRight } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { Slide, SlideTitle, SlideSubtitle, FundsItem, Grid, Callout } from '../shared';

const useOfFunds = [
  {
    icon: <Users className="w-5 h-5" />,
    title: 'Engineering',
    percentage: '60%',
    description: 'Hire 2-3 Engineers',
  },
  {
    icon: <Bot className="w-5 h-5" />,
    title: 'AI R&D',
    percentage: '25%',
    description: 'Agent Fleet MVP',
  },
  {
    icon: <Rocket className="w-5 h-5" />,
    title: 'GTM',
    percentage: '15%',
    description: 'Developer Advocacy',
  },
];

const milestones = [
  '50 Enterprise Clients',
  '$100K MRR',
  'Ready for Series A',
  'Planton DevOps AI Agents Used by Customers in Production',
];

// SAFE Explainer Modal Component
function SafeExplainerModal({ isOpen, onClose }: { isOpen: boolean; onClose: () => void }) {
  // Handle Escape key
  const handleKeyDown = useCallback((e: KeyboardEvent) => {
    if (e.key === 'Escape') {
      onClose();
    }
  }, [onClose]);

  useEffect(() => {
    if (isOpen) {
      document.addEventListener('keydown', handleKeyDown);
      // Prevent body scroll when modal is open
      document.body.style.overflow = 'hidden';
    }
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      document.body.style.overflow = '';
    };
  }, [isOpen, handleKeyDown]);

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50"
          />
          
          {/* Modal */}
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 20 }}
            transition={{ duration: 0.2 }}
            className="fixed inset-0 z-50 flex items-center justify-center p-4"
          >
            <div className="bg-[#111] border border-[#2a2a2a] rounded-xl p-3 sm:p-5 w-full max-w-3xl max-h-[85vh] overflow-y-auto relative">
              {/* Close button */}
              <button
                onClick={onClose}
                className="absolute top-2 right-2 sm:top-3 sm:right-3 p-1.5 text-[#666] hover:text-white hover:bg-[#1a1a1a] rounded-lg transition-colors duration-300"
                aria-label="Close modal"
              >
                <X className="w-4 h-4 sm:w-5 sm:h-5" />
              </button>

              <h4 className="text-sm sm:text-xl font-semibold text-white mb-2 sm:mb-3 text-center pr-6 sm:pr-8">
                SAFE = Simple Agreement for Future Equity
              </h4>
              
              {/* Two-column benefits - side by side even on mobile */}
              <div className="grid grid-cols-2 gap-2 sm:gap-3 mb-2 sm:mb-3">
                <div className="bg-[#10b981]/10 border border-[#10b981]/20 rounded-lg p-2 sm:p-3">
                  <h5 className="text-[#10b981] font-semibold text-[10px] sm:text-sm mb-1 sm:mb-1.5">For Planton</h5>
                  <ul className="space-y-0.5 sm:space-y-1.5 text-[10px] sm:text-xs text-[#a0a0a0]">
                    <li className="flex items-start gap-1">
                      <span className="text-[#10b981] shrink-0">✓</span>
                      <span>$500K to accelerate</span>
                    </li>
                    <li className="flex items-start gap-1">
                      <span className="text-[#10b981] shrink-0">✓</span>
                      <span>No valuation delays</span>
                    </li>
                    <li className="flex items-start gap-1">
                      <span className="text-[#10b981] shrink-0">✓</span>
                      <span>Focus on building</span>
                    </li>
                  </ul>
                </div>
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-2 sm:p-3">
                  <h5 className="text-white font-semibold text-[10px] sm:text-sm mb-1 sm:mb-1.5">For Investor</h5>
                  <ul className="space-y-0.5 sm:space-y-1.5 text-[10px] sm:text-xs text-[#a0a0a0]">
                    <li className="flex items-start gap-1">
                      <span className="text-white shrink-0">✓</span>
                      <span>Early discount</span>
                    </li>
                    <li className="flex items-start gap-1">
                      <span className="text-white shrink-0">✓</span>
                      <span>Equity at Series A</span>
                    </li>
                    <li className="flex items-start gap-1">
                      <span className="text-white shrink-0">✓</span>
                      <span>$7M cap protection</span>
                    </li>
                  </ul>
                </div>
              </div>

              {/* Conversion example - horizontal flow on mobile too */}
              <div className="bg-[#1a1a1a] rounded-lg p-2 sm:p-3">
                <h5 className="text-white font-semibold text-[10px] sm:text-sm mb-1.5 sm:mb-2 text-center">How Your $500K Converts</h5>
                <div className="flex items-center justify-center gap-1 sm:gap-2.5 text-[10px] sm:text-xs">
                  <div className="text-center p-1.5 sm:p-2.5 bg-[#151515] rounded-lg flex-1 sm:flex-none sm:w-auto">
                    <div className="text-[#666] text-[8px] sm:text-xs">Invest</div>
                    <div className="text-sm sm:text-xl font-bold text-white">$500K</div>
                  </div>
                  <ArrowRight className="w-3 h-3 sm:w-4 sm:h-4 text-[#3a3a3a] shrink-0" />
                  <div className="text-center p-1.5 sm:p-2.5 bg-[#151515] rounded-lg flex-1 sm:flex-none sm:w-auto">
                    <div className="text-[#666] text-[8px] sm:text-xs">@$20M</div>
                    <div className="text-sm sm:text-xl font-bold text-white">~7%</div>
                  </div>
                  <ArrowRight className="w-3 h-3 sm:w-4 sm:h-4 text-[#3a3a3a] shrink-0" />
                  <div className="text-center p-1.5 sm:p-2.5 bg-[#151515] rounded-lg flex-1 sm:flex-none sm:w-auto">
                    <div className="text-[#666] text-[8px] sm:text-xs">@$100M</div>
                    <div className="text-sm sm:text-xl font-bold text-[#10b981]">14x</div>
                  </div>
                </div>
                <p className="text-[8px] sm:text-xs text-[#666] text-center mt-1.5 sm:mt-2">$7M Cap • Terms Negotiable</p>
              </div>

              {/* $1000 Investment Example */}
              <div className="mt-5 sm:mt-8">
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-xl p-2 sm:p-3.5">
                  <h5 className="text-[#a0a0a0] font-semibold text-[10px] sm:text-sm mb-1.5 sm:mb-2.5 text-center">
                    What Could $1,000 Become?
                  </h5>
                  <div className="grid grid-cols-2 gap-2 sm:gap-3">
                    {/* Scenario 1: Exit at Cap */}
                    <div className="bg-[#151515] rounded-lg p-1.5 sm:p-3 text-center">
                      <div className="text-[#666] text-[8px] sm:text-xs mb-0.5 sm:mb-1">Exit at $7M (Cap)</div>
                      <div className="text-base sm:text-2xl font-bold text-white">$1,000</div>
                      <div className="text-[8px] sm:text-xs text-[#666]">1x (break even)</div>
                    </div>
                    {/* Scenario 2: Exit at 4x Cap */}
                    <div className="bg-[#10b981]/10 rounded-lg p-1.5 sm:p-3 text-center border border-[#10b981]/20">
                      <div className="text-[#666] text-[8px] sm:text-xs mb-0.5 sm:mb-1">Exit at $28M (4x Cap)</div>
                      <div className="text-base sm:text-2xl font-bold text-[#10b981]">$4,000</div>
                      <div className="text-[8px] sm:text-xs text-[#10b981]/70">4x (+$3K profit)</div>
                    </div>
                  </div>
                  {/* Cap Protection Explanation */}
                  <div className="mt-1.5 sm:mt-3 bg-[#151515] rounded-lg p-1.5 sm:p-3">
                    <p className="text-[8px] sm:text-xs text-[#666] text-center mb-1 sm:mb-2">
                      <span className="text-[#a0a0a0] font-medium">How the $7M cap protects you</span> — Imagine 1M shares, Series A at $20M:
                    </p>
                    <div className="grid grid-cols-2 gap-1.5 sm:gap-2.5 text-[10px] sm:text-xs">
                      <div className="text-center p-1 sm:p-2 bg-[#1a1a1a] rounded-lg">
                        <div className="text-[#666] text-[8px] sm:text-xs">Without cap</div>
                        <div className="text-sm sm:text-xl font-bold text-[#a0a0a0] my-0.5">50 shares</div>
                        <div className="text-[8px] sm:text-[10px] text-[#666]">$20/share</div>
                      </div>
                      <div className="text-center p-1 sm:p-2 bg-[#10b981]/10 rounded-lg border border-[#10b981]/20">
                        <div className="text-[#666] text-[8px] sm:text-xs">With $7M cap</div>
                        <div className="text-sm sm:text-xl font-bold text-[#10b981] my-0.5">143 shares</div>
                        <div className="text-[8px] sm:text-[10px] text-[#10b981]/60">$7/share</div>
                      </div>
                    </div>
                    <p className="text-[8px] sm:text-[10px] text-[#666] text-center mt-1 sm:mt-2">
                      Cap = <span className="text-[#10b981] font-medium">~3x more shares</span> for the same $1,000
                    </p>
                  </div>
                  {/* Market context reference */}
                  <p className="text-[8px] sm:text-[10px] text-[#666] text-center mt-1.5 sm:mt-3 italic">
                    DevTools companies typically raise Series A at $15-25M valuations. Our $7M cap gives you seed-stage pricing.
                  </p>
                </div>
              </div>

              {/* Escape hint - desktop only */}
              <p className="hidden sm:block text-xs text-[#666] text-center mt-3">
                Press <kbd className="px-1.5 py-0.5 bg-[#1a1a1a] border border-[#2a2a2a] rounded text-[#a0a0a0] text-[10px]">Esc</kbd> to close
              </p>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}

export default function SlideAsk() {
  const [showSafeExplainer, setShowSafeExplainer] = useState(false);

  return (
    <Slide className="!justify-start !pt-24 sm:!pt-28 md:!pt-32">
      <SlideTitle>The Ask</SlideTitle>
      <SlideSubtitle className="mb-3 sm:mb-6 text-xs sm:text-sm">
        Seed Round to Reach Next Milestones
      </SlideSubtitle>

      {/* The Number - more compact on mobile */}
      <div className="bg-[#1a1a1a] border border-[#3a3a3a] rounded-xl p-3 sm:p-6 md:p-8 max-w-[200px] sm:max-w-xs md:max-w-sm mx-auto mb-2 sm:mb-4">
        <div className="text-3xl sm:text-5xl md:text-6xl font-bold tracking-tight text-white mb-0.5 sm:mb-1">
          $500K
        </div>
        <div className="text-xs sm:text-sm md:text-base text-[#a0a0a0]">SAFE Note</div>
        <div className="text-[10px] sm:text-xs md:text-sm text-[#666] mt-0.5 sm:mt-1">~18 months runway</div>
      </div>

      {/* SAFE Explainer Toggle - styled as a pill button */}
      <button
        onClick={() => setShowSafeExplainer(true)}
        className="inline-flex items-center gap-1.5 sm:gap-2 px-3 sm:px-4 py-1.5 sm:py-2 
          bg-[#2a2a2a] 
          border border-[#3a3a3a] hover:border-white
          rounded-full text-[10px] sm:text-sm font-medium
          text-[#a0a0a0] hover:text-white
          hover:bg-white/5
          transition-all duration-300 mb-4 sm:mb-6 md:mb-8"
      >
        <HelpCircle className="w-3 h-3 sm:w-4 sm:h-4" />
        How does a SAFE work?
      </button>

      {/* SAFE Explainer Modal */}
      <SafeExplainerModal 
        isOpen={showSafeExplainer} 
        onClose={() => setShowSafeExplainer(false)} 
      />

      {/* Use of Funds - horizontal on mobile */}
      <Grid cols={3} gap="sm" className="mb-4 sm:mb-6 md:mb-8 sm:gap-5 md:gap-6">
        {useOfFunds.map((item) => (
          <FundsItem
            key={item.title}
            icon={item.icon}
            title={item.title}
            percentage={item.percentage}
            description={item.description}
          />
        ))}
      </Grid>

      {/* 18-Month Milestones - 2x2 grid on mobile */}
      <Callout className="max-w-xl sm:max-w-2xl p-2 sm:p-5 md:p-6">
        <div className="flex items-center justify-center gap-1 sm:gap-2 mb-1 sm:mb-3">
          <Target className="w-3 h-3 sm:w-5 sm:h-5 text-[#10b981]/70" />
          <h3 className="text-[10px] sm:text-base md:text-lg font-semibold text-white">18-Month Milestones</h3>
        </div>
        <div className="grid grid-cols-2 gap-x-3 sm:gap-x-8 gap-y-0.5 sm:gap-y-2 text-left">
          {milestones.map((milestone, index) => (
            <div key={index} className="flex items-start gap-1 sm:gap-2 text-[10px] sm:text-sm text-[#a0a0a0]">
              <span className="text-[#10b981] shrink-0 mt-0.5">→</span>
              <span className="leading-tight">{milestone}</span>
            </div>
          ))}
        </div>
      </Callout>
    </Slide>
  );
}

