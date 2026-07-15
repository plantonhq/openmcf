---
title: "Front Door Endpoint"
description: "Front Door Endpoint deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoorendpoint"
---

# Azure Front Door Endpoint

Creates an endpoint inside an AzureFrontDoorProfile -- the public entry point client traffic arrives at. Azure generates a globally unique `{name}-{hash}.z01.azurefd.net` hostname; routes attach to the endpoint and custom-domain DNS records CNAME onto its hostname.

## What Gets Created

When you deploy an AzureFrontDoorEndpoint resource, Planton provisions:

- **Front Door Endpoint** -- an `azurerm_cdn_frontdoor_endpoint` on the referenced profile, with its generated public hostname

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureFrontDoorProfile** to create the endpoint in (referenced through `profileId`)

## Quick Start

Create a file `endpoint.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorEndpoint
metadata:
  name: web-endpoint
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureFrontDoorEndpoint.web-endpoint
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: my-front-door
      fieldPath: status.outputs.profile_id
  endpointName: web
```

Deploy:

```shell
planton apply -f endpoint.yaml
```

The endpoint name becomes the prefix of the public hostname, so pick something recognizable -- and treat it as permanent: renaming replaces the endpoint and changes the hostname, breaking every DNS record pointing at the old one. Set `enabled: false` to provision dark and flip the switch at launch.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `endpoint_id` | The ARM id -- what AzureFrontDoorRoute references as its parent |
| `endpoint_name` | The endpoint's name inside the profile |
| `host_name` | The generated `*.azurefd.net` hostname -- the CNAME target for custom-domain DNS records |

## Related Resources

- [Azure Front Door Profile](/docs/catalog/azure/front-door-profile) -- the parent profile
- [Azure Front Door Route](/docs/catalog/azure/front-door-route) -- attaches traffic rules to this endpoint
- [Azure DNS Record](/docs/catalog/azure/dns-record) -- CNAMEs custom domains onto the endpoint hostname
