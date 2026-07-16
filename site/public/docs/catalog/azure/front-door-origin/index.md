---
title: "Front Door Origin"
description: "Front Door Origin deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoororigin"
---

# Azure Front Door Origin

Creates one backend inside an AzureFrontDoorOriginGroup: the backend's address, TLS validation posture, priority and weight in the pool, and optional Private Link connectivity (Premium profiles). Each region or deployment slot adds its own origin to a shared group -- the group never changes as backends come and go.

## What Gets Created

When you deploy an AzureFrontDoorOrigin resource, Planton provisions:

- **Front Door Origin** -- an `azurerm_cdn_frontdoor_origin` on the referenced origin group, with optional Private Link to the backend

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureFrontDoorOriginGroup** to join (referenced through `originGroupId`)
- **For Private Link**: a Premium-tier profile, and someone to approve the pending connection on the target resource

## Quick Start

Create a file `origin.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorOrigin
metadata:
  name: primary-app
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureFrontDoorOrigin.primary-app
spec:
  originGroupId:
    valueFrom:
      kind: AzureFrontDoorOriginGroup
      name: api-backends
      fieldPath: status.outputs.origin_group_id
  originName: primary-app
  hostName: myapp.azurewebsites.net
```

Deploy:

```shell
planton apply -f origin.yaml
```

The defaults suit App Service and other multi-tenant Azure backends: the Host header falls back to the origin's own hostname (which those services route by), certificate-name checking stays on, and ports stay 80/443. Use `priority` for active/passive failover tiers and `weight` for the traffic split within a tier -- a `weight: 50` origin beside a `weight: 950` one is a ~5% canary.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `origin_id` | The ARM id -- what AzureFrontDoorRoute lists in `originIds` to sequence deployment |
| `origin_name` | The origin's name inside its group |

## Related Resources

- [Azure Front Door Origin Group](/docs/catalog/azure/front-door-origin-group) -- the parent pool
- [Azure Front Door Route](/docs/catalog/azure/front-door-route) -- forwards traffic to the pool
- [Azure Front Door Profile](/docs/catalog/azure/front-door-profile) -- the grandparent container (PREMIUM required for Private Link)
