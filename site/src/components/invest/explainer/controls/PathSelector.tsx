'use client';

import React, { KeyboardEvent, useCallback, useRef, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { staggerContainer, fadeInUp, defaultTransition } from '../shared';

export interface PathOption<T extends string> {
  id: T;
  icon: string;
  title: string;
  description: string;
}

interface PathSelectorProps<T extends string> {
  /** Available path options */
  options: PathOption<T>[];
  /** Currently selected path */
  selectedPath: T | null;
  /** Selection handler */
  onSelect: (path: T) => void;
  /** Clear selection handler */
  onClear?: () => void;
  /** Section title */
  title?: string;
  /** Section subtitle */
  subtitle?: string;
  className?: string;
}

/**
 * Path selector with enhanced UX:
 * - Smooth collapse animation when a path is selected
 * - Selected path shows as a compact confirmation card
 * - "Change" button to reset selection
 * - Visual confirmation with checkmark
 */
export function PathSelector<T extends string>({
  options,
  selectedPath,
  onSelect,
  onClear,
  title,
  subtitle,
  className = '',
}: PathSelectorProps<T>): React.ReactElement {
  const containerRef = useRef<HTMLDivElement>(null);
  const [isExpanded, setIsExpanded] = useState(true);
  
  const selectedOption = selectedPath 
    ? options.find(o => o.id === selectedPath) 
    : null;

  const handleSelect = useCallback((path: T) => {
    onSelect(path);
    // Collapse after selection with a small delay for visual feedback
    setTimeout(() => setIsExpanded(false), 200);
  }, [onSelect]);

  const handleChange = useCallback(() => {
    setIsExpanded(true);
    onClear?.();
  }, [onClear]);

  const handleKeyDown = useCallback(
    (e: KeyboardEvent<HTMLButtonElement>, path: T, index: number) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        handleSelect(path);
        return;
      }

      // Arrow key navigation
      const buttons = containerRef.current?.querySelectorAll('button[role="option"]');
      if (!buttons) return;

      let nextIndex = index;
      if (e.key === 'ArrowDown' || e.key === 'ArrowRight') {
        nextIndex = (index + 1) % buttons.length;
      } else if (e.key === 'ArrowUp' || e.key === 'ArrowLeft') {
        nextIndex = (index - 1 + buttons.length) % buttons.length;
      }

      if (nextIndex !== index) {
        e.preventDefault();
        (buttons[nextIndex] as HTMLButtonElement).focus();
      }
    },
    [handleSelect]
  );

  return (
    <section className={`py-4 sm:py-6 ${className}`}>
      <div className="w-full max-w-3xl mx-auto px-4 sm:px-6">
        <AnimatePresence mode="wait">
          {/* Collapsed state: Show selected path as confirmation */}
          {selectedPath && !isExpanded && selectedOption ? (
            <motion.div
              key="collapsed"
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.3 }}
            >
              <SelectedConfirmation
                option={selectedOption}
                onChangeClick={handleChange}
              />
            </motion.div>
          ) : (
            /* Expanded state: Show all options */
            <motion.div
              key="expanded"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.3 }}
            >
              {(title || subtitle) && (
                <div className="text-center mb-4 sm:mb-5">
                  {title && (
                    <h2 className="text-base sm:text-lg md:text-xl font-bold text-white mb-1.5">
                      {title}
                    </h2>
                  )}
                  {subtitle && (
                    <p className="text-sm text-white/60">{subtitle}</p>
                  )}
                </div>
              )}

              <motion.div
                ref={containerRef}
                initial="hidden"
                animate="visible"
                variants={staggerContainer}
                className="space-y-2"
                role="listbox"
                aria-label="Select your background"
              >
                {options.map((option, index) => {
                  const isSelected = selectedPath === option.id;

                  return (
                    <motion.button
                      key={option.id}
                      variants={fadeInUp}
                      transition={defaultTransition}
                      type="button"
                      role="option"
                      aria-selected={isSelected}
                      tabIndex={0}
                      onClick={() => handleSelect(option.id)}
                      onKeyDown={(e) => handleKeyDown(e, option.id, index)}
                      className={`
                        w-full flex items-center gap-3 p-3 sm:p-4
                        rounded-xl border
                        text-left
                        transition-all duration-200
                        focus:outline-none focus:ring-2 focus:ring-white/50
                        ${isSelected
                          ? 'bg-gradient-to-r from-white/20 to-[#666]/20 border-white/50 scale-[1.02]'
                          : 'bg-white/5 border-white/10 hover:bg-white/10 hover:border-white/20'
                        }
                      `}
                    >
                      {/* Icon */}
                      <div
                        className={`
                          flex-shrink-0 w-9 h-9 sm:w-10 sm:h-10
                          flex items-center justify-center
                          text-lg sm:text-xl
                          rounded-lg transition-colors duration-200
                          ${isSelected ? 'bg-white/30' : 'bg-white/5'}
                        `}
                      >
                        {option.icon}
                      </div>

                      {/* Content */}
                      <div className="flex-1 min-w-0">
                        <div
                          className={`
                            text-sm sm:text-base font-semibold
                            ${isSelected ? 'text-white' : 'text-white/90'}
                          `}
                        >
                          {option.title}
                        </div>
                        <div className="text-xs sm:text-sm text-white/50 mt-0.5 line-clamp-1">
                          {option.description}
                        </div>
                      </div>

                      {/* Arrow / Check */}
                      <div
                        className={`
                          flex-shrink-0 text-lg sm:text-xl
                          transition-all duration-200
                          ${isSelected ? 'text-white scale-110' : 'text-white/30'}
                        `}
                      >
                        {isSelected ? '✓' : '→'}
                      </div>
                    </motion.button>
                  );
                })}
              </motion.div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </section>
  );
}

