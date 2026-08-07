---
title: "Presets"
description: "Ready-to-deploy configuration presets for DNS Record on Google Cloud"
type: "preset-list"
componentSlug: "dns-record-on-google-cloud"
componentTitle: "DNS Record on Google Cloud"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-a-record"
    rank: "01"
    title: "A Record"
    excerpt: "This preset creates a standard DNS A record pointing a domain name to an IPv4 address. This is the most common DNS record type, used for mapping hostnames to IP addresses."
  - slug: "02-cname-record"
    rank: "02"
    title: "CNAME Record"
    excerpt: "This preset creates a DNS CNAME record that aliases one hostname to another. CNAME records are used when you want a subdomain to resolve to the same address as another domain without hardcoding an IP..."
  - slug: "03-weighted-canary"
    rank: "03"
    title: "Weighted Canary"
    excerpt: "This preset creates an A record with a weighted round-robin routing policy — 95% of DNS answers point at the stable target, 5% at the canary. Shifting the weights progresses the rollout without..."
---

# DNS Record on Google Cloud Presets

Ready-to-deploy configuration presets for DNS Record on Google Cloud. Each preset is a complete manifest you can copy, customize, and deploy.
