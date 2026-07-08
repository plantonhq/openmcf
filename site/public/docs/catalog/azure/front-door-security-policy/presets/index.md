---
title: "Presets"
description: "Ready-to-deploy configuration presets for Front Door Security Policy"
type: "preset-list"
componentSlug: "front-door-security-policy"
componentTitle: "Front Door Security Policy"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-endpoint-association"
    rank: "01"
    title: "Endpoint WAF Association"
    excerpt: "This preset attaches a Front Door WAF policy to an endpoint's default `*.azurefd.net` hostname -- the association that actually turns the WAF on. Without a security policy, a WAF policy sits idle."
  - slug: "02-custom-domain-association"
    rank: "02"
    title: "Custom Domain WAF Association"
    excerpt: "This preset attaches a Front Door WAF policy to both the endpoint's default hostname and a validated custom domain -- the production shape once real domains serve traffic."
---

# Front Door Security Policy Presets

Ready-to-deploy configuration presets for Front Door Security Policy. Each preset is a complete manifest you can copy, customize, and deploy.
