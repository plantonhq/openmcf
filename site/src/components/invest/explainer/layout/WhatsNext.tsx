'use client';

import { FC } from 'react';
import Link from 'next/link';
import { motion } from 'framer-motion';
import { fadeInUp, staggerContainer, defaultTransition, Container } from '../shared';

// ============================================================================
// TYPES
// ============================================================================

interface WhatsNextProps {
  className?: string;
}

// ============================================================================
// WHAT'S NEXT SECTION
// ============================================================================

/**
 * A clean, elegant "What's Next?" section that provides two clear paths
 * for investors without demanding homework or proof of interest.
 *
 * - "Explore All Resources" → /invest (landing page hub)
 * - "Ready to Invest?" → /invest/process (step-by-step Carta process)
 */
export const WhatsNext: FC<WhatsNextProps> = ({ className = '' }) => {
  return (
    <section
      className={`
        py-16 sm:py-20 md:py-24
        bg-gradient-to-b from-transparent via-[#666]/5 to-[#666]/5
        ${className}
      `}
    >
      <Container>
        <motion.div
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true }}
          variants={staggerContainer}
          className="text-center"
        >
          {/* Section Title */}
          <motion.h2
            variants={fadeInUp}
            transition={defaultTransition}
            className="text-2xl sm:text-3xl md:text-4xl font-bold text-white mb-10 sm:mb-12"
          >
            What&apos;s Next?
          </motion.h2>

          {/* Two-card layout */}
          <motion.div
            variants={fadeInUp}
            transition={{ ...defaultTransition, delay: 0.1 }}
            className="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6 max-w-2xl mx-auto"
          >
            {/* Card 1: Explore All Resources */}
            <Link href="/invest" className="block group">
              <div
                className="
                  h-full rounded-xl sm:rounded-2xl border border-white/10
                  bg-white/5 p-6 sm:p-8
                  hover:border-white/20 hover:bg-white/[0.07]
                  transition-all duration-300
                  text-left
                "
              >
                <div className="text-3xl mb-4">🏠</div>
                <h3 className="text-lg sm:text-xl font-bold text-white mb-2 group-hover:text-white transition-colors">
                  Explore All Resources
                </h3>
                <p className="text-sm text-white/50 leading-relaxed mb-4">
                  See the full investor hub with deck, opportunity, updates, and more.
                </p>
                <div className="text-sm font-medium text-white/40 flex items-center gap-1 group-hover:text-white group-hover:gap-2 transition-all">
                  View Hub
                  <span className="group-hover:translate-x-1 transition-transform">→</span>
                </div>
              </div>
            </Link>

            {/* Card 2: Ready to Invest? */}
            <Link href="/invest/process" className="block group">
              <div
                className="
                  h-full rounded-xl sm:rounded-2xl border border-white/20
                  bg-gradient-to-br from-white/10 to-[#666]/10 p-6 sm:p-8
                  hover:border-white/50 hover:shadow-lg hover:shadow-black/10
                  transition-all duration-300
                  text-left
                "
              >
                <div className="text-3xl mb-4">✓</div>
                <h3 className="text-lg sm:text-xl font-bold text-white mb-2">
                  Ready to Invest?
                </h3>
                <p className="text-sm text-white/50 leading-relaxed mb-4">
                  View the step-by-step process to invest via Carta.
                </p>
                <div className="text-sm font-medium text-white flex items-center gap-1 group-hover:gap-2 transition-all">
                  See Process
                  <span className="group-hover:translate-x-1 transition-transform">→</span>
                </div>
              </div>
            </Link>
          </motion.div>
        </motion.div>
      </Container>
    </section>
  );
};

export default WhatsNext;
