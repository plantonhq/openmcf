import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';

/**
 * Investor Updates Content Utilities
 * 
 * Adapted from donepudi.me's fastlane.ts but with key differences:
 * - Date prefix is KEPT in the URL (e.g., /legal/investor-updates/2026-02-01-first-update)
 * - This makes chronological sorting obvious in shared links
 * 
 * Content lives in public/investor-updates/ with files named YYYY-MM-DD-slug.md
 */

// ============================================================================
// TYPES
// ============================================================================

export interface InvestorUpdateFrontmatter {
  title: string;
  date: string;
  tags?: string[];
  excerpt?: string;
}

export interface InvestorUpdate {
  slug: string;
  content: string;
  excerpt: string;
  frontmatter: InvestorUpdateFrontmatter;
}

// ============================================================================
// CONSTANTS
// ============================================================================

const UPDATES_DIR = 'public/investor-updates';

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

/**
 * Generate excerpt from markdown content
 */
function generateExcerpt(content: string, maxLength: number = 200): string {
  // Remove markdown formatting
  const plainText = content
    .replace(/#{1,6}\s+/g, '') // Remove headers
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1') // Remove links but keep text
    .replace(/[*_~`]/g, '') // Remove formatting
    .replace(/\n+/g, ' ') // Replace newlines with spaces
    .trim();

  if (plainText.length <= maxLength) {
    return plainText;
  }

  // Truncate at last complete sentence or word
  const truncated = plainText.substring(0, maxLength);
  const lastPeriod = truncated.lastIndexOf('.');
  const lastSpace = truncated.lastIndexOf(' ');

  if (lastPeriod > maxLength * 0.7) {
    return truncated.substring(0, lastPeriod + 1);
  }

  if (lastSpace > 0) {
    return truncated.substring(0, lastSpace) + '...';
  }

  return truncated + '...';
}

/**
 * Generate slug from filename
 * Unlike fastlane, we KEEP the date prefix in the slug
 */
function generateSlugFromFilename(filename: string): string {
  // Remove .md extension
  return filename.replace(/\.md$/, '');
}

/**
 * Sort updates by date (newest first)
 */
function sortByDate(updates: InvestorUpdate[]): InvestorUpdate[] {
  return updates.sort((a, b) => {
    const dateA = new Date(a.frontmatter.date).getTime();
    const dateB = new Date(b.frontmatter.date).getTime();
    return dateB - dateA; // Newest first
  });
}

// ============================================================================
// MAIN FUNCTIONS
// ============================================================================

/**
 * Get all investor updates sorted by date (newest first)
 */
export async function getAllInvestorUpdates(): Promise<InvestorUpdate[]> {
  const updatesPath = path.join(process.cwd(), UPDATES_DIR);

  if (!fs.existsSync(updatesPath)) {
    return [];
  }

  const files = fs.readdirSync(updatesPath);
  const updates: InvestorUpdate[] = [];

  for (const file of files) {
    // Skip README and non-markdown files
    if (!file.endsWith('.md') || file === 'README.md') {
      continue;
    }

    const update = await getInvestorUpdate(generateSlugFromFilename(file));
    if (update) {
      updates.push(update);
    }
  }

  return sortByDate(updates);
}

/**
 * Get a single investor update by slug
 * The slug includes the date prefix (e.g., "2026-02-01-first-update")
 */
export async function getInvestorUpdate(slug: string): Promise<InvestorUpdate | null> {
  try {
    const filePath = path.join(process.cwd(), UPDATES_DIR, `${slug}.md`);

    if (!fs.existsSync(filePath)) {
      return null;
    }

    const fileContents = fs.readFileSync(filePath, 'utf8');
    const { data, content } = matter(fileContents);

    // Validate required fields
    if (!data.title) {
      console.warn(`Investor update ${slug} is missing required 'title' field`);
      return null;
    }

    if (!data.date) {
      console.warn(`Investor update ${slug} is missing required 'date' field`);
      return null;
    }

    // Use frontmatter excerpt if provided, otherwise generate
    const excerpt = data.excerpt || generateExcerpt(content, 200);

    return {
      slug,
      content,
      excerpt,
      frontmatter: data as InvestorUpdateFrontmatter,
    };
  } catch (error) {
    console.error(`Error reading investor update ${slug}:`, error);
    return null;
  }
}

/**
 * Get recent investor updates (for widgets or previews)
 */
export async function getRecentInvestorUpdates(limit: number = 5): Promise<InvestorUpdate[]> {
  const allUpdates = await getAllInvestorUpdates();
  return allUpdates.slice(0, limit);
}

/**
 * Get all unique tags from investor updates
 */
export async function getAllInvestorUpdateTags(): Promise<string[]> {
  const allUpdates = await getAllInvestorUpdates();
  const tagsSet = new Set<string>();

  allUpdates.forEach(update => {
    update.frontmatter.tags?.forEach(tag => tagsSet.add(tag));
  });

  return Array.from(tagsSet).sort();
}

/**
 * Format date for display
 */
export function formatUpdateDate(dateString: string): string {
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
 * Format date as ISO string for datetime attributes
 */
export function formatDateISO(dateString: string): string {
  try {
    const date = new Date(dateString);
    return date.toISOString().split('T')[0];
  } catch {
    return dateString;
  }
}
