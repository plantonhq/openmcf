# Azure Publisher Domain Verification and Microsoft Sign-In Privacy Disclosures

**Date**: March 25, 2026
**Type**: Feature
**Components**: Legal Pages, Content Management, Domain Verification

## Summary

Added the Microsoft Entra publisher domain verification file at `/.well-known/microsoft-identity-association.json` and updated the Privacy Policy to disclose Microsoft as a sign-in provider. This unblocks Azure OAuth consent verification, the Azure counterpart to the Google OAuth consent work completed earlier today.

## Problem Statement / Motivation

Planton is adding Microsoft as a third-party sign-in provider alongside Google and GitHub. Microsoft Entra ID (formerly Azure AD) requires publisher domain verification before the OAuth consent screen can display a verified publisher. The verification process requires hosting a JSON file at a well-known URL on the publisher's domain.

Additionally, the Privacy Policy created earlier today only documented Google and GitHub sign-in data handling. With Microsoft sign-in being added, the policy needed to accurately reflect all three providers to maintain legal compliance and transparency.

### Pain Points

- Azure OAuth app registration showed "Unverified" publisher domain, which erodes user trust during consent
- Privacy Policy did not mention Microsoft as a sign-in provider, creating a legal gap
- Section 5 ("Google User Data") was provider-specific in a way that would not scale as new providers are added

## Solution / What's New

### New Files

| File | Purpose |
|------|---------|
| `public/.well-known/microsoft-identity-association.json` | Microsoft Entra publisher domain verification file associating application ID `596d6a3e-b228-4dc3-9a69-dc0971f84d12` with `planton.ai` |

### Modified Files

| File | Change |
|------|--------|
| `content/legal/privacy.md` | Added Microsoft as a sign-in provider in Section 1.A, 1.C, and Section 5; restructured Section 5 for multi-provider scalability |

### Privacy Policy Changes

Four targeted edits, all in `content/legal/privacy.md`:

1. **Section 1.A** — Sign-in provider list now reads "email/password, Google, GitHub, and Microsoft"
2. **Section 1.C intro** — Third-party provider enumeration updated to include Microsoft
3. **Section 1.C bullet** — New "Microsoft Sign-In" bullet documenting data received (name, email, profile picture) and exclusions (no Microsoft 365, OneDrive, Outlook access)
4. **Section 5** — Renamed from "Google User Data" to "Third-Party Sign-In Provider Data" with two subsections:
   - **Google** — Fully preserved, including Google API Services User Data Policy and Limited Use references (required by Google's verification)
   - **Microsoft** — Parallel structure with Microsoft-specific data handling commitments and a revocation link to `account.microsoft.com/privacy/app-access`

### Domain Verification File

The `microsoft-identity-association.json` file is a static JSON file placed in `public/.well-known/`. Since planton.ai uses Next.js `output: 'export'`, files in `public/` are copied verbatim to `out/` at build time. The existing `.nojekyll` file ensures dotfiles/directories are not stripped by the static hosting layer.

## Implementation Details

- The domain verification file follows Microsoft's exact specification: a single `associatedApplications` array containing the application ID of the "Planton OAuth" app registration in Microsoft Entra
- Section 5 restructuring uses `###` subsections under a `##` parent, keeping the document's numbered section structure intact (no renumbering of Sections 6-13)
- The Microsoft subsection in Section 5 mirrors the Google subsection's structure (data usage list, "we do not" list, revocation link) for consistency and legal parity
- Build verified: `make build` passes, `out/.well-known/microsoft-identity-association.json` is present in the static export output

## Benefits

- Unblocks Azure OAuth publisher domain verification in Microsoft Entra admin center
- Privacy Policy now accurately documents all three sign-in providers (Google, GitHub, Microsoft)
- Section 5's multi-provider structure is maintainable — future providers get a subsection without renumbering
- Consistent legal language across all provider subsections

## Impact

- **Azure OAuth**: Publisher domain can be verified by clicking "Verify and save domain" in the Entra admin center once the site is deployed
- **Users**: Privacy Policy transparently documents Microsoft sign-in data handling before the feature goes live
- **Build**: One new static file added to the export; build passes cleanly

## Related Work

- [Privacy Policy and Terms of Service Pages](2026-03-25-091626-privacy-and-terms-pages.md) — Created the privacy policy and terms pages earlier today
- [Rebrand and Google OAuth Fix](2026-03-25-095619-rebrand-planton-cloud-to-planton-and-google-oauth-fix.md) — Completed Google OAuth consent verification earlier today

---

**Status**: Live
**Timeline**: Single session
