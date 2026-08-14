---
title: "Presets"
description: "Ready-to-deploy configuration presets for Compute MIG"
type: "preset-list"
componentSlug: "compute-mig"
componentTitle: "Compute MIG"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-autoscaled-web-tier"
    rank: "01"
    title: "Autoscaled Web Tier"
    excerpt: "A regional, CPU-autoscaled serving fleet with zero-unavailability rolling updates — the canonical MIG shape behind a backend service and the HTTPS load-balancing family."
  - slug: "02-regional-ha-group"
    rank: "02"
    title: "Regional HA Group"
    excerpt: "A fixed-size, three-zone fleet with application-level auto-healing and hardened boot — the availability-first posture for serving tiers whose size is planned rather than demand-driven."
  - slug: "03-stateful-group"
    rank: "03"
    title: "Stateful Group"
    excerpt: "A fixed-size zonal fleet whose instances keep their names, data disks, and internal IPs through repairs and updates — brokers, quorum members, and databases-on-VM that peers address individually."
---

# Compute MIG Presets

Ready-to-deploy configuration presets for Compute MIG. Each preset is a complete manifest you can copy, customize, and deploy.
