---
title: "Presets"
description: "Ready-to-deploy configuration presets for SSL Certificate on Google Cloud"
type: "preset-list"
componentSlug: "ssl-certificate-on-google-cloud"
componentTitle: "SSL Certificate on Google Cloud"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-imported-cert"
    rank: "01"
    title: "Imported CA-Issued Certificate"
    excerpt: "The core self-managed pattern: upload a certificate chain issued by your own CA (or purchased commercially) with its private key, for a global external HTTPS load balancer."
  - slug: "02-regional-cert"
    rank: "02"
    title: "Regional ALB Certificate"
    excerpt: "A self-managed certificate scoped to one region — for regional external and internal Application Load Balancer proxies, which cannot reference global certificates."
  - slug: "03-rotation-versioned-name"
    rank: "03"
    title: "Rotation with a Versioned Name"
    excerpt: "Encode the issue year (or serial) into the GCP certificate name so rotations are explicit create-before-destroy steps instead of in-place mutations that cannot work."
---

# SSL Certificate on Google Cloud Presets

Ready-to-deploy configuration presets for SSL Certificate on Google Cloud. Each preset is a complete manifest you can copy, customize, and deploy.
