'use client';

import React, { useState, useCallback, useEffect, useRef } from 'react';
import { ArrowLeft, ArrowRight, ChevronRight } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

// Import slide components
import SlideCover from './slides/SlideCover';
import SlideProblem from './slides/SlideProblem';
import SlideSolution from './slides/SlideSolution';
import SlideProduct from './slides/SlideProduct';
import SlideComparison from './slides/SlideComparison';
import SlideTraction from './slides/SlideTraction';
import SlideCustomers from './slides/SlideCustomers';
import SlideWallOfLove from './slides/SlideWallOfLove';
import SlideMarket from './slides/SlideMarket';
import SlideRoadmap from './slides/SlideRoadmap';
import SlideTeam from './slides/SlideTeam';
import SlideAsk from './slides/SlideAsk';
import SlideWhy from './slides/SlideWhy';
import SlideClose from './slides/SlideClose';

// Slide configuration
const slideConfig = [
  { id: 'cover', name: 'Cover', component: SlideCover },
  { id: 'problem', name: 'Problem', component: SlideProblem },
  { id: 'solution', name: 'Solution', component: SlideSolution },
  { id: 'product', name: 'Product', component: SlideProduct },
  { id: 'comparison', name: 'Why Planton', component: SlideComparison },
  { id: 'traction', name: 'Traction', component: SlideTraction },
  { id: 'customers', name: 'Customers', component: SlideCustomers },
  { id: 'wall-of-love', name: 'Wall of Love', component: SlideWallOfLove },
  { id: 'market', name: 'Market', component: SlideMarket },
  { id: 'roadmap', name: 'Roadmap', component: SlideRoadmap },
  { id: 'team', name: 'Team', component: SlideTeam },
  { id: 'ask', name: 'The Ask', component: SlideAsk },
  { id: 'why-now', name: 'Why Now', component: SlideWhy },
  { id: 'close', name: 'Close', component: SlideClose },
] as const;

