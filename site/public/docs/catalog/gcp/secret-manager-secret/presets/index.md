---
title: "Presets"
description: "Ready-to-deploy configuration presets for Secret Manager Secret"
type: "preset-list"
componentSlug: "secret-manager-secret"
componentTitle: "Secret Manager Secret"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-app-secret-with-access"
    rank: "01"
    title: "App Secret with Access"
    excerpt: "The standard workload-credential shape: one manifest creates the secret, stores the payload as version 1, and grants the consuming service account read access — a READABLE secret with zero follow-up..."
  - slug: "02-regional-cmek-secret"
    rank: "02"
    title: "Regional CMEK Secret"
    excerpt: "The data-residency posture: a REGIONAL secret whose payloads never leave the region, encrypted with a customer-managed KMS key, with both destroy guards armed."
  - slug: "03-rotated-secret"
    rank: "03"
    title: "Rotated Secret"
    excerpt: "A credential with a rotation pipeline behind it: GCP publishes a reminder to Pub/Sub every 30 days, the subscriber rotates and adds a new version, and the `current` alias re-points consumers..."
---

# Secret Manager Secret Presets

Ready-to-deploy configuration presets for Secret Manager Secret. Each preset is a complete manifest you can copy, customize, and deploy.
