---
title: "Presets"
description: "Ready-to-deploy configuration presets for Service Connection Policy"
type: "preset-list"
componentSlug: "service-connection-policy"
componentTitle: "Service Connection Policy"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-memorystore-valkey"
    rank: "01"
    title: "Memorystore for Valkey Policy"
    excerpt: "The policy that unlocks Memorystore for Valkey on a network: authorizes the `gcp-memorystore` service class in one region and gives the connectivity automation a subnet to place PSC endpoints in."
  - slug: "02-shared-vpc-guarded"
    rank: "02"
    title: "Shared VPC with Connection Cap"
    excerpt: "A guarded policy for Shared VPC topologies: the platform team owns the policy in the host project, labels it for cost attribution, and caps how many managed-service instances can attach before the..."
  - slug: "03-producer-allowlist"
    rank: "03"
    title: "Producer Hierarchy Allowlist"
    excerpt: "A security-hardened policy that restricts which resource-hierarchy locations producer instances may live in — only producers under the listed organizations, folders, or projects can connect into the..."
---

# Service Connection Policy Presets

Ready-to-deploy configuration presets for Service Connection Policy. Each preset is a complete manifest you can copy, customize, and deploy.
