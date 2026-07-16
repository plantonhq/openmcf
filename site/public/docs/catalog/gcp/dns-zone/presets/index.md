---
title: "Presets"
description: "Ready-to-deploy configuration presets for DNS Zone"
type: "preset-list"
componentSlug: "dns-zone"
componentTitle: "DNS Zone"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-public-zone"
    rank: "01"
    title: "Public DNS Zone"
    excerpt: "Creates a public Cloud DNS managed zone for an internet-facing domain. DNS records are composed separately via [GcpDnsRecord](/docs/catalog/gcp/gcpdnsrecord)."
  - slug: "02-private-vpc"
    rank: "02"
    title: "Private VPC DNS Zone"
    excerpt: "Creates a private Cloud DNS managed zone visible only to resources on a VPC network. Use this for internal service discovery (`*.internal.example.com`) without exposing names to the public internet."
  - slug: "03-public-dnssec"
    rank: "03"
    title: "Public Zone with DNSSEC"
    excerpt: "Creates a production public Cloud DNS zone with DNSSEC signing and query logging enabled. DNSSEC applies to public zones only in Cloud DNS."
---

# DNS Zone Presets

Ready-to-deploy configuration presets for DNS Zone. Each preset is a complete manifest you can copy, customize, and deploy.
