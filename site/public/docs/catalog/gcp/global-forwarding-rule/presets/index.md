---
title: "Presets"
description: "Ready-to-deploy configuration presets for Global Forwarding Rule"
type: "preset-list"
componentSlug: "global-forwarding-rule"
componentTitle: "Global Forwarding Rule"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-https-frontend"
    rank: "01"
    title: "HTTPS Frontend VIP"
    excerpt: "The serving half of a production frontend: a reserved static IP bound on port 443 to a target HTTPS proxy, on the envoy-based `EXTERNAL_MANAGED` global external Application Load Balancer (advanced..."
  - slug: "02-http-redirect-frontend"
    rank: "02"
    title: "HTTP Redirect VIP (Shared IP)"
    excerpt: "The other half of the standard pair: port 80 on the SAME static IP as the 443 rule, pointing at a target HTTP proxy that serves an http→https redirect URL map. Two forwarding rules may share one..."
  - slug: "03-psc-google-apis"
    rank: "03"
    title: "Private Service Connect to Google APIs"
    excerpt: "The non-load-balancer face of the global forwarding rule: with `loadBalancingScheme: NONE` and the literal target `all-apis`, the rule becomes a Private Service Connect endpoint — VPC workloads reach..."
---

# Global Forwarding Rule Presets

Ready-to-deploy configuration presets for Global Forwarding Rule. Each preset is a complete manifest you can copy, customize, and deploy.
