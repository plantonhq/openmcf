# AWS Backup Framework

Continuous compliance auditing for your backups: controls that watch
whether resources are actually protected, retained long enough, and
encrypted — the evidence layer compliance reviews ask for.

## What Gets Managed

- The Audit Manager framework and its controls (coverage, retention,
  encryption, restore-time checks) with per-control parameters and
  scopes.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AWS Backup and Config permissions.

### AWS Account

- An ACTIVE AWS Config recorder in the region, recording the backup
  resource types
  ([AWS Config Recorder](/cloud-catalog/aws-config-recorder)) —
  without it, control evaluations cannot run and deployment lands
  FAILED.

## Deploy

### Console

Create the resource from the AWS catalog, pick controls and
parameters, and deploy. Note the AWS name forbids hyphens
(`backup_posture`, not `backup-posture`).

### CLI

```bash
planton apply -f backup-framework.yaml
```

## After Deploy

- Control evaluations appear in the Backup console's Audit Manager
  view as Config rules named after the framework.
- Wire the framework's ARN into a report plan for scheduled compliance
  evidence
  ([AWS Backup Report Plan](/cloud-catalog/aws-backup-report-plan)).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
