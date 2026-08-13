# CPU Threshold

The canonical infrastructure alert: any Compute Engine instance whose
mean CPU stays above 80% for five sustained minutes opens a WARNING
incident.

## What it configures

- A threshold condition with `ALIGN_MEAN` over 60s windows — the
  aggregation floor that keeps single-point blips from paging.
- `duration: 300s` — the flap guard; `0s` would page on one sample.

## Adjust before deploying

- **filter** — scope to your fleet (add `metadata.user_labels` or zone
  filters) so someone else's dev VM never pages you.
- **notificationChannels** — wire your channels via valueFrom; a policy
  without channels opens silent incidents.
- **thresholdValue / duration** — 0.8 for 5m is a starting point, not a
  law; latency-sensitive fleets often want 0.6 for 10m.

## When to choose something else

For "the service is down" rather than "the host is busy", the **Uptime
Check Failure** preset watches what customers actually experience.
