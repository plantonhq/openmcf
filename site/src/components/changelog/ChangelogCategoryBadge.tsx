import React from 'react';
import type { ChangelogCategory } from '@/lib/types-client';

interface ChangelogCategoryBadgeProps {
  category: ChangelogCategory;
  className?: string;
}

const CATEGORY_CONFIG: Record<
  ChangelogCategory,
  { label: string; bg: string; text: string; border: string }
> = {
  feature: {
    label: 'Feature',
    bg: 'bg-white/10',
    text: 'text-white',
    border: 'border-white/30',
  },
  improvement: {
    label: 'Improvement',
    bg: 'bg-white/5',
    text: 'text-white/70',
    border: 'border-white/20',
  },
  fix: {
    label: 'Fix',
    bg: 'bg-white/5',
    text: 'text-white/60',
    border: 'border-white/20',
  },
  breaking: {
    label: 'Breaking',
    bg: 'bg-white/5',
    text: 'text-white/50',
    border: 'border-white/20',
  },
};

const ChangelogCategoryBadge: React.FC<ChangelogCategoryBadgeProps> = ({
  category,
  className = '',
}) => {
  const config = CATEGORY_CONFIG[category];

  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 text-xs font-medium rounded border ${config.bg} ${config.text} ${config.border} ${className}`}
    >
      {config.label}
    </span>
  );
};

export default ChangelogCategoryBadge;
