---
title: "Presets"
description: "Ready-to-deploy configuration presets for Certificate Manager Certificate"
type: "preset-list"
componentSlug: "certificate-manager-certificate"
componentTitle: "Certificate Manager Certificate"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-managed-dns-auth"
    rank: "01"
    title: "Managed Certificate with DNS Authorization"
    excerpt: "This preset creates a Google-managed certificate validated through a referenced `GcpCertManagerDnsAuthorization`. DNS authorization lets the certificate reach ACTIVE before any traffic serves — the..."
  - slug: "02-wildcard-cert"
    rank: "02"
    title: "Wildcard Certificate"
    excerpt: "This preset creates a Google-managed certificate covering an apex domain and its wildcard under one certificate. Wildcards require DNS authorization — load-balancer authorization cannot validate them."
  - slug: "03-self-managed-pem"
    rank: "03"
    title: "Self-Managed (Uploaded) Certificate"
    excerpt: "This preset uploads a PEM certificate chain and its private key as a Certificate Manager certificate. Renewal before expiry is YOUR responsibility — Google serves the material but never rotates it."
---

# Certificate Manager Certificate Presets

Ready-to-deploy configuration presets for Certificate Manager Certificate. Each preset is a complete manifest you can copy, customize, and deploy.
