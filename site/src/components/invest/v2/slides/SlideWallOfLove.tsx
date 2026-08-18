'use client';

import React from 'react';
import Image from 'next/image';
import { Slide, SlideTitle, SlideSubtitle } from '../shared';

const row1 = [
  {
    name: 'Harsha Gurram',
    role: 'Solo Developer',
    company: 'Jai.CX',
    quote: "Weeks of Terraform → under 1 hr. As Planton's first user, I shaped Infra Charts.",
    avatar: '/_site/images/customers/people/harsha-ch.jpeg',
    logo: '/_site/images/customers/logos/jai-cx.svg',
  },
  {
    name: 'Balaji Borra',
    role: 'DevOps Engineer',
    company: 'TYNYBAY',
    quote: 'I handle 8+ client projects—no more rewriting Terraform between clients.',
    avatar: '/_site/images/customers/people/balaji-borra.png',
    logo: '/_site/images/customers/logos/tynybay.png',
  },
];

const row2 = [
  {
    name: 'Rakesh Kandhi',
    role: 'Senior Developer',
    company: 'TYNYBAY',
    quote: 'I deploy to dev, staging, and prod without waiting on DevOps. Game changer.',
    avatar: '/_site/images/customers/people/rakesh-kandhi.jpeg',
    logo: '/_site/images/customers/logos/tynybay.png',
  },
  {
    name: 'Sai Saketh',
    role: 'DevOps',
    company: 'iorta TechNext',
    quote: 'Mature developer experience for our 7-person team without deep AWS expertise.',
    avatar: null,
    logo: '/_site/images/customers/logos/iorta.svg',
  },
];

const row3 = [
  {
    name: 'Rohith Reddy Gopu',
    role: 'CEO',
    company: 'TYNYBAY',
    quote: 'Planton helps us deliver compliant infrastructure for regulated industries.',
    avatar: '/_site/images/customers/people/rohith-reddy-gopu.jpeg',
    logo: '/_site/images/customers/logos/tynybay.png',
  },
];

function AvatarFallback({ name, size = 'sm' }: { name: string; size?: 'xs' | 'sm' | 'md' }) {
  const sizeClass = {
    xs: 'w-5 h-5 text-[8px]',
    sm: 'w-8 h-8 sm:w-10 sm:h-10 text-xs sm:text-sm',
    md: 'w-7 h-7 sm:w-10 sm:h-10 text-xs sm:text-sm',
  }[size];

  return (
    <div className={`${sizeClass} rounded-full bg-[#2a2a2a] flex items-center justify-center text-[#a0a0a0] font-semibold shrink-0`}>
      {name.charAt(0)}
    </div>
  );
}

function MobileTestimonialCard({ 
  name, role, company, quote, avatar, logo,
}: {
  name: string; role: string; company: string; quote: string; avatar: string | null; logo: string;
}) {
  return (
    <div className="bg-[#151515] border border-[#2a2a2a] rounded-lg p-1.5 text-left">
      <div className="flex items-center gap-1.5 mb-1">
        {avatar ? (
          <Image src={avatar} alt={name} width={20} height={20} className="w-5 h-5 rounded-full object-cover shrink-0" />
        ) : (
          <AvatarFallback name={name} size="xs" />
        )}
        <div className="min-w-0 flex-1">
          <p className="text-[9px] font-medium text-white truncate">{name}</p>
          <p className="text-[8px] text-[#666] truncate">{role}, {company}</p>
        </div>
        <div className="w-4 h-4 shrink-0">
          <Image src={logo} alt={company} width={16} height={16} className="w-full h-full object-contain brightness-0 invert opacity-60" />
        </div>
      </div>
      <p className="text-[8px] text-[#a0a0a0] italic line-clamp-2 leading-tight">
        &ldquo;{quote}&rdquo;
      </p>
    </div>
  );
}

