---
title: "Error Counter"
description: "Count a service's error log entries, labeled by status code — the canonical log-based metric and the input to every error-rate alert."
type: "preset"
rank: "01"
presetSlug: "01-error-counter"
componentSlug: "log-metric"
componentTitle: "Log Metric"
provider: "gcp"
icon: "package"
order: 1
---

# Error Counter

Count a service's error log entries, labeled by status code — the
canonical log-based metric and the input to every error-rate alert.

## What it configures

- A DELTA/INT64 counter over `severity>=ERROR` entries from one Cloud
  Run service.
- A `status` label extracted from `httpRequest.status` — bounded
  cardinality (~60 values), safe for the bill.

## Adjust before deploying

- **filter** — point `resource.type` and `service_name` at your
  workload (GKE: `k8s_container` + namespace/container labels; GCE:
  `gce_instance`).
- **labels** — add `method` if you alert per-endpoint; never extract
  unbounded values (user IDs, request IDs) into labels — every distinct
  combination is a billed time series.

## After deploying

Chart and alert on `metric.type="logging.googleapis.com/user/error-counter"`.
Pair with a GcpMonitoringAlertPolicy threshold condition to page on
error spikes.

## When to choose something else

For "how slow", not "how many" — percentiles from access logs — start
from the **Latency Distribution** preset.
