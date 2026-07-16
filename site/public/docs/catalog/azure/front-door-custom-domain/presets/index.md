---
title: "Presets"
description: "Ready-to-deploy configuration presets for Front Door Custom Domain"
type: "preset-list"
componentSlug: "front-door-custom-domain"
componentTitle: "Front Door Custom Domain"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-managed-certificate"
    rank: "01"
    title: "Managed-Certificate Domain"
    excerpt: "This preset creates a custom domain with Azure's managed certificate -- the zero-maintenance TLS posture: Azure issues, hosts, and auto-rotates a free DV certificate for the exact hostname."
  - slug: "02-wildcard-byo-certificate"
    rank: "02"
    title: "Wildcard Domain with Bring-Your-Own Certificate"
    excerpt: "This preset creates a wildcard custom domain served by a customer certificate wrapped in an AzureFrontDoorSecret -- the shape for multi-tenant platforms where every tenant gets a subdomain."
  - slug: "03-hardened-ciphers"
    rank: "03"
    title: "Hardened Cipher Policy"
    excerpt: "This preset creates a managed-certificate domain with a hand-picked cipher policy: GCM-only TLS 1.2 suites (dropping the CBC pair) plus the mandatory TLS 1.3 suites -- for compliance baselines that..."
---

# Front Door Custom Domain Presets

Ready-to-deploy configuration presets for Front Door Custom Domain. Each preset is a complete manifest you can copy, customize, and deploy.
