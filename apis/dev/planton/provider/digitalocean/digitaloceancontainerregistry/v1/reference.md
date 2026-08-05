# DigitalOceanContainerRegistry

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1`

**DigitalOceanContainerRegistrySpec** defines the configuration for creating a DigitalOcean
Container Registry (DOCR). It exposes only the essential fields needed for the common 80 % use case.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.name` | `string` | yes |  |  |
| `spec.subscriptionTier` | `enum` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.garbageCollectionEnabled` | `bool` |  |  |  |

## Field Details

### spec.name

`string` · required

Registry name (must be unique within your DigitalOcean account).
1-63 characters, lowercase letters, numbers, and hyphens; must start and end with an alphanumeric.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"}}

### spec.subscriptionTier

`enum` · required

Subscription tier slug (defines storage limits and pricing).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digitalocean_container_registry_tier_unspecified`
- `starter`
- `basic`
- `professional`

### spec.region

`enum` · required

Required region slug where registry data is stored (e.g., "nyc3", "sfo3").
DigitalOcean container registries are single-region and the region cannot
be changed after creation, so a region must be specified explicitly.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digital_ocean_region_unspecified` -- 0: default / unspecified region
- `nyc3` -- new york 3
- `sfo3` -- san francisco 3
- `fra1` -- frankfurt 1
- `sgp1` -- singapore 1
- `lon1` -- london 1
- `tor1` -- toronto 1
- `blr1` -- bangalore 1
- `ams3` -- amsterdam 3

### spec.garbageCollectionEnabled

`bool`

Enable garbage collection of untagged images.
Default is false (no automatic GC).

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanContainerRegistry, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.registry_name` | `string` | The registry name. |
| `status.outputs.server_url` | `string` | Full server URL, e.g. "registry.digitalocean.com/<registry_name>". |
| `status.outputs.region` | `string` | Region slug where the registry is hosted. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| DigitalOceanAppPlatformService | `spec.imageSource.registry` | `status.outputs.server_url` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
