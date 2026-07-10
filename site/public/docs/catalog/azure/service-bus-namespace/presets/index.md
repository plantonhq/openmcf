---
title: "Presets"
description: "Ready-to-deploy configuration presets for Service Bus Namespace"
type: "preset-list"
componentSlug: "service-bus-namespace"
componentTitle: "Service Bus Namespace"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-namespace"
    rank: "01"
    title: "Standard Namespace"
    excerpt: "This preset creates a STANDARD-tier Service Bus namespace -- the full-featured multi-tenant tier (queues, topics, subscriptions, sessions, duplicate detection) that fits most production workloads."
  - slug: "02-premium-isolated"
    rank: "02"
    title: "Premium Isolated Namespace"
    excerpt: "This preset creates a PREMIUM namespace with the enterprise posture: dedicated messaging units, customer-managed-key encryption, and a deny-by-default VNet firewall admitting only the application..."
  - slug: "03-keyless-entra"
    rank: "03"
    title: "Keyless Namespace (Entra-Only)"
    excerpt: "This preset creates a namespace with SAS authentication disabled: no connection strings, no shared keys -- clients authenticate with Microsoft Entra identities holding Service Bus data-plane roles...."
---

# Service Bus Namespace Presets

Ready-to-deploy configuration presets for Service Bus Namespace. Each preset is a complete manifest you can copy, customize, and deploy.
