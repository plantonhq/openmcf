# Tutorial Detail Page Layout Redesign

**Date**: April 2, 2026
**Type**: Enhancement
**Components**: Tutorials, Blog System, Layout Components

## Summary

Replaced the three-column tutorial detail layout (left sidebar list + content + right TOC) with a cleaner two-column layout (content + right TOC) and a changelog-style back link. The left sidebar tutorials list was removed entirely to reduce noise while reading. Symmetric spacers on both sides give the content a centered, balanced feel.

## Problem Statement / Motivation

When reading a tutorial, the left sidebar listing all tutorials was visual noise. The reader has already chosen a tutorial -- showing the full list alongside the content competes for attention and adds no value to the reading experience.

### Pain Points

- Left sidebar occupied 320px of screen width with a list the reader doesn't need while reading
- Content area was squeezed between two sidebars
- No way to dismiss or collapse the list (previous iteration added a collapse button, but the list itself was still taking space even when collapsed)

## Solution / What's New

Removed the left sidebar entirely and replaced it with a simple back link (`← Tutorials` or `← Blog`) matching the existing changelog detail page pattern. Added symmetric `w-64` spacers on both sides to center the content + TOC visually.

### Layout Changes

**Before**: Left sidebar (w-80, tutorials list) | Content (flex-1) | Right sidebar (w-80, TOC/author)

**After**: Left spacer (w-64) | Content (flex-1) with back link | Right sidebar (w-80, TOC/author) | Right spacer (w-64)

### Back Navigation

The back link uses `ArrowLeft` from `lucide-react` (same icon the changelog detail page uses) and links to the `basePath` with `sectionTitle` as the label -- so tutorials show `← Tutorials` and blog posts show `← Blog`.

## Implementation Details

### Files Changed

**`src/components/common/content-layout.tsx`** -- Complete rewrite. Removed:
- `MdxRecordList` import and rendering (left sidebar)
- `SortContext` provider and `useSortContext` export
- Sidebar collapse/expand state, `localStorage` persistence, chevron icons
- `records` and `currentSlug` props

Added:
- `ArrowLeft` back link from `lucide-react`
- `w-64` spacers on both sides
- Simplified two-column layout

**`src/app/(root)/tutorials/[slug]/page.tsx`** -- Removed `records` and `currentSlug` props from `MdxContentLayout`.

**`src/app/(root)/blog/[slug]/page.tsx`** -- Same prop cleanup. Added explicit `sectionTitle="Blog"`.

**`src/components/tutorials/TutorialContent.tsx`** -- Removed `useNextArticle` hook dependency. Server-provided `nextArticle` passed directly to `MDXRenderer`.

**`src/components/blog/BlogPostContent.tsx`** -- Same simplification as TutorialContent.

**`src/hooks/useNextArticle.ts`** -- Simplified to a passthrough (no longer depends on sort context since there is no sidebar sort).

## Benefits

- Cleaner reading experience with no distracting list sidebar
- Content area gains ~320px of horizontal space
- Consistent navigation pattern with changelog pages (back link)
- Simplified component tree -- removed sort context, localStorage state, and collapse toggle
- Faster page loads -- no longer rendering 12+ tutorial links on every detail page

## Impact

- **Tutorial detail pages** (`/tutorials/[slug]`): Left sidebar removed, back link added
- **Blog detail pages** (`/blog/[slug]`): Same changes (shared layout component)
- **Tutorials index page** (`/tutorials`): Unchanged -- still shows the full list with category sidebar

## Related Work

- Tutorial rename and content refocus (same session) -- tutorials renamed to "How to..." verb-based titles
- Sitemap update with all 12 new tutorial URLs

---

**Status**: Live
**Timeline**: Same-day iteration through collapse button → hidden sidebar → final spacer tuning
