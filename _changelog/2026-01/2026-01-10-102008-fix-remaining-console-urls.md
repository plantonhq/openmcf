# Fix Remaining Console URLs Across Website

**Date**: January 10, 2026
**Type**: Bug Fix
**Components**: Landing Page, Solutions Pages, Agents Page

## Summary

Fixed 10 additional console URLs that were missed in the previous URL fix (2026-01-09). Two files used the wrong domain (`planton.ai` instead of `planton.ai`), and eight files were missing the `/signup` path for CTA buttons. All signup/CTA buttons now correctly redirect to `https://planton.ai/signup`.

## Problem Statement / Motivation

Following the URL fix on January 9th that established the pattern of `/login` for login buttons and `/signup` for CTAs, a scan of the repository revealed additional files that were not updated:

### Pain Points

- **Wrong domain**: The agents page components (`hero-section.tsx`, `cta-section.tsx`) were using `planton.ai` instead of `planton.ai`
- **Missing path**: Multiple landing page components across v2 and v3 versions were pointing to `planton.ai/` without the `/signup` path
- **Inconsistent user journey**: Users clicking "Get Started", "Start Free Trial", or "Start with Agents" buttons were landing on the console home instead of the signup page

## Solution / What's New

Updated all CTA buttons to use the correct URL pattern:
- **Signup/CTA buttons** → `https://planton.ai/signup`

### Files Updated

**Issue 1: Wrong Domain (2 files)**

| File | Button Text | Change |
|------|-------------|--------|
| `src/components/agents/hero-section.tsx` | "Get Started" | `planton.ai` → `planton.ai/signup` |
| `src/components/agents/cta-section.tsx` | "Start with Agents" | `planton.ai` → `planton.ai/signup` |

**Issue 2: Missing `/signup` Path (8 files)**

| File | Button Text |
|------|-------------|
| `src/components/landing-page/v3-2026-01-02-1000/FinalCTA.tsx` | "Start Free Trial" |
| `src/components/landing-page/v3-2026-01-02-1000/ROICalculator.tsx` | "Start Free Trial" |
| `src/components/landing-page/v3-2026-01-02-1000/HeroSection.tsx` | "Start Free Trial" |
| `src/components/landing-page/v2-2025-12-31-0900/ROICalculator.tsx` | "Start Free Trial" |
| `src/components/landing-page/v2-2025-12-31-0900/HeroSection.tsx` | "Start Free Trial" |
| `src/components/landing-page/v2-2025-12-31-0900/HowItWorks.tsx` | "Start Your First Deployment" |
| `src/components/landing-page/v2-2025-12-31-0900/FinalCTA.tsx` | "Start Free Trial" |
| `src/app/(root)/solutions/by-role/devops/page.tsx` | "Get Started Free" |

## Implementation Details

Simple string replacements across all files:

```typescript
// Before (wrong domain)
href="https://planton.ai"

// After
href="https://planton.ai/signup"
```

```typescript
// Before (missing path)
href="https://planton.ai/"

// After
href="https://planton.ai/signup"
```

### Verification

After changes, grep verification confirmed:
- `planton.ai` → 0 matches in `/src`
- `planton.ai/"` → Only 1 match remaining (commented-out code in `header.tsx`)

## Benefits

- **Correct user journey**: All CTA buttons now lead directly to the signup page
- **Consistent branding**: All URLs use the correct `planton.ai` domain
- **Improved conversion**: Reduced friction by eliminating extra navigation steps

## Impact

**Users**: Clicking any "Get Started", "Start Free Trial", or similar CTA button now takes users directly to signup
**Marketing**: Improved conversion funnel with correct signup flow across all pages
**Agents page**: Now properly integrated with the console authentication flow

## Related Work

- Continues the URL standardization from [2026-01-09-134439-fix-login-signup-urls-and-hydration-errors.md](2026-01-09-134439-fix-login-signup-urls-and-hydration-errors.md)

---

**Status**: ✅ Live
**Files Changed**: 10 files
**Lines Modified**: ~10 line changes (simple URL replacements)
