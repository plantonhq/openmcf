---
title: "Presets"
description: "Ready-to-deploy configuration presets for Container App Environment Managed Certificate"
type: "preset-list"
componentSlug: "container-app-environment-managed-certificate"
componentTitle: "Container App Environment Managed Certificate"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-cname-validated"
    rank: "01"
    title: "CNAME-Validated Certificate"
    excerpt: "The standard managed certificate for a subdomain: Azure validates ownership through the CNAME that already routes the hostname to the app, then issues and renews the certificate for free."
  - slug: "02-http-validated"
    rank: "02"
    title: "HTTP-Validated Certificate"
    excerpt: "A managed certificate for an apex domain (example.com itself): DNS forbids CNAME at the apex, so the domain routes by A record and Azure proves ownership by serving an HTTP token through it."
---

# Container App Environment Managed Certificate Presets

Ready-to-deploy configuration presets for Container App Environment Managed Certificate. Each preset is a complete manifest you can copy, customize, and deploy.
