'use client';

import { FC, ReactNode } from 'react';
import Link from 'next/link';
import { motion } from 'framer-motion';
import { fadeInUp, defaultTransition } from '../shared';

interface CTAButton {
  href: string;
  label: ReactNode;
  variant?: 'primary' | 'secondary' | 'ghost';
}

interface CTASectionProps {
  /** Section title */
  title: string;
  /** Section description */
  description: ReactNode;
  /** Primary CTA button */
  primaryButton: CTAButton;
  /** Secondary buttons */
  secondaryButtons?: CTAButton[];
  /** Email address to display */
  email?: string;
  /** Additional note below buttons */
  note?: ReactNode;
  className?: string;
}

export const CTASection: FC<CTASectionProps> = ({
  title,
  description,
  primaryButton,
  secondaryButtons,
  email,
  note,
  className = '',
}) => {
  const getButtonClasses = (variant: CTAButton['variant'] = 'primary') => {
    const base = 'inline-flex items-center justify-center px-6 py-3 rounded-full font-medium text-sm sm:text-base transition-all';
    
    switch (variant) {
      case 'primary':
        return `${base} bg-gradient-to-r from-white to-[#666] text-white hover:from-[#222] hover:to-[#333] shadow-lg shadow-black/25`;
      case 'secondary':
        return `${base} bg-white/10 text-white border border-white/20 hover:bg-white/15`;
      case 'ghost':
        return `${base} text-white/70 hover:text-white`;
      default:
        return base;
    }
  };

  return (
    <section
      className={`
        py-16 sm:py-20 md:py-24
        bg-gradient-to-b from-transparent via-[#666]/5 to-[#666]/5
        ${className}
      `}
    >
      <motion.div
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
        variants={fadeInUp}
        transition={defaultTransition}
        className="w-full max-w-3xl mx-auto px-4 sm:px-6 text-center"
      >
        <h2 className="text-2xl sm:text-3xl md:text-4xl font-bold text-white mb-4">
          {title}
        </h2>

        <p className="text-sm sm:text-base text-white/60 max-w-xl mx-auto mb-8">
          {description}
        </p>

        <div className="flex flex-col items-center gap-4">
          {/* Primary button */}
          <Link href={primaryButton.href} className={getButtonClasses(primaryButton.variant)}>
            {primaryButton.label}
          </Link>

          {/* Email */}
          {email && (
            <a
              href={`mailto:${email}`}
              className="flex items-center gap-2 text-sm sm:text-base text-white hover:text-white/70 transition-colors"
            >
              <span>📧</span>
              <span>{email}</span>
            </a>
          )}

          {/* Secondary buttons */}
          {secondaryButtons && secondaryButtons.length > 0 && (
            <div className="flex flex-wrap gap-3 justify-center mt-2">
              {secondaryButtons.map((button, index) => (
                <Link
                  key={index}
                  href={button.href}
                  className={getButtonClasses(button.variant || 'secondary')}
                >
                  {button.label}
                </Link>
              ))}
            </div>
          )}

          {/* Note */}
          {note && (
            <p className="text-xs sm:text-sm text-white/40 mt-4 max-w-md">
              {note}
            </p>
          )}
        </div>
      </motion.div>
    </section>
  );
};

export default CTASection;
