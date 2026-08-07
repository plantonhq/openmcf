# CivoBucket

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `civo.planton.dev/v1alpha1`

CivoObjectStoreBucketSpec defines the user configuration for a Civo object storage bucket.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.bucketName` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.versioningEnabled` | `bool` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |

## Field Details

### spec.bucketName

`string` · required

bucket name (DNS-compatible, 3–63 chars)

- rule: {"required":true,"string":{"minLen":"3","maxLen":"63","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"}}

### spec.region

`enum` · required

region code for the bucket

- rule: {"required":true}

Allowed values (use exactly as shown):

- `civo_region_unspecified` -- 0: default / unspecified region
- `lon1` -- london 1
- `lon2` -- london 2
- `fra1` -- frankfurt 1
- `nyc1` -- new york 1
- `phx1` -- phoenix 1
- `mum1` -- mumbai 1

### spec.versioningEnabled

`bool`

enable versioning for the bucket (disabled by default)

### spec.tags

`[]string`

tags to apply to the bucket (must be unique)

- rule: {"repeated":{"unique":true}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CivoBucket, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.bucket_id` | `string` | Unique identifier for the bucket (UUID format) |
| `status.outputs.endpoint_url` | `string` | Endpoint URL for the bucket (e.g., "https://objectstore.civo.com/<bucket-name>") |
| `status.outputs.access_key_secret_ref` | `string` | Reference to the secret storing the access key ID for the bucket |
| `status.outputs.secret_key_secret_ref` | `string` | Reference to the secret storing the secret key for the bucket |

## See Also

- [Overview](../README.md)
