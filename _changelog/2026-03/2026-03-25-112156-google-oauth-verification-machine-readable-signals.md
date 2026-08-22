# Machine-Readable Identity Signals for Google OAuth Verification

**Date**: March 25, 2026
**Type**: Enhancement
**Components**: SEO/Performance, Landing Page, Legal Pages, Content Management

## Summary

Added standard web infrastructure (`robots.txt`, `sitemap.xml`), machine-readable identity metadata (`application-name`, OpenGraph, JSON-LD structured data), and a homepage authentication disclosure to resolve persistent Google OAuth consent screen verification failures. Google's automated system was rejecting the submission despite correct human-readable content because it lacked the machine-readable signals it relies on for app identity matching and purpose verification.

## Problem Statement / Motivation

Google OAuth consent screen verification was rejected twice with two issues:

1. **"Your homepage does not explain the purpose of your app"** -- despite the homepage being extensively descriptive
2. **"The app name 'Planton' does not match the app name on your homepage"** -- despite "Planton" appearing in the H1 and `<title>`

Prior work (privacy/terms pages, site-wide rebrand) addressed the human-readable content. The root cause of the persistent failure was that Google's automated verification system relies on specific machine-readable signals that were completely absent from the site.

### Pain Points

- `https://planton.ai/robots.txt` returned 404 -- Google's bot may log a warning before crawling
- `https://planton.ai/sitemap.xml` returned 404 -- Google cannot discover `/privacy` and `/terms` programmatically
- No `<meta name="application-name">` tag -- the standard way to declare an app's name to user agents
- No JSON-LD structured data -- no machine-readable identity for Google's app name matcher
- No OpenGraph `og:site_name` -- another signal bots use for app name matching
- No `metadataBase` in Next.js metadata config -- metadata URLs may resolve incorrectly
- Homepage had zero mention of Google authentication despite the guideline requiring transparency about user data requests
- Footer had a "Complaince" typo visible during manual review

## Solution / What's New

### New Files

| File | Purpose |
|------|---------|
| `public/robots.txt` | Standard crawler allowance with sitemap reference |
| `public/sitemap.xml` | Static sitemap listing 95 public URLs with priority and changefreq |

### Modified Files

| File | Change |
|------|--------|
| `src/app/layout.tsx` | Added `metadataBase`, `applicationName`, `openGraph`, and JSON-LD structured data |
| `src/components/landing-page/v3-2026-01-02-1000/HeroSection.tsx` | Added authentication provider disclosure near CTA |
| `src/components/layout/footer.tsx` | Fixed "Complaince" typo to "Compliance" |

### Machine-Readable Metadata (layout.tsx)

The root layout now generates these `<head>` elements:

- `<meta name="application-name" content="Planton">` -- explicit app name for Google's matcher
- `<meta property="og:site_name" content="Planton">` -- OpenGraph app name signal
- `<meta property="og:url" content="https://planton.ai">` -- canonical identity
- `<meta property="og:type" content="website">` -- page type classification
- `<script type="application/ld+json">` -- WebApplication schema with `"name": "Planton"` and publisher `"Planton Cloud, Inc."`

The `metadataBase: new URL('https://planton.ai')` ensures all metadata URLs resolve correctly in the static export.

### JSON-LD Structured Data

A `WebApplication` schema provides the strongest machine-readable identity signal:

```json
{
  "@type": "WebApplication",
  "name": "Planton",
  "url": "https://planton.ai",
  "applicationCategory": "DeveloperApplication",
  "publisher": {
    "@type": "Organization",
    "name": "Planton Cloud, Inc."
  }
}
```

### Authentication Disclosure (HeroSection.tsx)

Added a small text line below the CTA's "100 automation minutes free" note:

> Sign in with Google, GitHub, or Microsoft. Privacy Policy

This directly addresses the Google guideline: *"Explain with transparency the purpose for which your app requests user data."* The link points to `/privacy`, matching the consent screen configuration.

### Sitemap Coverage

The static sitemap lists 95 URLs organized by priority:

- **Priority 1.0**: Homepage
- **Priority 0.8-0.9**: Core pages (privacy, terms, pricing, features, solutions, docs index)
- **Priority 0.6-0.7**: Section pages (docs sections, blog, changelog, tutorials, CLI, demo)
- **Priority 0.4-0.5**: Individual content pages (blog posts, changelog entries)

Excluded internal-facing pages: `/acme` (demo micro-app), `/invest` (investor content), `/meets` (meeting prep), `/legal/investor-updates`.

## Implementation Details

- `robots.txt` is a static file in `public/` -- copied verbatim to `out/` during static export
- `sitemap.xml` is a static file listing known routes -- appropriate for `output: 'export'` where the page set is deterministic at build time
- JSON-LD is injected via `dangerouslySetInnerHTML` with `JSON.stringify()` in the root layout `<head>` -- this ensures the structured data appears in every page's static HTML
- The `siteDescription` constant is shared between the Metadata export and the JSON-LD script to avoid description drift
- Next.js auto-generated Twitter Card metadata (`twitter:card`, `twitter:title`, `twitter:description`) from the OpenGraph config
- Build verified: `make build` passes cleanly (122 pages, 0 errors)
- All metadata confirmed present in `out/index.html` via grep

## Benefits

- Provides Google's verification system with unambiguous machine-readable app identity signals
- Adds standard web infrastructure (`robots.txt`, `sitemap.xml`) expected of any production website
- Improves SEO with OpenGraph metadata, structured data, and sitemap
- Homepage now transparently mentions authentication providers as required by Google's guidelines
- Fixes a visible typo in the footer that would look unprofessional during manual review

## Impact

- **Google OAuth**: Three layers of app name signals (application-name, og:site_name, JSON-LD name) should resolve the name-matching issue; authentication disclosure addresses the purpose-explanation issue
- **SEO**: Sitemap enables proper page discovery; OpenGraph metadata improves social sharing; JSON-LD enables rich results
- **Build**: Two new static files added to export; three modified source files; build passes cleanly

## Related Work

- [Privacy Policy and Terms of Service Pages](2026-03-25-091626-privacy-and-terms-pages.md) -- Created the pages that Google's verification requires
- [Rebrand and Google OAuth Fix](2026-03-25-095619-rebrand-planton-cloud-to-planton-and-google-oauth-fix.md) -- Addressed the human-readable signals; this changelog addresses the machine-readable signals
- [Azure Publisher Domain Verification](2026-03-25-101850-azure-publisher-domain-verification.md) -- Parallel OAuth verification work for Microsoft

---

**Status**: Live
**Timeline**: Single session
