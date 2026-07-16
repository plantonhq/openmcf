---
title: "Standard Rate Limit and IP Denylist"
description: "This preset creates a STANDARD-tier Front Door WAF policy with two custom rules: a per-client rate limit on the API path and an IP denylist -- edge protection without the Premium managed rule sets."
type: "preset"
rank: "01"
presetSlug: "01-standard-rate-limit"
componentSlug: "front-door-firewall-policy"
componentTitle: "Front Door Firewall Policy"
provider: "azure"
icon: "package"
order: 1
---

# Standard Rate Limit and IP Denylist

This preset creates a STANDARD-tier Front Door WAF policy with two
custom rules: a per-client rate limit on the API path and an IP
denylist -- edge protection without the Premium managed rule sets.

## When to Use

- Public APIs behind a Standard Front Door profile that need per-client
  rate limiting enforced at Microsoft's edge, before traffic reaches
  the origin
- Any deployment that needs to block known-bad networks (scrapers,
  abusive clients) globally
- When the Premium tier's managed rule sets are not (yet) justified --
  this policy upgrades in place to more custom rules, but moving to
  PREMIUM replaces the policy (the sku is fixed at creation)

## Key Configuration Choices

- **`sku: STANDARD`** -- custom rules only; must match the sku of the
  Front Door profile the policy gets associated with
- **`mode: PREVENTION`** -- matching requests are actually blocked; use
  DETECTION first if you want to watch a new rule against real traffic
- **`RATE_LIMIT_RULE` with a 5-minute window** -- the threshold counts
  per client; exceeding it blocks for the remainder of the window
- **`SOCKET_ADDR` over `REMOTE_ADDR`** -- matches the source IP Front
  Door actually sees; REMOTE_ADDR trusts X-Forwarded-For, which
  clients can spoof

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-resource-name>` | The AzureResourceGroup's Planton resource name | Your Azure composition |
| `203.0.113.0/24` | The CIDR(s) to block | Your abuse/denylist records |

## After Deploying

The policy enforces nothing until an AzureFrontDoorSecurityPolicy
associates it with your profile's endpoints or custom domains --
deploy one referencing this policy's `firewall_policy_id` output.
