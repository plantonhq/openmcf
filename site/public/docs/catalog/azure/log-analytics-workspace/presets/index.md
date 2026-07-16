---
title: "Presets"
description: "Ready-to-deploy configuration presets for Log Analytics Workspace"
type: "preset-list"
componentSlug: "log-analytics-workspace"
componentTitle: "Log Analytics Workspace"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-pay-as-you-go"
    rank: "01"
    title: "Pay-As-You-Go Workspace"
    excerpt: "This preset creates the everyday Log Analytics Workspace: PerGB2018 (pay-as-you-go) billing, 90-day retention, and a daily ingestion cap as a cost guard. It is the right starting point for nearly..."
  - slug: "02-commitment-tier"
    rank: "02"
    title: "Commitment-Tier Workspace"
    excerpt: "This preset creates a Log Analytics Workspace on the CapacityReservation sku -- Azure's discounted commitment billing for estates with sustained high ingestion. The tier discount grows with the..."
  - slug: "03-private-hardened"
    rank: "03"
    title: "Private, Hardened Workspace"
    excerpt: "This preset creates a Log Analytics Workspace in the regulated-estate posture: no public ingestion or query paths, no shared-key authentication, centralized query permissions, and no post-retention..."
---

# Log Analytics Workspace Presets

Ready-to-deploy configuration presets for Log Analytics Workspace. Each preset is a complete manifest you can copy, customize, and deploy.
