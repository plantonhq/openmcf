# AWS CloudWatch Synthetics

Outside-in monitoring: a canary runs your script on a schedule — hitting the health endpoint, walking the checkout flow, screenshotting the page — and alarms fire before your users file tickets.

## What Gets Managed

- The canary: its S3-staged code bundle, runtime, schedule (rate or cron, with retries), per-run sizing (memory, timeout, ephemeral storage), VPC placement for private endpoints, artifact storage and encryption, and retention windows.
- Owned groups (console aggregation containers) and the canary's group memberships — joins are by group name, so canaries across many instances share groups cleanly.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Synthetics, Lambda, S3, and IAM pass-role permissions.

### AWS Prerequisites

- An execution role trusting `lambda.amazonaws.com` with the canary permissions AWS documents (artifact writes, logs, metrics) — reference an AwsIamRole.
- An S3 bucket for artifacts and one holding the code zip (often the same bucket) — reference an AwsS3Bucket.
- The code zip in the runtime's layout: `nodejs/node_modules/<file>.js` with handler `<file>.handler` for Node.js runtimes.

## After You Deploy

- The canary is READY; with `start_canary: true` it runs on schedule and each run writes artifacts (screenshots, HAR, logs) to S3.
- Runs bill per canary run plus the underlying Lambda execution — cost scales linearly with schedule frequency, so the run rate is the cost knob.

## Common Changes

- New script version: upload a new zip (or object version), update `code.s3_key`/`s3_version` — in place, the provider stops/starts the canary around it.
- Runtime upgrades (`runtime_version`) are in-place; AWS deprecates runtimes on a published schedule.
