---
title: "Production Endpoint"
description: "This preset creates an enabled endpoint on an existing Front Door profile -- the public entry hostname client traffic arrives at. Routes attach to it; custom-domain DNS records CNAME onto its..."
type: "preset"
rank: "01"
presetSlug: "01-production-endpoint"
componentSlug: "front-door-endpoint"
componentTitle: "Front Door Endpoint"
provider: "azure"
icon: "package"
order: 1
---

# Production Endpoint

This preset creates an enabled endpoint on an existing Front Door
profile -- the public entry hostname client traffic arrives at. Routes
attach to it; custom-domain DNS records CNAME onto its generated
hostname.

## When to Use

- Every application a Front Door profile serves -- one endpoint per
  app is the common shape (one profile often fronts several endpoints)
- As the anchor for routes ("/api/*", "/static/*") that split one
  hostname's traffic across origin groups

## Key Configuration Choices

- **The name is the hostname prefix** -- Azure appends a hash
  (`{name}-{hash}.z01.azurefd.net`), which is why the name only needs
  per-profile uniqueness, and why renaming breaks every DNS record
  pointing at the old hostname
- **Enabled by default** -- the `enabled` switch exists for maintenance
  windows and staged cutovers, not day-one configuration

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `endpointName` (example value) | 2-46 chars; becomes the public hostname prefix -- rename to your convention | Your naming convention |

## Downstream Wiring

Routes attach to this endpoint, and DNS records point at its hostname:

```yaml
# On an AzureFrontDoorRoute
endpointId:
  valueFrom:
    kind: AzureFrontDoorEndpoint
    name: my-web-endpoint
    fieldPath: status.outputs.endpoint_id

# On an AzureDnsRecord (CNAME for a custom domain)
values:
  - valueFrom:
      kind: AzureFrontDoorEndpoint
      name: my-web-endpoint
      fieldPath: status.outputs.host_name
```
