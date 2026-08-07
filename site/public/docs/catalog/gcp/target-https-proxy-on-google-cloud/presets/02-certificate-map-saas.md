---
title: "Certificate-Map SaaS Frontend"
description: "The custom-domains-at-scale shape: a Certificate Manager certificate map selects the served certificate by SNI hostname, lifting the 15-certificate list limit — one proxy can serve thousands of..."
type: "preset"
rank: "02"
presetSlug: "02-certificate-map-saas"
componentSlug: "target-https-proxy-on-google-cloud"
componentTitle: "Target HTTPS Proxy on Google Cloud"
provider: "gcp"
icon: "package"
order: 2
---

# Certificate-Map SaaS Frontend

The custom-domains-at-scale shape: a Certificate Manager certificate map selects the served certificate by SNI hostname, lifting the 15-certificate list limit — one proxy can serve thousands of customer domains. QUIC is forced on for lower-latency mobile clients.

## When to Use

- SaaS platforms serving customer-owned custom domains from one load balancer
- More than ~15 certificates, or certificate lifecycles managed centrally in Certificate Manager

## Remix Notes

- `certificateMap` is mutually exclusive with `sslCertificates` and `certificateManagerCertificates` — the map is the single source of certificates.
- Certificate maps only work with external ALBs (`EXTERNAL` / `EXTERNAL_MANAGED` forwarding-rule schemes).
- Adding a customer domain becomes a Certificate Manager operation (new map entry) with zero proxy changes — the proxy reference never needs to move.
