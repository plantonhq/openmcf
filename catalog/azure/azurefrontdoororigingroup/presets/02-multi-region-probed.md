# Multi-Region Probed Backends

This preset creates the production multi-region shape: HTTPS health
probes against a dedicated endpoint, a latency window wide enough to
spread traffic across regions, and session affinity off for stateless
APIs.

## When to Use

- APIs or web apps with origins in two or more regions
- Any group where automatic failover matters -- probes are what take an
  unhealthy origin out of rotation

## Key Configuration Choices

- **`additionalLatencyInMilliseconds: 100`** -- regions within 100 ms of
  the fastest count as equally fast and share traffic by weight; tighten
  toward 0 to pin traffic to the closest region, widen for more even
  geographic spread
- **GET probes to `/healthz`** -- a dedicated health endpoint proves the
  application, not just the web server; HEAD is cheaper when the
  endpoint supports it
- **~2-minute ejection** -- 4 samples at 30 s intervals with 3 required
  successes; shorten the interval for faster failover at the cost of
  more probe traffic (every Front Door edge location probes every
  origin)
- **`sessionAffinityEnabled: false`** -- right for stateless APIs; keep
  Azure's default (on) for session-based web apps

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `originGroupName` (example value) | 2-90 chars -- rename to your convention | Your naming convention |

## Downstream Wiring

Each region's stamp adds its own AzureFrontDoorOrigin referencing this
group -- the group itself never changes as regions come and go:

```yaml
# On each region's AzureFrontDoorOrigin
originGroupId:
  valueFrom:
    kind: AzureFrontDoorOriginGroup
    name: my-api-backends
    fieldPath: status.outputs.origin_group_id
```
