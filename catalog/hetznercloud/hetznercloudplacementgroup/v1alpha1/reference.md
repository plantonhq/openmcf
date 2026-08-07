# HetznerCloudPlacementGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1alpha1`

HetznerCloudPlacementGroupSpec defines the specification for a Hetzner Cloud
placement group.

A placement group controls the physical distribution of servers across
Hetzner Cloud's infrastructure. Servers assigned to a spread placement group
are guaranteed to run on different physical hosts, providing fault tolerance
for high-availability workloads.

Other components (e.g., HetznerCloudServer) reference the placement group
via its placement_group_id output through StringValueOrRef.

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1alpha1
kind: HetznerCloudPlacementGroup
metadata:
  name: hetznercloudplacementgroup-demo
spec: {}
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.type` | `enum` |  | `spread` |  |

## Field Details

### spec.type

`enum` · optional (explicit presence)

Placement group strategy type.

Hetzner Cloud currently supports only "spread", which distributes
servers across different physical hosts for fault tolerance.

Default: spread

- default: `spread`

Allowed values (use exactly as shown):

- `type_unspecified`
- `spread`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudPlacementGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.placement_group_id` | `string` | The Hetzner Cloud numeric ID of the created placement group (as a string). Referenced by HetznerCloudServer.placement_group_id via StringValueOrRef. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| HetznerCloudServer | `spec.placementGroupId` | `status.outputs.placement_group_id` |

## See Also

- [Overview](../README.md)
