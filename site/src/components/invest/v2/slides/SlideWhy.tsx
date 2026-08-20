'use client';

import React from 'react';
import { Zap, Shield, TrendingUp } from 'lucide-react';
import { Slide, SlideTitle, SlideSubtitle, Card, CardTitle, CardText, Grid } from '../shared';

const reasons = [
  {
    icon: Zap,
    title: 'Proven Product',
    points: [
      '100% customer retention',
      '450+ production deployments',
      'Real revenue, real customers',
    ],
    color: 'text-[#10b981]',
  },
  {
    icon: Shield,
    title: 'Defensible Moat',
    points: [
      '120+ open source modules',
      'Multi-cloud expertise',
      '13+ years of domain knowledge',
    ],
    color: 'text-white',
  },
  {
    icon: TrendingUp,
    title: 'Massive Market',
    points: [
      '$20B+ DevOps tools market',
      'Every company needs this',
      'AI tailwind accelerating',
    ],
    color: 'text-white',
  },
];

export default function SlideWhy() {
  return (
    <Slide className="!justify-start !pt-24 sm:!pt-28 md:!pt-32">
      <SlideTitle>Why Invest Now</SlideTitle>
      <SlideSubtitle className="mb-3 sm:mb-8 md:mb-10 text-xs sm:text-sm">
        Early Stage Opportunity with Proven Foundation
      </SlideSubtitle>

      {/* Three Reasons - more compact on mobile, ~30% larger on desktop */}
      <Grid cols={3} gap="sm" className="mb-3 sm:mb-10 md:mb-12 sm:gap-6 md:gap-8">
        {reasons.map((reason) => (
          <Card key={reason.title} className="text-left p-2 sm:p-5 md:p-6">
            <div className="flex items-center gap-1.5 sm:gap-3 mb-1 sm:mb-4">
              <reason.icon className={`w-4 h-4 sm:w-7 sm:h-7 md:w-8 md:h-8 ${reason.color}`} />
              <CardTitle className="text-xs sm:text-lg md:text-xl">{reason.title}</CardTitle>
            </div>
            <ul className="space-y-0.5 sm:space-y-2 md:space-y-2.5">
              {reason.points.map((point, index) => (
                <li key={index} className="flex items-start gap-1 sm:gap-2">
                  <span className={`mt-0.5 text-[10px] sm:text-sm ${reason.color}`}>•</span>
                  <CardText className="text-[10px] sm:text-sm md:text-base leading-tight">{point}</CardText>
                </li>
              ))}
            </ul>
          </Card>
        ))}
      </Grid>

      {/* Bottom Line */}
      <div className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-2 sm:p-5 md:p-6 max-w-lg sm:max-w-xl md:max-w-2xl mx-auto">
        <p className="text-xs sm:text-lg md:text-xl text-white font-medium">
          $500K to Capture a $20B+ Market Opportunity
        </p>
        <p className="text-[10px] sm:text-sm md:text-base text-[#666] mt-0.5 sm:mt-2">
          Join Us at the Ground Floor of the Next Cloud Provider
        </p>
      </div>
    </Slide>
  );
}

