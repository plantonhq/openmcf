'use client';

import React from 'react';
import { motion } from 'framer-motion';

export default function SlideMarket() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-[#0a0a0a] via-[#1a1a1a] to-[#0a0a0a] p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center max-w-5xl w-full"
      >
        {/* Title */}
        <h2 className="text-4xl md:text-5xl font-bold text-white mb-4">
          Market Opportunity
        </h2>
        <p className="text-xl text-white/60 mb-12">
          Cloud spend is massive. Complexity is the problem.
        </p>

        {/* TAM/SAM */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="bg-white/5 border border-white/10 rounded-2xl p-8"
          >
            <div className="text-5xl md:text-6xl font-bold text-white mb-2">$150B+</div>
            <div className="text-white/60">Global Cloud Spend</div>
            <div className="text-xs text-white/40 mt-2">And growing 20%+ YoY</div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="bg-gradient-to-br from-white/10 to-[#666]/10 border border-white/20 rounded-2xl p-8"
          >
            <div className="text-5xl md:text-6xl font-bold text-white mb-2">$20B+</div>
            <div className="text-white/60">DevOps Tools Market</div>
            <div className="text-xs text-white/40 mt-2">Infrastructure automation</div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.4 }}
            className="bg-white/5 border border-white/10 rounded-2xl p-8"
          >
            <div className="text-5xl md:text-6xl font-bold text-white mb-2">~50K</div>
            <div className="text-white/60">Mid-size Companies</div>
            <div className="text-xs text-white/40 mt-2">Without dedicated DevOps</div>
          </motion.div>
        </div>

        {/* Why Now */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="grid grid-cols-1 md:grid-cols-3 gap-4 max-w-3xl mx-auto"
        >
          <div className="bg-white/5 border border-white/10 rounded-xl p-4">
            <h4 className="text-white font-semibold mb-1">AI Wave</h4>
            <p className="text-white/50 text-sm">
              LLMs enable new automation possibilities
            </p>
          </div>
          <div className="bg-white/5 border border-white/10 rounded-xl p-4">
            <h4 className="text-white font-semibold mb-1">Cloud Complexity</h4>
            <p className="text-white/50 text-sm">
              AWS has 200+ services—too many choices
            </p>
          </div>
          <div className="bg-white/5 border border-white/10 rounded-xl p-4">
            <h4 className="text-white font-semibold mb-1">Talent Shortage</h4>
            <p className="text-white/50 text-sm">
              Not enough DevOps engineers to go around
            </p>
          </div>
        </motion.div>
      </motion.div>
    </div>
  );
}
