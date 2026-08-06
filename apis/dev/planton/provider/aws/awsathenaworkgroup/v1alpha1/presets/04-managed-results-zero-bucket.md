# Preset: Managed Results, Zero Bucket

An Athena workgroup that stores query results in AWS-managed storage — no S3
bucket to create, secure, or lifecycle. Results are retained for 24 hours and
retrieved through Athena APIs (GetQueryResults), which is exactly what BI
tools and programmatic callers use.

## What This Configures

- AWS-managed query result storage (encrypted with an AWS-owned key).
- 10 GB per-query scan cutoff for cost control.
- Workgroup configuration enforcement and CloudWatch metrics.

## When to Use

- BI dashboards and applications that read results through the Athena API,
  never from result files in S3.
- Teams that want analytics without owning a results bucket at all.

## Customization Points

- Add `kmsKey` under `managedQueryResults` when compliance requires a
  customer-managed encryption key.
- Switch to `resultConfiguration.outputLocation` when downstream systems must
  read result files directly from S3 (the two are mutually exclusive).
