---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cert Manager DNS Authorization"
type: "preset-list"
componentSlug: "cert-manager-dns-authorization"
componentTitle: "Cert Manager DNS Authorization"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-standard-domain"
    rank: "01"
    title: "Standard Domain Authorization"
    excerpt: "This preset creates a global DNS authorization for one domain — the standard building block for issuing Google-managed certificates before traffic serves."
  - slug: "02-shared-per-project"
    rank: "02"
    title: "Shared Per-Project Authorization"
    excerpt: "This preset creates a DNS authorization whose validation record is scoped per (domain, project) — the shape that lets multiple teams and projects issue certificates for the same domain without..."
---

# Cert Manager DNS Authorization Presets

Ready-to-deploy configuration presets for Cert Manager DNS Authorization. Each preset is a complete manifest you can copy, customize, and deploy.
