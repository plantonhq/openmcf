# AWS CloudWatch Logs Anomaly Detector

Machine-learned log monitoring: the detector studies a log group's normal patterns and flags what breaks them — a new exception class, a volume spike, a format change — without you writing a single filter pattern.

## What Gets Managed

- The detector: which log group it trains on, how often it evaluates (one minute to one hour), how long surfaced anomalies stay visible (7–90 days), an optional training filter pattern, and an optional customer-managed KMS key for its findings.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with CloudWatch Logs anomaly-detection permissions.

### AWS Prerequisites

- The log group to train on (reference an AwsCloudwatchLogGroup's ARN output).

## After You Deploy

- The model trains on live traffic — anomalies start surfacing within hours (up to ~24h for a stable baseline).
- Pricing is per GB of log data analyzed; quiet groups cost effectively nothing.

## Common Changes

- Pause/resume: flip `enabled` — the trained model survives a pause.
- Tune noise: raise `evaluation_frequency` or narrow training with `filter_pattern` (both in place).
