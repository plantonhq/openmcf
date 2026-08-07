# Hetzner Cloud Placement Group

Creates a placement group that controls the physical distribution of servers across Hetzner Cloud infrastructure. Servers assigned to a spread placement group are guaranteed to run on different physical hosts, providing hardware-level fault tolerance for high-availability workloads. Hetzner Cloud currently supports only the "spread" strategy, with a maximum of 10 servers per placement group.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Placement Group** -- a single `hcloud_placement_group` resource with the spread strategy, ensuring assigned servers run on separate physical hosts

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and API token.

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud Placement Group**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudPlacementGroup
metadata:
  name: ha-group
  org: acme-corp
  env: prod
spec:
  type: spread
```

```shell
planton apply -f hetznercloud-placement-group.yaml
```

This creates a spread placement group. A Stack Job tracks the provisioning in real time. Reference the group in HetznerCloudServer manifests via `placement_group_id`.

### InfraChart

When deploying as part of an HA server cluster, use ValueFromRef to wire servers to this placement group:

```yaml
# In the HetznerCloudServer manifest:
spec:
  placementGroupId:
    valueFrom:
      kind: HetznerCloudPlacementGroup
      name: ha-group
      fieldPath: status.outputs.placement_group_id
```

The InfraPipeline resolves the dependency graph, creates the placement group first, then provisions servers with anti-affinity guarantees.

## Key Configuration

These are the most important decisions when configuring a placement group. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Type** -- The `type` field selects the placement strategy. Hetzner Cloud currently supports only `spread`, which distributes servers across distinct physical hosts. The type defaults to `spread` and cannot be changed after creation.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `placement_group_id` | Hetzner Cloud numeric ID of the placement group | HetznerCloudServer `placementGroupId` field |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HA server cluster** -- Create a spread placement group before deploying multiple servers that must tolerate single-host failure. Combine with a load balancer for automatic failover.

## Works With

- [**Hetzner Cloud Server**](/cloud-catalog/hetznercloud-server) -- servers reference this placement group via `placementGroupId` for anti-affinity guarantees
