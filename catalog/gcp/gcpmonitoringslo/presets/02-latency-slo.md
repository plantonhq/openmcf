# Latency SLO

95% of requests under 500ms, measured over a calendar month — the
"fast enough" objective, from the load balancer's own latency
distribution.

## What it configures

- A `distributionCut` on the external HTTPS LB's `total_latencies`
  DISTRIBUTION metric: values landing in [unbounded, 500] count as
  good. The unset `min` is deliberate — an unset bound is open, and 0
  vs unset are different statements.
- `calendarPeriod: MONTH` — matches how latency commitments are usually
  written into contracts (budget resets at month boundaries).

## Adjust before deploying

- **distributionFilter** — point it at YOUR latency distribution (Cloud
  Run: `run.googleapis.com/request_latencies`; custom histograms from a
  distribution-valued GcpLogMetric).
- **range.max** — the metric's unit is the resource's own (LB latencies
  are milliseconds). 500 means 500ms here, not seconds.
- **goal** — p95-under-bound expressed as a ratio; tighten toward 0.99
  only with measured headroom.

## When to choose something else

For "is it up" rather than "is it fast", start from the **Availability
SLO** preset.
