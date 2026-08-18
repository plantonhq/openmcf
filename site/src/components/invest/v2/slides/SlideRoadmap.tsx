'use client';

import React from 'react';
import { CheckCircle, Activity, Bot, Cloud } from 'lucide-react';
import { Slide, SlideTitle, SlideSubtitle, RoadmapItem } from '../shared';

const roadmapItems = [
  {
    phase: 'Now',
    title: 'The Self-Service Cloud Platform',
    icon: <CheckCircle className="w-4 h-4" />,
    color: 'emerald' as const,
    items: ['Infra Hub + Service Hub Live', 'Multi-Cloud Support'],
    status: 'Live',
  },
  {
    phase: 'Q1 2026',
    title: 'Observability',
    icon: <Activity className="w-4 h-4" />,
    color: 'cyan' as const,
    items: ['Full-Stack Observability', 'Strong Customer Demand'],
    status: 'In Design',
  },
  {
    phase: '2026',
    title: 'AI Agent Fleet',
    icon: <Bot className="w-4 h-4" />,
    color: 'violet' as const,
    items: ['Cursor for Cloud Infrastructure', 'Specialized DevOps Agents'],
    status: 'In R&D',
  },
  {
    phase: 'Long Term',
    title: 'Cloud Provider',
    icon: <Cloud className="w-4 h-4" />,
    color: 'pink' as const,
    items: ['Join AWS, GCP, Azure Tier', 'Simple Service Catalog'],
    status: 'Vision',
  },
];

export default function SlideRoadmap() {
  return (
    <Slide className="!justify-start !pt-24 sm:!pt-28 md:!pt-32 sm:px-4 md:px-6">
      <SlideTitle>Roadmap</SlideTitle>
      <SlideSubtitle className="mb-2 sm:mb-16 md:mb-20 text-xs sm:text-sm">
        Building Essential Infrastructure, Step by Step
      </SlideSubtitle>

      {/* Roadmap Items - 4 column on desktop, single column centered on mobile */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-1.5 sm:gap-4 md:gap-5 lg:gap-6 mb-3 sm:mb-8 md:mb-10 max-w-xs sm:max-w-none mx-auto w-full px-0 sm:px-2">
        {roadmapItems.map((item) => (
          <RoadmapItem
            key={item.phase}
            phase={item.phase}
            title={item.title}
            status={item.status}
            items={item.items}
            color={item.color}
            icon={item.icon}
          />
        ))}
      </div>

      {/* Vision Statement */}
      <p className="text-xs sm:text-base md:text-lg text-[#666] max-w-xl sm:max-w-2xl mx-auto text-center">
        Vision: Planton Becomes a Cloud Provider Alongside AWS, GCP, Azure
        <br />
        Our Moat is the Control Plane and DX
      </p>
    </Slide>
  );
}

