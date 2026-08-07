---
title: "Managed-Certificate HTTPS Frontend"
description: "The standard production HTTPS frontend: a Google-managed certificate terminates TLS (no key material to handle), and an SSL policy replaces GCP's permissive default (TLS 1.0, COMPATIBLE profile) with..."
type: "preset"
rank: "01"
presetSlug: "01-managed-cert-frontend"
componentSlug: "target-https-proxy-on-google-cloud"
componentTitle: "Target HTTPS Proxy on Google Cloud"
provider: "gcp"
icon: "package"
order: 1
---

# Managed-Certificate HTTPS Frontend

The standard production HTTPS frontend: a Google-managed certificate terminates TLS (no key material to handle), and an SSL policy replaces GCP's permissive default (TLS 1.0, COMPATIBLE profile) with modern minimums.

## When to Use

- Any internet-facing application on a global external Application Load Balancer
- Up to 15 domains via one or more `GcpManagedSslCertificate` resources in `sslCertificates`

## Remix Notes

- Reference the `GcpManagedSslCertificate` and `GcpUrlMap` via `valueFrom` instead of literal self-links — the cert's `self_link` output exists immediately even while it is still PROVISIONING.
- Certificate rotation is zero-downtime: attach the replacement cert to the list first, wait for it to go ACTIVE, then remove the old one — the proxy swaps the list in place.
- Pair with the port-80 redirect frontend (`GcpTargetHttpProxy` preset 01) and two `GcpGlobalForwardingRule`s sharing one static `GcpGlobalAddress`.
