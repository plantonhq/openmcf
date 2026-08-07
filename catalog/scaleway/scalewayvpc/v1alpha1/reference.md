# ScalewayVpc

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayVpcSpec defines the specification for a Scaleway Virtual Private Cloud (VPC).

A Scaleway VPC is a regional, logical container that groups Private Networks.
Unlike VPCs in some other cloud providers, Scaleway VPCs do not define IP ranges
or CIDR blocks -- IP planning happens at the Private Network level.

The VPC's primary purpose is to provide network isolation and, when routing is enabled,
allow communication between Private Networks attached to the same VPC.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.enableRouting` | `bool` |  | `false` |  |
| `spec.enableCustomRoutesPropagation` | `bool` |  | `false` |  |

## Field Details

### spec.region

`string` · required

The Scaleway region where the VPC will be created.
Examples: "fr-par", "nl-ams", "pl-waw"

This field is required and cannot be changed after creation.

- rule: {"required":true}

### spec.enableRouting

`bool`

Whether to enable routing between Private Networks in this VPC.

When enabled, resources in different Private Networks attached to this VPC
can communicate with each other. This is essential for multi-tier architectures
(e.g., a Kapsule cluster in one Private Network talking to an RDB instance
in another).

IMPORTANT: Once enabled, routing cannot be disabled. Plan accordingly.

Default: false

- default: `false`

### spec.enableCustomRoutesPropagation

`bool`

Whether to enable custom routes propagation between Private Networks in this VPC.

When enabled, custom routes from one Private Network are advertised to other
Private Networks in the same VPC. This is useful for advanced networking
scenarios such as VPN gateways or network appliances.

IMPORTANT: Once enabled, custom routes propagation cannot be disabled. Plan accordingly.

Default: false

- default: `false`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayVpc, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vpc_id` | `string` | The unique identifier (UUID) of the created Scaleway VPC. This is the primary output referenced by downstream resources (e.g., ScalewayPrivateNetwork) via StringValueOrRef. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| ScalewayPrivateNetwork | `spec.vpcId` | `status.outputs.vpc_id` |

## See Also

- [Overview](./README.md)
