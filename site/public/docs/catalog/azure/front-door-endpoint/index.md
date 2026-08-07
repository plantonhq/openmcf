---
title: "Front Door Endpoint"
description: "Front Door Endpoint deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoorendpoint"
---

# Azure Front Door Endpoint

Deploys a public Front Door endpoint -- the `*.azurefd.net` hostname clients hit at Microsoft's edge. The endpoint belongs to an AzureFrontDoorProfile and is the attachment point for routes (and later custom domains). Renaming the endpoint replaces it and changes the generated hostname, so every DNS record pointing at the old name breaks. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring to the parent profile.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Front Door Endpoint** -- a named child of the profile with a generated public hostname (`{endpointName}-{hash}.z01.azurefd.net`)
- **Enabled state** -- whether the endpoint accepts traffic (default on; disable for maintenance without deleting the hostname)
- **Azure Tags** -- resource metadata tags merged with user tags and applied to the endpoint

## The Endpoint in the Front Door Family

The endpoint is the public face of the delivery surface:

- **AzureFrontDoorProfile** -- the parent container, referenced by `profileId`; its SKU tier and identity govern every child
- **AzureFrontDoorRoute** -- attaches URL patterns to this endpoint and forwards them to an origin group
- **AzureFrontDoorCustomDomain** -- later, CNAMEs a branded hostname at this endpoint's `host_name` output (or aliases it via Azure DNS)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Azure Subscription

- **An Azure Front Door Profile** the endpoint nests under. Reference an AzureFrontDoorProfile Cloud Resource via ValueFromRef, or provide the profile ARM ID directly.

## Deploy

### Console

Open the deployment store, find **Azure Front Door Endpoint**, and click **Deploy**. The wizard walks you through the parent profile, the endpoint name (2–46 characters -- shorter than other Front Door names because Azure builds the public hostname from it), and the enabled switch. Start from the **Production Endpoint** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorEndpoint
metadata:
  name: my-web-endpoint
  org: acme-corp
  env: prod
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  endpointName: my-web
```

```shell
planton apply -f front-door-endpoint.yaml
```

This creates an enabled endpoint under the profile. Routes attach next, referencing the `endpoint_id` output; custom-domain DNS points at the `host_name` output.

### InfraChart

```yaml
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  endpointName: my-web
```

## Key Configuration

**Endpoint name** -- 2–46 characters; letters, digits, and hyphens; must start and end with a letter or digit. Becomes the prefix of the generated public hostname. ForceNew: renaming replaces the endpoint and changes the hostname.

**Enabled** -- the traffic switch. Leave on for production; set false for a maintenance hold without deleting the DNS target. Defaults to true when omitted.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureFrontDoorProfile** | `profileId` | `status.outputs.profile_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `endpoint_id` | ARM resource ID of the endpoint | AzureFrontDoorRoute.`endpointId` |
| `endpoint_name` | The endpoint's name | Operator tooling |
| `host_name` | The generated `*.azurefd.net` hostname | AzureDnsRecord CNAME/alias targets; AzureFrontDoorCustomDomain validation |

## Presets

| Preset | Rank | Description |
|--------|------|-------------|
| Production Endpoint | 1 | Enabled endpoint for live traffic |
| Maintenance Disabled | 2 | Endpoint created with traffic off |

## Related Components

- **AzureFrontDoorProfile** -- the parent container
- **AzureFrontDoorRoute** -- attaches patterns to this endpoint
- **AzureDnsRecord** -- CNAMEs or aliases a custom hostname at `host_name`
