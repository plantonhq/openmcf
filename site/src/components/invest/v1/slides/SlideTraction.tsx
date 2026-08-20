'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { TrendingUp } from 'lucide-react';

const metrics = [
  { value: '3+', label: 'Years Building', highlight: false },
  { value: '$500K+', label: 'Self-Funded', highlight: false },
  { value: '$800', label: 'MRR', highlight: true },
  { value: '5', label: 'Team Members', highlight: false },
];

const milestones = [
  'Production platform with real paying customers',
  'Full Infra Hub + Service Hub live',
  'Multi-cloud support (AWS, GCP, Azure, K8s)',
  'Planton open source on GitHub',
];

export default function SlideTraction() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-[#0a0a0a] via-[#1a1a1a] to-[#0a0a0a] p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center max-w-5xl w-full"
      >
        {/* Title */}
        <h2 className="text-4xl md:text-5xl font-bold text-white mb-4">
          Traction
        </h2>
        <p className="text-xl text-white/60 mb-12">
          Pre-revenue with production platform
        </p>

        {/* Metrics */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-6 mb-12">
          {metrics.map((metric, index) => (
            <motion.div
              key={metric.label}
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: 0.1 * (index + 1) }}
              className={`${
                metric.highlight
                  ? 'bg-gradient-to-br from-white/20 to-[#666]/20 border-white/20'
                  : 'bg-white/5 border-white/10'
              } border rounded-2xl p-6`}
            >
              <div className={`text-4xl md:text-5xl font-bold mb-2 ${
                metric.highlight ? 'text-white' : 'text-white'
              }`}>
                {metric.value}
              </div>
              <div className="text-sm text-white/60">{metric.label}</div>
            </motion.div>
          ))}
        </div>

        {/* Milestones */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="bg-white/5 border border-white/10 rounded-2xl p-6 max-w-2xl mx-auto"
        >
          <div className="flex items-center justify-center gap-2 mb-4">
            <TrendingUp className="w-5 h-5 text-emerald-400" />
            <h3 className="text-lg font-semibold text-white">What We&apos;ve Built</h3>
          </div>
          <ul className="space-y-3 text-left">
            {milestones.map((milestone, index) => (
              <motion.li
                key={index}
                initial={{ opacity: 0, x: -10 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 0.6 + index * 0.1 }}
                className="flex items-start gap-3 text-white/70"
              >
                <span className="text-emerald-400 mt-1">✓</span>
                {milestone}
              </motion.li>
            ))}
          </ul>
        </motion.div>

        {/* Stage Context */}
        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1 }}
          className="text-white/40 mt-8 text-sm"
        >
          Early stage: Focused on product-market fit before aggressive growth
        </motion.p>
      </motion.div>
    </div>
  );
}
