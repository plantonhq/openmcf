'use client';

import { FC, ReactNode } from 'react';
import { motion, Variants } from 'framer-motion';

// ============================================================================
// DESIGN TOKENS - Monochrome palette
// ============================================================================

export const colors = {
  // Semantic (functional indicators only)
  emerald: '#10b981',
  amber: '#f59e0b',
  red: '#ef4444',

  // Background layers (dark theme)
  bgPrimary: '#0a0a0a',
  bgSecondary: '#111111',
  bgTertiary: '#1a1a1a',

  // Text
  textPrimary: '#ededed',
  textSecondary: '#a0a0a0',
  textMuted: '#666666',

  // Borders
  borderDefault: '#2a2a2a',
  borderHover: '#3a3a3a',
} as const;

// ============================================================================
// ANIMATION PRESETS
// ============================================================================

export const fadeInUp: Variants = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0 },
};

export const fadeIn: Variants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1 },
};

export const staggerContainer: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
      delayChildren: 0.1,
    },
  },
};

export const defaultTransition = { duration: 0.4, ease: 'easeOut' };

// ============================================================================
// CONSTANTS
// ============================================================================

/** SAFE valuation cap in USD */
export const VALUATION_CAP = 7_000_000;

/** USD to INR exchange rate */
export const USD_TO_INR = 83;

/** Exit scenarios for calculator */
export const EXIT_SCENARIOS = {
  conservative: { value: 7_000_000, label: 'Conservative', description: 'Exit at cap' },
  base: { value: 28_000_000, label: 'Base Case', description: 'Things went well' },
  optimistic: { value: 100_000_000, label: 'Optimistic', description: 'We crushed it' },
  best: { value: 500_000_000, label: 'Best Case', description: 'Home run' },
} as const;

/** Investment preset amounts in USD */
export const INVESTMENT_PRESETS_USD = [10_000, 25_000, 50_000, 100_000] as const;

/** Investment preset amounts in INR */
export const INVESTMENT_PRESETS_INR = [2_500_000, 5_000_000, 8_300_000] as const;

// ============================================================================
// TYPES
// ============================================================================

export type Currency = 'USD' | 'INR';

export type PathAndYouGet = 'beginner' | 'intermediate' | 'advanced' | 'friend';
export type PathIfYouAre = 'vc-angel' | 'technical' | 'friend' | 'customer' | 'general';

export interface ListItemData {
  icon?: ReactNode;
  iconColor?: 'pink' | 'emerald' | 'red' | 'amber' | 'cyan' | 'default';
  text: ReactNode;
}

export interface FAQItemData {
  question: string;
  answer: ReactNode;
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

export function formatUSD(value: number): string {
  if (value >= 1_000_000) {
    const millions = value / 1_000_000;
    return `$${millions % 1 === 0 ? millions.toFixed(0) : millions.toFixed(1)}M`;
  }
  if (value >= 1_000) {
    const thousands = value / 1_000;
    return `$${thousands % 1 === 0 ? thousands.toFixed(0) : thousands.toFixed(1)}K`;
  }
  return `$${value.toLocaleString('en-US')}`;
}

export function formatINR(value: number): string {
  if (value >= 10_000_000) {
    const crores = value / 10_000_000;
    return `₹${crores % 1 === 0 ? crores.toFixed(0) : crores.toFixed(1)} Cr`;
  }
  if (value >= 100_000) {
    const lakhs = value / 100_000;
    return `₹${lakhs % 1 === 0 ? lakhs.toFixed(0) : lakhs.toFixed(1)} L`;
  }
  return `₹${value.toLocaleString('en-IN')}`;
}

export function formatCurrency(usdValue: number, currency: Currency): string {
  if (currency === 'INR') {
    return formatINR(usdValue * USD_TO_INR);
  }
  return formatUSD(usdValue);
}

export function formatPercent(value: number): string {
  return `${(value * 100).toFixed(2)}%`;
}

export function formatMultiple(value: number): string {
  return `${value.toFixed(1)}x`;
}

/** Revenue multiple used for valuation (industry standard for DevTools) */
export const REVENUE_MULTIPLE = 12;

export function valuationToARR(valuation: number): number {
  return valuation / REVENUE_MULTIPLE;
}

export function valuationToMRR(valuation: number): number {
  return valuationToARR(valuation) / 12;
}

export function formatARR(valuation: number, currency: Currency): string {
  const arr = valuationToARR(valuation);
  const mrr = valuationToMRR(valuation);
  if (currency === 'INR') {
    return `${formatINR(arr * USD_TO_INR)} ARR / ${formatINR(mrr * USD_TO_INR)} MRR`;
  }
  return `${formatUSD(arr)} ARR / ${formatUSD(mrr)} MRR`;
}

export function getValuationWithARR(valuation: number, currency: Currency): {
  valuation: string;
  arrMrr: string;
} {
  return {
    valuation: formatCurrency(valuation, currency),
    arrMrr: formatARR(valuation, currency),
  };
}

// ============================================================================
// SECTION WRAPPER
// ============================================================================

interface SectionProps {
  children: ReactNode;
  className?: string;
  visible?: boolean;
  background?: 'default' | 'gradient-subtle' | 'danger-subtle';
  id?: string;
}

export const Section: FC<SectionProps> = ({
  children,
  className = '',
  visible = true,
  id,
}) => {
  if (!visible) return null;

  return (
    <motion.section
      id={id}
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, margin: '-50px' }}
      variants={fadeInUp}
      transition={defaultTransition}
      className={`py-12 sm:py-16 md:py-20 ${className}`}
    >
      {children}
    </motion.section>
  );
};

