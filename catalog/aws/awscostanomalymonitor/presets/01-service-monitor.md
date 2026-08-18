# Service Spend Monitor

This preset creates AWS's recommended baseline: one dimensional
monitor segmenting spend by SERVICE, with a daily email summary for
anomalies whose absolute impact reaches $100.

## When to Use

- The first anomaly monitor every account should carry
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
