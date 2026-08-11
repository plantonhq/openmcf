---
title: "Block Destructive APIs"
description: "The org-wide invariant shape: project deletion is denied for everyone in the organization — no role grant overrides it — except on resources tagged as sandboxes, where experimentation stays cheap."
type: "preset"
rank: "02"
presetSlug: "02-block-destructive-apis"
componentSlug: "iam-deny-policy"
componentTitle: "IAM Deny Policy"
provider: "gcp"
icon: "package"
order: 2
---

# Block Destructive APIs

The org-wide invariant shape: project deletion is denied for everyone in
the organization — no role grant overrides it — except on resources
tagged as sandboxes, where experimentation stays cheap.

## What it configures

- `parent.organizationId` — the policy attaches at the organization and
  applies to every folder and project below it.
- One rule denying `cloudresourcemanager.googleapis.com/projects.delete`
  to `principalSet://goog/public:all`.
- A `denialCondition` on resource tags:
  `!resource.matchTag('12345678/env', 'sandbox')` — the denial applies
  everywhere EXCEPT sandbox-tagged resources. Environments go in the
  condition; people go in exception lists.
- `deletionPolicy: PREVENT` — an org guardrail must not vanish as a side
  effect of a teardown.

## Adjust before deploying

- **organizationId** — replace with the real numeric organization ID.
- **denialCondition.expression** — the tag key is namespaced by YOUR
  org's numeric ID (`{org-id}/env`); replace `12345678` and make sure
  the sandbox hierarchy actually carries the tag before relying on the
  exemption.
- **deniedPermissions** — extend with other destructive operations worth
  guarding (only permissions on Google's supported-permissions list can
  be denied).

## When to choose something else

If the goal is protecting a specific data path (e.g. break-glass
secrets) inside one project rather than an org-wide invariant, start
from the **Guard Secret Access** preset.
