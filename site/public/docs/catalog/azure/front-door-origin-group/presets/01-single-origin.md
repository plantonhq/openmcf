---
title: "Single-Origin Group"
description: "This preset creates an origin group with no health probe -- the honest shape for a single-backend deployment, where probing adds origin load without buying failover (there is nowhere to fail over to)."
type: "preset"
rank: "01"
presetSlug: "01-single-origin"
componentSlug: "front-door-origin-group"
componentTitle: "Front Door Origin Group"
provider: "azure"
icon: "package"
order: 1
---

# Single-Origin Group

This preset creates an origin group with no health probe -- the honest
shape for a single-backend deployment, where probing adds origin load
without buying failover (there is nowhere to fail over to).

## When to Use

- One backend behind Front Door (a single App Service, storage static
  site, or API) where the CDN value is caching and edge TLS, not
  multi-origin load balancing
- The starting point that grows into a probed multi-origin group when a
  second region arrives

## Key Configuration Choices

- **No `healthProbe`** -- deliberately absent: Front Door treats missing
  probe settings as probing disabled, and a single origin is served
  regardless of probe results anyway
- **Azure's load-balancing defaults** -- meaningful only once the group
  has multiple origins; safe to leave until then
- **Session affinity stays on** (Azure's default) -- harmless with one
  origin, correct if a stateful second origin appears later

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `originGroupName` (example value) | 2-90 chars -- rename to your convention | Your naming convention |

## Downstream Wiring

Origins join this group; routes point at it:

```yaml
# On an AzureFrontDoorOrigin
originGroupId:
  valueFrom:
    kind: AzureFrontDoorOriginGroup
    name: my-app-backends
    fieldPath: status.outputs.origin_group_id
```
