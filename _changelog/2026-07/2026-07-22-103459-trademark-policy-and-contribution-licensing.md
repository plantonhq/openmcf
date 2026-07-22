# Trademark Policy and Contribution Licensing Terms

**Date**: July 22, 2026
**Type**: Feature
**Components**: Legal/Licensing, Documentation

## Summary

The repository now documents its trademark boundary and its contribution
licensing terms. `TRADEMARKS.md` explains — generously — what Apache-2.0 §6
already says: the code is licensed, the identity is not. `CONTRIBUTING.md`
gained a licensing section restating the license's own inbound=outbound
clause (§5), and the README's Licensing section carries a one-sentence
pointer to the policy.

## Problem Statement / Motivation

Apache-2.0 deliberately withholds trademark rights (§6), but nothing in the
repository explained what that means in practice. That ambiguity cuts both
ways: adopters writing "built on Planton" had no assurance the reference was
welcome, and nothing told a fork-and-rebrand actor where the line actually
is. Similarly, contributors had no stated licensing terms for their
submissions beyond the license text itself.

## Solution / What's New

1. **`TRADEMARKS.md`** (new, repo root) — anchored in the license: it opens
   from §6 and explains the boundary permissions-first. Truthful references,
   teaching, forking, and redistribution need no permission, ever. What needs
   permission: using the name/logo in your own product, domain, or logo, or
   presenting a fork as official. Forks get their own section: the code is
   yours to fork, the identity is not. Contact: `trademarks@planton.dev`
   (following the repo's per-purpose contact-alias convention).
2. **`CONTRIBUTING.md` `## Licensing`** — inbound=outbound Apache-2.0,
   explicitly framed as what license §5 already provides; a right-to-submit
   affirmation; pointers to LICENSE, NOTICE, and TRADEMARKS.md.
3. **README** — one sentence appended to the existing Licensing section:
   the name and logo are trademarks of Planton Cloud, Inc.

Two deliberate design constraints:

- **No registration claims** — no ®, no "registered". The wording is
  registration-neutral, truthful today and needing no edit after a USPTO
  filing.
- **Domain-agnostic** — the policy covers the name and logo only and never
  enumerates domains, so it cannot go stale as web surfaces converge.

## Validation

- All relative links resolve (LICENSE, NOTICE, TRADEMARKS.md,
  CODE-OF-CONDUCT.md).
- README diff verified to be exactly the one added sentence.
- Entity name ("Planton Cloud, Inc.") verified consistent with the NOTICE
  and the site's terms content.
- License §5/§6 text read from the root LICENSE — both documents restate the
  license rather than adding terms.

## Impact

- **Adopters** get explicit, generous permission for every truthful use of
  the name — lower friction for writing, teaching, and building on the
  catalog.
- **The project** gets its documented identity boundary: the enforceable
  line against fork-and-rebrand, stated before it is ever needed.
- **Contributors** see the terms of submission in plain words at the point
  of contribution.

## Related Work

- The Apache NOTICE and its traveling attribution across release artifacts
  (CLI archives, terraform module zips) — the attribution counterpart of
  this identity boundary.
- A contributor license agreement with a repository-native signing check is
  planned follow-up work; the CONTRIBUTING licensing section is structured
  to take that paragraph without restructuring.

---

**Status**: ✅ Production Ready
