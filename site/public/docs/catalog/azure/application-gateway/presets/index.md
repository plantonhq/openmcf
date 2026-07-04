---
title: "Presets"
description: "Ready-to-deploy configuration presets for Application Gateway"
type: "preset-list"
componentSlug: "application-gateway"
componentTitle: "Application Gateway"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-https"
    rank: "01"
    title: "Standard HTTPS Gateway with Key Vault TLS"
    excerpt: "This preset creates the production L7 baseline: a zone-redundant Standard_v2 gateway that terminates TLS with a Key Vault certificate, redirects all HTTP to HTTPS, autoscales 2-10 instances,..."
  - slug: "02-waf-path-routing"
    rank: "02"
    title: "WAF Gateway with Path-Based Routing"
    excerpt: "This preset creates the protected microservice front door: a WAF_v2 gateway enforcing a referenced Web Application Firewall policy, with path-based routing that sends `/api/*` to a dedicated backend..."
  - slug: "03-internal-gateway"
    rank: "03"
    title: "Internal (Private-Only) Gateway"
    excerpt: "This preset creates an east-west L7 gateway with no public exposure: a single private frontend pinned to a static address in the gateway subnet, fixed two-instance capacity, and connection draining..."
---

# Application Gateway Presets

Ready-to-deploy configuration presets for Application Gateway. Each preset is a complete manifest you can copy, customize, and deploy.
