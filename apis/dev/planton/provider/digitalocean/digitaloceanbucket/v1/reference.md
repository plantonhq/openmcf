# DigitalOceanBucket

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1`

DigitalOceanBucketSpec defines the user configuration for a DigitalOcean Spaces bucket.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.bucketName` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.accessControl` | `enum` |  |  |  |
| `spec.versioningEnabled` | `bool` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |

## Field Details

### spec.bucketName

`string` · required

bucket name (DNS-compatible, 3–63 chars)

- rule: {"required":true,"string":{"minLen":"3","maxLen":"63","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"}}

### spec.region

`enum` · required

region slug (datacenter location for the bucket)

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

### spec.accessControl

`enum`

access control setting for the bucket (private or public-read)

Allowed values (use exactly as shown):

- `PRIVATE`
- `PUBLIC_READ`

### spec.versioningEnabled

`bool`

enable versioning for the bucket (disabled by default)

### spec.tags

`[]string`

tags to apply to the bucket (must be unique)

- rule: {"repeated":{"unique":true}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanBucket, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.bucket_id` | `string` | Unique identifier for the bucket (UUID format) |
| `status.outputs.endpoint` | `string` | Regional endpoint URL for the bucket (e.g., "https://<region>.digitaloceanspaces.com") |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
