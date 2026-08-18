/**
 * Pre-build script: generates public/changelog/recent.json
 *
 * Reads all *.md files in public/changelog/, parses frontmatter, selects the
 * 10 most recent entries by date, and writes a JSON file consumed by the
 * console dashboard widget.
 *
 * Usage:  node scripts/generate-recent-changelog.mjs
 * Called: automatically before `next build` via the "build" script in package.json
 */

import fs from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';

// gray-matter is a CJS module; use createRequire to import it from ESM.
const require = createRequire(import.meta.url);
const matter = require('gray-matter');

const CHANGELOG_DIR = path.join(process.cwd(), 'public/changelog');
const OUTPUT_PATH = path.join(CHANGELOG_DIR, 'recent.json');
const MAX_ENTRIES = 10;
const SITE_ORIGIN = 'https://planton.ai';

function generateExcerpt(content, maxLength = 200) {
  return content
    .replace(/^---[\s\S]*?---/, '')
    .replace(/```[\s\S]*?```/g, '')
    .replace(/`([^`]+)`/g, '$1')
    .replace(/^#{1,6}\s+.*$/gm, '')
    .replace(/\*\*([^*]+)\*\*/g, '$1')
    .replace(/\*([^*]+)\*/g, '$1')
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
    .replace(/!\[([^\]]*)\]\([^)]+\)/g, '')
    .replace(/<[^>]*>/g, '')
    .replace(/^[-*_]{3,}$/gm, '')
    .replace(/^>\s+/gm, '')
    .replace(/^[-*+]\s+/gm, '')
    .replace(/^\d+\.\s+/gm, '')
    .replace(/\|.*\|/g, '')
    .replace(/\n\s*\n/g, '\n')
    .replace(/\s+/g, ' ')
    .trim()
    .slice(0, maxLength);
}

function main() {
  if (!fs.existsSync(CHANGELOG_DIR)) {
    console.log('[recent.json] No changelog directory found. Writing empty file.');
    fs.writeFileSync(OUTPUT_PATH, JSON.stringify({ entries: [], generatedAt: new Date().toISOString() }, null, 2));
    return;
  }

  const files = fs.readdirSync(CHANGELOG_DIR).filter((f) => f.endsWith('.md'));

  if (files.length === 0) {
    console.log('[recent.json] No changelog markdown files found. Writing empty file.');
    fs.writeFileSync(OUTPUT_PATH, JSON.stringify({ entries: [], generatedAt: new Date().toISOString() }, null, 2));
    return;
  }

  const entries = files
    .map((fileName) => {
      const slug = fileName.replace(/\.md$/, '');
      const raw = fs.readFileSync(path.join(CHANGELOG_DIR, fileName), 'utf8');
      const { data, content } = matter(raw);

      if (!data.title || !data.date || !data.category) {
        console.warn(`[recent.json] Skipping ${fileName}: missing required frontmatter.`);
        return null;
      }

      // gray-matter parses YAML dates as JS Date objects; normalise to YYYY-MM-DD.
      const dateValue = data.date instanceof Date
        ? data.date.toISOString().split('T')[0]
        : String(data.date);

      return {
        title: data.title,
        date: dateValue,
        category: data.category,
        slug,
        excerpt: data.excerpt || generateExcerpt(content),
        url: `${SITE_ORIGIN}/changelog/${slug}`,
      };
    })
    .filter(Boolean)
    .sort((a, b) => (a.date > b.date ? -1 : 1))
    .slice(0, MAX_ENTRIES);

  const output = {
    entries,
    generatedAt: new Date().toISOString(),
  };

  fs.writeFileSync(OUTPUT_PATH, JSON.stringify(output, null, 2));
  console.log(`[recent.json] Generated with ${entries.length} entries.`);
}

main();
