---
title: "Presets"
description: "Ready-to-deploy configuration presets for Application Security Group"
type: "preset-list"
componentSlug: "application-security-group"
componentTitle: "Application Security Group"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-workload-tier"
    rank: "01"
    title: "Workload Tier Group"
    excerpt: "This preset creates a single application security group named after a workload role (\"web-tier\"). It is the building block of address-independent micro-segmentation: instead of writing NSG rules..."
  - slug: "02-tagged-governance"
    rank: "02"
    title: "Governed Data-Tier Group"
    excerpt: "This preset creates a data-tier application security group carrying a full governance tag set -- cost center, owning team, and data classification. Tags are Azure's governance surface: Azure Policy..."
---

# Application Security Group Presets

Ready-to-deploy configuration presets for Application Security Group. Each preset is a complete manifest you can copy, customize, and deploy.
