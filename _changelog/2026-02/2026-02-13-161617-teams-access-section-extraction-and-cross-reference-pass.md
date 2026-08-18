# Teams & Access Section Extraction and Cross-Reference Pass

**Date**: February 13, 2026
**Type**: Feature
**Components**: Documentation

## Summary

Extracted teams-and-access and billing pages from the platform/ section into their own top-level teams-and-access/ section, with clean-slate rewrites of both pages. Followed by a comprehensive cross-reference pass across all 53 documentation pages to fix broken links, remove references to the deferred Security section, and correct terminology inconsistencies.

## Problem Statement / Motivation

The existing teams-and-access and billing pages had been sitting inside platform/ since the docs were first created, despite the information architecture calling for them to be their own section. Both pages also had significant quality issues inherited from the original documentation: marketing language, "Planton" terminology, fabricated YAML examples not verified against source code, author attribution in frontmatter, and unverified claims about pricing details and special programs.

Additionally, after 10 sessions of docs work creating new sections and rewriting pages, cross-references had accumulated some inconsistencies: links pointing to old paths, references to a Security section that was deferred, and minor terminology mismatches.

### Pain Points

- teams-and-access.md (430 lines) used "Planton", marketing tone, fabricated team permission patterns, unverified troubleshooting claims
- billing.md (455 lines) contained specific pricing that may be outdated, unverified special programs, and marketing language
- 3 pages still linked to `/docs/platform/teams-and-access` (old path)
- 2 pages linked to `/docs/security` which does not exist (Phase 7b deferred)
- cli.md used "InfraHub" instead of "Infra Hub" and "Planton" instead of "Planton"

## Solution / What's New

### Phase 8: Teams & Access Section

Two new pages created from scratch with source-verified content:

**teams-and-access/index.md** — Members, invitations (with verified lifecycle: pending/accepted/removed), teams (with nesting support: teams can contain identity accounts or other teams), roles and permissions (OpenFGA-based with scoped resource-kind roles), API keys, and a complete CLI reference table covering 10 commands.

**teams-and-access/billing.md** — Structural coverage of the three subscription plans (Free, Plus, Pro), automation minutes tracking, and Stripe-powered billing management. Deliberately omits specific pricing (changes frequently, belongs on a pricing page) while accurately documenting what the billing console shows and how subscription management works.

Both old pages deleted: `platform/teams-and-access.md` (430 lines) and `platform/billing.md` (455 lines).

### Phase 9: Cross-Reference Pass

- 3 links updated from `/docs/platform/teams-and-access` to `/docs/teams-and-access`
- Security section removed from root `index.md` and `platform/index.md` (Phase 7b deferred — no security docs exist yet)
- "InfraHub" corrected to "Infra Hub" in `cli.md`
- "Planton" corrected to "Planton" in `cli.md`
- All 53 pages verified: zero broken internal links remain

## Implementation Details

### Source Verification

All content in the new pages was verified against:

- **11 protobuf files**: identity account, user invitation, team, IAM role, IAM policy/permission, billing account, subscription plan enum
- **13 CLI command files**: iam invite, iam-policy add/get/remove, role list, lookup-invitations, remove-invitation, api-key list/new, get team
- **8 web console files**: settings layout (General/Members/Teams tabs), billing layout (Plans/Subscription tabs), invite-members modal, invitations list, teams tab

### Key Decisions

- Billing page uses structural-only approach: documents that 3 plans exist and how billing works, without specific prices that would require maintenance as pricing changes
- OpenFGA explicitly named in teams-and-access page for open-source credibility
- Teams nesting documented (teams can contain other teams as members) — verified in protobuf `TeamMember.member_type` constraint
- Security section references removed entirely rather than redirected to runner/security-model — keeps the navigation clean and avoids suggesting a section exists when it does not

## Benefits

- Docs now have 8 sections in the sidebar (up from 7), matching the information architecture
- 883 lines of unverified marketing content replaced with 255 lines of source-verified documentation
- Zero broken internal links across the entire documentation site
- Platform section reduced to its intended 5 pages (core concepts, not administrative functions)

## Impact

- Documentation site now has a clean `teams-and-access/` section at sidebar position 80 (after Runner)
- All cross-references across 53 pages resolve correctly
- Only Phase 7b (Security section) and Phase 10 (screenshot placeholders) remain for Week 1 structural work

## Related Work

- Phase 4 (Session 4): Connect section extraction — same pattern of extracting from platform/ to top-level section
- Phase 6b (Session 7): Runner IP sanitization — informed the decision to omit implementation details
- Phase 4.5/4.5b (Sessions 9-10): Retrospective audit — established the clean-slate principle used here

---

**Status**: Live
**Timeline**: Single session
