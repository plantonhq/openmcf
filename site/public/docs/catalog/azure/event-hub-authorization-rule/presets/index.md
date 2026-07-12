---
title: "Presets"
description: "Ready-to-deploy configuration presets for Event Hub Authorization Rule"
type: "preset-list"
componentSlug: "event-hub-authorization-rule"
componentTitle: "Event Hub Authorization Rule"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-hub-producer"
    rank: "01"
    title: "Hub-Scoped Producer Credential"
    excerpt: "This preset mints a send-only SAS credential on ONE event hub -- the least-privilege connection string a producer service should hold."
  - slug: "02-namespace-operator"
    rank: "02"
    title: "Namespace-Scoped Operator Credential"
    excerpt: "This preset mints a full-rights (listen+send+manage) SAS credential over a whole namespace -- for operational tooling that creates and inspects hubs and consumer groups, distinct from the built-in..."
---

# Event Hub Authorization Rule Presets

Ready-to-deploy configuration presets for Event Hub Authorization Rule. Each preset is a complete manifest you can copy, customize, and deploy.
