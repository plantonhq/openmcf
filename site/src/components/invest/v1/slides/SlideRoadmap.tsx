'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { CheckCircle, Cloud, Bot, Activity } from 'lucide-react';

const roadmapItems = [
  {
    phase: 'Now',
    title: 'The Self-Service Cloud Platform',
    icon: CheckCircle,
    color: 'emerald',
    items: [
      'Infra Hub - Infrastructure provisioning',
      'Service Hub — Vercel for Backend, In Your Own Cloud',
      'Multi-cloud support (AWS, GCP, Azure)',
    ],
    status: 'Live',
  },
  {
    phase: 'Q1 2026',
    title: 'Observability',
    icon: Activity,
    color: 'cyan',
    items: [
      'Full-stack observability built-in',
      'Strong demand from early customers',
      'Design complete, development underway',
    ],
    status: 'In Design',
  },
  {
    phase: '2026',
    title: 'AI Agent Fleet',
    icon: Bot,
    color: 'violet',
    items: [
      'Specialized agents per DevOps function',
      'Pipeline troubleshooter agent',
      'Infrastructure optimization agent',
    ],
    status: 'In R&D',
  },
  {
    phase: 'Long Term',
    title: 'Become a Cloud Provider',
    icon: Cloud,
    color: 'pink',
    items: [
      'Join AWS, GCP, Azure, DigitalOcean',
      'Simple, focused service catalog',
      'Software layer on Tier-2 data centers',
    ],
    status: 'Vision',
  },
];

export default function SlideRoadmap() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-[#0a0a0a] via-[#1a1a1a] to-[#0a0a0a] p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center max-w-5xl w-full"
      >
        {/* Title */}
        <h2 className="text-4xl md:text-5xl font-bold text-white mb-4">
          Roadmap
        </h2>
        <p className="text-xl text-white/60 mb-12">
          Building essential infrastructure, step by step
        </p>

        {/* Roadmap Items */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {roadmapItems.map((item, index) => {
            const colorClasses = {
              emerald: {
                border: 'border-emerald-500/30',
                icon: 'text-emerald-400',
                badge: 'bg-emerald-500/20 text-emerald-400',
              },
              cyan: {
                border: 'border-white/30',
                icon: 'text-white',
                badge: 'bg-white/20 text-white',
              },
              violet: {
                border: 'border-white/20',
                icon: 'text-white',
                badge: 'bg-white/10 text-white',
              },
              pink: {
                border: 'border-white/20',
                icon: 'text-white',
                badge: 'bg-white/10 text-white',
              },
            };
            const colors = colorClasses[item.color as keyof typeof colorClasses];

            return (
              <motion.div
                key={item.phase}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.1 * (index + 1) }}
                className={`bg-white/5 border rounded-2xl p-5 text-left ${colors.border}`}
              >
                {/* Header */}
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2">
                    <item.icon className={`w-4 h-4 ${colors.icon}`} />
                    <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${colors.badge}`}>
                      {item.status}
                    </span>
                  </div>
                </div>

                {/* Phase */}
                <div className="text-xs text-white/50 mb-1">{item.phase}</div>
                <h3 className="text-lg font-bold text-white mb-3">{item.title}</h3>

                {/* Items */}
                <ul className="space-y-1.5">
                  {item.items.map((listItem, i) => (
                    <li key={i} className="text-xs text-white/60 flex items-start gap-2">
                      <span className={`mt-0.5 ${colors.icon}`}>•</span>
                      {listItem}
                    </li>
                  ))}
                </ul>
              </motion.div>
            );
          })}
        </div>

        {/* Vision Statement */}
        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="text-white/40 mt-10 text-sm max-w-2xl mx-auto"
        >
          Vision: Planton becomes a cloud provider alongside AWS, GCP, Azure—our moat is the control plane and DX
        </motion.p>
      </motion.div>
    </div>
  );
}
