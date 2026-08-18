'use client';

import React from 'react';
import TableOfContents from '@/components/blog/TableOfContents';

const SITE_HEADER_HEIGHT = '70px';

interface BrandingContentLayoutProps {
  children: React.ReactNode;
  content: string;
}

export const BrandingContentLayout: React.FC<BrandingContentLayoutProps> = ({
  children,
  content,
}) => {
  return (
    <div className="min-h-screen font-inter antialiased">
      <div className="flex">
        {/* Main content area */}
        <div className="flex-1 min-h-screen min-w-0 overflow-x-hidden">
          <div className="px-4 sm:px-6 md:px-8 lg:px-12 py-4 md:py-8 w-full">
            {children}
          </div>
        </div>

        {/* Right sidebar — Table of Contents, sticks below the site header */}
        <div
          className="hidden xl:block sticky flex-shrink-0 w-80"
          style={{ top: SITE_HEADER_HEIGHT, height: `calc(100vh - ${SITE_HEADER_HEIGHT})` }}
        >
          <div className="h-full overflow-y-auto border-l border-[#2a2a2a]">
            <TableOfContents content={content} />
          </div>
        </div>
      </div>
    </div>
  );
};
