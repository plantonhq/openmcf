# DigitalOceanVpc

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1`

DigitalOceanVpcSpec defines the specification required to deploy a DigitalOcean Virtual Private Cloud (VPC).
A DigitalOcean VPC allows you to create a private, isolated network for your Droplets and other resources,
enabling secure communication within your infrastructure.
This specification focuses on the essential parameters for creating a VPC, adhering to the 80/20 principle.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.description` | `string` |  |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.ipRangeCidr` | `string` |  |  |  |
| `spec.isDefaultForRegion` | `bool` |  | `false` |  |

## Field Details

### spec.description

`string`

A human-readable description for the VPC.
Constraints: Maximum 100 characters.

- rule: {"string":{"maxLen":"100"}}

### spec.region

`enum` · required

The DigitalOcean region where the VPC will be created.
This determines the geographical location of the VPC.

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

### spec.ipRangeCidr

`string`

The IP range for the VPC in CIDR notation (optional).
Only /16, /20, or /24 CIDR blocks are supported for VPCs on DigitalOcean.
Example: "10.10.0.0/16"

80/20 Principle: When omitted, DigitalOcean auto-generates a non-conflicting /20 CIDR block (4,096 IPs).
This is the recommended approach for dev/test environments and when explicit IP planning is not required.
For production environments with specific IPAM requirements, explicitly specify the CIDR block.

- rule: {"string":{"pattern":"^([0-9]{1,3}\\.){3}[0-9]{1,3}/(16|20|24)$"}}

### spec.isDefaultForRegion

`bool`

A boolean indicating whether this VPC should be set as the default for the specified region.
Only one VPC can be the default for a given region.
Default: false

- default: `false`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanVpc, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vpc_id` | `string` | The unique identifier (UUID) of the created DigitalOcean VPC. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| DigitalOceanDatabaseCluster | `spec.vpc` | `status.outputs.vpc_id` |
| DigitalOceanDroplet | `spec.vpc` | `status.outputs.vpc_id` |
| DigitalOceanKubernetesCluster | `spec.vpc` | `status.outputs.vpc_id` |
| DigitalOceanLoadBalancer | `spec.vpc` | `status.outputs.vpc_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
