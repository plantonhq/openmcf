'use client';

import { FC, useState, useCallback, useId } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { FAQItemData, ChevronDownIcon } from '../shared';

interface FAQProps {
  /** Array of FAQ items */
  items: FAQItemData[];
  /** Allow multiple items open at once */
  allowMultiple?: boolean;
  className?: string;
}

export const FAQ: FC<FAQProps> = ({ items, allowMultiple = true, className = '' }) => {
  const [openItems, setOpenItems] = useState<Set<number>>(new Set());
  const baseId = useId();

  const toggleItem = useCallback(
    (index: number) => {
      setOpenItems((prev) => {
        const next = new Set(prev);
        if (next.has(index)) {
          next.delete(index);
        } else {
          if (!allowMultiple) {
            next.clear();
          }
          next.add(index);
        }
        return next;
      });
    },
    [allowMultiple]
  );

  return (
    <div className={`space-y-3 ${className}`}>
      {items.map((item, index) => {
        const isOpen = openItems.has(index);
        const questionId = `${baseId}-question-${index}`;
        const answerId = `${baseId}-answer-${index}`;

        return (
          <div
            key={index}
            className="border border-white/10 rounded-xl overflow-hidden bg-white/[0.02]"
          >
            <button
              type="button"
              id={questionId}
              aria-expanded={isOpen}
              aria-controls={answerId}
              onClick={() => toggleItem(index)}
              className="w-full flex items-center justify-between gap-3 p-4 sm:p-5 text-left hover:bg-white/5 transition-colors"
            >
              <span className="text-sm sm:text-base font-medium text-white">
                {item.question}
              </span>
              <motion.span
                animate={{ rotate: isOpen ? 180 : 0 }}
                transition={{ duration: 0.2 }}
                className="flex-shrink-0 text-white/50"
              >
                <ChevronDownIcon />
              </motion.span>
            </button>

            <AnimatePresence initial={false}>
              {isOpen && (
                <motion.div
                  id={answerId}
                  role="region"
                  aria-labelledby={questionId}
                  initial={{ height: 0, opacity: 0 }}
                  animate={{ height: 'auto', opacity: 1 }}
                  exit={{ height: 0, opacity: 0 }}
                  transition={{ duration: 0.2 }}
                  className="overflow-hidden"
                >
                  <div className="px-4 pb-4 sm:px-5 sm:pb-5 text-sm sm:text-base text-white/60">
                    {item.answer}
                  </div>
                </motion.div>
              )}
            </AnimatePresence>
          </div>
        );
      })}
    </div>
  );
};

export default FAQ;
