import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';
import { TUTORIALS_DIRECTORY } from '@/lib/constants';

export interface TutorialAuthor {
  name: string;
  title: string;
  bio?: string;
  profilePicture?: string;
  twitter?: string;
  github?: string;
  linkedin?: string;
  website?: string;
}

export interface Tutorial {
  slug: string;
  title: string;
  date: string;
  author: TutorialAuthor[];
  featuredImage?: string;
  featuredImageType?: string;
  featuredVideo?: string;
  tags: string[];
  category: string;
  excerpt?: string;
  content: string;
}

export function getTutorialContentBySlug(slug: string): string {
  try {
    const fullPath = path.join(TUTORIALS_DIRECTORY, `${slug}.md`);
    
    if (!fs.existsSync(fullPath)) {
      return '';
    }

    return fs.readFileSync(fullPath, 'utf8');
  } catch (error) {
    console.error('Error reading tutorial content for slug: ', slug, error);
    return '';
  }
}

export function getAllTutorials(): Tutorial[] {
  try {
    if (!fs.existsSync(TUTORIALS_DIRECTORY)) {
      return [];
    }

    const fileNames = fs.readdirSync(TUTORIALS_DIRECTORY);
    const allTutorialsData = fileNames
      .filter((fileName) => fileName.endsWith('.md'))
      .map((fileName) => {
        const slug = fileName.replace(/\.md$/, '');
        const fullPath = path.join(TUTORIALS_DIRECTORY, fileName);
        const fileContents = fs.readFileSync(fullPath, 'utf8');
        const { data, content } = matter(fileContents);
        return {
          slug,
          title: data.title,
          date: data.date,
          author: data.author || [],
          featuredImage: data.featuredImage,
          featuredImageType: data.featuredImageType,
          featuredVideo: data.featuredVideo,
          tags: data.tags || [],
          category: data.category || 'general',
          excerpt: data.excerpt,
          content,
        };
      })
      .filter((tutorial) => tutorial !== null);

    return allTutorialsData;
  } catch (error) {
    console.error('Error reading tutorials:', error);
    return [];
  }
}

export function getTutorialsByCategory(category: string): Tutorial[] {
  const allTutorials = getAllTutorials();
  return allTutorials.filter((tutorial) => tutorial.category === category);
}

export function getTutorialsByTag(tag: string): Tutorial[] {
  const allTutorials = getAllTutorials();
  return allTutorials.filter((tutorial) => tutorial.tags.includes(tag));
}

export function getTutorialsByAuthor(author: string): Tutorial[] {
  const allTutorials = getAllTutorials();
  return allTutorials.filter((tutorial) => tutorial.author.some(a => a.name === author));
}

export function getAllCategories(): string[] {
  const allTutorials = getAllTutorials();
  const categories = allTutorials.map(tutorial => tutorial.category);
  return [...new Set(categories)].sort();
}

export function getTutorialsBySearch(query: string): Tutorial[] {
  const allTutorials = getAllTutorials();
  const lowercaseQuery = query.toLowerCase();
  
  return allTutorials.filter(tutorial => 
    tutorial.title.toLowerCase().includes(lowercaseQuery) ||
    tutorial.excerpt?.toLowerCase().includes(lowercaseQuery) ||
    tutorial.tags.some(tag => tag.toLowerCase().includes(lowercaseQuery)) ||
    tutorial.category.toLowerCase().includes(lowercaseQuery)
  );
}

export function getTutorialsWithPagination(page: number = 1, limit: number = 6): { tutorials: Tutorial[]; total: number; hasMore: boolean } {
  const allTutorials = getAllTutorials();
  const startIndex = (page - 1) * limit;
  const endIndex = startIndex + limit;
  
  return {
    tutorials: allTutorials.slice(startIndex, endIndex),
    total: allTutorials.length,
    hasMore: endIndex < allTutorials.length
  };
}
