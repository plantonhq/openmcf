# AWS Cost Anomaly Monitor

Machine-learning spend surveillance: AWS learns your normal spend
pattern and flags anomalies — with alert subscriptions deciding who
hears, how often, and above what impact.

## What Gets Managed

- The monitor: by-service segmentation (the recommended posture) or a
  custom Cost Explorer expression selecting exactly the slice to
  watch (accounts, tags, cost categories).
- Alert subscriptions: immediate individual alerts to SNS, or
  daily/weekly email summaries, each with an optional impact
  threshold ("only alert above $100 of impact").

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Cost Explorer permissions.

### AWS Account

- Cost Explorer must be enabled (first visit to the console enables
  it; data appears within 24 hours).
- Anomaly detection is free.

## Deploy

### Console

Create the resource from the AWS catalog, pick the monitor shape, add
a subscription, and deploy.

### CLI

```bash
planton apply -f monitor.yaml
```

## After Deploy

- The monitor trains on ~10 days of history before flagging; alerts
  then flow per subscription.
- Outputs publish the monitor ARN and each subscription's ARN.
