---
title: "Presets"
description: "Ready-to-deploy configuration presets for WAF Regex Pattern Set"
type: "preset-list"
componentSlug: "waf-regex-pattern-set"
componentTitle: "WAF Regex Pattern Set"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-scanner-probes"
    rank: "01"
    title: "Scanner Probe Patterns"
    excerpt: "A REGIONAL set of common automated-scanner URI paths. Pair with a web ACL block rule whose `regex_pattern_set_reference` inspects `uriPath`."
  - slug: "02-internal-admin-paths"
    rank: "02"
    title: "Internal Admin Path Patterns"
    excerpt: "Regexes for internal admin and debug URL prefixes — the shape for applications that expose operator endpoints under predictable paths."
---

# WAF Regex Pattern Set Presets

Ready-to-deploy configuration presets for WAF Regex Pattern Set. Each preset is a complete manifest you can copy, customize, and deploy.
