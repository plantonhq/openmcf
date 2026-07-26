---
title: "GHA Runner Scale Set Controller"
description: "GHA Runner Scale Set Controller deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesgharunnerscalesetcontroller"
---

# GHA Runner Scale Set Controller

GitHub's official manager for self-hosted Actions runners on
Kubernetes: install it once, then declare runner fleets as resources —
the controller long-polls GitHub and turns queued jobs into ephemeral
runner pods that live exactly one job.

## Highlights

- **One controller, many fleets** — a single cluster-wide install
  serves every runner scale set on the cluster; a namespace fence
  exists for hard multi-tenancy.
- **The official line** — GitHub's supported scale-set architecture
  from the official OCI charts, chart and image in lockstep at a
  pinned version.
- **Production knobs, typed** — hot-standby replicas with automatic
  leader election, `eventual` upgrades that wait for running jobs,
  reconcile-throughput and API rate tuning, health probes, metrics on
  the controller and every listener.
- **No credentials here** — the controller holds no GitHub secret;
  each fleet brings its own, so blast radius stays per-registration.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
