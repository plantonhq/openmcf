---
title: "Presets"
description: "Ready-to-deploy configuration presets for HTTP API Domain"
type: "preset-list"
componentSlug: "http-api-domain"
componentTitle: "HTTP API Domain"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-single-api-domain"
    rank: "01"
    title: "Preset: Single API Domain"
    excerpt: "Use this preset to give one HTTP API a branded production URL: `https://api.example.com/` serving the API's `$default` stage at the domain root."
  - slug: "02-multi-api-mtls-domain"
    rank: "02"
    title: "Preset: Multi-API mTLS Domain"
    excerpt: "Use this preset for B2B / machine-to-machine API surfaces: several APIs composed under one hostname, with mutual TLS requiring partner clients to present certificates chaining to your CA truststore."
---

# HTTP API Domain Presets

Ready-to-deploy configuration presets for HTTP API Domain. Each preset is a complete manifest you can copy, customize, and deploy.
