---
title: "Presets"
description: "Ready-to-deploy configuration presets for App Runner Auto Scaling Configuration"
type: "preset-list"
componentSlug: "app-runner-auto-scaling-configuration"
componentTitle: "App Runner Auto Scaling Configuration"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-latency-sensitive-api"
    rank: "01"
    title: "Latency-Sensitive API"
    excerpt: "A scaling posture for user-facing APIs where tail latency matters more than baseline cost: three warm instances absorb traffic without cold starts, and a lowered concurrency ceiling gives every..."
  - slug: "02-scale-conservative"
    rank: "02"
    title: "Cost-Conscious"
    excerpt: "A scaling posture for internal tools and low-traffic services: one warm instance, maximum request packing, and a tight scale-out cap so a traffic anomaly can never triple the bill."
  - slug: "03-org-default-policy"
    rank: "03"
    title: "Organization Default Scaling Policy"
    excerpt: "This preset registers a balanced auto scaling configuration and claims it as the ACCOUNT-WIDE default for its region: every App Runner service created afterwards without an explicit..."
---

# App Runner Auto Scaling Configuration Presets

Ready-to-deploy configuration presets for App Runner Auto Scaling Configuration. Each preset is a complete manifest you can copy, customize, and deploy.
