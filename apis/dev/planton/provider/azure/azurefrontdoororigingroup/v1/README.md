# AzureFrontDoorOriginGroup

An origin group inside an Azure Front Door profile: the load-balanced
pool of backends a route sends traffic to. The group carries pool-level
behavior -- health probing, latency-aware selection, session affinity,
and the traffic-restore ramp -- while the backends themselves are
first-class AzureFrontDoorOrigin resources referencing the group.

That split is what keeps regional stamps composable: a new region adds
its origin to a shared group without touching the group or any other
region's resources.

## When to Use

Use AzureFrontDoorOriginGroup when you need:

- **A backend pool for a route** -- every AzureFrontDoorRoute forwards
  to exactly one origin group
- **Health-probed failover** -- unhealthy origins leave rotation
  automatically and traffic ramps back gently when they recover
- **Distinct pools per content type** -- e.g. an API pool and a
  static-assets pool behind different routes on one endpoint

## Key Configuration

- `profile_id` -- the parent profile, referenced from an
  AzureFrontDoorProfile's output; fixed at creation
- `origin_group_name` -- 2-90 characters, unique within the profile;
  ForceNew (replaces the group and every origin under it)
- `load_balancing` -- sample size / required successes / latency
  window; unset deploys Azure's defaults (4 / 3 / 50 ms)
- `health_probe` -- protocol, interval, method, path; ABSENT means
  probing disabled (right for single-origin groups)
- `session_affinity_enabled` -- default true; disable for stateless APIs
- `restore_traffic_time_to_healed_or_new_endpoint_in_minutes` -- 0-50,
  default 10; the recovery ramp

## Composition

```yaml
profileId:
  valueFrom:
    kind: AzureFrontDoorProfile
    name: my-front-door
    fieldPath: status.outputs.profile_id
```

Origins join the group and routes point at it through its
`origin_group_id` output.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
