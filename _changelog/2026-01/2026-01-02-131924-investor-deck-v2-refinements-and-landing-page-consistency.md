# Investor Deck V2 Refinements and Landing Page Consistency

**Date**: January 2, 2026
**Type**: Enhancement
**Components**: Investor Deck, Landing Page, Wall of Love, How It Works, Customers

## Summary

Refined the investor deck v2 with improved visual consistency, updated customer information, and better testimonial layouts. Applied consistent styling and vocabulary across both the investor deck and landing page, including icon updates, logo sizing adjustments, and role title corrections.

## Problem Statement / Motivation

After the initial investor deck v2 implementation, several refinements were needed:

- Customer logos had inconsistent sizing and unnecessary backgrounds
- The "How It Works" section on the landing page didn't match the investor deck styling
- Testimonial layouts on the Wall of Love slide didn't optimize screen space
- Customer names and roles needed corrections (Harsha's full name, Sai Saketh's role)
- TYNYBAY logo appeared too small due to internal padding in the PNG

### Pain Points

- Customer logos on landing page had distracting background boxes
- Jai.CX logo appeared too large relative to other logos
- TYNYBAY logo appeared too small due to internal padding
- Wall of Love testimonials were cramped in a single row on desktop
- "Junior DevOps" title was inappropriate for Sai Saketh
- Landing page "How It Works" section had outdated styling with step badges and backgrounds

## Solution / What's New

### Customer Logos Cleanup

**Landing Page (`SocialProofBar.tsx`)**:

- Removed background boxes and borders from customer logos
- Clean, borderless logos with hover opacity effect
- Increased gap spacing for better visual breathing room

**Investor Deck (`SlideCustomers.tsx`)**:

- Removed card backgrounds, showing only logos
- Per-logo scaling to balance visual weight:
  - Jai.CX: `scale-50` (logo is naturally larger)
  - TYNYBAY: `scale-[1.8]` (compensates for PNG padding)
  - Others: default scale

### How It Works Consistency

Updated landing page `HowItWorks.tsx` to match investor deck styling:

- Removed gradient backgrounds and borders from icon containers
- Added proper brand icons:
  - AWS icon for "Connect Cloud"
  - Layers icon for "Choose Infra"
  - Terraform icon for "Deploy"
  - GitHub icon (inline SVG) for "Connect Git"
  - Tekton icon for "Deploy Services"
- Removed "Step 1", "Step 2" badges
- Updated section labels: "Infrastructure" and "Services" (not "InfraHub: Deploy Infrastructure")
- Updated result callout with proper capitalization

### Wall of Love Layout Improvements

Restructured testimonial layout for better readability:

- **Desktop**: 2-2-1 grid layout
  - Row 1: 2 testimonials side by side
  - Row 2: 2 testimonials side by side
  - Row 3: 1 testimonial full width, center-aligned
- **Mobile**: Single column stack
- Full-width card shows quote on top, attribution below (centered)

### Customer Information Updates

- **Harsha CH → Harsha Gurram**: Corrected full name across all files
- **Junior DevOps → DevOps**: Updated Sai Saketh's role title across:
  - `SlideWallOfLove.tsx`
  - `WallOfLove.tsx`
  - `FinalCTA.tsx`
  - `CustomerStories.tsx`
  - `ROICalculator.tsx`

### Rohith's Testimonial Update

Changed Rohith Reddy Gopu's quote to focus on regulated industries:

- Old: "For a BFSI client on GCP, Planton delivered despite zero GCP experience on our team."
- New: "Planton helps us deliver compliant infrastructure for regulated industries."

## Implementation Details

### Logo Scaling Approach

Used Tailwind's `scale-*` classes for per-logo adjustments:

```tsx
const customers = [
  {
    name: 'Jai.CX',
    logo: '/images/customers/logos/jai-cx.svg',
    className: 'scale-50', // Logo is naturally larger
  },
  {
    name: 'TYNYBAY',
    logo: '/images/customers/logos/tynybay.png',
    className: 'scale-[1.8]', // Compensates for PNG padding
  },
  // ...
];
```

### Wall of Love 2-2-1 Layout

Separated testimonials into row arrays for explicit layout control:

```tsx
const row1 = [
  /* 2 testimonials */
];
const row2 = [
  /* 2 testimonials */
];
const row3 = [
  /* 1 testimonial - full width */
];
```

Full-width card uses different layout with centered content:

```tsx
if (fullWidth) {
  return (
    <div className="... flex flex-col items-center text-center">
      <p className="italic mb-3">{quote}</p>
      <div className="flex items-center gap-2">{/* avatar, name, role, logo */}</div>
    </div>
  );
}
```

### How It Works Icon System

Replaced generic icons with brand-specific images:

```tsx
const infraSteps = [
  {
    title: 'Connect Cloud',
    iconType: 'image',
    iconSrc: '/images/providers/aws.svg',
  },
  {
    title: 'Choose Infra',
    iconType: 'lucide',
    icon: Layers,
  },
  {
    title: 'Deploy',
    iconType: 'image',
    iconSrc: '/images/providers/terraform.svg',
  },
];
```

## Benefits

- **Visual Consistency**: Investor deck and landing page now share the same icon and styling language
- **Cleaner Aesthetics**: Removed unnecessary backgrounds and borders from logos
- **Better Readability**: Wall of Love testimonials are easier to read with 2-2-1 layout
- **Accurate Information**: Customer names and roles are now correct
- **Professional Appearance**: Balanced logo sizing creates more polished look

## Impact

### Files Modified

**Investor Deck (11 files)**:

- `InvestorDeckV2.tsx` - Added SlideWallOfLove to slide config
- `SlideCustomers.tsx` - Logo-only display with per-logo scaling
- `SlideWallOfLove.tsx` - New 2-2-1 layout with centered full-width card
- Various slide refinements

**Landing Page (8 files)**:

- `HowItWorks.tsx` - New icon system, removed step badges
- `SocialProofBar.tsx` - Removed logo backgrounds
- `WallOfLove.tsx` - Updated Harsha's name, Sai's role
- `FinalCTA.tsx` - Updated Sai's role
- `CustomerStories.tsx` - Updated Sai's role
- `ROICalculator.tsx` - Updated role label

### User Experience

- Investors see a more polished, professional deck
- Landing page visitors see consistent branding with the deck
- Testimonials are more readable and impactful

## Related Work

- Previous: `2026-01-02-105743-investor-deck-v2-mobile-optimized-redesign.md` - Initial v2 implementation
- Related: Landing page v2 marketing optimization work

---

**Status**: ✅ Live
**Timeline**: ~2 hours of iterative refinements
