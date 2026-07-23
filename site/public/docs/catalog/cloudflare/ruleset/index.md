---
title: "Ruleset"
description: "Ruleset deployment documentation"
icon: "package"
order: 100
componentName: "cloudflareruleset"
---

# Cloudflare Ruleset

Provision a Cloudflare Ruleset — an ordered collection of rules evaluated
during a specific phase of Cloudflare's HTTP request pipeline. Rulesets are
the unified engine behind WAF custom and managed rules, rate limiting,
cache rules, origin rules, redirects, transforms, and configuration rules:
one component covers every phase.

## What Gets Created

- A `cloudflare_ruleset` at zone scope (most common) or account scope,
  bound to one processing phase (for example
  `http_request_firewall_custom` for WAF custom rules,
  `http_ratelimit` for rate limiting, or
  `http_request_dynamic_redirect` for redirects).
- The ruleset's ordered rules: each pairs a wirefilter match expression
  with an action (block, skip, execute, redirect, rewrite, set cache
  settings, log, and more) and that action's parameters. Rate-limiting
  rules carry their counting characteristics; WAF rules can enable
  exposed-credential checks and per-rule logging overrides.

## Prerequisites

- A Cloudflare zone (reference a `CloudflareDnsZone` or pass a literal
  zone ID) for zone-level rulesets, or an account ID for account-level
  (`custom`/`root`) rulesets — exactly one of the two scopes.
- Rules that reference reusable Lists (for example
  `ip.src in $my_list`) need the `CloudflareList` to exist first.

## Configuration Reference

**Required**

- `phase` — when the ruleset runs in the pipeline (immutable; one ruleset
  per phase per scope).
- `rules` — the ordered rule list; each rule needs an `expression` and an
  `action`.
- `zoneId` or `accountId` — exactly one (the ruleset's scope).

**Optional**

- `rulesetKind` — `zone` (default), `custom` (reusable, invoked via the
  `execute` action), `managed`, or `root`.
- `name`, `description` — identity and documentation.
- Per rule: `actionParameters` (action-specific configuration — redirect
  targets, cache settings, managed-ruleset overrides, header rewrites),
  `ratelimit` (counting characteristics, period, mitigation),
  `exposedCredentialCheck`, `logging`, `enabled`.

Not every action is valid in every phase — Cloudflare enforces the
pairing at the API, and the component's validation mirrors the common
contracts before anything reaches a deploy.

## Stack Outputs

| Output | Description |
|---|---|
| `ruleset_id` | The Cloudflare-assigned ruleset ID |
| `version` | The ruleset version (increments on each update) |
| `zone_id` | Pass-through zone scope for downstream references |
| `phase` | Pass-through phase for infra-chart conditionals |
| `last_updated` | RFC3339 timestamp of the last change |

## Related Components

- CloudflareDnsZone
- CloudflareList / CloudflareListItem
- CloudflareWorker
