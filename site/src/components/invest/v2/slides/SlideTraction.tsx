'use client';

import React from 'react';
import { TrendingUp } from 'lucide-react';
import { Slide, SlideTitle, SlideSubtitle, Metric, Callout } from '../shared';

const metrics = [
  { value: '450+', label: 'Deployments', sublabel: 'Production infra' },
  { value: '100%', label: 'Retention', sublabel: 'Since launch', highlight: true },
  { value: '5', label: 'Customers', sublabel: 'Paying' },
  { value: '$800+', label: 'MRR', sublabel: 'Growing' },
];

const milestones = [
  'Production Platform with Paying Customers',
  'Full Infra Hub + Service Hub Live',
  'Multi-Cloud: AWS, GCP, Azure, K8s',
  'Open Source on GitHub',
];

export default function SlideTraction() {
  return (
    <Slide className="!justify-start !pt-24 sm:!pt-28 md:!pt-32">
      <SlideTitle>Traction</SlideTitle>
      <SlideSubtitle className="mb-4 sm:mb-6 md:mb-8">
        Early Revenue with Production Platform
      </SlideSubtitle>

      {/* Metrics Grid */}
      <div className="flex flex-wrap justify-center gap-4 sm:gap-8 md:gap-12 mb-4 sm:mb-8 md:mb-10">
        {metrics.map((metric) => (
          <Metric
            key={metric.label}
            value={metric.value}
            label={metric.label}
            sublabel={metric.sublabel}
            highlight={metric.highlight}
          />
        ))}
      </div>

      {/* Milestones */}
      <Callout className="max-w-3xl sm:max-w-4xl sm:p-6 md:p-8">
        <div className="flex items-center justify-center gap-2 sm:gap-3 mb-3 sm:mb-5">
          <TrendingUp className="w-4 h-4 sm:w-5 sm:h-5 text-[#10b981]/70" />
          <h3 className="text-sm sm:text-lg md:text-xl font-semibold text-white">What We&apos;ve Built</h3>
        </div>
        <div className="inline-grid grid-cols-1 sm:grid-cols-2 gap-x-8 sm:gap-x-12 gap-y-1.5 sm:gap-y-3 text-left">
          {milestones.map((milestone, index) => (
            <div key={index} className="flex items-center gap-2 sm:gap-3 text-xs sm:text-sm md:text-base text-[#a0a0a0] whitespace-nowrap">
              <span className="text-[#10b981]/70 shrink-0">✓</span>
              <span>{milestone}</span>
            </div>
          ))}
        </div>
      </Callout>

      {/* Stage Context */}
      <div className="mt-4 sm:mt-6 md:mt-8">
        <p className="sm:hidden text-xs text-[#666]">
          Early Stage: Focused on Product-Market Fit Before Aggressive Growth
        </p>
        <div className="hidden sm:block bg-[#1a1a1a] border border-[#2a2a2a] rounded-xl px-6 py-3 md:px-8 md:py-4">
          <p className="text-sm md:text-base text-[#a0a0a0] font-medium">
            <span className="text-white">Early Stage:</span> Focused on Product-Market Fit Before Aggressive Growth
          </p>
        </div>
      </div>
    </Slide>
  );
}

