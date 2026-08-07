# Canary Weighted Forward

This preset splits one route's traffic between two target groups -- 95% to
the stable version, 5% to the canary -- with group stickiness so a client
stays on whichever version served its first request. Promoting the canary is
a weight edit (5 → 25 → 50 → 100), and rolling back is setting its weight to
0; the rule, listener, and DNS never change.

## When to Use

- Validating a new service version against a slice of real traffic before
  full rollout
- Blue/green deployments where the shift should be gradual rather than
  all-at-once
- Per-route experiments -- the split applies to exactly the requests this
  rule matches, leaving every other route alone

## Key Configuration Choices

- **Weights are relative** -- 95/5 means 95 parts of 100; the numbers need
  not sum to anything in particular (0-999 each, 0 drains a group)
- **Stickiness is group-level** -- a client pins to the *group* (not the
  individual target) for `durationSeconds` (here 1 hour), so sessions do not
  flap between versions mid-canary; drop the block for stateless APIs where
  per-request splitting is fine
- **Two independent target groups** -- each version keeps its own health
  checks and drain behavior, so a failing canary marks itself unhealthy
  without touching stable capacity

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<service-name>` | Name prefix for the rule resource | Your service's name |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<listener-resource-name>` | Name of the AwsLbListener to attach to | Your AwsLbListener manifest's `metadata.name` |
| `<service-hostname>` | Hostname to match (e.g., `api.example.com`) | Your DNS zone |
| `<stable-target-group-name>` | AwsLbTargetGroup running the current version | Your stable target group manifest's `metadata.name` |
| `<canary-target-group-name>` | AwsLbTargetGroup running the new version | Your canary target group manifest's `metadata.name` |

## Common Additions

- Swap the `hostHeader` condition for a `pathPattern` to canary a single
  path
- Add an `httpHeader` condition (e.g. a `x-canary: always` header) as a
  second rule with higher priority, so testers can opt into the canary
  deterministically

## Related Presets

- **01-path-based-routing** -- the single-group forward this graduates from
- **02-host-based-routing** -- split by domain instead of path
