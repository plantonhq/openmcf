---
title: "Presets"
description: "Ready-to-deploy configuration presets for Ingress"
type: "preset-list"
componentSlug: "ingress"
componentTitle: "Ingress"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-single-host"
    rank: "01"
    title: "Single Host"
    excerpt: "This preset creates the simplest useful Ingress: one hostname routed to one Service backend on a prefix path covering everything (`/`). It is the standard \"expose this app on this domain\" shape."
  - slug: "02-tls-cert-manager"
    rank: "02"
    title: "TLS with cert-manager"
    excerpt: "This preset serves one host over HTTPS with a certificate issued automatically by cert-manager. The `tls` block names the certificate Secret; the `cert-manager.io/cluster-issuer` annotation tells..."
  - slug: "03-fanout-paths"
    rank: "03"
    title: "Path Fan-Out"
    excerpt: "This preset routes one hostname to multiple Service backends by URL path: `/` to the frontend, `/api` to the API. The classic single-domain, multi-service layout — one certificate, one DNS record,..."
  - slug: "04-default-backend"
    rank: "04"
    title: "Default Backend Only"
    excerpt: "This preset declares no host or path rules at all — just a `default_backend`. Every request the controller routes to this Ingress, regardless of hostname or path, goes to one Service. It is the..."
---

# Ingress Presets

Ready-to-deploy configuration presets for Ingress. Each preset is a complete manifest you can copy, customize, and deploy.
