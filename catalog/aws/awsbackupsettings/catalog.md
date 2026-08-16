# AWS Backup Settings

The account-level switches behind AWS Backup: which resource types the
service protects in each region, which it fully manages, and whether
cross-account backup is enabled organization-wide.

## What Gets Managed

- Region settings: resource-type opt-in and management preferences
  (one instance per region).
- Global settings: cross-account backup (account-wide — set it in
  exactly one instance).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AWS Backup permissions.

### AWS Account

- Nothing else. Cross-account backup additionally assumes an AWS
  Organizations management context to be meaningful.

## Deploy

### Console

Create the resource from the AWS catalog, list every resource type you
intend to manage (AWS returns the full set on read — a type missing
from your map shows as a perpetual difference), and deploy.

### CLI

```bash
planton apply -f backup-settings.yaml
```

## After Deploy

- Backup plans only protect resource types the region has opted in —
  this component is the reason a selection silently skips a type.
- Destroy changes nothing at AWS (both arms are no-op deletes): to
  turn a setting off, apply it as false first.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
