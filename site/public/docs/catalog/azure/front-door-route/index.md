---
title: "Front Door Route"
description: "Front Door Route deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoorroute"
---

# Azure Front Door Route

Deploys a Front Door route -- the URL-pattern glue that connects an endpoint to an origin group, with supported protocols, forwarding protocol, HTTPS redirect, custom domains, and optional edge caching. `originIds` exists only for deploy sequencing (ARM rejects a route whose origin group is empty); Azure never receives the list. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring across the Front Door serving spine.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Front Door Route** -- a named child of an endpoint matching URL patterns and forwarding to an origin group
- **Protocol posture** -- supported protocols (HTTP/HTTPS), forwarding protocol, and HTTPS redirect
- **Domain binding** -- link to the default `*.azurefd.net` hostname and/or custom domains
- **Edge cache** (optional) -- query-string caching behavior, compression, and content types to compress
- **Rule sets** (optional) -- references to AzureFrontDoorRuleSet for request/response transforms

## The Route in the Front Door Family

- **AzureFrontDoorEndpoint** -- the public entry, referenced by `endpointId` (ForceNew)
- **AzureFrontDoorOriginGroup** -- the backend pool, referenced by `originGroupId` (editable in place)
- **AzureFrontDoorOrigin** -- listed in `originIds` so the route waits for backends to exist
- **AzureFrontDoorCustomDomain** / **AzureFrontDoorRuleSet** -- optional attachments (slice-two kinds)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Azure Subscription

- **An endpoint and an origin group** under the same Front Door profile, with at least one origin in the group before the route applies (or list them in `originIds` for sequencing).

## Deploy

### Console

Open the deployment store, find **Azure Front Door Route**, and click **Deploy**. The wizard walks you through attachment (endpoint, origin group, sequencing origins), matching and protocols, domains, and optional edge caching. Start from the **Catch-All HTTPS** preset for the default route shape.

### CLI

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorRoute
metadata:
  name: my-default-route
  org: acme-corp
  env: prod
spec:
  endpointId:
    valueFrom:
      kind: AzureFrontDoorEndpoint
      name: my-web-endpoint
      fieldPath: status.outputs.endpoint_id
  routeName: default-route
  originGroupId:
    valueFrom:
      kind: AzureFrontDoorOriginGroup
      name: my-app-backends
      fieldPath: status.outputs.origin_group_id
  originIds:
    - valueFrom:
        kind: AzureFrontDoorOrigin
        name: my-app-origin
        fieldPath: status.outputs.origin_id
  patternsToMatch:
    - /*
  supportedProtocols:
    - HTTP
    - HTTPS
```

```shell
planton apply -f front-door-route.yaml
```

This creates a catch-all route with both protocols and the default HTTPS redirect (HTTP arrives only to be 301'd to HTTPS). Caching stays off until a `cache` block is declared.

### InfraChart

```yaml
spec:
  endpointId:
    valueFrom:
      kind: AzureFrontDoorEndpoint
      name: my-web-endpoint
      fieldPath: status.outputs.endpoint_id
  routeName: static-assets
  originGroupId:
    valueFrom:
      kind: AzureFrontDoorOriginGroup
      name: my-app-backends
      fieldPath: status.outputs.origin_group_id
  patternsToMatch:
    - /assets/*
  supportedProtocols:
    - HTTPS
  httpsRedirectEnabled: false
  cache:
    queryStringCachingBehavior: IGNORE_QUERY_STRING
    compressionEnabled: true
```

## Key Configuration

**Patterns** -- at least one; each must start with `/`. `/*` is the catch-all.

**HTTPS redirect** -- requires both HTTP and HTTPS in `supportedProtocols` (the redirect needs HTTP to arrive and HTTPS to land on). Disable the redirect, or list both.

**Default domain vs custom domains** -- `linkToDefaultDomain` can only be false when at least one custom domain is attached; otherwise the route would answer on no hostname.

**Cache** -- absence means caching off. When present, choose query-string behavior and optionally enable compression with Front Door's supported content-type list.

**originIds** -- deploy sequencing only. List origins so the route provisions after the group is non-empty; Azure never sees this field.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureFrontDoorEndpoint** | `endpointId` | `status.outputs.endpoint_id` |
| **AzureFrontDoorOriginGroup** | `originGroupId` | `status.outputs.origin_group_id` |
| **AzureFrontDoorOrigin** | `originIds[]` | `status.outputs.origin_id` (sequencing) |
| **AzureFrontDoorCustomDomain** | `customDomainIds[]` | `status.outputs.custom_domain_id` |
| **AzureFrontDoorRuleSet** | `ruleSetIds[]` | `status.outputs.rule_set_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `route_id` | ARM resource ID of the route | Operator tooling / diagnostics |
| `route_name` | The route's name | Operator tooling |

## Presets

| Preset | Rank | Description |
|--------|------|-------------|
| Catch-All HTTPS | 1 | `/*` with both protocols and HTTPS redirect |
| Static Assets Cached | 2 | Path-scoped route with edge caching and compression |
| API HTTPS-Only | 3 | HTTPS-only API path without redirect |

## Related Components

- **AzureFrontDoorEndpoint** -- the public entry
- **AzureFrontDoorOriginGroup** / **AzureFrontDoorOrigin** -- the backend pool
- **AzureFrontDoorCustomDomain** -- branded hostnames (slice two)
- **AzureFrontDoorRuleSet** -- request/response transforms (slice two)
