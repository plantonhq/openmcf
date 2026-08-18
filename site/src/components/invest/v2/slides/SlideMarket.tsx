'use client';

import React from 'react';
import { Slide, SlideTitle, SlideSubtitle, Metric, Grid, Card, CardTitle, CardText } from '../shared';

const whyNow = [
  {
    title: 'AI Wave',
    description: 'LLMs Enable New Automation',
  },
  {
    title: 'Cloud Complexity',
    description: 'AWS Has 200+ Services',
  },
  {
    title: 'Talent Shortage',
    description: 'Not Enough DevOps Engineers',
  },
];

export default function SlideMarket() {
  return (
    <Slide className="!justify-start !pt-24 sm:!pt-28 md:!pt-32">
      <SlideTitle>Market Opportunity</SlideTitle>
      <SlideSubtitle className="mb-4 sm:mb-16 md:mb-20">
        Cloud Spend Is Massive. Complexity Is the Problem.
      </SlideSubtitle>

      {/* TAM/SAM Numbers */}
      <Grid cols={3} gap="sm" className="mb-4 sm:mb-8 md:mb-10 sm:gap-6 md:gap-8">
        <Metric
          value="$150B+"
          label="Global Cloud Spend"
          sublabel="20%+ YoY growth"
        />
        <Metric
          value="$20B+"
          label="DevOps Tools Market"
          sublabel="Infrastructure automation"
          highlight
        />
        <Metric
          value="~50K"
          label="Mid-size Companies"
          sublabel="Without dedicated DevOps"
        />
      </Grid>

      {/* Why Now */}
      <Grid cols={3} gap="sm" className="max-w-2xl sm:max-w-3xl md:max-w-4xl mx-auto sm:gap-6 md:gap-8">
        {whyNow.map((item) => (
          <Card key={item.title} className="text-center sm:p-5 md:p-6">
            <CardTitle className="text-sm sm:text-base md:text-lg mb-1 sm:mb-2">{item.title}</CardTitle>
            <CardText className="text-xs sm:text-sm">{item.description}</CardText>
          </Card>
        ))}
      </Grid>
    </Slide>
  );
}

