---
title: "Presets"
description: "Ready-to-deploy configuration presets for Container App Custom Domain"
type: "preset-list"
componentSlug: "container-app-custom-domain"
componentTitle: "Container App Custom Domain"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-managed-certificate-domain"
    rank: "01"
    title: "Managed-Certificate Domain"
    excerpt: "The standard custom-domain flow: bind the hostname certificate-less, then let a free Azure-managed certificate attach automatically once issued."
  - slug: "02-byo-certificate-domain"
    rank: "02"
    title: "BYO-Certificate Domain"
    excerpt: "A custom domain served with a certificate you brought: an EV/OV chain, an org-mandated CA, or a wildcard certificate stored on the app's environment."
---

# Container App Custom Domain Presets

Ready-to-deploy configuration presets for Container App Custom Domain. Each preset is a complete manifest you can copy, customize, and deploy.
