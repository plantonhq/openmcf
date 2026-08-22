// Client-safe types (no Node.js imports)

export interface Author {
  name: string;
  title: string;
  bio?: string;
  profilePicture?: string;
  twitter?: string;
  github?: string;
  linkedin?: string;
  website?: string;
}

export interface BlogPost {
  slug: string;
  title: string;
  date: string;
  author: Author[];
  featuredImage?: string;
  featuredImageType?: string;
  featuredVideo?: string;
  tags: string[];
  excerpt?: string;
  content: string;
}

export interface Tutorial {
  slug: string;
  title: string;
  date: string;
  author: Author[];
  featuredImage?: string;
  featuredImageType?: string;
  featuredVideo?: string;
  tags: string[];
  category?: string;
  excerpt?: string;
  content: string;
}

export type ChangelogCategory = 'feature' | 'improvement' | 'fix' | 'breaking';

export interface ChangelogEntry {
  slug: string;
  title: string;
  date: string;
  category: ChangelogCategory;
  tags: string[];
  excerpt?: string;
  author: Author[];
  content: string;
}
