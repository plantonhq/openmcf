# AzureFrontDoorProfile

An Azure Front Door (Standard/Premium) profile: the container for a
global CDN and application-delivery deployment on Microsoft's edge
network. The profile owns the SKU tier, the origin response timeout,
the managed identity, access-log scrubbing, and tags -- the delivery
surface composes against it from first-class resources:

- **AzureFrontDoorEndpoint** -- the public entry hostname (*.azurefd.net)
- **AzureFrontDoorOriginGroup** -- a load-balanced backend pool
- **AzureFrontDoorOrigin** -- one backend inside an origin group
- **AzureFrontDoorRoute** -- connects an endpoint to an origin group by
  URL pattern

This mirrors Azure's own resource model (each is a separate ARM child
resource with an independent lifecycle) and keeps regional stamps
composable: a new region adds its origin to a shared group without
touching the profile or any other region's resources.

## When to Use

Use AzureFrontDoorProfile when you need:

- **Global HTTP(S) delivery** -- edge caching, TLS at the edge,
  latency-aware load balancing across regions, instant failover
- **The PREMIUM tier** -- Private Link to origins (backends with public
  access disabled) and the managed WAF rule sets
- **The anchor of a Front Door composition** -- every endpoint, origin
  group, and route references this profile's outputs

## Key Configuration

- `resource_group` -- referenced from an AzureResourceGroup's output;
  Front Door is global (no region) but lives in a resource group for
  ARM organization
- `profile_name` -- 2-90 characters, unique within the resource group;
  ForceNew, and it replaces everything nested under the profile
- `sku` -- STANDARD (default) or PREMIUM; ForceNew, and Azure refuses a
  PREMIUM -> STANDARD downgrade outright
- `response_timeout_seconds` -- 16-240, default 120; how long Front
  Door waits for the origin before returning 504
- `identity` -- system/user-assigned managed identity for keyless Key
  Vault certificate access (bring-your-own TLS on custom domains)
- `log_scrubbing_variables` -- request parts masked out of access logs
  (client IP, request URI, query-string arguments); presence enables
  scrubbing
- `tags` -- merged over the Planton-derived resource tags; user tags win

## Composition

```yaml
resourceGroup:
  valueFrom:
    kind: AzureResourceGroup
    name: my-rg
    fieldPath: status.outputs.resource_group_name
```

Delivery resources reference the profile through its `profile_id`
output; the `identity_principal_id` output is the principal to grant
Key Vault access to for bring-your-own certificates.

## Documentation

- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
