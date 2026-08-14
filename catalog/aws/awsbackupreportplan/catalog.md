# AWS Backup Report Plan

Scheduled evidence for your backup posture: daily reports of backup,
copy, and restore jobs — or control-compliance reports over your audit
frameworks — delivered to S3 as CSV/JSON.

## What Gets Managed

- The report plan: template (job reports or compliance reports),
  framework references, account/OU/region coverage.
- The S3 delivery channel: bucket, prefix, formats.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AWS Backup permissions.

### AWS Account

- A destination S3 bucket whose policy allows the AWS Backup report
  service to write
  ([AWS S3 Bucket](/cloud-catalog/aws-s3-bucket)).
- For compliance templates, deployed frameworks
  ([AWS Backup Framework](/cloud-catalog/aws-backup-framework)).

## Deploy

### Console

Create the resource from the AWS catalog, pick the template, point it
at the bucket, and deploy. Note the AWS name forbids hyphens
(`daily_backup_jobs`, not `daily-backup-jobs`).

### CLI

```bash
planton apply -f backup-report-plan.yaml
```

## After Deploy

- Reports land in the bucket daily (AWS-managed cadence); run one
  on-demand with `aws backup start-report-job`.
- Changing the report template replaces the whole report plan —
  deliberate template choice up front saves churn.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
