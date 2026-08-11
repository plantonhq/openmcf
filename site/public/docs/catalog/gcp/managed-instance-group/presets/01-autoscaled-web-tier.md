---
title: "Autoscaled Web Tier"
description: "A regional, CPU-autoscaled serving fleet with zero-unavailability rolling updates — the canonical MIG shape behind a backend service and the HTTPS load-balancing family."
type: "preset"
rank: "01"
presetSlug: "01-autoscaled-web-tier"
componentSlug: "managed-instance-group"
componentTitle: "Managed Instance Group"
provider: "gcp"
icon: "package"
order: 1
---

# Autoscaled Web Tier

A regional, CPU-autoscaled serving fleet with zero-unavailability
rolling updates — the canonical MIG shape behind a backend service and
the HTTPS load-balancing family.

## What it configures

- A REGIONAL group (zone-outage resilient) of private e2-small VMs on
  your subnetwork, app started from the boot script.
- CPU-target autoscaling 2–10 with a scale-in cap (at most one instance
  removed per 10-minute window).
- Proactive `REPLACE` rollouts with a surge budget of 3 — template
  changes roll with zero unavailability.
- A named port (`http` → 8080) for the backend service's `portName`.

## Adjust before deploying

- **networkInterfaces** — point the subnetwork reference at your
  GcpSubnetwork (its region must match the group's).
- **startupScript / sourceImage** — your app start and your baked image;
  every change rotates the template and rolls the fleet.
- **cpuTarget / replicas** — 0.6 keeps headroom for zone loss; size
  maxReplicas against quota (rollout surge counts too).

## After deploying

Put the `instance_group` output behind a GcpBackendService backend and
the fleet serves the load-balancer family. Wire a conservative
GcpHealthCheck into `autoHealing` once the app exposes a health
endpoint.

## When to choose something else

For per-instance identity and preserved disks (brokers, databases),
start from the **Stateful Group** preset. For maximum availability with
application-level repairs, start from **Regional HA Group**.
