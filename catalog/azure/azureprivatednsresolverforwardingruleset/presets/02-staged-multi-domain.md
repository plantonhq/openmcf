# Staged Multi-Domain

This preset carries two namespaces on one rule book: the live corporate domain, plus an acquired company's domain PARKED with `enabled: false` until its connectivity is ready -- flipping one field cuts it over, with no redeploys and no rule rewrites.

## When to Use

- Estates resolving more than one private namespace (acquisitions, subsidiaries, partner networks)
- Staged migrations where the rule should exist -- reviewed, correct, attributed -- before the tunnel it depends on is live

## Key Configuration Choices

- **A disabled rule keeps its configuration but forwards nothing** -- the domain resolves normally inside Azure until you enable it
- **Per-rule `metadata` records ownership** -- who owns the rule and why it is parked lives ON the rule, not in a wiki
- **A non-standard port rides the target entry** (`port: 5353`) -- unset targets default to 53, the standard DNS port

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The AzureResourceGroup the ruleset is created in | The resource-group component's name |
| `<your-dns-resolver>` | The AzurePrivateDnsResolver whose outbound endpoint the rules steer | The resolver component's name |

Replace the target `ipAddress` values with your actual DNS servers. Everything on a rule except its domain updates in place -- enabling the parked rule later is a one-field change.
