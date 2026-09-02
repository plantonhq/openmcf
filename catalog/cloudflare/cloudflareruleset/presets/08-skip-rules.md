# Skip Specific Rules in Another Ruleset

Surgically exempt a path from individual rules inside another ruleset (typically noisy managed WAF rules that false-positive on your payloads) without turning off the whole ruleset or product.

## When to use

- A managed WAF rule blocks legitimate webhook or API payloads and you want to skip exactly that rule, keeping the rest of the protection active.
- Finer-grained than `action_parameters.ruleset` (which skips a whole ruleset) or `products`.

## Key choices

- `action_parameters.rules` maps a target ruleset ID to the rule IDs to skip inside it — both are 32-character hex IDs from the Cloudflare API (`GET /zones/<zone_id>/rulesets` and the ruleset's own rules list).
- Incompatible with the single `ruleset` option on the same rule — pick one grain.
- The skip rule's `expression` bounds WHERE the exemption applies; keep it as narrow as possible.

## Placeholders

| Placeholder | Description |
|---|---|
| `<cloudflare-zone-id>` | The zone the skip ruleset lives in |
| `efb7b8c949ac4650a09736fc376e9aee` | The ruleset containing the rules to skip (32-char hex; the shipped value is Cloudflare's managed WAF ruleset id — replace with your target) |
| `5de7edfa648c4d6891dc3e7f84534ffa` | A rule ID inside the target ruleset (32-char hex, from the ruleset's own rules list) |
| `75a0060762034a6cb663fd51a02344cb` | Another rule ID to skip (remove if only one) |
