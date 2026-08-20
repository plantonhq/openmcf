'use client';

import { FC, ReactNode } from 'react';
import { motion } from 'framer-motion';
import { fadeInUp, defaultTransition } from '../shared';

interface HeroProps {
  /** Page title */
  title: ReactNode;
  /** Page subtitle/description */
  subtitle: ReactNode;
  className?: string;
}

export const Hero: FC<HeroProps> = ({ title, subtitle, className = '' }) => {
  return (
    <section
      className={`
        relative py-6 sm:py-8 md:py-10
        bg-[#0a0a0a]
        ${className}
      `}
    >
      <div className="relative w-full max-w-3xl mx-auto px-4 sm:px-6 text-center">
        <motion.h1
          initial="hidden"
          animate="visible"
          variants={fadeInUp}
          transition={defaultTransition}
          className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl font-semibold text-white leading-snug tracking-tight mb-4 sm:mb-6"
        >
          {title}
        </motion.h1>

        <motion.p
          initial="hidden"
          animate="visible"
          variants={fadeInUp}
          transition={{ ...defaultTransition, delay: 0.1 }}
          className="text-sm sm:text-base md:text-lg text-[#a0a0a0] max-w-2xl mx-auto"
        >
          {subtitle}
        </motion.p>
      </div>
    </section>
  );
};

export default Hero;