// ============================================================================
// CONTAINER
// ============================================================================

interface ContainerProps {
  children: ReactNode;
  className?: string;
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl' | '4xl' | '5xl';
}

export const Container: FC<ContainerProps> = ({
  children,
  className = '',
  maxWidth = '3xl',
}) => {
  const maxWidthClass = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl',
    '2xl': 'max-w-2xl',
    '3xl': 'max-w-3xl',
    '4xl': 'max-w-4xl',
    '5xl': 'max-w-5xl',
  }[maxWidth];

  return (
    <div className={`w-full ${maxWidthClass} mx-auto px-4 sm:px-6 ${className}`}>
      {children}
    </div>
  );
};

// ============================================================================
// TYPOGRAPHY
// ============================================================================

interface TypographyProps {
  children: ReactNode;
  className?: string;
}

export const SectionTitle: FC<TypographyProps & { gradient?: boolean; color?: string }> = ({
  children,
  className = '',
}) => (
  <h2
    className={`text-xl sm:text-2xl md:text-3xl font-semibold leading-snug tracking-tight mb-4 text-white ${className}`}
  >
    {children}
  </h2>
);

export const SectionSubtitle: FC<TypographyProps> = ({ children, className = '' }) => (
  <p className={`text-sm sm:text-base text-[#a0a0a0] max-w-2xl ${className}`}>{children}</p>
);

export const BodyText: FC<TypographyProps> = ({ children, className = '' }) => (
  <div className={`text-sm sm:text-base text-[#a0a0a0] leading-relaxed ${className}`}>{children}</div>
);

export const SmallText: FC<TypographyProps> = ({ children, className = '' }) => (
  <p className={`text-xs sm:text-sm text-[#666] ${className}`}>{children}</p>
);

export const GradientText: FC<TypographyProps & { variant?: 'default' | 'pink' }> = ({
  children,
  className = '',
}) => (
  <span className={`text-white ${className}`}>{children}</span>
);

// ============================================================================
// CARD
// ============================================================================

interface CardProps {
  children: ReactNode;
  className?: string;
  variant?: 'default' | 'highlight' | 'success' | 'warning' | 'danger' | 'cyan';
}

export const Card: FC<CardProps> = ({ children, className = '', variant = 'default' }) => {
  const variantClass = {
    default: 'bg-[#151515] border-[#2a2a2a] hover:border-[#3a3a3a] hover:bg-[#1a1a1a]',
    highlight: 'bg-[#1a1a1a] border-[#3a3a3a]',
    success: 'bg-[#10b981]/10 border-[#10b981]/30',
    warning: 'bg-[#f59e0b]/10 border-[#f59e0b]/30',
    danger: 'bg-[#ef4444]/10 border-[#ef4444]/30',
    cyan: 'bg-[#151515] border-[#2a2a2a] hover:border-[#3a3a3a] hover:bg-[#1a1a1a]',
  }[variant];

  return (
    <div className={`rounded-xl border p-4 sm:p-5 md:p-6 transition-all duration-300 ${variantClass} ${className}`}>
      {children}
    </div>
  );
};

export const CardTitle: FC<TypographyProps & { color?: string }> = ({
  children,
  className = '',
}) => (
  <h3 className={`text-base sm:text-lg font-semibold text-white ${className}`}>
    {children}
  </h3>
);

