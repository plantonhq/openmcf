---
title: "Presets"
description: "Ready-to-deploy configuration presets for Dataproc Autoscaling Policy"
type: "preset-list"
componentSlug: "dataproc-autoscaling-policy"
componentTitle: "Dataproc Autoscaling Policy"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-balanced-autoscaling"
    rank: "01"
    title: "Balanced Autoscaling"
    excerpt: "A moderate, general-purpose autoscaling policy: half-strength scale factors, an even primary/secondary split, and a 30-minute graceful decommission window. A sensible default for mixed interactive..."
  - slug: "02-aggressive-batch"
    rank: "02"
    title: "Aggressive Batch"
    excerpt: "A maximum-speed scaling policy for retryable batch workloads: full scale factors in both directions, a spot-heavy 4:1 weight split, a scale-to-zero secondary group, and a short drain window."
  - slug: "03-conservative-production"
    rank: "03"
    title: "Conservative Production"
    excerpt: "A smooth, stability-first scaling policy for SLA-bound clusters: small scale factors, minimum-change fractions that filter out scaling noise, a long cooldown, and a 2-hour graceful decommission..."
---

# Dataproc Autoscaling Policy Presets

Ready-to-deploy configuration presets for Dataproc Autoscaling Policy. Each preset is a complete manifest you can copy, customize, and deploy.
