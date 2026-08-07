# Hetzner Cloud Snapshot

Captures a point-in-time disk image from a Hetzner Cloud server. Snapshots are stored as Hetzner Cloud Images and can be used to create new servers from the captured state. They persist independently of the source server -- deleting the server does not remove existing snapshots. Useful for golden image creation, pre-upgrade baselines, and server cloning.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Snapshot** -- a single `hcloud_snapshot` resource that captures the disk state of the referenced server as a Hetzner Cloud Image

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and API token.
- **An existing server** -- provide the server ID directly or reference a HetznerCloudServer resource via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud Snapshot**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudSnapshot
metadata:
  name: pre-upgrade-baseline
  org: acme-corp
  env: prod
spec:
  serverId:
    value: "12345678"
  description: "Baseline before v2.0 upgrade"
```

```shell
planton apply -f hetznercloud-snapshot.yaml
```

This captures the disk state of the specified server. A Stack Job tracks the provisioning in real time. The resulting snapshot ID can be used as the `image` field in a new HetznerCloudServer manifest.

### InfraChart

When deploying as part of a backup workflow, use ValueFromRef to reference the source server:

```yaml
spec:
  serverId:
    valueFrom:
      kind: HetznerCloudServer
      name: app-server
      fieldPath: status.outputs.server_id
  description: "Automated pre-deployment snapshot"
```

The InfraPipeline resolves the dependency graph and creates the server before capturing the snapshot.

## Key Configuration

These are the most important decisions when configuring a snapshot. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Server ID** -- The `serverId` field references the server to snapshot. Accepts a literal server ID or a ValueFromRef reference to a HetznerCloudServer output. Changing this value forces replacement of the snapshot.

**Description** -- The `description` field provides a human-readable label for identifying the snapshot's purpose. Can be updated after creation without replacing the snapshot.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **HetznerCloudServer** | `serverId` | `status.outputs.server_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `snapshot_id` | Hetzner Cloud image ID of the snapshot | HetznerCloudServer `image` field for booting from this snapshot |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Pre-upgrade baseline** -- Capture a snapshot before a major upgrade so the server can be restored to a known good state if the upgrade fails.

**Golden image** -- Create a snapshot of a fully configured server to use as a template for cloning new servers with identical configuration.

## Works With

- [**Hetzner Cloud Server**](/cloud-catalog/hetznercloud-server) -- the source server for the snapshot, and the consumer of snapshot IDs via the `image` field
