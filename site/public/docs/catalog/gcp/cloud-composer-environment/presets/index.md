---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud Composer Environment"
type: "preset-list"
componentSlug: "cloud-composer-environment"
componentTitle: "Cloud Composer Environment"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-dev-small"
    rank: "01"
    title: "Dev Small"
    excerpt: "A minimal Cloud Composer environment for developing and testing Airflow DAGs. Small allocations, public endpoint, no private networking — the cheapest footprint that still runs real pipelines."
  - slug: "02-production-private"
    rank: "02"
    title: "Production Private"
    excerpt: "A production-grade Cloud Composer environment with private networking, multi-zone resilience, VPC-native ranges, control-plane allowlisting, and scaled workloads."
  - slug: "03-enterprise-encrypted"
    rank: "03"
    title: "Enterprise Encrypted"
    excerpt: "A large Cloud Composer environment with the full security and operations surface: CMEK encryption, a bring-your-own DAG bucket, private networking, web UI allowlisting, disaster-recovery snapshots,..."
---

# Cloud Composer Environment Presets

Ready-to-deploy configuration presets for Cloud Composer Environment. Each preset is a complete manifest you can copy, customize, and deploy.
