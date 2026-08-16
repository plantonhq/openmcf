# AWS SSM Association

Keep your fleet in a known state: bind an SSM document — AWS-managed
or your own — to instances by tag, ID, or resource group, on a
schedule, with compliance reporting when a machine drifts.

## What Gets Managed

- The document binding: AWS-managed documents (patch scans, CloudWatch
  agent management) as plain names, or your own documents by
  reference, pinned to a version.
- Targeting and scheduling: tag/instance/resource-group targets, cron
  or rate expressions, and only-at-interval semantics.
- Operations posture: compliance severity, rate controls across large
  fleets, Change Calendar gating, and S3 output delivery.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SSM permissions.

### AWS Account

- Managed nodes (instances with the SSM agent and an instance profile)
  for the association to act on — the association itself deploys fine
  before they exist.
- Optionally a customer document
  ([AWS SSM Document](/cloud-catalog/aws-ssm-document)) and an output
  bucket ([AWS S3 Bucket](/cloud-catalog/aws-s3-bucket)).

## Deploy

### Console

Create the resource from the AWS catalog, name the document, target by
tag, set the schedule, and deploy.

### CLI

```bash
planton apply -f ssm-association.yaml
```

## After Deploy

- State Manager applies the association on schedule; per-target status
  and compliance land in the Systems Manager console.
- Adding matching instances later needs no change — tag-based targets
  pick them up on the next interval.
- Associations are free; the commands they run consume only what they
  touch.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
