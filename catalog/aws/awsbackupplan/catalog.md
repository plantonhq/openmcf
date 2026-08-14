# AWS Backup Plan

Scheduled backups for your AWS resources: when they run, where
recovery points land, how long they live, and which resources are
covered — one plan, centrally managed.

## What Gets Managed

- Backup rules: cron schedules, backup windows, continuous
  (point-in-time) backups, lifecycle to cold storage and expiry, copy
  actions to other vaults/regions, GuardDuty malware scans, and
  air-gapped vault targeting.
- Resource selections: which resources each plan covers (by ARN, tag,
  or condition) and the IAM role AWS Backup assumes.
- Windows VSS application-consistent backups for EC2.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AWS Backup permissions.

### AWS Account

- A backup vault to target
  ([AWS Backup Vault](/cloud-catalog/aws-backup-vault)).
- For selections, an IAM role trusting `backup.amazonaws.com` (AWS's
  managed `AWSBackupServiceRolePolicyForBackup` policy is the usual
  grant) ([AWS IAM Role](/cloud-catalog/aws-iam-role)).

## Deploy

### Console

Create the resource from the AWS catalog, point each rule at a vault,
add a selection covering your resources, and deploy.

### CLI

```bash
planton apply -f backup-plan.yaml
```

## After Deploy

- Backup jobs appear in the AWS Backup console per rule schedule;
  recovery points land in the target vault.
- Lifecycle expiry is the cost lever — recovery-point storage is what
  AWS bills.
- Prove restores actually work with
  [AWS Backup Restore Testing Plan](/cloud-catalog/aws-backup-restore-testing-plan).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
