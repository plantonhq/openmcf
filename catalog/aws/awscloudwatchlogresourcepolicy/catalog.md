# AWS CloudWatch Logs Resource Policy

The permission slip AWS services need before they can write logs into your account — Route53 query logs, EventBridge delivery, OpenSearch slow logs. One policy, and the service's logging feature turns on.

## What Gets Managed

- The resource policy: its scope (account-wide by name, or pinned to one log group's ARN) and its IAM document granting service principals `logs:CreateLogStream` / `logs:PutLogEvents` on the target groups.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with `logs:PutResourcePolicy` / `logs:DescribeResourcePolicies` / `logs:DeleteResourcePolicy`.

### Quota

AWS caps account-scope resource policies at 10 per region — prefer one shared policy per service class over many near-duplicates.

## After You Deploy

- The granted service can write immediately; `revision_id` in outputs is the concurrency token guarding future edits.

## Common Changes

- Document edits apply in place (revision-guarded); scope changes (name ↔ ARN) replace the policy.
