'use client';

import React, { useState, useCallback, useEffect, useRef, ComponentType } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { Navigation } from './shared/navigation';
import { PresenterNotes } from './shared/presenter-notes';

// ============================================================================
// TYPES
// ============================================================================

export interface SlideConfig {
  id: string;
  name: string;
  component: ComponentType<SlideComponentProps>;
  presenterNotes?: string[];
}

export interface SlideComponentProps {
  notesVisible?: boolean;
}

export interface MeetsDeckProps {
  slides: SlideConfig[];
  guest: string;
  meetingDate: string;
  presenter?: string;
  company?: string;
  location?: string;
}

// ============================================================================
// MEETS DECK - Main Presentation Engine
// ============================================================================

/**
 * MeetsDeck - Core presentation engine for prospect meetings
 *
 * Features:
 * - Hash-based navigation (#slide-id)
 * - AnimatePresence slide transitions
 * - Keyboard navigation (arrows, space, N for notes, F for fullscreen)
 * - Touch swipe support
 * - Presenter notes state management
 * - Responsive layout
 */
export default function MeetsDeck({ slides }: MeetsDeckProps) {
  const [currentSlideIndex, setCurrentSlideIndex] = useState(0);
  const [notesVisible, setNotesVisible] = useState(false);
  const touchStartX = useRef<number | null>(null);
  const touchEndX = useRef<number | null>(null);

  const currentSlide = slides[currentSlideIndex];
  const isFirstSlide = currentSlideIndex === 0;
  const isLastSlide = currentSlideIndex === slides.length - 1;

  // ============================================================================
  // HASH NAVIGATION
  // ============================================================================

  // Sync slide index from URL hash (initial load + browser back/forward).
  useEffect(() => {
    const syncFromHash = () => {
      const hash = window.location.hash.slice(1);
      if (hash) {
        const index = slides.findIndex((s) => s.id === hash);
        if (index !== -1) setCurrentSlideIndex(index);
      }
    };

    const raf = requestAnimationFrame(syncFromHash);
    window.addEventListener('hashchange', syncFromHash);
    return () => {
      cancelAnimationFrame(raf);
      window.removeEventListener('hashchange', syncFromHash);
    };
  }, [slides]);

  // Update hash when slide changes
  const updateHash = useCallback((slideId: string) => {
    window.history.replaceState(null, '', `#${slideId}`);
  }, []);

  // ============================================================================
  // NAVIGATION HANDLERS
  // ============================================================================

  const goToSlide = useCallback(
    (index: number) => {
      if (index >= 0 && index < slides.length) {
        setCurrentSlideIndex(index);
        updateHash(slides[index].id);
      }
    },
    [slides, updateHash]
  );

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

  const toggleNotes = useCallback(() => {
    setNotesVisible((prev) => !prev);
  }, []);

  const toggleFullscreen = useCallback(() => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen();
    } else {
      document.exitFullscreen();
    }
  }, []);

  // ============================================================================
  // KEYBOARD NAVIGATION
  // ============================================================================

  const handleKeyPress = useCallback(
    (e: KeyboardEvent) => {
      switch (e.key) {
        case 'ArrowRight':
        case ' ':
        case 'PageDown':
          e.preventDefault();
          navigateNext();
          break;
        case 'ArrowLeft':
        case 'PageUp':
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
          goToSlide(slides.length - 1);
          break;
        case 'n':
        case 'N':
          e.preventDefault();
          toggleNotes();
          break;
        case 'f':
        case 'F':
          e.preventDefault();
          toggleFullscreen();
          break;
      }
    },
    [navigateNext, navigatePrev, goToHome, goToSlide, slides.length, toggleNotes, toggleFullscreen]
  );

  useEffect(() => {
    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [handleKeyPress]);

  // ============================================================================
  // TOUCH/SWIPE NAVIGATION
  // ============================================================================

  const handleTouchStart = (e: React.TouchEvent) => {
    touchStartX.current = e.touches[0].clientX;
    touchEndX.current = e.touches[0].clientX;
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

  // ============================================================================
  // RENDER
  // ============================================================================

  const CurrentSlideComponent = currentSlide.component;
  const slideInfos = slides.map((s) => ({ id: s.id, name: s.name }));

  return (
    <div
      className="h-[100dvh] bg-[#0a0a0a] flex flex-col relative overflow-hidden touch-pan-x"
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleTouchEnd}
      style={{ touchAction: 'pan-x' }}
    >
      {/* Navigation UI */}
      <Navigation
        slides={slideInfos}
        currentIndex={currentSlideIndex}
        onNavigate={goToSlide}
        onHome={goToHome}
        onPrev={navigatePrev}
        onNext={navigateNext}
        notesVisible={notesVisible}
        onToggleNotes={toggleNotes}
      />

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
            <CurrentSlideComponent notesVisible={notesVisible} />
          </motion.div>
        </AnimatePresence>
      </div>

      {/* Presenter Notes Overlay */}
      <PresenterNotes
        notes={currentSlide.presenterNotes || []}
        visible={notesVisible}
      />
    </div>
  );
}
