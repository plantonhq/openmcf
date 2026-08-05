# HetznerCloudSnapshot

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1`

HetznerCloudSnapshotSpec defines the specification for a Hetzner Cloud
server snapshot.

A snapshot is a point-in-time disk image captured from a running or stopped
server. Snapshots are stored as Hetzner Cloud Images (type "snapshot") and
can be used to create new servers from the captured state. Snapshots persist
independently of the source server -- deleting the server does not remove
existing snapshots.

This is a single-resource component. Changing the source server_id forces
replacement of the snapshot (the provider marks server_id as ForceNew).

Bundled provider resources:
  - hcloud_snapshot:  The server snapshot (stored as an Image).

Fields not exposed in this spec (hardcoded or derived in IaC modules):
  - labels: Derived from metadata (CG01 pattern). Standard labels take
            precedence over user-specified metadata.labels.

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudSnapshot
metadata:
  name: hetznercloudsnapshot-demo
spec:
  serverId:
    value: "12345678"
  description: demo snapshot
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.serverId` | `string \| valueFrom` | yes |  | HetznerCloudServer (`status.outputs.server_id`) |
| `spec.description` | `string` |  |  |  |

## Field Details

### spec.serverId

`string | valueFrom` · required

Server to create the snapshot from. Required.

Accepts a literal Hetzner Cloud server ID (as a string) or a reference
to a HetznerCloudServer resource's output via valueFrom.

Changing this value forces replacement of the snapshot (the existing
snapshot is destroyed and a new one is created from the new server).

Example (literal):
  serverId:
    value: "12345678"

Example (reference):
  serverId:
    valueFrom:
      kind: HetznerCloudServer
      name: my-app-server
      fieldPath: status.outputs.server_id

- references: HetznerCloudServer (`status.outputs.server_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudServer, name: <that resource's name>, fieldPath: status.outputs.server_id}} -- a bare string does not parse

### spec.description

`string`

Human-readable description of the snapshot. Optional.

Useful for identifying the purpose of the snapshot (e.g.,
"pre-upgrade baseline", "golden image v2.1"). Can be updated after
creation without replacing the snapshot.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudSnapshot, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.snapshot_id` | `string` | The Hetzner Cloud image ID of the created snapshot (as a string). Snapshots are stored as Images in the Hetzner Cloud API. This ID can be used as the `image` parameter when creating a new HetznerCloudServer to boot from this snapshot. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.serverId` | HetznerCloudServer | `status.outputs.server_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
