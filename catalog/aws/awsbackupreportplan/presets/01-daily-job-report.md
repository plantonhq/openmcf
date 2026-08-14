# Daily Job Report

This preset creates the operational evidence stream: every backup
job's outcome, delivered daily to S3 as CSV and JSON.

## When to Use

- Feeding backup-success dashboards and alerting pipelines
- The evidence file audits ask for ("show me your backup job history")

## What You Get

- A daily report (AWS-managed cadence) of backup job outcomes under
  the `backup-reports/` prefix
- Both formats: CSV for humans, JSON for pipelines

## Customize

- Swap `reportTemplate` to `RESTORE_JOB_REPORT` or `COPY_JOB_REPORT`
  for the other job streams — note a template change REPLACES the
  report plan
- Widen coverage with `accounts: ["*"]` and a `regions` list for
  organization-wide evidence
- The bucket policy must allow the AWS Backup report service to write
  — the grant lives on the bucket, not here
