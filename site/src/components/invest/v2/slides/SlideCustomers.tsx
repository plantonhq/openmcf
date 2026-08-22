'use client';

import React from 'react';
import Image from 'next/image';
import { Slide, SlideTitle, SlideSubtitle, Callout } from '../shared';

const customers = [
  {
    name: 'Jai.CX',
    logo: '/_site/images/customers/logos/jai-cx.svg',
    className: 'scale-50', // Reduce size - logo appears larger than others
  },
  {
    name: 'TYNYBAY',
    logo: '/_site/images/customers/logos/tynybay.png',
    className: 'scale-[1.8]', // Larger to compensate for logo padding
  },
  {
    name: 'iorta TechNext',
    logo: '/_site/images/customers/logos/iorta.svg',
    className: '',
  },
  {
    name: 'Odwen',
    logo: '/_site/images/customers/logos/odwen.svg',
    className: '',
  },
  // TODO: Add Aadvil logo to /_site/images/customers/logos/aadvil.svg (First GCP Customer)
];

export default function SlideCustomers() {
  return (
    <Slide className="!justify-start !pt-24 sm:!pt-28 md:!pt-32">
      <SlideTitle>Customers</SlideTitle>
      <SlideSubtitle className="mb-6 sm:mb-10 md:mb-12">
        Trusted by IT Consulting Firms and Startups
      </SlideSubtitle>

      {/* Customer Logos */}
      <div className="flex flex-wrap items-center justify-center gap-8 sm:gap-20 md:gap-24 mb-8 sm:mb-10 md:mb-12">
        {customers.map((customer) => (
          <div key={customer.name} className="w-14 h-14 sm:w-32 sm:h-32 md:w-36 md:h-36 flex items-center justify-center">
            <Image 
              src={customer.logo} 
              alt={customer.name} 
              width={144} 
              height={144} 
              className={`object-contain brightness-0 invert opacity-80 hover:opacity-100 transition-opacity duration-300 ${customer.className}`}
            />
          </div>
        ))}
      </div>

      {/* Key Quote with avatar */}
      <Callout variant="highlight" className="max-w-2xl sm:max-w-3xl">
        <div className="flex items-center gap-3 sm:gap-4 mb-2 sm:mb-3">
          <Image 
            src="/_site/images/customers/people/balaji-borra.png" 
            alt="Balaji Borra" 
            width={48} 
            height={48} 
            className="w-8 h-8 sm:w-12 sm:h-12 rounded-full object-cover"
          />
          <div className="text-left">
            <p className="text-xs sm:text-sm md:text-base text-white font-medium">Balaji Borra</p>
            <p className="text-xs sm:text-sm text-[#666]">DevOps Engineer, TYNYBAY</p>
          </div>
        </div>
        <p className="text-sm sm:text-base md:text-lg text-[#a0a0a0] italic sm:whitespace-nowrap">
          &ldquo;Developers Self-Onboard Their Own Services. That&apos;s a Huge Win.&rdquo;
        </p>
      </Callout>

      {/* Retention highlight */}
      <div className="mt-6 sm:mt-8 md:mt-10 flex items-center justify-center gap-2 sm:gap-3">
        <span className="text-[#10b981]/70 text-base sm:text-lg">✓</span>
        <span className="text-sm sm:text-base md:text-lg text-white font-medium">
          100% Customer Retention Since Launch
        </span>
      </div>
    </Slide>
  );
}

