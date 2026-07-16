---
title: "Presets"
description: "Ready-to-deploy configuration presets for WAF IP Set"
type: "preset-list"
componentSlug: "waf-ip-set"
componentTitle: "WAF IP Set"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-office-allowlist"
    rank: "01"
    title: "Office Allow-List"
    excerpt: "A REGIONAL IPv4 set carrying a /24 office range and a /32 single-host entry — the two CIDR shapes AWS accepts. Pair with a web ACL whose default action is block and an early-priority allow rule..."
  - slug: "02-placeholder-set"
    rank: "02"
    title: "Placeholder IP Set"
    excerpt: "An empty REGIONAL IPv4 set — valid in AWS and useful when the web ACL and its rules must exist before the address list is finalized."
---

# WAF IP Set Presets

Ready-to-deploy configuration presets for WAF IP Set. Each preset is a complete manifest you can copy, customize, and deploy.
