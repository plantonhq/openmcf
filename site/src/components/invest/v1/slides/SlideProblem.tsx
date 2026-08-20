'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { AlertTriangle, Clock, DollarSign, Users } from 'lucide-react';

const problems = [
  {
    icon: Clock,
    title: 'Weeks to Deploy',
    description: 'Setting up cloud infrastructure takes weeks, not minutes',
  },
  {
    icon: Users,
    title: 'Deployment Expertise Is Scarce',
    description: 'Every company needs ops talent that is hard to hire',
  },
  {
    icon: AlertTriangle,
    title: 'Multi-Cloud Chaos',
    description: 'AWS, GCP, Azure - each with different tools and APIs',
  },
  {
    icon: DollarSign,
    title: 'Repetitive Work',
    description: 'Same infrastructure patterns rebuilt company after company',
  },
];

export default function SlideProblem() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-[#1a0a0a] via-[#2d1515] to-[#1a0a0a] p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center max-w-5xl w-full"
      >
        {/* Title */}
        <h2 className="text-4xl md:text-5xl font-bold text-white mb-4">
          The Problem
        </h2>
        <p className="text-xl text-white/60 mb-12 max-w-2xl mx-auto">
          Cloud infrastructure is complex. Every company faces the same struggles.
        </p>

        {/* Problem Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-12">
          {problems.map((problem, index) => (
            <motion.div
              key={problem.title}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 * (index + 1) }}
              className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-2xl p-6 text-left hover:bg-white/10 transition-colors"
            >
              <problem.icon className="w-10 h-10 text-red-400 mb-4" />
              <h3 className="text-xl font-bold text-white mb-2">
                {problem.title}
              </h3>
              <p className="text-white/60">
                {problem.description}
              </p>
            </motion.div>
          ))}
        </div>

        {/* Bottom Line */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="bg-red-500/10 border border-red-500/30 rounded-xl p-6 max-w-2xl mx-auto"
        >
          <p className="text-xl text-white font-medium">
            80% of cloud services are the same patterns—
            <br />
            <span className="text-red-400">rebuilt from scratch at every company.</span>
          </p>
        </motion.div>
      </motion.div>
    </div>
  );
}