export const CardText: FC<TypographyProps> = ({ children, className = '' }) => (
  <p className={`text-xs sm:text-sm text-[#a0a0a0] ${className}`}>{children}</p>
);

// ============================================================================
// CALLOUT
// ============================================================================

interface CalloutProps {
  children: ReactNode;
  className?: string;
  variant?: 'default' | 'highlight' | 'success' | 'warning' | 'danger' | 'cyan';
  title?: ReactNode;
  icon?: ReactNode;
}

export const Callout: FC<CalloutProps> = ({
  children,
  className = '',
  variant = 'default',
  title,
  icon,
}) => {
  const variantClass = {
    default: 'bg-[#151515] border-[#2a2a2a]',
    highlight: 'bg-[#151515] border-[#2a2a2a]',
    success: 'bg-[#10b981]/10 border-[#10b981]/30',
    warning: 'bg-[#f59e0b]/10 border-[#f59e0b]/30',
    danger: 'bg-[#ef4444]/10 border-[#ef4444]/30',
    cyan: 'bg-[#151515] border-[#2a2a2a]',
  }[variant];

  return (
    <div className={`border rounded-xl p-4 sm:p-5 md:p-6 ${variantClass} ${className}`}>
      {title && (
        <div className="flex items-center gap-2 mb-2 font-semibold text-white">
          {icon && <span>{icon}</span>}
          <span>{title}</span>
        </div>
      )}
      <div className="text-sm sm:text-base text-[#a0a0a0]">{children}</div>
    </div>
  );
};

// ============================================================================
// METRIC
// ============================================================================

interface MetricProps {
  value: ReactNode;
  label: string;
  sublabel?: string;
  highlight?: boolean;
  className?: string;
  valueColor?: string;
}

export const Metric: FC<MetricProps> = ({
  value,
  label,
  sublabel,
  highlight = false,
  className = '',
  valueColor,
}) => (
  <div
    className={`
      text-center p-4 sm:p-6
      ${highlight ? 'bg-[#1a1a1a] border border-[#3a3a3a] rounded-xl' : ''}
      ${className}
    `}
  >
    <div
      className="text-2xl sm:text-4xl md:text-5xl font-bold tracking-tight mb-1 sm:mb-2 text-white"
      style={valueColor ? { color: valueColor } : undefined}
    >
      {value}
    </div>
    <div className="text-xs sm:text-base text-[#a0a0a0]">{label}</div>
    {sublabel && <div className="text-xs sm:text-sm text-[#666] mt-1">{sublabel}</div>}
  </div>
);

// ============================================================================
// BADGE
// ============================================================================

interface BadgeProps {
  children: ReactNode;
  variant?: 'default' | 'success' | 'warning' | 'danger' | 'purple' | 'pink' | 'cyan';
  className?: string;
}

export const Badge: FC<BadgeProps> = ({ children, variant = 'default', className = '' }) => {
  const neutral = 'bg-[#2a2a2a] text-[#a0a0a0] border-[#3a3a3a]';
  const variantClass = {
    default: neutral,
    success: 'bg-[#10b981]/10 text-[#10b981] border-[#10b981]/30',
    warning: 'bg-[#f59e0b]/10 text-[#f59e0b] border-[#f59e0b]/30',
    danger: 'bg-[#ef4444]/10 text-[#ef4444] border-[#ef4444]/30',
    purple: neutral,
    pink: neutral,
    cyan: neutral,
  }[variant];

  return (
    <span
      className={`inline-flex items-center px-3 py-1.5 rounded-full border text-xs font-medium ${variantClass} ${className}`}
    >
      {children}
    </span>
  );
};

// ============================================================================
// GRID
// ============================================================================

interface GridProps {
  children: ReactNode;
  cols?: 2 | 3 | 4;
  gap?: 'sm' | 'md' | 'lg';
  className?: string;
}

export const Grid: FC<GridProps> = ({ children, cols = 2, gap = 'md', className = '' }) => {
  const colsClass = {
    2: 'grid-cols-1 sm:grid-cols-2',
    3: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3',
    4: 'grid-cols-2 lg:grid-cols-4',
  }[cols];

  const gapClass = {
    sm: 'gap-2 sm:gap-3',
    md: 'gap-4 sm:gap-6',
    lg: 'gap-6 sm:gap-8',
  }[gap];

  return <div className={`grid ${colsClass} ${gapClass} ${className}`}>{children}</div>;
};

