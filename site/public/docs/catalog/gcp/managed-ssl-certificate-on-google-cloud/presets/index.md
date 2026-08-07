---
title: "Presets"
description: "Ready-to-deploy configuration presets for Managed SSL Certificate on Google Cloud"
type: "preset-list"
componentSlug: "managed-ssl-certificate-on-google-cloud"
componentTitle: "Managed SSL Certificate on Google Cloud"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-single-domain"
    rank: "01"
    title: "Single-Domain Load Balancer Certificate"
    excerpt: "The most common Google-managed SSL certificate: one fully-qualified domain name for an application served through a global external HTTPS load balancer."
  - slug: "02-multi-domain"
    rank: "02"
    title: "Multi-Domain Certificate (Apex + WWW)"
    excerpt: "One Google-managed SSL certificate covering multiple hostnames — typically the apex domain and its `www` alias, or several related subdomains on the same load balancer."
  - slug: "03-explicit-name"
    rank: "03"
    title: "Explicit GCP Certificate Name"
    excerpt: "Use when the Planton resource name (`metadata.name`) should differ from the certificate name that appears in GCP — common during rotation workflows where you create `prod-lb-tls-2026` while the..."
---

# Managed SSL Certificate on Google Cloud Presets

Ready-to-deploy configuration presets for Managed SSL Certificate on Google Cloud. Each preset is a complete manifest you can copy, customize, and deploy.
