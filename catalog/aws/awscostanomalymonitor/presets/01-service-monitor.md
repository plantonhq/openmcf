# Service Spend Monitor

This preset creates AWS's recommended baseline: one dimensional
monitor segmenting spend by SERVICE, with a daily email summary for
anomalies whose absolute impact reaches $100.

## Check first: your account probably already has this monitor

AWS permits exactly ONE services monitor per account, and it
auto-creates one (named `Default-Services-Monitor`) for every account
that enabled Cost Explorer on or after 2023-03-27. If yours exists
(`aws ce get-anomaly-monitors` shows it), creating this preset fails
with "Limit exceeded on dimensional spend monitor creation" — import
the existing monitor into this kind instead (the import map derives it
from the monitor ARN), or start from the custom-slice preset, which
has no singleton limit.

## When to Use

- Accounts WITHOUT the auto-created default monitor (Cost Explorer
  enabled before 2023-03-27, or the default was deleted)
- Accounts where per-service anomaly streams (an EC2 spike vs an S3
  spike) are the right granularity

## What You Get

- A free ML monitor that learns each service's normal spend and flags
  deviations (training takes ~10 days of history)
- A daily email digest, filtered to anomalies worth reading ($100+
  impact)

## Customize

- Tighten or loosen the threshold: compose
  `ANOMALY_TOTAL_IMPACT_ABSOLUTE` with
  `ANOMALY_TOTAL_IMPACT_PERCENTAGE` under `and` for "$100 AND 10%
  above normal"
- Switch to `frequency: IMMEDIATE` with SNS subscribers for real-time
  alerts into chat/ops tooling (the topic policy must allow
  costalerts.amazonaws.com to publish)
- Watch a slice instead of services: the custom-slice preset's
  `monitorSpecification`