// ============================================================================
// LIST
// ============================================================================

interface ListProps {
  items: ListItemData[];
  className?: string;
}

export const List: FC<ListProps> = ({ items, className = '' }) => {
  const getIconColor = (color?: ListItemData['iconColor']) => {
    switch (color) {
      case 'emerald':
        return 'text-[#10b981]';
      case 'red':
        return 'text-[#ef4444]';
      case 'amber':
        return 'text-[#f59e0b]';
      default:
        return 'text-[#a0a0a0]';
    }
  };

  return (
    <ul className={`space-y-2 sm:space-y-3 ${className}`}>
      {items.map((item, index) => (
        <li key={index} className="flex items-start gap-2 sm:gap-3">
          <span className={`mt-0.5 flex-shrink-0 ${getIconColor(item.iconColor)}`}>
            {item.icon ?? '→'}
          </span>
          <span className="text-sm sm:text-base text-[#a0a0a0]">{item.text}</span>
        </li>
      ))}
    </ul>
  );
};

// ============================================================================
// TABLE
// ============================================================================

interface TableColumn {
  header: ReactNode;
  accessor: string;
  className?: string;
  headerClassName?: string;
}

interface TableProps {
  columns: TableColumn[];
  data: Record<string, ReactNode>[];
  className?: string;
  highlightRow?: number;
}

export const Table: FC<TableProps> = ({ columns, data, className = '', highlightRow }) => (
  <div className={`overflow-x-auto -mx-4 px-4 ${className}`}>
    <table className="w-full text-sm">
      <thead>
        <tr className="border-b border-[#2a2a2a]">
          {columns.map((col, i) => (
            <th
              key={i}
              className={`text-left py-3 px-2 sm:px-4 font-semibold text-white ${col.headerClassName ?? ''}`}
            >
              {col.header}
            </th>
          ))}
        </tr>
      </thead>
      <tbody>
        {data.map((row, rowIndex) => (
          <tr
            key={rowIndex}
            className={`border-b border-[#1a1a1a] ${highlightRow === rowIndex ? 'bg-[#1a1a1a]' : ''}`}
          >
            {columns.map((col, colIndex) => (
              <td key={colIndex} className={`py-3 px-2 sm:px-4 text-[#a0a0a0] ${col.className ?? ''}`}>
                {row[col.accessor]}
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  </div>
);

// ============================================================================
// STEP (for numbered step lists)
// ============================================================================

interface StepProps {
  number: number;
  title: string;
  children: ReactNode;
  className?: string;
}

export const Step: FC<StepProps> = ({ number, title, children, className = '' }) => (
  <div className={`flex gap-4 sm:gap-6 ${className}`}>
    <div className="flex-shrink-0 w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-[#2a2a2a] flex items-center justify-center text-white font-bold text-sm sm:text-base">
      {number}
    </div>
    <div className="flex-1 pt-1">
      <h4 className="text-base sm:text-lg font-semibold text-white mb-2">{title}</h4>
      <div className="text-sm sm:text-base text-[#a0a0a0]">{children}</div>
    </div>
  </div>
);

export const Steps: FC<{ children: ReactNode; className?: string }> = ({ children, className = '' }) => (
  <div className={`space-y-6 sm:space-y-8 ${className}`}>{children}</div>
);

// ============================================================================
// ICONS (semantic colors preserved)
// ============================================================================

export const CheckIcon: FC<{ className?: string }> = ({ className = '' }) => (
  <svg
    className={`w-3.5 h-3.5 ${className}`}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
  </svg>
);

export const XIcon: FC<{ className?: string }> = ({ className = '' }) => (
  <svg
    className={`w-3.5 h-3.5 ${className}`}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M6 18L18 6M6 6l12 12" />
  </svg>
);

export const WarningIcon: FC<{ className?: string }> = ({ className = '' }) => (
  <svg
    className={`w-3.5 h-3.5 ${className}`}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2.5}
      d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
    />
  </svg>
);

export const ChevronDownIcon: FC<{ className?: string }> = ({ className = '' }) => (
  <svg
    className={`w-4 h-4 sm:w-5 sm:h-5 ${className}`}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
  </svg>
);

export const ArrowRightIcon: FC<{ className?: string }> = ({ className = '' }) => (
  <svg
    className={`w-4 h-4 sm:w-5 sm:h-5 ${className}`}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
  </svg>
);
