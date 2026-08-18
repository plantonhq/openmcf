'use client';

import React, { useState, useMemo, useEffect, useRef } from 'react';
import Link from 'next/link';
import { Author } from '@/lib/mdx';

// Custom type with only the properties actually used by the component
interface PostRecord {
  slug: string;
  title: string;
  date: string;
  author: Author[];
  content: string;
}

interface MdxRecordListProps {
  records: PostRecord[];
  currentSlug?: string;
  onSortChange: (sortOption: SortOption) => void;
  currentSort: SortOption;
  sectionTitle?: string;
  basePath?: string;
}

export type SortOption = 'date-desc' | 'date-asc' | 'title-asc' | 'title-desc';

const MdxRecordList: React.FC<MdxRecordListProps> = ({ 
  records, 
  currentSlug, 
  onSortChange, 
  currentSort, 
  sectionTitle = "Content",
  basePath = "/blog"
}) => {
  const [showSortMenu, setShowSortMenu] = useState(false);
  const sortMenuRef = useRef<HTMLDivElement>(null);
  const sortButtonRef = useRef<HTMLButtonElement>(null);

  // Close sort menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        showSortMenu &&
        sortMenuRef.current &&
        !sortMenuRef.current.contains(event.target as Node) &&
        sortButtonRef.current &&
        !sortButtonRef.current.contains(event.target as Node)
      ) {
        setShowSortMenu(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [showSortMenu]);

  const sortedRecords = useMemo(() => {
    const sorted = [...records];
    
    switch (currentSort) {
      case 'date-desc':
        return sorted.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
      case 'date-asc':
        return sorted.sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());
      case 'title-asc':
        return sorted.sort((a, b) => a.title.localeCompare(b.title));
      case 'title-desc':
        return sorted.sort((a, b) => b.title.localeCompare(a.title));
      default:
        return sorted;
    }
  }, [records, currentSort]);

  const getTimeAgo = (date: string) => {
    const now = new Date();
    const recordDate = new Date(date);
    const diffTime = Math.abs(now.getTime() - recordDate.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 1) return '1 day ago';
    if (diffDays < 7) return `${diffDays} days ago`;
    if (diffDays < 30) return `${Math.floor(diffDays / 7)} week${Math.floor(diffDays / 7) === 1 ? '' : 's'} ago`;
    if (diffDays < 365) return `${Math.floor(diffDays / 30)} month${Math.floor(diffDays / 30) === 1 ? '' : 's'} ago`;
    return `${Math.floor(diffDays / 365)} year${Math.floor(diffDays / 365) === 1 ? '' : 's'} ago`;
  };

  const getSortLabel = (option: SortOption) => {
    switch (option) {
      case 'date-desc': return 'Newest first';
      case 'date-asc': return 'Oldest first';
      case 'title-asc': return 'A-Z';
      case 'title-desc': return 'Z-A';
      default: return 'Sort by';
    }
  };

  return (
    <div className="w-80 border-r border-[#2a2a2a] h-full flex flex-col">
      {/* Content Section Header with Sort Controls - Fixed */}
      <div className="p-4 border-b border-[#2a2a2a] flex-shrink-0">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-[#b0b0b0] font-medium">{sectionTitle}</span>
          </div>
          
          {/* Sort Controls */}
          <div className="flex items-center gap-2 relative">
            <button
              onClick={() => setShowSortMenu(!showSortMenu)}
              className="p-1 text-[#a0a0a0] hover:text-white transition-colors"
              title="Sort options"
              ref={sortButtonRef}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
              </svg>
            </button>
            {/* Sort Menu */}
            {showSortMenu && (
              <div className="absolute top-0 right-4 mt-2 w-48 bg-[#0a0a0a] border border-[#3a3a3a] rounded-md shadow-lg z-10" ref={sortMenuRef}>
                {(['date-desc', 'date-asc', 'title-asc', 'title-desc'] as SortOption[]).map((option) => (
                  <button
                    key={option}
                    onClick={() => {
                      onSortChange(option);
                      setShowSortMenu(false);
                    }}
                    className={`w-full text-left px-4 py-2 text-sm hover:bg-white/5 transition-colors ${
                      currentSort === option ? 'text-white bg-white/10' : 'text-[#a0a0a0]'
                    }`}
                  >
                    {getSortLabel(option)}
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Content Records List - Scrollable */}
      <div className="flex-1 overflow-y-auto p-4">
        {sortedRecords.map((record) => (
          <Link
            key={record.slug}
            href={`${basePath}/${record.slug}`}
            className={`block py-3 px-2 rounded-md transition-colors ${
              currentSlug === record.slug
                ? 'bg-white/10 text-white'
                : 'hover:bg-white/5'
            }`}
          >
            <div className="flex items-start justify-between">
              <h3 className={`text-sm font-medium leading-tight ${
                currentSlug === record.slug ? 'text-white' : 'text-[#b0b0b0]'
              }`}>
                {record.title}
              </h3>
              <span className={`text-xs ml-2 flex-shrink-0 ${
                currentSlug === record.slug ? 'text-white' : 'text-[#a0a0a0]'
              }`}>
                {getTimeAgo(record.date)}
              </span>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
};

export { MdxRecordList };
