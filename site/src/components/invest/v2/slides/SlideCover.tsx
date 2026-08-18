'use client';

import React from 'react';
import { Slide, SlideTitle, Badge } from '../shared';

export default function SlideCover() {
  return (
    <Slide>
      {/* Brand */}
      <div className="mb-4 sm:mb-6 md:mb-8">
        <span className="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-semibold text-white tracking-tight">
          Planton
        </span>
      </div>

      {/* Tagline */}
      <SlideTitle className="mb-3 sm:mb-4 md:mb-6">
        The Self-Service Cloud Platform
      </SlideTitle>

      {/* One-liner */}
      <p className="text-sm sm:text-base md:text-lg text-[#a0a0a0] mb-6 sm:mb-8 md:mb-10 max-w-2xl mx-auto px-4">
        Deploy Infrastructure and Services to Any Cloud.
        <br className="hidden sm:block" />
        <span className="sm:hidden"> </span>
        Self-Service by Design.
      </p>

      {/* Key Metrics */}
      <div className="flex flex-wrap justify-center gap-4 sm:gap-6 md:gap-8 mb-6 sm:mb-8 md:mb-10">
        <div className="text-center">
          <div className="text-2xl sm:text-3xl md:text-4xl font-bold tracking-tight text-white">3+</div>
          <div className="text-xs sm:text-sm text-[#666]">Years Building</div>
        </div>
        <div className="text-center">
          <div className="text-2xl sm:text-3xl md:text-4xl font-bold tracking-tight text-white">$500K+</div>
          <div className="text-xs sm:text-sm text-[#666]">Self-Funded</div>
        </div>
        <div className="text-center">
          <div className="text-2xl sm:text-3xl md:text-4xl font-bold tracking-tight text-white">100%</div>
          <div className="text-xs sm:text-sm text-[#666]">Retention</div>
        </div>
      </div>

      {/* Stage Badge */}
      <Badge className="mb-4">
        Seed Stage • Raising $500K
      </Badge>

      {/* Date */}
      <p className="text-xs sm:text-sm text-[#666]">
        January 2026
      </p>
    </Slide>
  );
}

