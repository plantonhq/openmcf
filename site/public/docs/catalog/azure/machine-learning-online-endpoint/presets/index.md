---
title: "Presets"
description: "Ready-to-deploy configuration presets for Machine Learning Online Endpoint"
type: "preset-list"
componentSlug: "machine-learning-online-endpoint"
componentTitle: "Machine Learning Online Endpoint"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-key-auth-endpoint"
    rank: "01"
    title: "Key-Auth Endpoint"
    excerpt: "This preset creates the everyday online endpoint: static-key authentication, a system-assigned identity for pulling images and models, and all traffic routed to a first deployment named `blue`."
  - slug: "02-entra-auth-private-endpoint"
    rank: "02"
    title: "Entra-Auth Private Endpoint"
    excerpt: "This preset creates the hardened posture: Microsoft Entra token authentication (no scoring secret exists at all) with public network access disabled, for internal services scoring over private..."
---

# Machine Learning Online Endpoint Presets

Ready-to-deploy configuration presets for Machine Learning Online Endpoint. Each preset is a complete manifest you can copy, customize, and deploy.
