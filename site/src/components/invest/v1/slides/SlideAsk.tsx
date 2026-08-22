'use client';

import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Target, Users, Rocket, Bot, HelpCircle, X, ArrowRight } from 'lucide-react';

const useOfFunds = [
  {
    icon: Users,
    title: 'Engineering',
    percentage: '60%',
    description: 'Hire 2-3 engineers to accelerate platform development',
  },
  {
    icon: Bot,
    title: 'AI R&D',
    percentage: '25%',
    description: 'Build Agent Fleet MVP—specialized DevOps AI assistants',
  },
  {
    icon: Rocket,
    title: 'GTM',
    percentage: '15%',
    description: 'Developer advocacy and initial customer acquisition',
  },
];

export default function SlideAsk() {
  const [showSafeExplainer, setShowSafeExplainer] = useState(false);

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-[#0a0a0a] via-[#1a1a1a] to-[#0a0a0a] p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center max-w-5xl w-full"
      >
        {/* Title */}
        <h2 className="text-4xl md:text-5xl font-bold text-white mb-4">
          The Ask
        </h2>
        <p className="text-xl text-white/60 mb-10">
          Seed round to reach next milestones
        </p>

        {/* The Number */}
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.2 }}
          className="bg-gradient-to-br from-white/20 to-[#666]/20 border border-white/20 rounded-3xl p-10 max-w-md mx-auto mb-4"
        >
          <div className="text-6xl md:text-7xl font-bold text-white mb-2">
            $500K
          </div>
          <div className="text-white/60">SAFE Note</div>
          <div className="text-sm text-white/40 mt-2">~18 months runway</div>
        </motion.div>

        {/* SAFE Explainer Toggle Button */}
        <button
          onClick={() => setShowSafeExplainer(!showSafeExplainer)}
          className="inline-flex items-center gap-2 text-sm text-white/50 hover:text-white/80 transition-colors mb-8"
        >
          {showSafeExplainer ? (
            <>
              <X className="w-4 h-4" />
              Hide SAFE Explanation
            </>
          ) : (
            <>
              <HelpCircle className="w-4 h-4" />
              How does a SAFE work?
            </>
          )}
        </button>

        {/* SAFE Explainer Panel */}
        <AnimatePresence>
          {showSafeExplainer && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
              className="overflow-hidden mb-8"
            >
              <div className="bg-white/5 border border-white/20 rounded-2xl p-6 max-w-3xl mx-auto text-left">
                <h4 className="text-lg font-semibold text-white mb-4 text-center">
                  SAFE = Simple Agreement for Future Equity
                </h4>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                  {/* For Planton */}
                  <div className="bg-emerald-500/10 border border-emerald-500/30 rounded-xl p-4">
                    <h5 className="text-emerald-400 font-semibold mb-3">For Planton</h5>
                    <ul className="space-y-2 text-sm text-white/70">
                      <li className="flex items-start gap-2">
                        <span className="text-emerald-400 mt-0.5">✓</span>
                        Receive $500K now to build
                      </li>
                      <li className="flex items-start gap-2">
                        <span className="text-emerald-400 mt-0.5">✓</span>
                        No valuation negotiation today
                      </li>
                      <li className="flex items-start gap-2">
                        <span className="text-emerald-400 mt-0.5">✓</span>
                        Focus on building, not fundraising
                      </li>
                    </ul>
                  </div>

                  {/* For Investor */}
                  <div className="bg-white/10 border border-white/20 rounded-xl p-4">
                    <h5 className="text-white font-semibold mb-3">For Investor</h5>
                    <ul className="space-y-2 text-sm text-white/70">
                      <li className="flex items-start gap-2">
                        <span className="text-white mt-0.5">✓</span>
                        Invest early at a discount
                      </li>
                      <li className="flex items-start gap-2">
                        <span className="text-white mt-0.5">✓</span>
                        Converts to equity at Series A
                      </li>
                      <li className="flex items-start gap-2">
                        <span className="text-white mt-0.5">✓</span>
                        Protected by valuation cap
                      </li>
                    </ul>
                  </div>
                </div>

                {/* Simple Math Example */}
                <div className="bg-white/5 rounded-xl p-4">
                  <h5 className="text-white font-semibold mb-3 text-center">Example: How Your $500K Converts</h5>
                  <div className="flex flex-col md:flex-row items-center justify-center gap-4 text-sm">
                    <div className="text-center p-3 bg-white/5 rounded-lg">
                      <div className="text-white/50 text-xs mb-1">You invest</div>
                      <div className="text-xl font-bold text-white">$500K</div>
                    </div>
                    <ArrowRight className="w-5 h-5 text-white/30 rotate-90 md:rotate-0" />
                    <div className="text-center p-3 bg-white/5 rounded-lg">
                      <div className="text-white/50 text-xs mb-1">If Series A at $20M</div>
                      <div className="text-xl font-bold text-white">~7% equity</div>
                      <div className="text-xs text-white/40">($500K ÷ $7M cap = 7.1%)</div>
                    </div>
                    <ArrowRight className="w-5 h-5 text-white/30 rotate-90 md:rotate-0" />
                    <div className="text-center p-3 bg-white/5 rounded-lg">
                      <div className="text-white/50 text-xs mb-1">If we hit $100M valuation</div>
                      <div className="text-xl font-bold text-emerald-400">$7M value</div>
                      <div className="text-xs text-white/40">14x return</div>
                    </div>
                  </div>
                  <p className="text-xs text-white/40 text-center mt-4">
                    * Assumes $7M valuation cap. Actual terms negotiable.
                  </p>
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>

        {/* Use of Funds */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          {useOfFunds.map((item, index) => (
            <motion.div
              key={item.title}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 + index * 0.1 }}
              className="bg-white/5 border border-white/10 rounded-2xl p-6 text-left"
            >
              <div className="flex items-center justify-between mb-4">
                <item.icon className="w-6 h-6 text-white/60" />
                <span className="text-2xl font-bold text-white">{item.percentage}</span>
              </div>
              <h3 className="text-lg font-semibold text-white mb-2">{item.title}</h3>
              <p className="text-sm text-white/50">{item.description}</p>
            </motion.div>
          ))}
        </div>

        {/* Milestones */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="bg-white/5 border border-white/10 rounded-xl p-6 max-w-2xl mx-auto"
        >
          <div className="flex items-center justify-center gap-2 mb-4">
            <Target className="w-5 h-5 text-emerald-400" />
            <h3 className="text-lg font-semibold text-white">18-Month Milestones</h3>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-left">
            <div className="flex items-center gap-2 text-white/60 text-sm">
              <span className="text-emerald-400">→</span>
              5-10 enterprise pilots
            </div>
            <div className="flex items-center gap-2 text-white/60 text-sm">
              <span className="text-emerald-400">→</span>
              Agent Fleet MVP launched
            </div>
            <div className="flex items-center gap-2 text-white/60 text-sm">
              <span className="text-emerald-400">→</span>
              $50K+ MRR target
            </div>
            <div className="flex items-center gap-2 text-white/60 text-sm">
              <span className="text-emerald-400">→</span>
              Ready for Series A
            </div>
          </div>
        </motion.div>
      </motion.div>
    </div>
  );
}
