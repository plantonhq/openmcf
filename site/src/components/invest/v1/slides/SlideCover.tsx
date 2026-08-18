'use client';

import React from 'react';
import { motion } from 'framer-motion';

export default function SlideCover() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-[#0a0a0a] via-[#1a1a1a] to-[#0a0a0a] p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="text-center max-w-4xl"
      >
        {/* Logo/Brand */}
        <div className="mb-8">
          <span className="text-5xl md:text-6xl font-extrabold text-white">
            Planton
          </span>
          <span className="text-5xl md:text-6xl font-extrabold text-white">
            Cloud
          </span>
        </div>

        {/* Tagline */}
        <h1 className="text-3xl md:text-5xl font-bold text-white mb-6 leading-tight">
          The Self-Service Cloud Platform
        </h1>

        {/* One-liner */}
        <p className="text-xl md:text-2xl text-white/70 mb-12 max-w-2xl mx-auto">
          Deploy infrastructure and services to any cloud.
          <br />
          Self-service by design.
        </p>

        {/* Key Metrics Preview */}
        <div className="flex flex-wrap justify-center gap-8 mb-12">
          <div className="text-center">
            <div className="text-4xl font-bold text-white">3+</div>
            <div className="text-sm text-white/50">Years Building</div>
          </div>
          <div className="text-center">
            <div className="text-4xl font-bold text-white">$500K+</div>
            <div className="text-sm text-white/50">Self-Funded</div>
          </div>
          <div className="text-center">
            <div className="text-4xl font-bold text-white">Seed</div>
            <div className="text-sm text-white/50">Stage</div>
          </div>
        </div>

        {/* Date */}
        <p className="text-sm text-white/40">
          December 2025
        </p>
      </motion.div>
    </div>
  );
}
