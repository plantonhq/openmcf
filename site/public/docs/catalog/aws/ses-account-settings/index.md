---
title: "SES Account Settings"
description: "SES Account Settings deployment documentation"
icon: "package"
order: 100
componentName: "awssesaccountsettings"
---

# AWS SES Account Settings

The region's SES account-level settings — the suppression list and
Virtual Deliverability Manager (VDM) posture. A settings singleton:
one SES account object per region, at most one instance deployed per
region.

## What Gets Managed

- Which events (bounces, complaints) auto-suppress recipients
  account-wide.
- VDM: engagement dashboards and Guardian delivery optimization
  (VDM carries its own AWS pricing).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SES account permissions.

### AWS Account

- Nothing else — the settings target the account object itself.

## Deploy

### Console

Create the resource from the AWS catalog, pick the posture, and
deploy.

### CLI

```bash
planton apply -f ses-account-settings.yaml
```

## After Deploy

- Suppression applies to every send from the account in the region —
  including sends from configuration sets and identities managed
  elsewhere.
- Destroy resets VDM to disabled; suppression settings PERSIST (apply
  an empty reasons list first if you mean to stop suppressing).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
