# Header Luminance Hierarchy

**Date**: April 3, 2026
**Type**: Design
**Components**: Navigation, Design System

## Summary

Introduced a three-tier luminance hierarchy in the planton.ai website header. The logo and navigation items (Product, Solutions, Resources, Pricing) are now rendered at `#ededed` (`grey[100]`) instead of the body-text `#a0a0a0`, with hover states that brighten to `#fff`. Utility actions (Discord, Sign in) remain at `#a0a0a0`. This creates a clear visual hierarchy: CTA > navigation > utility.

## Problem Statement / Motivation

Every element in the header — logo, nav items, Discord, Sign in — was rendered at the same `#a0a0a0` body-text gray. This created a flat header with no visual differentiation between primary navigation and secondary utility actions. Only the Sign up CTA button stood out.

### Pain Points

- Navigation items were at body-text brightness, requiring unnecessary cognitive effort to scan
- The logo appeared timid at `#a0a0a0` — brand marks need confidence
- No luminance hierarchy between primary wayfinding and utility actions
- Did not match the pattern used by comparable devtools marketing sites (Linear, Vercel, Stripe, Cursor)

## Solution / What's New

Three-tier luminance hierarchy in the header:

| Tier | Elements | Color | Rationale |
|------|----------|-------|-----------|
| **Tier 1 — CTA** | Sign up button | `#fff` bg / `#000` text | Maximum attention. No change. |
| **Tier 2 — Navigation** | Logo, Product, Solutions, Resources, Pricing | `#ededed` (`grey[100]`) | Primary wayfinding. Hover brightens to `#fff`. |
| **Tier 3 — Utility** | Discord, Sign in | `#a0a0a0` | Secondary actions. No change. |

## Implementation Details

All changes in `packages/website-shell/src/components/`:

- **`WebsiteLogo.tsx`** — Changed `SvgIcon` color from `text.primary` to `grey.100`
- **`MegaMenu.tsx`** — Added `color: 'grey.100'` with `transition: 'color 150ms ease'` to trigger Typography; added `&:hover` to parent Stack so both label and chevron brighten to `#fff`
- **`DesktopNav.tsx`** — Added `color: 'grey.100'` with hover-to-white and transition to Pricing Typography

No new design tokens were created. `grey[100]` (`#ededed`) already exists in the website palette as the canonical marketing-page primary text brightness. The `#fff` hover matches the existing pattern in `MegaMenuItem.tsx` dropdown labels.

## Benefits

- Header navigation is immediately scannable — users find Product/Solutions/Resources/Pricing faster
- Logo has brand-appropriate presence at `#ededed`
- CTA (Sign up) retains the strongest visual pull
- Hero headline below still dominates — header supports without competing

## Impact

- **All marketing pages**: Every page with the header benefits from improved navigation scanability
- **Design system**: No new tokens, no new patterns — uses existing `grey[100]` and established hover-to-white convention

---

**Status**: Live
