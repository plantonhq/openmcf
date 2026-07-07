# Azure Front Door Profile

Creates an Azure Front Door (Standard/Premium) profile -- the container for a global CDN deployment on Microsoft's edge network. Endpoints, origin groups, origins, and routes compose against the profile as first-class resources, so one profile anchors many applications' delivery.

## What Gets Created

When you deploy an AzureFrontDoorProfile resource, Planton provisions:

- **Front Door Profile** -- an `azurerm_cdn_frontdoor_profile` with the chosen SKU tier, origin response timeout, optional managed identity, optional access-log scrubbing, and tags

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureResourceGroup** to create the profile in (referenced through `resourceGroup`)

## Quick Start

Create a file `front-door.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorProfile
metadata:
  name: my-front-door
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureFrontDoorProfile.my-front-door
spec:
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: my-rg
      fieldPath: status.outputs.resource_group_name
  profileName: my-front-door
```

Deploy:

```shell
planton apply -f front-door.yaml
```

This deploys a Standard-tier profile. Choose `sku: PREMIUM` when origins need Private Link (public access disabled on backends) or when the managed WAF rule sets will be attached -- the tier is fixed at creation and Azure refuses a downgrade. Add an `identity` block when custom domains will carry bring-your-own Key Vault certificates, and `logScrubbingVariables` to mask client IPs, URIs, or query strings out of access logs.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `profile_id` | The ARM id -- what AzureFrontDoorEndpoint and AzureFrontDoorOriginGroup reference as their parent |
| `profile_name` | The ARM namespace every child resource nests under |
| `resource_guid` | Front Door's own GUID for the profile (traffic-ownership validation, e.g. apex-domain afdverify DNS records) |
| `identity_principal_id` | The system-assigned identity's principal -- the Key Vault grant target for bring-your-own TLS |

## Related Resources

- [Azure Front Door Endpoint](/docs/catalog/azure/azurefrontdoorendpoint) -- the public entry hostname
- [Azure Front Door Origin Group](/docs/catalog/azure/azurefrontdoororigingroup) -- the load-balanced backend pool
- [Azure Front Door Origin](/docs/catalog/azure/azurefrontdoororigin) -- one backend in a group
- [Azure Front Door Route](/docs/catalog/azure/azurefrontdoorroute) -- connects endpoints to origin groups
