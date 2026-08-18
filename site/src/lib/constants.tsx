import path from 'path';

const CONTENT_DIRECTORIES = {
  BLOG: path.join(process.cwd(), 'public/blog'),
  BRANDING: path.join(process.cwd(), 'public/branding'),
  CHANGELOG: path.join(process.cwd(), 'public/changelog'),
  DOCS: path.join(process.cwd(), 'public/docs'),
  TUTORIALS: path.join(process.cwd(), 'public/tutorials'),
} as const;

export const BLOG_DIRECTORY = CONTENT_DIRECTORIES.BLOG;
export const BRANDING_DIRECTORY = CONTENT_DIRECTORIES.BRANDING;
export const CHANGELOG_DIRECTORY = CONTENT_DIRECTORIES.CHANGELOG;
export const DOCS_DIRECTORY = CONTENT_DIRECTORIES.DOCS;
export const TUTORIALS_DIRECTORY = CONTENT_DIRECTORIES.TUTORIALS;