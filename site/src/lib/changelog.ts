import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';
import { CHANGELOG_DIRECTORY } from '@/lib/constants';
import { generateExcerptFromContent } from '@/lib/utils';

export type ChangelogCategory = 'feature' | 'improvement' | 'fix' | 'breaking';

/** All valid changelog categories in display order. */
export const CHANGELOG_CATEGORIES: ChangelogCategory[] = [
  'feature',
  'improvement',
  'fix',
  'breaking',
];

export interface ChangelogAuthor {
  name: string;
  title?: string;
}

export interface ChangelogEntry {
  slug: string;
  title: string;
  date: string;
  category: ChangelogCategory;
  tags: string[];
  excerpt: string;
  author: ChangelogAuthor[];
  content: string;
}

/**
 * Read raw markdown content for a single changelog entry by slug.
 * Returns an empty string if the file does not exist.
 */
export function getChangelogContentBySlug(slug: string): string {
  try {
    const fullPath = path.join(CHANGELOG_DIRECTORY, `${slug}.md`);

    if (!fs.existsSync(fullPath)) {
      return '';
    }

    return fs.readFileSync(fullPath, 'utf8');
  } catch (error) {
    console.error('Error reading changelog content for slug:', slug, error);
    return '';
  }
}

/**
 * Load and parse all changelog entries, sorted by date descending (newest first).
 *
 * Gracefully handles a missing or empty directory.
 */
export function getAllChangelogEntries(): ChangelogEntry[] {
  try {
    if (!fs.existsSync(CHANGELOG_DIRECTORY)) {
      return [];
    }

    const fileNames = fs.readdirSync(CHANGELOG_DIRECTORY);

    const entries = fileNames
      .filter((fileName) => fileName.endsWith('.md'))
      .map((fileName): ChangelogEntry | null => {
        const slug = fileName.replace(/\.md$/, '');
        const fullPath = path.join(CHANGELOG_DIRECTORY, fileName);
        const fileContents = fs.readFileSync(fullPath, 'utf8');
        const { data, content } = matter(fileContents);

        // Validate required fields
        if (!data.title || !data.date || !data.category) {
          console.warn(
            `Changelog entry ${fileName} is missing required frontmatter (title, date, category). Skipping.`,
          );
          return null;
        }

        // Validate category value
        if (!CHANGELOG_CATEGORIES.includes(data.category as ChangelogCategory)) {
          console.warn(
            `Changelog entry ${fileName} has invalid category "${data.category}". ` +
              `Must be one of: ${CHANGELOG_CATEGORIES.join(', ')}. Skipping.`,
          );
          return null;
        }

        const excerpt =
          data.excerpt || generateExcerptFromContent(content, 200);

        // gray-matter parses YAML dates as JS Date objects; normalise to YYYY-MM-DD.
        const dateValue =
          data.date instanceof Date
            ? data.date.toISOString().split('T')[0]
            : String(data.date);

        return {
          slug,
          title: data.title,
          date: dateValue,
          category: data.category as ChangelogCategory,
          tags: data.tags || [],
          excerpt,
          author: data.author || [],
          content,
        };
      })
      .filter((entry): entry is ChangelogEntry => entry !== null)
      .sort((a, b) => (a.date > b.date ? -1 : 1));

    return entries;
  } catch (error) {
    console.error('Error reading changelog entries:', error);
    return [];
  }
}

/**
 * Get the next changelog entry in date-descending order.
 *
 * Returns `null` when the current entry is the last (oldest) one.
 */
export function getNextChangelogEntry(
  currentSlug: string,
  allEntries: ChangelogEntry[],
): { title: string; excerpt: string; slug: string } | null {
  const sorted = [...allEntries].sort((a, b) =>
    a.date > b.date ? -1 : 1,
  );

  const currentIndex = sorted.findIndex((e) => e.slug === currentSlug);

  if (currentIndex === -1 || currentIndex === sorted.length - 1) {
    return null;
  }

  const next = sorted[currentIndex + 1];
  return {
    title: next.title,
    excerpt: next.excerpt,
    slug: next.slug,
  };
}