// ============================================================================
// SELECTED CONFIRMATION
// ============================================================================

interface SelectedConfirmationProps<T extends string> {
  option: PathOption<T>;
  onChangeClick: () => void;
}

/**
 * Compact confirmation card showing the selected path.
 * Provides clear visual feedback that selection is complete.
 */
function SelectedConfirmation<T extends string>({
  option,
  onChangeClick,
}: SelectedConfirmationProps<T>): React.ReactElement {
  return (
    <motion.div
      initial={{ scale: 0.95, opacity: 0 }}
      animate={{ scale: 1, opacity: 1 }}
      transition={{ type: 'spring', stiffness: 300, damping: 25 }}
      className="
        flex items-center gap-3 sm:gap-4 p-3 sm:p-4
        rounded-xl sm:rounded-2xl
        bg-gradient-to-r from-white/10 to-[#666]/10
        border border-white/20
      "
    >
      {/* Success indicator */}
      <div className="
        flex-shrink-0 w-8 h-8 sm:w-10 sm:h-10
        flex items-center justify-center
        rounded-full bg-[#10b981]/20
        text-[#10b981] text-sm sm:text-base
      ">
        ✓
      </div>

      {/* Icon */}
      <div className="
        flex-shrink-0 w-8 h-8 sm:w-10 sm:h-10
        flex items-center justify-center
        text-lg sm:text-xl
        rounded-lg bg-white/10
      ">
        {option.icon}
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0">
        <div className="text-xs text-white/40 uppercase tracking-wide">
          Your Background
        </div>
        <div className="text-sm sm:text-base font-medium text-white truncate">
          {option.title}
        </div>
      </div>

      {/* Change button */}
      <button
        type="button"
        onClick={onChangeClick}
        className="
          flex-shrink-0 px-3 py-1.5
          text-xs font-medium text-white/50
          rounded-lg border border-white/10
          hover:text-white hover:border-white/30
          transition-all duration-200
        "
      >
        Change
      </button>
    </motion.div>
  );
}

export default PathSelector;
