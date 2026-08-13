---
title: "Stateful Group"
description: "A fixed-size zonal fleet whose instances keep their names, data disks, and internal IPs through repairs and updates — brokers, quorum members, and databases-on-VM that peers address individually."
type: "preset"
rank: "03"
presetSlug: "03-stateful-group"
componentSlug: "managed-instance-group"
componentTitle: "Managed Instance Group"
provider: "gcp"
icon: "package"
order: 3
---

# Stateful Group

A fixed-size zonal fleet whose instances keep their names, data disks,
and internal IPs through repairs and updates — brokers, quorum members,
and databases-on-VM that peers address individually.

## What it configures

- A per-instance 100 GB data disk marked stateful (`deleteRule: NEVER` —
  the disk survives even permanent instance deletion) and stateful
  internal IPs on `nic0`.
- `RECREATE` replacement: instance NAMES are preserved through template
  updates, one instance at a time (`maxUnavailableFixed: 1`).
- `OPPORTUNISTIC` rollout — template changes wait for deliberate
  refreshes; repairs never apply a new template as a side effect.
- `PREVENT` teardown posture.

## Adjust before deploying

- **startupScript** — mount the `data` device and start the service;
  the disk arrives formatted-or-not exactly as the instance left it.
- **targetSize** — stateful groups scale by deliberate resize, not by
  autoscaler (the two models conflict on identity).
- Add **perInstanceConfigs** when individual instances need pinned
  metadata or pre-existing disks (the config name IS the instance name).

## After deploying

Rehearse the teardown path before trusting it with data: `deleteRule:
NEVER` keeps the disks when instances die permanently, and the group's
`deletionPolicy` governs the rest of the stack.

## When to choose something else

For stateless serving, start from **Autoscaled Web Tier** — stateful
machinery costs rollout speed and autoscaling, and a stateless tier
should pay neither.
