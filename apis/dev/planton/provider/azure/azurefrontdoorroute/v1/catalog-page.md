# Azure Front Door Route

Creates a route inside an AzureFrontDoorEndpoint -- the rule that matches client requests by URL pattern and forwards them to an origin group, with protocol policy and edge caching. Routes are the glue of the Front Door graph: endpoint -> route -> origin group -> origins.

## What Gets Created

When you deploy an AzureFrontDoorRoute resource, Planton provisions:

- **Front Door Route** -- an `azurerm_cdn_frontdoor_route` on the referenced endpoint, forwarding to the referenced origin group

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureFrontDoorEndpoint** to attach to (referenced through `endpointId`)
- **An AzureFrontDoorOriginGroup with at least one origin** -- ARM rejects a route whose origin group is empty (list the origins in `originIds` when deploying the whole chain together)

## Quick Start

Create a file `route.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorRoute
metadata:
  name: default-route
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureFrontDoorRoute.default-route
spec:
  endpointId:
    valueFrom:
      kind: AzureFrontDoorEndpoint
      name: web-endpoint
      fieldPath: status.outputs.endpoint_id
  routeName: default
  originGroupId:
    valueFrom:
      kind: AzureFrontDoorOriginGroup
      name: api-backends
      fieldPath: status.outputs.origin_group_id
  originIds:
    - valueFrom:
        kind: AzureFrontDoorOrigin
        name: primary-app
        fieldPath: status.outputs.origin_id
  patternsToMatch:
    - /*
  supportedProtocols:
    - HTTP
    - HTTPS
```

Deploy:

```shell
planton apply -f route.yaml
```

This serves everything on the endpoint's hostname with the default HTTPS redirect (HTTP arrives only to be 301'd). Add a `cache` block for static content -- its absence means caching is off, a deliberate switch, not a defaults bundle. Front Door picks the most specific pattern across an endpoint's routes, so "/api/*" and "/*" cleanly split traffic.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `route_id` | The ARM id of the route |
| `route_name` | The route's name inside its endpoint |

The client-facing hostname deliberately lives on the ENDPOINT's outputs, not here -- the route is policy attached to that hostname.

## Related Resources

- [Azure Front Door Endpoint](/docs/catalog/azure/azurefrontdoorendpoint) -- the parent entry hostname
- [Azure Front Door Origin Group](/docs/catalog/azure/azurefrontdoororigingroup) -- the destination pool
- [Azure Front Door Origin](/docs/catalog/azure/azurefrontdoororigin) -- the backends the pool serves
