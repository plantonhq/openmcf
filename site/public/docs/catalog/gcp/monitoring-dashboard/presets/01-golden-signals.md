---
title: "Golden Signals"
description: "The SRE starter dashboard: traffic, errors, latency, and saturation on one page — the four questions every incident starts with, answered for an external HTTPS load balancer."
type: "preset"
rank: "01"
presetSlug: "01-golden-signals"
componentSlug: "monitoring-dashboard"
componentTitle: "Monitoring Dashboard"
provider: "gcp"
icon: "package"
order: 1
---

# Golden Signals

The SRE starter dashboard: traffic, errors, latency, and saturation on
one page — the four questions every incident starts with, answered for
an external HTTPS load balancer.

## What it configures

- A 2×2 mosaic layout: request rate, 5xx rate, p95 latency, and backend
  request volume, each aligned over 60s windows.
- Rates use `ALIGN_RATE` + `REDUCE_SUM` (fleet totals); latency uses
  `ALIGN_PERCENTILE_95` (tail truth, not averages).

## Adjust before deploying

- **filters** — these charts watch `https_lb_rule` resources (external
  ALB). Point them at your service's resource type (Cloud Run:
  `run.googleapis.com/request_count` on `cloud_run_revision`; GKE:
  Istio/gateway metrics) and scope by `resource.labels`.
- **displayName** — name it for the service, not the team; dashboards
  outlive re-orgs.

## When to choose something else

For host-level fleet health (CPU/memory/disk) rather than service
signals, start from the **Infrastructure Overview** preset.
