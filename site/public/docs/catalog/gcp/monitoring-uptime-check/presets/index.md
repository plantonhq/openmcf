---
title: "Presets"
description: "Ready-to-deploy configuration presets for Monitoring Uptime Check"
type: "preset-list"
componentSlug: "monitoring-uptime-check"
componentTitle: "Monitoring Uptime Check"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-public-https-check"
    rank: "01"
    title: "Public HTTPS Check"
    excerpt: "The canonical availability monitor: probe a public URL over TLS from all regions every five minutes, with certificate validation so an expired certificate fails the probe instead of passing silently."
  - slug: "02-tcp-port-check"
    rank: "02"
    title: "TCP Port Check"
    excerpt: "Asserts a non-HTTP service accepts connections — message brokers, databases with public endpoints, custom TCP protocols. The probe passes when the port completes a TCP handshake."
  - slug: "03-authenticated-api-check"
    rank: "03"
    title: "Authenticated API Check"
    excerpt: "Probes a health endpoint that sits behind HTTP basic auth, asserting both transport health (2xx over valid TLS) and body truth (the health JSON actually says ok)."
---

# Monitoring Uptime Check Presets

Ready-to-deploy configuration presets for Monitoring Uptime Check. Each preset is a complete manifest you can copy, customize, and deploy.
