'use client';

import React from 'react';
import { ArrowRight, Mail, Globe } from 'lucide-react';
import { Slide } from '../shared';

// GitHub icon as inline SVG for authentic look
const GitHubIcon = ({ className }: { className?: string }) => (
  <svg className={className} viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
  </svg>
);

export default function SlideClose() {
  return (
    <Slide>
      {/* Main Message */}
      <h2 className="text-2xl sm:text-3xl md:text-4xl font-semibold text-white tracking-tight mb-3 sm:mb-4 md:mb-6 leading-snug">
        <span className="sm:whitespace-nowrap">
          DevOps Shouldn&apos;t Be This Hard
        </span>
      </h2>

      <p className="text-sm sm:text-base md:text-lg text-[#a0a0a0] mb-6 sm:mb-8 max-w-lg mx-auto px-4">
        <span className="sm:whitespace-nowrap">We&apos;ve Spent 3+ Years and $500K Building the Platform We Wished Existed.</span>
        <br />
        <span className="sm:whitespace-nowrap">Now We Need Partners to Take It Global.</span>
      </p>

      {/* CTA */}
      <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-xl p-4 sm:p-6 max-w-sm mx-auto mb-6 sm:mb-8">
        <h3 className="text-lg sm:text-xl font-semibold text-white mb-2 sm:mb-3">
          Let&apos;s Talk
        </h3>
        <p className="text-sm text-[#a0a0a0] mb-3 sm:mb-4">
          Interested in learning more?
        </p>
        <a
          href="mailto:swarup@planton.ai"
          className="inline-flex items-center justify-center gap-2 bg-[#fff] text-black font-medium px-4 sm:px-5 py-2 sm:py-2.5 rounded-lg hover:bg-gray-200 transition-all duration-300 text-sm"
        >
          <Mail className="w-4 h-4" />
          swarup@planton.ai
        </a>
      </div>

      {/* Links */}
      <div className="flex flex-wrap justify-center gap-4 sm:gap-6 text-xs sm:text-sm">
        <a
          href="https://planton.ai"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center gap-1.5 text-[#a0a0a0] hover:text-white transition-colors duration-300"
        >
          <Globe className="w-4 h-4" />
          planton.ai
        </a>
        <a
          href="https://github.com/plantonhq/planton"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center gap-1.5 text-[#a0a0a0] hover:text-white transition-colors duration-300"
        >
          <GitHubIcon className="w-4 h-4" />
          GitHub
        </a>
        <a
          href="https://planton.ai/demo"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center gap-1.5 text-[#a0a0a0] hover:text-white transition-colors duration-300"
        >
          <ArrowRight className="w-4 h-4" />
          Demo
        </a>
      </div>

      {/* Tagline */}
      <p className="text-xs text-[#666] mt-6 sm:mt-8">
        Planton • The Self-Service Cloud Platform • January 2026
      </p>
    </Slide>
  );
}

