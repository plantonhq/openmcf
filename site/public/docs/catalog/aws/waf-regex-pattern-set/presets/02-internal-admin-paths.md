---
title: "Internal Admin Path Patterns"
description: "Regexes for internal admin and debug URL prefixes — the shape for applications that expose operator endpoints under predictable paths."
type: "preset"
rank: "02"
presetSlug: "02-internal-admin-paths"
componentSlug: "waf-regex-pattern-set"
componentTitle: "WAF Regex Pattern Set"
provider: "aws"
icon: "package"
order: 2
---

# Internal Admin Path Patterns

Regexes for internal admin and debug URL prefixes — the shape for applications that expose operator endpoints under predictable paths.

## When to Use

- Applications with `/internal/admin/*` or `/debug/*` routes that should be blocked at the edge unless traffic comes from an allow-listed IP set
- Pairing with an IP-set allow rule: allow office IPs OR block admin paths for everyone else

## What It Configures

- **Two path-prefix patterns** — broad internal admin tree plus a specific debug status endpoint
- **REGIONAL scope**

## What to Customize

- Replace `<aws-region>` with your target region
- Adjust patterns to your application's actual admin URL structure
- Combine with [AwsWafIpSet](../../awswafipset/presets/01-office-allowlist.yaml) for allow-list + block-probes layering

## Layering Example

1. Priority 1 — allow rule referencing office IP set
2. Priority 10 — block rule referencing this pattern set on `uriPath`
3. Managed rule groups at higher priorities for OWASP coverage
