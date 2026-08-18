# AWS CloudWatch Logs Delivery

Route AWS service logs where they belong: CloudFront access logs to S3 for Athena, Bedrock knowledge-base logs to a central archive, or the whole organization's log streams into one account's Kinesis pipeline.

## What Gets Managed

- **Vended pipeline**: the delivery source (which resource, which log type), owned destinations (S3, CloudWatch Logs, Firehose, X-Ray — with output format and cross-account policies), and the deliveries joining them (record fields, delimiter, Hive-compatible S3 layout).
- **Cross-account destination**: the legacy Kinesis-backed endpoint (with its access policy) that subscription filters in other accounts send to.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with CloudWatch Logs delivery permissions (plus IAM pass-role for the cross-account arm).

### AWS Prerequisites

- Vended arm: the producing resource must exist (its service must support vended log delivery — each service's docs list its log types), plus the receiving bucket/group/stream.
- Cross-account arm: a role trusting `logs.amazonaws.com` with write access to the Kinesis stream.

## After You Deploy

- Owned destination ARNs (keyed by name) are what other instances' deliveries reference; delivery IDs are the import handles.
- Vended delivery itself is free; the destination storage (S3/Firehose/CWL ingestion) bills normally.

## Common Changes

- Adding a second destination type for the same source: add a destinations entry + a deliveries entry — AWS allows one delivery per destination TYPE per source.
- Sharing a destination across teams: keep it in one owning instance and hand consumers its ARN output.
