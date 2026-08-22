'use client';

import React from 'react';
import { MDXRenderer } from '@/lib/MDXRenderer';
import { MDXParserClient } from '@/lib/mdx-client';

interface NextArticle {
  title: string;
  excerpt?: string;
  slug: string;
}

interface BlogPostContentProps {
  slug: string;
  post: string;
  allPosts: any[];
  nextArticle?: NextArticle;
}

export function BlogPostContent({ slug, post, nextArticle }: BlogPostContentProps) {
  const mdxContent = MDXParserClient.reconstructMDX(post);

  return (
    <div className="p-8">
      <MDXRenderer 
        mdxContent={mdxContent} 
        markdownContent={post}
        nextArticle={nextArticle}
        path={slug}
      />
    </div>
  );
}
