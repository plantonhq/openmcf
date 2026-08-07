---
title: "Presets"
description: "Ready-to-deploy configuration presets for DNS Zone on AWS Route53"
type: "preset-list"
componentSlug: "dns-zone-on-aws-route53"
componentTitle: "DNS Zone on AWS Route53"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-public-zone"
    rank: "01"
    title: "Public DNS Zone"
    excerpt: "This preset creates a public Route53 hosted zone for managing DNS records that resolve globally on the internet. Public zones are the most common type and are used for any domain that needs to be..."
  - slug: "02-private-vpc-zone"
    rank: "02"
    title: "Private VPC DNS Zone"
    excerpt: "This preset creates a private Route53 hosted zone that resolves DNS queries only within associated VPCs. Private zones enable split-horizon DNS, where internal services use private domain names..."
---

# DNS Zone on AWS Route53 Presets

Ready-to-deploy configuration presets for DNS Zone on AWS Route53. Each preset is a complete manifest you can copy, customize, and deploy.
