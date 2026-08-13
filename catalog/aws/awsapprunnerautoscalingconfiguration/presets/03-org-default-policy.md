# Organization Default Scaling Policy

This preset registers a balanced auto scaling configuration and claims it as the ACCOUNT-WIDE default for its region: every App Runner service created afterwards without an explicit `autoScalingConfigurationArn` adopts this policy automatically. Use it to give your whole account a sane scaling baseline instead of AWS's built-in default (100 concurrency / 25 max / 1 min).

## When to Use

- Platform teams setting a governed scaling baseline for every team's services
- Accounts where most services should share one posture and only exceptions reference their own configuration
- Replacing AWS's built-in account default with an explicitly owned, versioned policy

## Key Configuration Choices

- **Balanced dials** (`maxConcurrency: 80`, `maxSize: 10`, `minSize: 1`) -- Headroom per instance without a runaway cost ceiling. Tune to your org's traffic profile; every change registers a new revision and the designation follows the configuration name.
- **Account-default claim** (`setAsAccountDefault: true`) -- The load-bearing choice. Three truths to understand before applying:
  - **One per account/region.** Claiming silently displaces whichever configuration held the designation. Coordinate so exactly one resource in your account sets this flag.
  - **Only future services are affected.** Existing services keep their current associations; the default applies to services created after the claim.
  - **One-way at AWS.** Destroying this resource (or dropping the flag) does NOT restore the previous default -- AWS has no restore API. The designation stays until another configuration claims it.
- **`status.outputs.is_default`** confirms the claim after deployment.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `region` | AWS region the configuration (and its designation) lives in | Your deployment region |
| Scaling dials | `maxConcurrency` / `maxSize` / `minSize` for your org's baseline | Your traffic profile |

## Related Presets

- **01-latency-sensitive-api** -- A per-service posture with warm instances; reference it explicitly from the services that need it.
- **02-scale-conservative** -- A cost-capped posture for background services.
