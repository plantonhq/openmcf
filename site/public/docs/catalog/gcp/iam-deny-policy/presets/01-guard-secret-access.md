---
title: "Guard Secret Access"
description: "The protect-break-glass-secrets shape: nobody in the project can read secret versions — not even project owners — except the one break-glass service account. Deny outranks every role grant, which is..."
type: "preset"
rank: "01"
presetSlug: "01-guard-secret-access"
componentSlug: "iam-deny-policy"
componentTitle: "IAM Deny Policy"
provider: "gcp"
icon: "package"
order: 1
---

# Guard Secret Access

The protect-break-glass-secrets shape: nobody in the project can read
secret versions — not even project owners — except the one break-glass
service account. Deny outranks every role grant, which is exactly why
this works where an allow-policy audit never quite does.

## What it configures

- One rule denying `secretmanager.googleapis.com/versions.access` to
  `principalSet://goog/public:all` — everyone.
- An `exceptionPrincipals` carve-out for the break-glass service
  account — the only identity that can still read the guarded secrets.
- `deletionPolicy: PREVENT` — destroying the guardrail must be a
  deliberate two-step, because its silent removal re-opens secret access
  for every role holder with no symptom.

## Adjust before deploying

- **exceptionPrincipals** — replace with the real break-glass account,
  and verify it works BEFORE applying: a wrong exception locks out the
  recovery path itself.
- **parent** — omitted here (provider default project); set
  `parent.projectId`, `parent.folderId`, or `parent.organizationId` to
  attach elsewhere. Folder/org attachment guards everything below.
- Consider scoping `deniedPermissions` if only some secret operations
  should be guarded — only permissions on Google's supported list can be
  denied.

## When to choose something else

If the goal is org-wide protection against destructive operations rather
than guarding a data path, start from the **Block Destructive APIs**
preset — it attaches at the organization and uses a tag condition to
exempt sandboxes.