export default function InvestorDeckV2() {
  const [currentSlideIndex, setCurrentSlideIndex] = useState(0);
  const touchStartX = useRef<number | null>(null);
  const touchEndX = useRef<number | null>(null);

  const currentSlide = slideConfig[currentSlideIndex];
  const isFirstSlide = currentSlideIndex === 0;
  const isLastSlide = currentSlideIndex === slideConfig.length - 1;

  // Sync slide index from URL hash (initial load + browser back/forward).
  // The setState calls happen inside event callbacks / rAF to satisfy the
  // react-hooks/set-state-in-effect rule.
  useEffect(() => {
    const syncFromHash = () => {
      const hash = window.location.hash.slice(1);
      if (hash) {
        const index = slideConfig.findIndex(s => s.id === hash);
        if (index !== -1) setCurrentSlideIndex(index);
      }
    };

    const raf = requestAnimationFrame(syncFromHash);
    window.addEventListener('hashchange', syncFromHash);
    return () => {
      cancelAnimationFrame(raf);
      window.removeEventListener('hashchange', syncFromHash);
    };
  }, []);

  // Update hash when slide changes
  const updateHash = useCallback((slideId: string) => {
    window.history.replaceState(null, '', `#${slideId}`);
  }, []);

  const goToSlide = useCallback((index: number) => {
    if (index >= 0 && index < slideConfig.length) {
      setCurrentSlideIndex(index);
      updateHash(slideConfig[index].id);
    }
  }, [updateHash]);

  const goToHome = useCallback(() => {
    goToSlide(0);
  }, [goToSlide]);

  const navigateNext = useCallback(() => {
    if (!isLastSlide) {
      goToSlide(currentSlideIndex + 1);
    }
  }, [isLastSlide, currentSlideIndex, goToSlide]);

  const navigatePrev = useCallback(() => {
    if (!isFirstSlide) {
      goToSlide(currentSlideIndex - 1);
    }
  }, [isFirstSlide, currentSlideIndex, goToSlide]);

  // Keyboard navigation
  const handleKeyPress = useCallback(
    (e: KeyboardEvent) => {
      switch (e.key) {
        case 'ArrowRight':
        case ' ':
          e.preventDefault();
          navigateNext();
          break;
        case 'ArrowLeft':
          e.preventDefault();
          navigatePrev();
          break;
        case 'Home':
        case 'h':
        case 'H':
          e.preventDefault();
          goToHome();
          break;
        case 'End':
          e.preventDefault();
          goToSlide(slideConfig.length - 1);
          break;
        // Note: Escape is intentionally not bound here to allow modals to use it
      }
    },
    [navigateNext, navigatePrev, goToHome, goToSlide]
  );

  useEffect(() => {
    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [handleKeyPress]);

  // Touch/swipe navigation
  const handleTouchStart = (e: React.TouchEvent) => {
    touchStartX.current = e.touches[0].clientX;
    touchEndX.current = e.touches[0].clientX; // Initialize to same position
  };

  const handleTouchMove = (e: React.TouchEvent) => {
    touchEndX.current = e.touches[0].clientX;
  };

  const handleTouchEnd = () => {
    if (touchStartX.current === null || touchEndX.current === null) return;

    const diff = touchStartX.current - touchEndX.current;
    const threshold = 50; // Minimum swipe distance in pixels

    if (Math.abs(diff) > threshold) {
      if (diff > 0) {
        // Swipe left -> next slide
        navigateNext();
      } else {
        // Swipe right -> previous slide
        navigatePrev();
      }
    }

    touchStartX.current = null;
    touchEndX.current = null;
  };

  const CurrentSlideComponent = currentSlide.component;

  return (
    <div
      className="h-[100dvh] bg-[#0a0a0a] flex flex-col relative overflow-hidden touch-pan-x"
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleTouchEnd}
      style={{ touchAction: 'pan-x' }}
    >
      {/* Progress bar */}
      <div className="absolute top-0 left-0 right-0 h-1 bg-[#2a2a2a] z-50">
        <motion.div
          className="h-full bg-white"
          initial={{ width: 0 }}
          animate={{ width: `${((currentSlideIndex + 1) / slideConfig.length) * 100}%` }}
          transition={{ duration: 0.3 }}
        />
      </div>

      {/* Slide navigation dots */}
      <div className="absolute top-4 left-1/2 -translate-x-1/2 z-50 flex gap-1 sm:gap-1.5">
        {slideConfig.map((slide, index) => (
          <button
            key={slide.id}
            onClick={() => goToSlide(index)}
            className={`h-1.5 sm:h-2 rounded-full transition-all duration-300 ${
              index === currentSlideIndex
                ? 'bg-white w-4 sm:w-6'
                : 'bg-[#3a3a3a] hover:bg-[#666] w-1.5 sm:w-2'
            }`}
            title={slide.name}
            aria-label={`Go to ${slide.name} slide`}
          />
        ))}
      </div>

      {/* Main content */}
      <div className="flex-1 flex items-center justify-center">
        <AnimatePresence mode="wait">
          <motion.div
            key={currentSlide.id}
            initial={{ opacity: 0, x: 50 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -50 }}
            transition={{ duration: 0.3 }}
            className="w-full h-full"
          >
            <CurrentSlideComponent />
          </motion.div>
        </AnimatePresence>
      </div>

      {/* Mobile swipe hint - only on cover slide */}
      {isFirstSlide && (
        <motion.div
          className="sm:hidden absolute bottom-20 left-0 right-0 text-center z-50"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1, duration: 0.6 }}
        >
          <div className="inline-flex items-center gap-1.5 px-4 py-2 bg-[#1a1a1a] border border-[#2a2a2a] rounded-full">
            <span className="text-xs font-medium text-[#a0a0a0]">Swipe to Navigate</span>
            <ChevronRight className="w-3.5 h-3.5 text-[#a0a0a0]" />
          </div>
        </motion.div>
      )}

      {/* Navigation footer */}
      <div className="absolute bottom-0 left-0 right-0 bg-[#0a0a0a]/90 backdrop-blur-sm border-t border-[#2a2a2a] pb-3 sm:pb-4 md:pb-6 pt-3 sm:pt-4">
        <div className="max-w-7xl mx-auto px-3 sm:px-4 md:px-8">
          <div className="flex items-center justify-between">
            {/* Slide counter */}
            <div className="flex items-center gap-2 sm:gap-4">
              <span className="text-xs sm:text-sm font-medium text-white/60">
                {currentSlideIndex + 1}/{slideConfig.length}
              </span>
              <span className="hidden sm:inline text-xs sm:text-sm text-white/40">
                {currentSlide.name}
              </span>
            </div>

            {/* Navigation buttons - hidden on mobile */}
            <div className="hidden sm:flex gap-2 sm:gap-3">
              <button
                onClick={navigatePrev}
                disabled={isFirstSlide}
                className="flex items-center gap-1 sm:gap-2 px-2 sm:px-3 md:px-4 py-1.5 sm:py-2 bg-transparent text-white text-xs sm:text-sm border border-[#3a3a3a] hover:border-white hover:bg-white/5 disabled:opacity-30 disabled:cursor-not-allowed rounded-lg transition-all duration-300"
              >
                <ArrowLeft className="w-3 h-3 sm:w-4 sm:h-4" />
                <span className="hidden sm:inline">Back</span>
              </button>
              <button
                onClick={navigateNext}
                disabled={isLastSlide}
                className="flex items-center gap-1 sm:gap-2 px-2 sm:px-3 md:px-4 py-1.5 sm:py-2 bg-[#fff] text-black text-xs sm:text-sm font-medium disabled:opacity-30 disabled:cursor-not-allowed hover:bg-gray-200 rounded-lg transition-all duration-300"
              >
                <span className="hidden sm:inline">Next</span>
                <ArrowRight className="w-3 h-3 sm:w-4 sm:h-4" />
              </button>
            </div>
          </div>

          {/* Keyboard shortcuts hint */}
          <div className="hidden md:block text-center mt-3">
            <p className="text-xs text-[#666]">
              <kbd className="px-1.5 py-0.5 bg-[#1a1a1a] border border-[#2a2a2a] rounded text-[#a0a0a0]">←</kbd>{' '}
              <kbd className="px-1.5 py-0.5 bg-[#1a1a1a] border border-[#2a2a2a] rounded text-[#a0a0a0]">→</kbd>{' '}
              to navigate •{' '}
              <kbd className="px-1.5 py-0.5 bg-[#1a1a1a] border border-[#2a2a2a] rounded text-[#a0a0a0]">Home</kbd>{' '}
              to start
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

