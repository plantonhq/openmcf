'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { Code, Palette, Layout, Briefcase } from 'lucide-react';

const teamMembers = [
  {
    name: 'Swarup Donepudi',
    role: 'Founder & CEO',
    icon: Code,
    description: '10+ years DevOps. Built entire platform. $500K+ personal investment.',
    highlight: true,
  },
  {
    name: 'Suresh Attaluri',
    role: 'Co-Founder & Backend Engineer',
    icon: Code,
    description: 'College classmate. 10 years collaboration. Backend systems.',
    highlight: false,
  },
  {
    name: 'Irshad Ahmed',
    role: 'Lead UX Designer',
    icon: Palette,
    description: '5 years with team. Full-time from day one. All product design.',
    highlight: false,
  },
  {
    name: 'Satish Lakhani',
    role: 'Full-Stack Engineer',
    icon: Layout,
    description: 'Built entire web console. Transitioning to full-stack.',
    highlight: false,
  },
  {
    name: 'Avinash Sana',
    role: 'Operations & BD',
    icon: Briefcase,
    description: 'College classmate. Handles all non-technical operations.',
    highlight: false,
  },
];

export default function SlideTeam() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-[#0a0a0a] via-[#1a1a1a] to-[#0a0a0a] p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center max-w-5xl w-full"
      >
        {/* Title */}
        <h2 className="text-4xl md:text-5xl font-bold text-white mb-4">
          Team
        </h2>
        <p className="text-xl text-white/60 mb-10">
          Small, focused, committed for 3+ years
        </p>

        {/* Team Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
          {teamMembers.map((member, index) => (
            <motion.div
              key={member.name}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 * (index + 1) }}
              className={`p-5 rounded-2xl text-left ${
                member.highlight
                  ? 'bg-gradient-to-br from-white/20 to-[#666]/20 border border-white/20 md:col-span-2 lg:col-span-1'
                  : 'bg-white/5 border border-white/10'
              }`}
            >
              <div className="flex items-start gap-3">
                <div className={`p-2 rounded-lg ${
                  member.highlight ? 'bg-white/10' : 'bg-white/10'
                }`}>
                  <member.icon className={`w-5 h-5 ${
                    member.highlight ? 'text-white' : 'text-white/60'
                  }`} />
                </div>
                <div className="flex-1">
                  <h3 className="text-white font-semibold">{member.name}</h3>
                  <p className={`text-sm ${
                    member.highlight ? 'text-white/70' : 'text-white/50'
                  }`}>{member.role}</p>
                  <p className="text-sm text-white/40 mt-2">{member.description}</p>
                </div>
              </div>
            </motion.div>
          ))}
        </div>

        {/* Team Strengths */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="bg-white/5 border border-white/10 rounded-xl p-6 max-w-2xl mx-auto"
        >
          <h3 className="text-white font-semibold mb-4">Why This Team Wins</h3>
          <div className="grid grid-cols-2 gap-4 text-left">
            <div className="flex items-center gap-2 text-white/60 text-sm">
              <span className="text-emerald-400">✓</span>
              Deep domain expertise
            </div>
            <div className="flex items-center gap-2 text-white/60 text-sm">
              <span className="text-emerald-400">✓</span>
              3+ years building together
            </div>
            <div className="flex items-center gap-2 text-white/60 text-sm">
              <span className="text-emerald-400">✓</span>
              $500K+ skin in the game
            </div>
            <div className="flex items-center gap-2 text-white/60 text-sm">
              <span className="text-emerald-400">✓</span>
              Production platform shipped
            </div>
          </div>
        </motion.div>
      </motion.div>
    </div>
  );
}
