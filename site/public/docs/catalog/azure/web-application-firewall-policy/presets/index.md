---
title: "Presets"
description: "Ready-to-deploy configuration presets for Web Application Firewall Policy"
type: "preset-list"
componentSlug: "web-application-firewall-policy"
componentTitle: "Web Application Firewall Policy"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-owasp-baseline"
    rank: "01"
    title: "OWASP 3.2 Baseline"
    excerpt: "This preset creates the policy almost every gateway should start with: the OWASP 3.2 core rule set in Prevention mode (Azure's default), no custom rules, no overrides. It blocks SQL injection,..."
  - slug: "02-rate-limit-and-geo"
    rank: "02"
    title: "Edge Protection: Rate Limits, Geo Fencing, Bot Challenges"
    excerpt: "This preset layers custom edge rules over the OWASP baseline: a health-check allowlist, a geo fence, per-client API rate limiting, and a JavaScript challenge for script-shaped user agents -- plus..."
  - slug: "03-detection-tuning"
    rank: "03"
    title: "Detection-Mode Tuning with Exclusions and Overrides"
    excerpt: "This preset is the false-positive workflow: the policy runs in Detection mode (logging matches without blocking), a scoped exclusion skips the session cookie that trips SQL-injection rules, two..."
---

# Web Application Firewall Policy Presets

Ready-to-deploy configuration presets for Web Application Firewall Policy. Each preset is a complete manifest you can copy, customize, and deploy.
