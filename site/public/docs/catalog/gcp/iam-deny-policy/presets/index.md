---
title: "Presets"
description: "Ready-to-deploy configuration presets for IAM Deny Policy"
type: "preset-list"
componentSlug: "iam-deny-policy"
componentTitle: "IAM Deny Policy"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-guard-secret-access"
    rank: "01"
    title: "Guard Secret Access"
    excerpt: "The protect-break-glass-secrets shape: nobody in the project can read secret versions — not even project owners — except the one break-glass service account. Deny outranks every role grant, which is..."
  - slug: "02-block-destructive-apis"
    rank: "02"
    title: "Block Destructive APIs"
    excerpt: "The org-wide invariant shape: project deletion is denied for everyone in the organization — no role grant overrides it — except on resources tagged as sandboxes, where experimentation stays cheap."
---

# IAM Deny Policy Presets

Ready-to-deploy configuration presets for IAM Deny Policy. Each preset is a complete manifest you can copy, customize, and deploy.
