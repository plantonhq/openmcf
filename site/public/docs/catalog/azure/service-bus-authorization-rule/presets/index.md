---
title: "Presets"
description: "Ready-to-deploy configuration presets for Service Bus Authorization Rule"
type: "preset-list"
componentSlug: "service-bus-authorization-rule"
componentTitle: "Service Bus Authorization Rule"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-queue-sender"
    rank: "01"
    title: "Queue Sender Credential"
    excerpt: "This preset mints a send-only SAS credential scoped to one queue -- the least-privilege shape for a producer service: it can send to its queue and do nothing else in the namespace."
  - slug: "02-namespace-operator"
    rank: "02"
    title: "Namespace Operator Credential"
    excerpt: "This preset mints a namespace-wide manage credential -- for platform tooling that creates and deletes entities at runtime (dynamic queue provisioning, tenant onboarding automation). A deliberate,..."
---

# Service Bus Authorization Rule Presets

Ready-to-deploy configuration presets for Service Bus Authorization Rule. Each preset is a complete manifest you can copy, customize, and deploy.
