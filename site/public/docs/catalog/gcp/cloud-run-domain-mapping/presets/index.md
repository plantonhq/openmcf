---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud Run Domain Mapping"
type: "preset-list"
componentSlug: "cloud-run-domain-mapping"
componentTitle: "Cloud Run Domain Mapping"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-custom-domain"
    rank: "01"
    title: "Custom Domain"
    excerpt: "The default shape: a verified domain mapped onto a Cloud Run service with a managed TLS certificate — no load balancer, no certificate handling, just DNS."
  - slug: "02-migration-no-cert"
    rank: "02"
    title: "Migration Without Certificate"
    excerpt: "The zero-downtime cutover shape: create the mapping with NO managed certificate, publish the DNS records while the old host still serves, then flip to `AUTOMATIC` once traffic has moved."
---

# Cloud Run Domain Mapping Presets

Ready-to-deploy configuration presets for Cloud Run Domain Mapping. Each preset is a complete manifest you can copy, customize, and deploy.
