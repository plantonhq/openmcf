---
title: "Standard Delivery"
description: "This preset creates a Standard-tier Front Door profile -- the container for a global CDN deployment. Endpoints, origin groups, origins, and routes compose against it as their own resources, so the..."
type: "preset"
rank: "01"
presetSlug: "01-standard-delivery"
componentSlug: "front-door-profile"
componentTitle: "Front Door Profile"
provider: "azure"
icon: "package"
order: 1
---

# Standard Delivery

This preset creates a Standard-tier Front Door profile -- the container
for a global CDN deployment. Endpoints, origin groups, origins, and
routes compose against it as their own resources, so the profile itself
stays a small, stable anchor.

## When to Use

- Public web applications and APIs that need global load balancing,
  edge caching, and HTTPS at the edge
- Any deployment that does NOT need Private Link to origins or the
  managed WAF rule sets (those are PREMIUM-only)

## Key Configuration Choices

- **`sku: STANDARD`** -- the production default. The sku is fixed at
  creation, and Azure refuses a PREMIUM -> STANDARD downgrade outright,
  so upgrade paths are one-way; pick PREMIUM only when its features are
  actually needed
- **Default 120 s response timeout** -- raise it for slow APIs or large
  downloads, lower it when fast failover matters more
- **No identity** -- a managed identity is only needed when custom
  domains carry bring-your-own Key Vault certificates

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The AzureResourceGroup's Planton resource name | Your Azure composition |
| `profileName` (example value) | 2-90 chars, letters/digits/hyphens -- rename to your convention | Your naming convention |

## Downstream Wiring

Every Front Door delivery resource references this profile's ARM id:

```yaml
# On an AzureFrontDoorEndpoint or AzureFrontDoorOriginGroup
profileId:
  valueFrom:
    kind: AzureFrontDoorProfile
    name: my-front-door
    fieldPath: status.outputs.profile_id
```
