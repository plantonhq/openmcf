---
title: "Presets"
description: "Ready-to-deploy configuration presets for Compute Instance"
type: "preset-list"
componentSlug: "compute-instance"
componentTitle: "Compute Instance"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-dev-spot-vm"
    rank: "01"
    title: "Dev Spot VM"
    excerpt: "The cheapest way to get a full Linux machine on GCP: an `e2-medium` Spot VM booting the latest Debian 12, on the default network with an ephemeral external IP, deleted outright when GCP reclaims the..."
  - slug: "02-hardened-web-server"
    rank: "02"
    title: "Hardened Web Server"
    excerpt: "The production security posture for an internet-facing workload: a Shielded VM on a custom-mode subnetwork with no external IP, a dedicated least-privilege service account, OS Login instead of static..."
  - slug: "03-stateful-data-vm"
    rank: "03"
    title: "Stateful Data VM"
    excerpt: "A database-on-VM pattern where everything durable outlives the instance: data on a referenced `GcpComputeDisk`, a stable internal address from a referenced `GcpAddress`, deletion protection on the VM..."
---

# Compute Instance Presets

Ready-to-deploy configuration presets for Compute Instance. Each preset is a complete manifest you can copy, customize, and deploy.
