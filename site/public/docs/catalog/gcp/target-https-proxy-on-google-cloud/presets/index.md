---
title: "Presets"
description: "Ready-to-deploy configuration presets for Target HTTPS Proxy on Google Cloud"
type: "preset-list"
componentSlug: "target-https-proxy-on-google-cloud"
componentTitle: "Target HTTPS Proxy on Google Cloud"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-managed-cert-frontend"
    rank: "01"
    title: "Managed-Certificate HTTPS Frontend"
    excerpt: "The standard production HTTPS frontend: a Google-managed certificate terminates TLS (no key material to handle), and an SSL policy replaces GCP's permissive default (TLS 1.0, COMPATIBLE profile) with..."
  - slug: "02-certificate-map-saas"
    rank: "02"
    title: "Certificate-Map SaaS Frontend"
    excerpt: "The custom-domains-at-scale shape: a Certificate Manager certificate map selects the served certificate by SNI hostname, lifting the 15-certificate list limit — one proxy can serve thousands of..."
  - slug: "03-mtls-server-tls-policy"
    rank: "03"
    title: "Mutual-TLS Frontend"
    excerpt: "An HTTPS frontend that authenticates CLIENTS, not just the server: a network security `ServerTlsPolicy` demands and validates client certificates during the handshake, on top of the normal server..."
---

# Target HTTPS Proxy on Google Cloud Presets

Ready-to-deploy configuration presets for Target HTTPS Proxy on Google Cloud. Each preset is a complete manifest you can copy, customize, and deploy.
