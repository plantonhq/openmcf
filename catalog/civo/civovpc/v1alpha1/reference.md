# CivoVpc

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `civo.planton.dev/v1alpha1`

CivoVpcSpec defines the specification for an isolated private network on Civo.

## Example

```yaml
apiVersion: civo.planton.dev/v1alpha1
kind: CivoVpc
metadata:
  name: example-vpc
spec:
  network_name: example-vpc
  region: lon1
  ip_range_cidr: "10.10.1.0/24"
  description: "Example VPC network"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.networkName` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.ipRangeCidr` | `string` |  |  |  |
| `spec.isDefaultForRegion` | `bool` |  |  |  |
| `spec.description` | `string` |  |  |  |

## Field Details

### spec.networkName

`string` · required

The name of the network (DNS-friendly label).

- rule: {"required":true}

### spec.region

`enum` · required

The Civo region where this network will be created.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `civo_region_unspecified` -- 0: default / unspecified region
- `lon1` -- london 1
- `lon2` -- london 2
- `fra1` -- frankfurt 1
- `nyc1` -- new york 1
- `phx1` -- phoenix 1
- `mum1` -- mumbai 1

### spec.ipRangeCidr

`string`

The IPv4 CIDR range for the network (max /24). If omitted, an available range will be auto-allocated.

### spec.isDefaultForRegion

`bool`

Whether this network should be the default for the region (only one default network per region).

### spec.description

`string`

An optional description for the network (up to 100 characters).

- rule: {"string":{"maxLen":"100"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CivoVpc, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.network_id` | `string` | the unique ID of the network created on Civo |
| `status.outputs.cidr_block` | `string` | the IPv4 CIDR block of the created network |
| `status.outputs.created_at_rfc3339` | `string` | timestamp when the network was created (RFC 3339 format) |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| CivoComputeInstance | `spec.network` | `status.outputs.network_id` |
| CivoDatabase | `spec.networkId` | `status.outputs.network_id` |
| CivoFirewall | `spec.networkId` | `status.outputs.network_id` |
| CivoKubernetesCluster | `spec.network` | `status.outputs.network_id` |

## See Also

- [Overview](../README.md)
