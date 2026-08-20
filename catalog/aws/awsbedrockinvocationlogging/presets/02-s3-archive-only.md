# S3 Archive Only

This preset archives text and image invocation logs to S3 without a
CloudWatch arm — the low-cost retention posture when nobody queries
logs interactively.

## When to Use

- Regions where invocation logs exist for retention/compliance, not
  live debugging
- Cost-sensitive workloads (S3 storage is cheaper than CloudWatch
  ingestion)

## What You Get

- Text and image payloads archived under the key prefix
- Video and embedding payloads excluded (explicit false on the
  presence-typed toggles)

## Customize

- The bucket policy must let `bedrock.amazonaws.com` put objects
  (with an `aws:SourceAccount` condition in production)
- Add the `cloudwatch` arm later without replacing anything — the
  configuration updates in place
