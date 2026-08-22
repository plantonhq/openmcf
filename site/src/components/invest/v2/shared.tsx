'use client';

import { FC, ReactNode } from 'react';
import { motion } from 'framer-motion';
import Image from 'next/image';

// ============================================================================
// DESIGN TOKENS - Monochrome palette
// ============================================================================

export const colors = {
  // Semantic (functional indicators only)
  accentGreen: '#10b981',
  accentRed: '#ef4444',
  accentAmber: '#f59e0b',
  
  // Background layers
  bgPrimary: '#0a0a0a',
  bgSecondary: '#111111',
  bgCard: '#151515',
  bgCardHover: '#1a1a1a',
  
  // Text
  textPrimary: '#ededed',
  textSecondary: '#a0a0a0',
  textMuted: '#666666',
  
  // Borders
  border: '#2a2a2a',
  borderHover: '#3a3a3a',
};

// ============================================================================
// SLIDE WRAPPER - Mobile-first, single screen
// ============================================================================

interface SlideProps {
  children: ReactNode;
  className?: string;
  variant?: 'default' | 'problem' | 'solution' | 'gradient';
}

export const Slide: FC<SlideProps> = ({ 
  children, 
  className = '',
  variant: _variant = 'default'
}) => {
  return (
    <div 
      className={`
        h-[100dvh] overflow-hidden
        flex flex-col items-center justify-center
        px-4 py-6 sm:px-6 sm:py-8 md:px-8 md:py-12
        bg-[#0a0a0a]
        ${className}
      `}
    >
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4 }}
        className="w-full max-w-5xl mx-auto text-center"
      >
        {children}
      </motion.div>
    </div>
  );
};

// ============================================================================
// TYPOGRAPHY - Responsive scaling
// ============================================================================

interface TypographyProps {
  children: ReactNode;
  className?: string;
}

export const SlideTitle: FC<TypographyProps> = ({ children, className = '' }) => (
  <h2 className={`
    text-2xl sm:text-3xl md:text-4xl
    font-semibold text-white leading-snug tracking-tight
    ${className}
  `}>
    {children}
  </h2>
);

export const SlideSubtitle: FC<TypographyProps> = ({ children, className = '' }) => (
  <p className={`
    text-sm sm:text-base md:text-lg
    text-[#a0a0a0] font-normal mt-2 sm:mt-3 md:mt-4
    max-w-2xl mx-auto
    ${className}
  `}>
    {children}
  </p>
);

export const CardTitle: FC<TypographyProps> = ({ children, className = '' }) => (
  <h3 className={`
    text-base sm:text-lg md:text-xl
    font-semibold text-white
    ${className}
  `}>
    {children}
  </h3>
);

export const CardText: FC<TypographyProps> = ({ children, className = '' }) => (
  <p className={`
    text-xs sm:text-sm
    text-white/60
    ${className}
  `}>
    {children}
  </p>
);

// ============================================================================
// CARDS - Compact for mobile
// ============================================================================

interface CardProps {
  children: ReactNode;
  className?: string;
  variant?: 'default' | 'highlight' | 'success' | 'danger';
}

export const Card: FC<CardProps> = ({ 
  children, 
  className = '',
  variant = 'default'
}) => {
  const variantClass = {
    default: 'bg-[#151515] border-[#2a2a2a] hover:border-[#3a3a3a] hover:bg-[#1a1a1a]',
    highlight: 'bg-[#1a1a1a] border-[#3a3a3a]',
    success: 'bg-[#10b981]/10 border-[#10b981]/30',
    danger: 'bg-[#ef4444]/10 border-[#ef4444]/30',
  }[variant];

  return (
    <div className={`
      rounded-xl
      border p-4 sm:p-5 md:p-6
      transition-all duration-300
      ${variantClass}
      ${className}
    `}>
      {children}
    </div>
  );
};

// ============================================================================
// METRICS - Large numbers
// ============================================================================

interface MetricProps {
  value: string;
  label: string;
  sublabel?: string;
  highlight?: boolean;
  className?: string;
}

export const Metric: FC<MetricProps> = ({ 
  value, 
  label, 
  sublabel,
  highlight = false,
  className = ''
}) => (
  <div className={`
    text-center p-2 sm:p-4 md:p-6
    ${highlight ? 'bg-[#1a1a1a] border border-[#3a3a3a] rounded-xl' : ''}
    ${className}
  `}>
    <div className={`
      text-2xl sm:text-3xl md:text-4xl lg:text-5xl
      font-bold tracking-tight mb-1 sm:mb-2 text-white
    `}>
      {value}
    </div>
    <div className="text-xs sm:text-sm md:text-base text-[#a0a0a0]">{label}</div>
    {sublabel && (
      <div className="text-xs sm:text-sm text-[#666] mt-0.5 sm:mt-1">{sublabel}</div>
    )}
  </div>
);

