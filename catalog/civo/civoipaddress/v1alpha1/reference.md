# CivoIpAddress

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `civo.planton.dev/v1alpha1`

CivoIpAddressSpec defines the specification for provisioning a static reserved IP in Civo Cloud.
A reserved IP is a persistent public IPv4 address that can be attached to an instance or load balancer in the same region.
This spec focuses on the essential fields, following the 80/20 principle: specifying the region and an optional descriptive label.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.description` | `string` |  |  |  |
| `spec.region` | `enum` | yes |  |  |

## Field Details

### spec.description

`string`

An optional human-readable name or description for the reserved IP.
If not provided, Civo may default to using the IP address as the name.

- rule: {"string":{"maxLen":"100"}}

### spec.region

`enum` · required

Reserved IPs are region-specific; the IP can only be attached to resources in this region.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `civo_region_unspecified` -- 0: default / unspecified region
- `lon1` -- london 1
- `lon2` -- london 2
- `fra1` -- frankfurt 1
- `nyc1` -- new york 1
- `phx1` -- phoenix 1
- `mum1` -- mumbai 1

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CivoIpAddress, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.reserved_ip_id` | `string` | The unique identifier of the reserved IP in Civo (UUID). |
| `status.outputs.ip_address` | `string` | The static IPv4 address allocated for this reserved IP. |
| `status.outputs.attached_resource_id` | `string` | The ID of the Civo resource (instance or load balancer) currently attached to this IP, if any. Will be empty if the reserved IP is not attached to any resource. |
| `status.outputs.created_at_rfc3339` | `string` | The timestamp when the reserved IP was created, in RFC 3339 format. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| CivoComputeInstance | `spec.reservedIpId` | `status.outputs.reserved_ip_id` |

## See Also

- [Overview](../README.md)
