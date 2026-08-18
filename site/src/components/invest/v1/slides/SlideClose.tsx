'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { ArrowRight, Mail, Globe, Github } from 'lucide-react';

export default function SlideClose() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-[#0a0a0a] via-[#1a1a1a] to-[#0a0a0a] p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center max-w-4xl w-full"
      >
        {/* Main Message */}
        <h2 className="text-4xl md:text-6xl font-bold text-white mb-6 leading-tight">
          DevOps shouldn&apos;t be
          <br />
          <span className="text-white">
            this hard
          </span>
        </h2>

        <p className="text-xl md:text-2xl text-white/70 mb-12 max-w-2xl mx-auto">
          We&apos;ve spent 3+ years and $500K building the platform we wished existed.
          <br />
          Now we need partners to take it global.
        </p>

        {/* CTA */}
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.3 }}
          className="bg-gradient-to-r from-white to-[#666] rounded-2xl p-8 max-w-lg mx-auto mb-12"
        >
          <h3 className="text-2xl font-bold text-white mb-4">
            Let&apos;s Talk
          </h3>
          <p className="text-white/80 mb-6">
            Interested in learning more about Planton?
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <a
              href="mailto:swarup@planton.ai"
              className="inline-flex items-center justify-center gap-2 bg-white text-white font-semibold px-6 py-3 rounded-xl hover:bg-white/90 transition-colors"
            >
              <Mail className="w-5 h-5" />
              swarup@planton.ai
            </a>
          </div>
        </motion.div>

        {/* Links */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="flex flex-wrap justify-center gap-8"
        >
          <a
            href="https://planton.ai"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-2 text-white/60 hover:text-white transition-colors"
          >
            <Globe className="w-5 h-5" />
            planton.ai
          </a>
          <a
            href="https://planton.dev"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-2 text-white/60 hover:text-white transition-colors"
          >
            <Github className="w-5 h-5" />
            planton.dev
          </a>
          <a
            href="https://planton.ai/demo"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-2 text-white/60 hover:text-white transition-colors"
          >
            <ArrowRight className="w-5 h-5" />
            Interactive Demo
          </a>
        </motion.div>

        {/* Tagline */}
        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.7 }}
          className="text-white/30 mt-16 text-sm"
        >
          Planton • The Self-Service Cloud Platform • December 2025
        </motion.p>
      </motion.div>
    </div>
  );
}
