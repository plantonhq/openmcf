---
title: "Front Door Origin Group"
description: "Front Door Origin Group deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoororigingroup"
---

# Azure Front Door Origin Group

Creates an origin group inside an AzureFrontDoorProfile -- the load-balanced pool of backends a route sends traffic to. The group owns health probing, latency-aware origin selection, session affinity, and the traffic-restore ramp; the backends are AzureFrontDoorOrigin resources that reference the group.

## What Gets Created

When you deploy an AzureFrontDoorOriginGroup resource, Planton provisions:

- **Front Door Origin Group** -- an `azurerm_cdn_frontdoor_origin_group` on the referenced profile, with load balancing and optional health probe settings

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureFrontDoorProfile** to create the group in (referenced through `profileId`)

## Quick Start

Create a file `origin-group.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFrontDoorOriginGroup
metadata:
  name: api-backends
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureFrontDoorOriginGroup.api-backends
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: my-front-door
      fieldPath: status.outputs.profile_id
  originGroupName: api-backends
  healthProbe:
    protocol: HTTPS
    intervalInSeconds: 30
    path: /healthz
```

Deploy:

```shell
planton apply -f origin-group.yaml
```

Omit `healthProbe` entirely for single-origin groups -- Front Door treats absent probe settings as probing disabled, and with one origin probes only add load. Leave `loadBalancing` unset to deploy Azure's defaults; tune `additionalLatencyInMilliseconds` upward to spread traffic across geographically dispersed regions.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `origin_group_id` | The ARM id -- what AzureFrontDoorOrigin references as parent and AzureFrontDoorRoute references as destination |
| `origin_group_name` | The group's name inside the profile |

## Related Resources

- [Azure Front Door Profile](/docs/catalog/azure/front-door-profile) -- the parent profile
- [Azure Front Door Origin](/docs/catalog/azure/front-door-origin) -- the backends in this group
- [Azure Front Door Route](/docs/catalog/azure/front-door-route) -- sends traffic to this group
