# Standalone Detector

This preset enables threat detection for one account+region: the
detector with the two highest-value protection plans, a
low-severity noise filter, and findings archived to S3.

## When to Use

- A standalone account (not under an organization's GuardDuty
  administration)
- The first region of a new security posture — findings within
  minutes, no agents required

## What You Get

- Foundational detection (CloudTrail, VPC flow, DNS) plus S3 data
  events and runtime monitoring with the EKS agent managed by AWS
- Findings below severity 4 auto-archived (the console shows what
  matters)
- Updated findings re-published every 15 minutes; everything exported
  to S3 for retention beyond GuardDuty's 90 days

## Customize

- The export bucket policy must grant `guardduty.amazonaws.com`
  PutObject, and the KMS key policy `kms:GenerateDataKey` — AWS
  rejects the destination otherwise
- Add `RDS_LOGIN_EVENTS` / `LAMBDA_NETWORK_LOGS` / `AI_PROTECTION` as
  the estate grows; agent sub-toggles install software — enable
  deliberately
- Drop `publishingDestination` to keep findings console-only

## Composing

Pair with AwsCloudTrail (GuardDuty consumes CloudTrail events
natively) and point security tooling at the export bucket.
