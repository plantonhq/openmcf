---
title: "Edge Protection: Rate Limits, Geo Fencing, Bot Challenges"
description: "This preset layers custom edge rules over the OWASP baseline: a health-check allowlist, a geo fence, per-client API rate limiting, and a JavaScript challenge for script-shaped user agents -- plus..."
type: "preset"
rank: "02"
presetSlug: "02-rate-limit-and-geo"
componentSlug: "web-application-firewall-policy"
componentTitle: "Web Application Firewall Policy"
provider: "azure"
icon: "package"
order: 2
---

# Edge Protection: Rate Limits, Geo Fencing, Bot Challenges

This preset layers custom edge rules over the OWASP baseline: a health-check
allowlist, a geo fence, per-client API rate limiting, and a JavaScript
challenge for script-shaped user agents -- plus Microsoft's bot-manager rule
set alongside OWASP 3.2.

## When to Use

- Public APIs and web applications that need abuse protection beyond
  attack-signature matching
- Applications that legally or operationally serve a bounded set of
  countries

## Key Configuration Choices

- **Priority ordering is load-bearing** -- the health-check ALLOW runs
  first (5) so probes are never geo-fenced or throttled; the geo fence (10)
  runs before the rate limit (20) so blocked geos never consume counters
- **`GEO_MATCH` + `negationCondition`** -- "block everything except these
  countries" is the negated form; list codes without negation to block
  specific countries instead
- **`RATE_LIMIT_RULE` grouped by `CLIENT_ADDR`** -- 300 requests/minute per
  client IP on `/api/`; rate-limit rules cannot ALLOW (Azure's contract)
- **`JS_CHALLENGE`** -- real browsers solve it invisibly; scripts fail;
  the challenge validity window is tunable via
  `policySettings.jsChallengeCookieExpirationInMinutes`
- **Bot manager 1.1** -- classifies good/bad/unknown bots with Microsoft
  threat intelligence, complementing the signature-based OWASP set

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the policy in | The resource group's `status.outputs.resource_group_name` |
| `<policy-name>` | The policy's name, unique within the resource group | Your naming convention |
| `<allowed-country-code>` | Two-letter ISO codes the app serves (repeat the list entry) | Your compliance/market requirements |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
