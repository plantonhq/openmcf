'use client';

import React, { useState } from 'react';
import { ChevronDown, ExternalLink, Link as LinkIcon, Check } from 'lucide-react';

/**
 * Client-safe date formatter
 */
function formatDate(dateString: string): string {
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  } catch {
    return dateString;
  }
}

/**
 * Format date as YYYY-MM-DD for display
 */
function formatDateCompact(dateString: string): string {
  try {
    const date = new Date(dateString);
    return date.toISOString().split('T')[0];
  } catch {
    return dateString;
  }
}

interface InvestorUpdate {
  slug: string;
  content: string;
  excerpt: string;
  frontmatter: {
    title: string;
    date: string;
    tags?: string[];
  };
}

interface InvestorUpdatesTimelineProps {
  updates: InvestorUpdate[];
}

/**
 * Simple Markdown renderer for investor updates
 * Handles basic markdown without importing heavy dependencies
 */
function SimpleMarkdownRenderer({ content }: { content: string }) {
  // Convert markdown to HTML (basic support)
  const html = content
    // Headers
    .replace(/^### (.*$)/gim, '<h3 class="text-lg font-semibold text-white mt-6 mb-2">$1</h3>')
    .replace(/^## (.*$)/gim, '<h2 class="text-xl font-bold text-white mt-8 mb-3">$1</h2>')
    .replace(/^# (.*$)/gim, '<h1 class="text-2xl font-bold text-white mt-8 mb-4">$1</h1>')
    // Bold
    .replace(/\*\*(.+?)\*\*/g, '<strong class="text-white font-semibold">$1</strong>')
    // Italic
    .replace(/\*(.+?)\*/g, '<em>$1</em>')
    // Links
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" class="text-white hover:underline" target="_blank" rel="noopener noreferrer">$1</a>')
    // Unordered lists
    .replace(/^\s*[-*]\s+(.*)$/gim, '<li class="ml-4 text-white/60">$1</li>')
    // Paragraphs (double newlines)
    .replace(/\n\n/g, '</p><p class="text-white/60 mb-4">')
    // Single newlines within paragraphs
    .replace(/\n/g, '<br />');

  // Wrap list items in ul tags (simple approach)
  const processedHtml = html.replace(/(<li[^>]*>.*?<\/li>\s*)+/g, '<ul class="list-disc list-inside mb-4 space-y-1">$&</ul>');

  return (
    <div 
      className="prose prose-invert prose-sm max-w-none"
      dangerouslySetInnerHTML={{ __html: `<p class="text-white/60 mb-4">${processedHtml}</p>` }}
    />
  );
}

export default function InvestorUpdatesTimeline({ updates }: InvestorUpdatesTimelineProps) {
  const [expandedSlug, setExpandedSlug] = useState<string | null>(null);
  const [copiedSlug, setCopiedSlug] = useState<string | null>(null);

  const toggleUpdate = (slug: string) => {
    setExpandedSlug(expandedSlug === slug ? null : slug);
  };

  const copyLink = async (slug: string, e: React.MouseEvent) => {
    e.stopPropagation();
    const url = `${window.location.origin}/legal/investor-updates/${slug}`;
    await navigator.clipboard.writeText(url);
    setCopiedSlug(slug);
    setTimeout(() => setCopiedSlug(null), 2000);
  };

  const openInNewTab = (slug: string, e: React.MouseEvent) => {
    e.stopPropagation();
    window.open(`/legal/investor-updates/${slug}`, '_blank');
  };

  return (
    <div className="space-y-1">
      {updates.map((update) => {
        const isExpanded = expandedSlug === update.slug;
        const isCopied = copiedSlug === update.slug;

        return (
          <article key={update.slug} className="relative">
            {/* Clickable header */}
            <div
              onClick={() => toggleUpdate(update.slug)}
              className="relative py-8 px-6 -mx-6 rounded-lg hover:bg-white/[0.02] transition-all duration-200 cursor-pointer group"
            >
              {/* Date and action icons row */}
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-3">
                  {/* Compact date badge */}
                  <span className="px-2 py-1 bg-white/5 border border-white/10 rounded text-xs font-mono text-white/60">
                    {formatDateCompact(update.frontmatter.date)}
                  </span>
                  <time
                    dateTime={update.frontmatter.date}
                    className="text-sm text-white/40"
                  >
                    {formatDate(update.frontmatter.date)}
                  </time>
                </div>

                {/* Action icons */}
                <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                  {/* Copy link */}
                  <button
                    onClick={(e) => copyLink(update.slug, e)}
                    className="p-1.5 rounded-md hover:bg-white/10 text-white/40 hover:text-white transition-colors"
                    title="Copy link"
                  >
                    {isCopied ? (
                      <Check className="w-4 h-4 text-emerald-400" />
                    ) : (
                      <LinkIcon className="w-4 h-4" />
                    )}
                  </button>

                  {/* Open in new tab */}
                  <button
                    onClick={(e) => openInNewTab(update.slug, e)}
                    className="p-1.5 rounded-md hover:bg-white/10 text-white/40 hover:text-white transition-colors"
                    title="Open in new tab"
                  >
                    <ExternalLink className="w-4 h-4" />
                  </button>

                  {/* Expand/collapse indicator */}
                  <div className={`p-1.5 text-white/40 transition-transform duration-200 ${isExpanded ? 'rotate-180' : ''}`}>
                    <ChevronDown className="w-4 h-4" />
                  </div>
                </div>
              </div>

              {/* Title */}
              <h2 className="text-xl font-semibold text-white group-hover:text-white/80 transition-colors mb-3">
                {update.frontmatter.title}
              </h2>

              {/* Excerpt (only when collapsed) */}
              {!isExpanded && (
                <p className="text-white/50 text-base leading-relaxed mb-4">
                  {update.excerpt}
                </p>
              )}

              {/* Tags */}
              {update.frontmatter.tags && update.frontmatter.tags.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {update.frontmatter.tags.map((tag) => (
                    <span
                      key={tag}
                      className="px-2 py-0.5 bg-white/5 text-white/40 border border-white/10 rounded text-xs"
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              )}
            </div>

            {/* Expanded content */}
            {isExpanded && (
              <div className="px-6 -mx-6 pb-8 pt-2 border-l-2 border-white/20 ml-0">
                <SimpleMarkdownRenderer content={update.content} />
              </div>
            )}
          </article>
        );
      })}
    </div>
  );
}
