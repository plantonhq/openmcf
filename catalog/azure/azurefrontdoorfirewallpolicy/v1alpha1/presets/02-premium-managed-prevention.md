# Premium Managed Rules Prevention

This preset creates a PREMIUM-tier Front Door WAF policy running
Microsoft's managed rule sets in blocking mode -- the default posture
for production workloads on a Premium profile: OWASP-class attack
coverage and bot classification, maintained by Microsoft, with log
scrubbing for sensitive request data.

## When to Use

- Production web applications and APIs behind a PREMIUM Front Door
  profile -- the managed rule sets are the main reason the Premium
  tier exists
- When you want attack-signature maintenance outsourced: Microsoft
  updates the rule sets server-side; you keep only the tuning
  (exclusions, overrides) in code
- Compliance postures that require sensitive data (auth headers,
  client IPs) scrubbed from security logs

## Key Configuration Choices

- **`Microsoft_DefaultRuleSet 2.1` on anomaly scoring** -- individual
  rules contribute to a score instead of acting alone, which cuts
  false positives; per-rule overrides on 2.x sets may only use
  `OVERRIDE_ANOMALY_SCORING` or `OVERRIDE_LOG` (the spec enforces it)
- **Cookie exclusion over rule disabling** -- the `session-token`
  exclusion shows the surgical false-positive tool: the rules still
  run, they just skip that one cookie
- **JS challenge for unknown bots** -- `OVERRIDE_JS_CHALLENGE` is valid
  only on the bot manager set; real browsers solve it invisibly while
  scripted clients fail
- **Log scrubbing enabled** -- the Authorization header and client IPs
  never reach the logs

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-resource-name>` | The AzureResourceGroup's Planton resource name | Your Azure composition |
| `session-token` | Your app's session cookie name | Your application's auth configuration |

## After Deploying

Associate the policy with your profile's domains through an
AzureFrontDoorSecurityPolicy referencing this policy's
`firewall_policy_id` output. The profile must also be PREMIUM -- Azure
rejects a sku mismatch at association time.
