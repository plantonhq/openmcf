# AWS Backup Restore Testing Plan

A backup you've never restored is a hope, not a backup. Restore
testing runs scheduled, automated restore drills against your recovery
points and reports whether they actually come back.

## What Gets Managed

- The testing plan: schedule, recovery-point selection (latest or
  random within a lookback window, per vault, per point type).
- Per-resource-type selections: the IAM role, coverage by ARN or tag,
  restore metadata overrides, and how long each restored copy stays up
  for validation.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AWS Backup permissions.

### AWS Account

- Recovery points to test (a vault filled by a backup plan —
  [AWS Backup Plan](/cloud-catalog/aws-backup-plan)).
- An IAM role trusting `backup.amazonaws.com` with restore permissions
  for the tested types ([AWS IAM Role](/cloud-catalog/aws-iam-role)).

## Deploy

### Console

Create the resource from the AWS catalog, schedule the drills, add a
selection per resource type, and deploy. Note the AWS name forbids
hyphens (`weekly_restore_tests`, not `weekly-restore-tests`).

### CLI

```bash
planton apply -f restore-testing-plan.yaml
```

## After Deploy

- Test results appear in the Backup console's restore testing view;
  restore-time metrics feed Audit Manager's restore-time controls.
- Each test creates (and deletes) a temporary restored resource — the
  drill bills as a real restore plus that resource's brief runtime.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
