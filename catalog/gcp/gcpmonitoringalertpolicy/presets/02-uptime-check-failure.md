# Uptime Check Failure

The other half of every availability monitor: pages CRITICAL when a
GcpMonitoringUptimeCheck stops passing. This is Google's own recommended
shape for uptime alerting — count the FALSE probe results per host and
trigger when any exist.

## What it configures

- A threshold on `uptime_check_passed` with `REDUCE_COUNT_FALSE` over a
  20-minute alignment window, grouped by host — more than one failed
  probe result opens the incident.
- `severity: CRITICAL` — a customer-facing outage is the escalation
  tier.

## Adjust before deploying

- **check_id** — replace with the real check's `uptime_check_id`; in a
  chart, wire the whole filter through the check's output rather than
  hand-editing (recreated checks get new ids, and a stale id leaves this
  policy permanently green).
- **notificationChannels** — wire the escalation-tier channel
  (PagerDuty) here.

## When to choose something else

For host-level saturation signals, the **CPU Threshold** preset; for
error-driven paging on application logs, the **Error Log Match** preset.
