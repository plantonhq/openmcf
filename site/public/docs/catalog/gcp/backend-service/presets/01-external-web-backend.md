---
title: "External Web Backend"
description: "The workhorse backend service: an instance-group pool behind the global external Application Load Balancer, health-checked, balanced on CPU utilization, draining gracefully on rollouts, with sampled..."
type: "preset"
rank: "01"
presetSlug: "01-external-web-backend"
componentSlug: "backend-service"
componentTitle: "Backend Service"
provider: "gcp"
icon: "package"
order: 1
---

# External Web Backend

The workhorse backend service: an instance-group pool behind the global external Application Load Balancer, health-checked, balanced on CPU utilization, draining gracefully on rollouts, with sampled request logging.

## When to Use

- The dynamic half of a public web application (static paths go to a backend bucket)
- Any VM-based service pool that needs a global anycast frontend

## Remix Notes

- `portName` must match a named port defined on every instance group in `backends` — that is how each group maps the logical port to its own number.
- Add a second backend with `capacityScaler: 0` to pre-stage a blue/green pool, then flip the scalers to shift traffic.
- Switch `balancingMode` to `RATE` with `maxRatePerInstance` when request cost is uneven and CPU is a lagging signal.
- Attach a `securityPolicy` reference (a GcpCloudArmorPolicy of type `CLOUD_ARMOR`) before exposing anything sensitive.
