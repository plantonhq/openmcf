# Security Headers Policy

This preset creates a rule set with one condition-less rule that stamps
the standard security headers on every response and strips the
backend's technology fingerprint -- the baseline every production site
behind Front Door should attach.

## When to Use

- Any Front Door deployment serving browsers -- attach this set to every
  route so the security posture is enforced at the edge regardless of
  which backend answered
- As the first rule set in a profile: one shared policy beats
  re-implementing headers in every application

## Key Configuration Choices

- **No `conditions`** -- deliberately absent: a rule without conditions
  applies to every request, which is exactly right for a baseline
- **`OVERWRITE` rather than `APPEND`** -- the edge is authoritative for
  these headers; overwriting prevents duplicate or conflicting values
  when a backend also sets them
- **`DELETE X-Powered-By`** -- header DELETE actions carry no value (the
  spec enforces it)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `ruleSetName` (example value) | 1-60 letters/digits, no hyphens -- rename to your convention | Your naming convention |

## Downstream Wiring

Routes attach the set by ARM ID:

```yaml
# On an AzureFrontDoorRoute
ruleSetIds:
  - valueFrom:
      kind: AzureFrontDoorRuleSet
      name: my-security-headers
      fieldPath: status.outputs.rule_set_id
```
