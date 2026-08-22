# CloudflareRuleset guide

Operational judgment for rulesets. The README covers phases, actions, and fields; this covers how the pieces interact.

## One ruleset per zone per phase

Cloudflare allows exactly one custom ruleset per `{scope, phase}`. Deploying a second resource for the same zone and phase does not add rules — it fights the first for the same underlying ruleset. Keep each phase's rules in one manifest and let list order (or explicit ordering by rule) express precedence.

## The skip grains, narrowest first

The skip action has five surfaces, and the narrowest that works is the right one:

- `rules` (map of ruleset ID → rule IDs) — skip INDIVIDUAL rules inside another ruleset. The surgical tool for managed-WAF false positives: the rest of the ruleset keeps protecting. Incompatible with `ruleset` on the same rule.
- `ruleset` — skip one whole ruleset (`"current"` skips the remainder of this one).
- `rulesets` — skip several whole rulesets.
- `phases` / `products` — skip entire phases or products; the bluntest grain, easy to over-exempt.

Both IDs in the `rules` map are 32-char hex from the API (`GET /zones/<zone_id>/rulesets`, then the target ruleset's rules). They are environment-specific — parameterize them per environment rather than hardcoding one environment's IDs into a shared manifest.

## Entitlements bound what a phase can do

Custom firewall rules (`http_request_firewall_custom`) work on the free plan. Executing managed WAF rulesets (`http_request_firewall_managed`) and rate limiting (`http_ratelimit`) are paid-plan features — a deploy into those phases on a free-plan zone fails at the API, not at validation.

## Expressions fail at deploy, not at validation

Rule expressions are Cloudflare's wirefilter language, parsed server-side. A typo deploys nothing: the API rejects the whole ruleset. Keep expressions short, test complicated ones in the dashboard's expression editor first, and prefer several narrow rules over one giant expression — per-rule enable flags then give you a kill switch per behavior.
