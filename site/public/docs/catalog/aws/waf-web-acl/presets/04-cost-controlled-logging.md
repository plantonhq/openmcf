---
title: "Cost-Controlled Logging"
description: "This preset creates a Web ACL with baseline managed-rule protection, IP rate limiting, and request logging that keeps ONLY the records where WAF acted -- blocked requests, counted matches, and rules..."
type: "preset"
rank: "04"
presetSlug: "04-cost-controlled-logging"
componentSlug: "waf-web-acl"
componentTitle: "WAF Web ACL"
provider: "aws"
icon: "package"
order: 4
---

# Cost-Controlled Logging

This preset creates a Web ACL with baseline managed-rule protection, IP rate limiting, and request logging that keeps ONLY the records where WAF acted -- blocked requests, counted matches, and rules demoted to count by a tuning override. Allowed traffic (the overwhelming majority on a healthy application) is dropped before it reaches the log destination, which is where WAF logging costs actually accrue.

## When to Use

- High-traffic applications where logging every inspected request is prohibitively expensive
- Security monitoring that only needs enforcement events, not full traffic capture
- WAF tuning workflows -- `EXCLUDED_AS_COUNT` keeps the records for rules you have overridden to count while you evaluate false positives

## Key Configuration Choices

- **`filter.defaultBehavior: DROP`** -- records matching no filter are discarded; the KEEP filter is the allowlist
- **`MEETS_ANY` with three action conditions** -- a record is kept when WAF blocked it, counted it, or would have blocked it but for a count override
- **Header and query-string redaction** -- `authorization` and `cookie` values and the query string appear as `REDACTED` in kept records, so credentials never land in logs
- **Destination name must start with `aws-waf-logs-`** -- AWS rejects any other destination name at the API
- **One destination per web ACL** -- AWS's contract allows exactly one logging destination

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `<acl-name>` | Unique name for the Web ACL (lowercase, alphanumeric, hyphens) |
| `<aws-region>` | Region for the ACL and its log group (e.g. `us-west-2`) |
| `<account-id>` | The AWS account id owning the log group |

## Common Additions

- A `labelName` condition to drop a specific managed group's noise, e.g. keep everything except records labeled `awswaf:managed:aws:core-rule-set:NoUserAgent_Header`
- `redactUriPath: true` / `redactMethod: true` for stricter log hygiene
- `dataProtectionConfig` when fields must be masked in ALL WAF outputs (sampled requests and rule match details, not just the log destination)

## Related Presets

- **01-managed-rules-basic** -- the protection baseline without logging
- **03-production-web-app** -- full production configuration with unfiltered logging
