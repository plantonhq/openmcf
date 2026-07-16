---
title: "Presets"
description: "Ready-to-deploy configuration presets for Front Door Firewall Policy"
type: "preset-list"
componentSlug: "front-door-firewall-policy"
componentTitle: "Front Door Firewall Policy"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-rate-limit"
    rank: "01"
    title: "Standard Rate Limit and IP Denylist"
    excerpt: "This preset creates a STANDARD-tier Front Door WAF policy with two custom rules: a per-client rate limit on the API path and an IP denylist -- edge protection without the Premium managed rule sets."
  - slug: "02-premium-managed-prevention"
    rank: "02"
    title: "Premium Managed Rules Prevention"
    excerpt: "This preset creates a PREMIUM-tier Front Door WAF policy running Microsoft's managed rule sets in blocking mode -- the default posture for production workloads on a Premium profile: OWASP-class..."
  - slug: "03-detection-rollout"
    rank: "03"
    title: "Detection-Mode Rollout"
    excerpt: "This preset creates a Premium WAF policy in DETECTION mode -- the managed rule sets run and log everything they would have blocked, but no traffic is affected. It is the safe first step of every WAF..."
---

# Front Door Firewall Policy Presets

Ready-to-deploy configuration presets for Front Door Firewall Policy. Each preset is a complete manifest you can copy, customize, and deploy.