function TestimonialCard({ 
  name, role, company, quote, avatar, logo, fullWidth = false,
}: {
  name: string; role: string; company: string; quote: string; avatar: string | null; logo: string; fullWidth?: boolean;
}) {
  if (fullWidth) {
    return (
      <div className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-4 sm:p-5 md:p-6 h-full flex flex-col items-center text-center">
        <p className="text-xs sm:text-base md:text-lg text-[#a0a0a0] italic mb-3 sm:mb-4">
          &ldquo;{quote}&rdquo;
        </p>
        <div className="flex items-center gap-2 sm:gap-3">
          {avatar ? (
            <Image src={avatar} alt={name} width={40} height={40} className="w-7 h-7 sm:w-10 sm:h-10 rounded-full object-cover shrink-0" />
          ) : (
            <AvatarFallback name={name} size="md" />
          )}
          <div className="min-w-0 text-left">
            <p className="text-xs sm:text-sm md:text-base font-medium text-white">{name}</p>
            <p className="text-xs sm:text-sm text-[#666]">{role}, {company}</p>
          </div>
          <div className="w-5 h-5 sm:w-7 sm:h-7 shrink-0">
            <Image src={logo} alt={company} width={28} height={28} className="w-full h-full object-contain brightness-0 invert opacity-60" />
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-4 sm:p-5 md:p-6 text-left h-full flex flex-col">
      <div className="flex items-center gap-2 sm:gap-3 mb-2 sm:mb-3">
        {avatar ? (
          <Image src={avatar} alt={name} width={40} height={40} className="w-8 h-8 sm:w-10 sm:h-10 rounded-full object-cover shrink-0" />
        ) : (
          <AvatarFallback name={name} size="sm" />
        )}
        <div className="min-w-0">
          <p className="text-xs sm:text-sm md:text-base font-medium text-white">{name}</p>
          <p className="text-xs sm:text-sm text-[#666]">{role}, {company}</p>
        </div>
        <div className="w-5 h-5 sm:w-7 sm:h-7 shrink-0 ml-auto">
          <Image src={logo} alt={company} width={28} height={28} className="w-full h-full object-contain brightness-0 invert opacity-60" />
        </div>
      </div>
      <p className="text-xs sm:text-base md:text-lg text-[#a0a0a0] italic flex-1">
        &ldquo;{quote}&rdquo;
      </p>
    </div>
  );
}

export default function SlideWallOfLove() {
  const allTestimonials = [...row1, ...row2, ...row3];
  
  return (
    <Slide className="!justify-start !pt-24 sm:!pt-28 md:!pt-32">
      <SlideTitle>They Shipped. We Listened.</SlideTitle>
      <SlideSubtitle className="mb-2 sm:mb-14 md:mb-16 text-[10px] sm:text-sm">
        Voices from Teams Who Moved to Production with Planton
      </SlideSubtitle>

      {/* Mobile: Compact 2-column grid */}
      <div className="sm:hidden grid grid-cols-2 gap-1.5 mx-auto mb-2">
        {allTestimonials.map((testimonial, index) => (
          <MobileTestimonialCard
            key={index}
            name={testimonial.name}
            role={testimonial.role}
            company={testimonial.company}
            quote={testimonial.quote}
            avatar={testimonial.avatar}
            logo={testimonial.logo}
          />
        ))}
      </div>

      {/* Desktop: 2-2-1 Layout */}
      <div className="hidden sm:flex flex-col gap-5 md:gap-6 max-w-4xl md:max-w-5xl mx-auto mb-6">
        <div className="grid grid-cols-2 gap-5 md:gap-6">
          {row1.map((testimonial, index) => (
            <TestimonialCard key={index} {...testimonial} />
          ))}
        </div>
        <div className="grid grid-cols-2 gap-5 md:gap-6">
          {row2.map((testimonial, index) => (
            <TestimonialCard key={index} {...testimonial} />
          ))}
        </div>
        <div className="w-full">
          {row3.map((testimonial, index) => (
            <TestimonialCard key={index} {...testimonial} fullWidth />
          ))}
        </div>
      </div>
    </Slide>
  );
}
