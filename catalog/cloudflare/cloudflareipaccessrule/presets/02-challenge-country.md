# Challenge a country

A zone-scoped IP Access rule that presents Cloudflare's managed challenge to visitors from one country. `managed_challenge` picks the least intrusive challenge that confirms the visitor is human.

## When to Use

- Challenge (not block) traffic from a country on one zone
- Tor exit nodes -- set `value: T1` instead of a country code
- A zone override of a looser account-wide rule

## Key Configuration Choices

- **zone_id, not account_id** -- exactly one scope. This rule applies only to the named zone.
- **mode: managed_challenge** -- recommended over `challenge` or `js_challenge`. Updates in place if you later switch to `block`.
- **target: country** -- two-character ISO 3166-1 alpha-2. Changing the country later does not stick; recreate the rule.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone this rule applies to | Cloudflare Dashboard -> zone Overview -> API section (right sidebar), or reference a CloudflareDnsZone via `value_from` instead |
| `configuration.value` | Country code (`US`) or `T1` for Tor | Your geo policy |

## Related Presets

- **01-block-ip** -- an account-wide block of a single IPv4 address
