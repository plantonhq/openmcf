---
title: "Presets"
description: "Ready-to-deploy configuration presets for GKE Cluster"
type: "preset-list"
componentSlug: "gke-cluster"
componentTitle: "GKE Cluster"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-private-standard"
    rank: "01"
    title: "Private Production Cluster — Standard"
    excerpt: "This preset creates a regional Standard cluster with private nodes, Dataplane V2, planned secondary ranges, a control-plane access allowlist, a daily maintenance window, and the Security Posture..."
  - slug: "02-autopilot"
    rank: "02"
    title: "Autopilot Cluster"
    excerpt: "This preset creates a private Autopilot cluster: GKE provisions and manages the nodes, bills per running pod, and enforces a hardened security posture out of the box."
  - slug: "03-dev-zonal"
    rank: "03"
    title: "Development Zonal Cluster"
    excerpt: "This preset creates the smallest, cheapest GKE control plane: a zonal Standard cluster with GKE-managed IP ranges, public nodes, the RAPID channel, and deletion protection off."
---

# GKE Cluster Presets

Ready-to-deploy configuration presets for GKE Cluster. Each preset is a complete manifest you can copy, customize, and deploy.
