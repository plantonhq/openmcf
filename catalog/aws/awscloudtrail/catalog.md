# AWS CloudTrail

The account's API audit log: a trail that records AWS API calls and
delivers them to S3, with optional CloudWatch Logs mirroring, SNS
delivery notices, Insights anomaly detection, and organization-wide
capture — the first checkbox in every security review.

## What Gets Managed

- The trail: multi-region capture, log-file validation, logging
  on/off.
- Delivery to S3 (optionally SSE-KMS encrypted) with per-file SNS
  notices.
- Event scope via classic or advanced selectors (one style per
  trail).
- CloudTrail Insights anomaly engines.
- Organization trails, with the optional delegated-administrator
  registration.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with CloudTrail permissions.

### AWS Account

- A delivery bucket carrying the `cloudtrail.amazonaws.com` bucket
  policy ([AWS S3 Bucket](/cloud-catalog/aws-s3-bucket) with
  `s3:GetBucketAcl` on the bucket and `s3:PutObject` under it) — AWS
  rejects trail creation without it.
- For CloudWatch mirroring: a log group
  ([AWS CloudWatch Log Group](/cloud-catalog/aws-cloudwatch-log-group))
  and a role trusting `cloudtrail.amazonaws.com` with write access to
  it ([AWS IAM Role](/cloud-catalog/aws-iam-role)).

## Deploy

### Console

Create the resource from the AWS catalog, pick the delivery bucket,
choose the event scope, and deploy.

### CLI

```bash
planton apply -f cloudtrail.yaml
```

## After Deploy

- Log files land under `s3://<bucket>/<prefix>/AWSLogs/<account-id>/`
  within minutes.
- With CloudWatch mirroring, query live events with Logs Insights.
- Destroying the component deletes the trail; delivered files stay in
  the bucket.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
