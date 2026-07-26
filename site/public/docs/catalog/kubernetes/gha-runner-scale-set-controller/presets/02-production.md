---
title: "Production controller preset"
description: "The controller hardened for a fleet the business depends on: a hot standby behind leader election, `eventual` update strategy (controller upgrades wait for running jobs instead of overprovisioning..."
type: "preset"
rank: "02"
presetSlug: "02-production"
componentSlug: "gha-runner-scale-set-controller"
componentTitle: "GHA Runner Scale Set Controller"
provider: "kubernetes"
icon: "package"
order: 2
---

# Production controller preset

The controller hardened for a fleet the business depends on: a hot
standby behind leader election, `eventual` update strategy (controller
upgrades wait for running jobs instead of overprovisioning runners),
structured logs at info, health probes, metrics on the controller AND
every listener it creates, and `system-cluster-critical` priority so
node pressure evicts workloads before it evicts the thing that
schedules your CI.

Change first: `runner_max_concurrent_reconciles` moves with fleet size
— it trades GitHub-API and API-server load for runner startup
throughput. Watch the listener metrics (queued vs started jobs) to
know when to raise it.

See [02-production.yaml](./02-production.yaml) for the manifest.