// ============================================================================
// BADGES
// ============================================================================

interface BadgeProps {
  children: ReactNode;
  variant?: 'default' | 'success' | 'warning';
  className?: string;
}

export const Badge: FC<BadgeProps> = ({ 
  children, 
  variant = 'default',
  className = '' 
}) => {
  const variantClass = {
    default: 'bg-[#2a2a2a] text-[#a0a0a0] border-[#3a3a3a]',
    success: 'bg-[#10b981]/10 text-[#10b981] border-[#10b981]/30',
    warning: 'bg-[#f59e0b]/10 text-[#f59e0b] border-[#f59e0b]/30',
  }[variant];

  return (
    <span className={`
      inline-flex items-center
      px-3 py-1.5
      rounded-full border
      text-xs md:text-sm font-medium
      ${variantClass}
      ${className}
    `}>
      {children}
    </span>
  );
};

// ============================================================================
// ICONS (semantic colors preserved)
// ============================================================================

export const CheckIcon: FC<{ className?: string }> = ({ className = '' }) => (
  <svg className={`w-3.5 h-3.5 text-[#10b981]/70 ${className}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
  </svg>
);

export const XIcon: FC<{ className?: string }> = ({ className = '' }) => (
  <svg className={`w-3.5 h-3.5 text-[#ef4444]/70 ${className}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M6 18L18 6M6 6l12 12" />
  </svg>
);

export const WarningIcon: FC<{ className?: string }> = ({ className = '' }) => (
  <svg className={`w-3.5 h-3.5 text-[#f59e0b]/70 ${className}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
  </svg>
);

// ============================================================================
// GRID LAYOUTS
// ============================================================================

interface GridProps {
  children: ReactNode;
  cols?: 2 | 3 | 4;
  gap?: 'sm' | 'md';
  className?: string;
}

export const Grid: FC<GridProps> = ({ 
  children, 
  cols = 3,
  gap = 'md',
  className = '' 
}) => {
  const colsClass = {
    2: 'grid-cols-1 sm:grid-cols-2',
    3: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3',
    4: 'grid-cols-2 sm:grid-cols-2 lg:grid-cols-4',
  }[cols];

  const gapClass = {
    sm: 'gap-2 sm:gap-3',
    md: 'gap-3 sm:gap-4 md:gap-6',
  }[gap];

  return (
    <div className={`grid ${colsClass} ${gapClass} ${className}`}>
      {children}
    </div>
  );
};

// ============================================================================
// CALLOUT BOX
// ============================================================================

interface CalloutProps {
  children: ReactNode;
  variant?: 'default' | 'success' | 'highlight';
  className?: string;
}

export const Callout: FC<CalloutProps> = ({ 
  children, 
  variant = 'default',
  className = ''
}) => {
  const variantClass = {
    default: 'bg-[#151515] border-[#2a2a2a]',
    success: 'bg-[#10b981]/10 border-[#10b981]/30',
    highlight: 'bg-[#1a1a1a] border-[#3a3a3a]',
  }[variant];

  return (
    <div className={`
      border rounded-xl
      p-4 sm:p-5 md:p-6
      max-w-2xl mx-auto
      ${variantClass}
      ${className}
    `}>
      {children}
    </div>
  );
};

// ============================================================================
// TEAM MEMBER CARD
// ============================================================================

interface TeamMemberProps {
  name: string;
  role: string;
  description: ReactNode;
  highlight?: boolean;
  icon?: ReactNode;
  avatar?: string;
  badge?: ReactNode;
}

export const TeamMember: FC<TeamMemberProps> = ({
  name,
  role,
  description,
  highlight = false,
  icon,
  avatar,
  badge,
}) => (
  <Card variant={highlight ? 'highlight' : 'default'} className="text-left relative md:p-6 lg:p-7 w-full h-full">
    {badge && (
      <div className="absolute -top-1 -right-1 sm:top-2 sm:right-2">
        {badge}
      </div>
    )}
    <div className="flex items-start gap-2 sm:gap-3 md:gap-5">
      {avatar ? (
        <Image 
          src={avatar} 
          alt={name} 
          width={96} 
          height={96} 
          className="w-11 h-11 sm:w-16 sm:h-16 md:w-20 md:h-20 lg:w-[84px] lg:h-[84px] rounded-full object-cover object-[center_25%] shrink-0"
        />
      ) : (
        <div className={`
          p-1.5 sm:p-2 md:p-3 rounded-lg shrink-0
          ${highlight ? 'bg-white/20' : 'bg-white/10'}
        `}>
          {icon}
        </div>
      )}
      <div className="min-w-0">
        <h3 className="text-sm sm:text-base md:text-xl font-semibold text-white truncate">{name}</h3>
        <p className={`text-xs sm:text-sm md:text-base ${highlight ? 'text-white/70' : 'text-white/50'}`}>{role}</p>
        <div className="text-xs md:text-base text-white/40 mt-1 md:mt-2">{description}</div>
      </div>
    </div>
  </Card>
);

// ============================================================================
// CUSTOMER CARD
// ============================================================================

interface CustomerCardProps {
  name: string;
  metric: string;
  metricLabel: string;
  className?: string;
}

export const CustomerCard: FC<CustomerCardProps> = ({
  name,
  metric,
  metricLabel,
  className = '',
}) => (
  <Card className={`text-center ${className}`}>
    <div className="text-sm sm:text-base font-semibold text-white mb-1">{name}</div>
    <div className="text-xl sm:text-2xl font-bold text-white">{metric}</div>
    <div className="text-xs text-white/50">{metricLabel}</div>
  </Card>
);

// ============================================================================
// ROADMAP ITEM (brightness-based hierarchy)
// ============================================================================

interface RoadmapItemProps {
  phase: string;
  title: string;
  status: string;
  items: string[];
  color: 'emerald' | 'cyan' | 'violet' | 'pink';
  icon: ReactNode;
}

export const RoadmapItem: FC<RoadmapItemProps> = ({
  phase,
  title,
  status,
  items,
  icon,
}) => (
  <div className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-3 sm:p-4 md:p-5 text-left">
    <div className="flex items-center justify-between mb-1 sm:mb-2 md:mb-3">
      <div className="flex items-center gap-1 sm:gap-2">
        <span className="text-white/60 sm:scale-110 md:scale-125">{icon}</span>
        <span className="text-[10px] sm:text-xs md:text-sm font-medium px-1.5 sm:px-2.5 py-0.5 sm:py-1 rounded-full bg-[#2a2a2a] text-[#a0a0a0] border border-[#3a3a3a]">
          {status}
        </span>
      </div>
    </div>
    <div className="text-[10px] sm:text-xs md:text-sm text-white/50">{phase}</div>
    <h3 className="text-xs sm:text-lg md:text-xl font-semibold text-white mb-1 sm:mb-2 md:mb-3">{title}</h3>
    <ul className="space-y-0.5 sm:space-y-1.5 md:space-y-2">
      {items.slice(0, 2).map((item, i) => (
        <li key={i} className="text-[10px] sm:text-sm md:text-base text-white/60 flex items-start gap-1 sm:gap-1.5">
          <span className="mt-0.5 text-white/40">•</span>
          <span>{item}</span>
        </li>
      ))}
    </ul>
  </div>
);

// ============================================================================
// USE OF FUNDS ITEM
// ============================================================================

interface FundsItemProps {
  title: string;
  percentage: string;
  description: string;
  icon: ReactNode;
}

export const FundsItem: FC<FundsItemProps> = ({
  title,
  percentage,
  description,
  icon,
}) => (
  <Card className="text-left">
    <div className="flex items-center justify-between mb-2">
      <span className="text-white/60">{icon}</span>
      <span className="text-lg sm:text-xl font-bold text-white">{percentage}</span>
    </div>
    <h3 className="text-sm sm:text-base font-semibold text-white mb-1">{title}</h3>
    <p className="text-xs text-white/50 line-clamp-2">{description}</p>
  </Card>
);

// ============================================================================
// COMPARISON TABLE
// ============================================================================

interface ComparisonRowProps {
  feature: string;
  planton: ReactNode;
  competitor1: ReactNode;
  competitor2: ReactNode;
}

export const ComparisonRow: FC<ComparisonRowProps> = ({
  feature,
  planton,
  competitor1,
  competitor2,
}) => (
  <div className="grid grid-cols-4 gap-2 py-2 border-b border-[#2a2a2a] text-xs sm:text-sm">
    <div className="text-white/70 text-left">{feature}</div>
    <div className="text-center font-medium text-white">{planton}</div>
    <div className="text-center text-white/60">{competitor1}</div>
    <div className="text-center text-white/60">{competitor2}</div>
  </div>
);
