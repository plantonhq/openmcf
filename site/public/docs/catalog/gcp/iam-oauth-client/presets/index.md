---
title: "Presets"
description: "Ready-to-deploy configuration presets for IAM OAuth Client"
type: "preset-list"
componentSlug: "iam-oauth-client"
componentTitle: "IAM OAuth Client"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-web-app-client"
    rank: "01"
    title: "Web App Client"
    excerpt: "The standard server-side shape: a confidential client with one managed credential whose secret GCP generates — the app authenticates users through the code flow and keeps sessions alive with refresh..."
  - slug: "02-spa-public-client"
    rank: "02"
    title: "SPA Public Client"
    excerpt: "The browser-app shape: a public client using the authorization code flow with PKCE. No credentials exist because none can — a secret shipped in a browser bundle is public by definition, and GCP..."
---

# IAM OAuth Client Presets

Ready-to-deploy configuration presets for IAM OAuth Client. Each preset is a complete manifest you can copy, customize, and deploy.
