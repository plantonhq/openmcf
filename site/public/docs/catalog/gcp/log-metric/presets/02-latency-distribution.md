---
title: "Latency Distribution"
description: "Percentile latency charts from access logs alone: extract each request's latency into a histogram — no instrumentation change, no tracing dependency."
type: "preset"
rank: "02"
presetSlug: "02-latency-distribution"
componentSlug: "log-metric"
componentTitle: "Log Metric"
provider: "gcp"
icon: "package"
order: 2
---

# Latency Distribution

Percentile latency charts from access logs alone: extract each
request's latency into a histogram — no instrumentation change, no
tracing dependency.

## What it configures

- A DISTRIBUTION metric with `EXTRACT(jsonPayload.latency_ms)` pulling
  the number from each matching entry.
- Exponential buckets (64 buckets, √2 growth from 1ms) — sub-ms to
  ~90 minutes of range with constant relative precision, the right
  shape for latencies.

## Adjust before deploying

- **valueExtractor** — point it at YOUR latency field; use
  `REGEXP_EXTRACT(textPayload, "took (\\d+)ms")` when the value lives
  inside message text.
- **unit + scale** — keep them consistent: if the field is seconds, say
  `unit: s` and scale the buckets accordingly.
- **filter** — the `>0` guard drops entries where the field is missing
  or zero-valued so they never skew the histogram.

## After deploying

Chart p50/p95/p99 with percentile aligners on
`metric.type="logging.googleapis.com/user/latency-distribution"`, or
feed a GcpMonitoringSlo `distributionCut` SLI.

## When to choose something else

For counting events rather than measuring values, start from the
**Error Counter** preset.
