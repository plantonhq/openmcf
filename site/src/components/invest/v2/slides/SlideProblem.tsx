'use client';

import React from 'react';
import { Clock, DollarSign, AlertTriangle, Puzzle } from 'lucide-react';
import { Slide, SlideTitle, SlideSubtitle, Card, CardTitle, CardText, Callout } from '../shared';

const problems = [
  {
    icon: Clock,
    title: 'Weeks to Deploy',
    description: 'Infrastructure Setup Takes Weeks, Not Minutes',
  },
  {
    icon: DollarSign,
    title: 'Deployment Expertise Is Scarce',
    description: 'Every Company Needs Ops Talent That Is Hard to Hire',
  },
  {
    icon: AlertTriangle,
    title: 'Multi-Cloud Chaos',
    description: 'AWS, GCP, Azure—Each with Different Tools',
  },
  {
    icon: Puzzle,
    title: 'Tool Sprawl',
    description: 'Multiple Platforms for Infra, CI/CD, and Deployments',
  },
];

export default function SlideProblem() {
  return (
    <Slide>
      <SlideTitle>The Problem</SlideTitle>
      <SlideSubtitle className="mb-4 sm:mb-6">
        Cloud Infrastructure Is Complex. Every Company Faces the Same Struggles.
      </SlideSubtitle>

      {/* Problem Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4 mb-4 sm:mb-6">
        {problems.map((problem) => (
          <Card key={problem.title} className="text-left">
            <problem.icon className="w-5 h-5 sm:w-6 sm:h-6 text-[#ef4444]/70 mb-2" />
            <CardTitle className="mb-1 text-sm sm:text-base">{problem.title}</CardTitle>
            <CardText className="text-xs">{problem.description}</CardText>
          </Card>
        ))}
      </div>

      {/* Bottom Line */}
      <Callout variant="default" className="!border-[#ef4444]/20 !bg-[#ef4444]/5 max-w-3xl">
        <p className="text-sm sm:text-base md:text-lg text-white font-medium sm:whitespace-nowrap">
          80% of cloud services are the same patterns—<span className="text-[#ef4444]/80">rebuilt from scratch at every company.</span>
        </p>
      </Callout>
    </Slide>
  );
}

