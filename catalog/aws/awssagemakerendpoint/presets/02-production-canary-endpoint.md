# Production Canary Endpoint

This preset runs two weighted instance-backed variants — 90% of
traffic on the stable model, 10% on the candidate — with data capture
feeding Model Monitor and every capacity change rolled blue/green
through a canary batch guarded by an alarm.

## When to Use

- Production endpoints where a bad model version must back itself out
- A/B testing model versions on live traffic with a controlled split

## Before You Deploy

- Request the `ml.m5.large for endpoint usage` Service Quota first: it
  defaults to ZERO on fresh AWS accounts (as do nearly all instance
  endpoint quotas), so this preset fails with ResourceLimitExceeded
  until the increase is granted. Size it for rollouts, not steady
  state — blue/green transiently doubles the fleet, so two variants
  need a quota of 4.

## What You Get

- Two `ml.m5.large` variants splitting traffic 90/10 by weight, routed
  to free capacity (`LEAST_OUTSTANDING_REQUESTS`)
- 20% of payloads captured to S3 — the Model Monitor feed
- Configuration rolls that shift 20% of capacity first, bake for five
  minutes while the 5xx alarm is watched, and roll back automatically
  if it fires; the old fleet lingers five more minutes as a safety
  window

## Customize

- Promote the canary by shifting weights (0 keeps a variant deployed
  but takes no traffic — an instant rollback target); weight changes
  are themselves guarded configuration rolls
- Point `autoRollbackAlarmNames` at your real CloudWatch alarms (1–10)
  — the preset names a placeholder
- Add `managedInstanceScaling` per variant to let the endpoint itself
  scale between a floor and ceiling
