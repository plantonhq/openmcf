# Full Recording Posture

This preset records EVERY supported resource type — including global
ones (IAM) — with a daily-snapshot override taming the noisiest
compute types. The compliance-first posture; the bill follows the
account's activity.

## When to Use

- Audit regimes that require complete configuration history
- The ONE region designated to record global types

## What You Get

- Everything AWS Config supports, recorded continuously
- EC2 instances and ENIs as daily snapshots instead of per-change
  items (autoscaling churn is the classic bill amplifier)
- Full history in S3 with daily snapshots

## Customize

- Record global types in exactly ONE region — drop
  `includeGlobalResourceTypes` everywhere else or IAM items multiply
  per region
- Extend the override's `resourceTypes` with types that churn in your
  account
- Add `retentionPeriodInDays` to bound the queryable window

## Composing

Pair with organization-scoped AwsConfigRule instances for org-wide
compliance over the recorded surface.
