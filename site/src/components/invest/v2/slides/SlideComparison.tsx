'use client';

import React from 'react';
import { Slide, SlideTitle, SlideSubtitle, CheckIcon, XIcon, WarningIcon } from '../shared';

const comparisonData = [
  { feature: 'Setup Time', planton: '<1 hour', terraform: '1-2 weeks', vercel: 'N/A', heroku: 'minutes' },
  { feature: 'Backend CI/CD', planton: true, terraform: 'DIY', vercel: 'frontend', heroku: true },
  { feature: 'Your Cloud', planton: true, terraform: true, vercel: false, heroku: false },
  { feature: 'Infra Control', planton: true, terraform: true, vercel: false, heroku: false },
  { feature: 'Open Source', planton: true, terraform: 'partial', vercel: false, heroku: false },
];

type StatusType = 'yes' | 'no' | 'partial' | 'text';

function getStatus(value: boolean | string): { type: StatusType; display: string | null } {
  // Boolean true/false - just show icon, no text
  if (value === true) {
    return { type: 'yes', display: null };
  }
  if (value === false) {
    return { type: 'no', display: null };
  }
  // String values - show with appropriate styling
  if (value === '<1 hour') {
    return { type: 'yes', display: value };
  }
  if (value === 'N/A') {
    return { type: 'no', display: value };
  }
  // Partial/warning values
  if (value === 'DIY' || value === 'partial' || value === '1-2 weeks' || value === 'minutes' || value === 'frontend') {
    return { type: 'partial', display: value };
  }
  return { type: 'text', display: value };
}

function StatusCell({ value }: { value: boolean | string }) {
  const { type, display } = getStatus(value);
  
  // All cells use the same layout: icon at fixed position, optional text to the right
  // On mobile: just center the icon (no text shown)
  // On desktop: icon starts at consistent position, text flows right
  return (
    <div className="flex items-center justify-center sm:justify-start sm:pl-[calc(50%-10px)]">
      <span className="w-3.5 h-3.5 flex items-center justify-center shrink-0">
        {type === 'yes' && <CheckIcon />}
        {type === 'no' && <XIcon />}
        {type === 'partial' && <WarningIcon />}
      </span>
      {display && (
        <span className={`hidden sm:inline text-sm ml-1.5 whitespace-nowrap ${
          type === 'yes' ? 'text-[#10b981] font-medium' : 
          type === 'no' ? 'text-red-400' : 
          type === 'partial' ? 'text-[#a0a0a0]' : 'text-white/60'
        }`}>
          {display}
        </span>
      )}
    </div>
  );
}

export default function SlideComparison() {
  return (
    <Slide>
      <SlideTitle>Why Planton Wins</SlideTitle>
      <SlideSubtitle className="mb-4 sm:mb-6 sm:whitespace-nowrap">
        The Only Platform That&apos;s Open Source, Multi-Cloud, and No Lock-In
      </SlideSubtitle>

      {/* Comparison Table */}
      <div className="w-full max-w-3xl sm:max-w-4xl mx-auto bg-[#151515] border border-[#2a2a2a] rounded-xl overflow-hidden">
        {/* Header */}
        <div className="grid grid-cols-5 gap-1 sm:gap-4 p-2 sm:p-4 bg-[#1a1a1a] border-b border-[#2a2a2a]">
          <div className="text-xs sm:text-sm text-white/50 text-left">Feature</div>
          <div className="text-xs sm:text-sm text-white font-semibold text-center">Planton</div>
          <div className="text-xs sm:text-sm text-white/50 text-center">Terraform</div>
          <div className="text-xs sm:text-sm text-white/50 text-center">Vercel</div>
          <div className="text-xs sm:text-sm text-white/50 text-center">Heroku</div>
        </div>
        
        {/* Rows */}
        {comparisonData.map((row, index) => (
          <div 
            key={row.feature}
            className={`grid grid-cols-5 gap-1 sm:gap-4 p-2 sm:p-4 ${
              index < comparisonData.length - 1 ? 'border-b border-[#2a2a2a]' : ''
            }`}
          >
            <div className="text-xs sm:text-sm text-white/70 text-left">{row.feature}</div>
            <StatusCell value={row.planton} />
            <StatusCell value={row.terraform} />
            <StatusCell value={row.vercel} />
            <StatusCell value={row.heroku} />
          </div>
        ))}
      </div>

      {/* Key differentiator */}
      <p className="text-xs sm:text-sm text-[#666] mt-4 sm:mt-6 max-w-3xl sm:max-w-4xl mx-auto sm:whitespace-nowrap">
        <span className="text-[#10b981]/70">✓</span> Only Platform Combining PaaS Developer Experience + IaC Infrastructure + Your Cloud
      </p>
    </Slide>
  );
}

