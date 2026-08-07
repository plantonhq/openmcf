# OpenStack Volume Attach

Attaches an OpenStack Cinder volume to a compute instance, making the volume appear as a block device (e.g., `/dev/vdb`) inside the instance. This is a join resource that connects two independently managed resources -- a volume and an instance -- and makes the attachment relationship explicit in the deployment graph. Supports ValueFromRef for wiring both instance and volume dependencies in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Volume Attachment** -- an `openstack_compute_volume_attach_v2` resource that connects a Cinder volume to a compute instance via the Nova API. The volume transitions from "available" to "in-use" state, and the hypervisor presents the block device to the instance at the specified or auto-assigned device path.

## Before You Deploy

### OpenStack Account

- **A compute instance** in running state -- provide the instance UUID directly or reference an OpenStackInstance Cloud Resource via ValueFromRef.
- **A Cinder volume** in "available" state -- provide the volume UUID directly or reference an OpenStackVolume Cloud Resource via ValueFromRef.
- **Same availability zone** for both the volume and instance, as required by most OpenStack deployments.

## Deploy

### Console

Open the deployment store, find **OpenStack Volume Attach**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab for a basic volume attachment.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackVolumeAttach
metadata:
  name: app-data-attach
  org: acme-corp
  env: prod
spec:
  instanceId:
    value: "3b4c5d6e-7f8a-9b0c-1d2e-3f4a5b6c7d8e"
  volumeId:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
```

```shell
planton apply -f openstack-volume-attach.yaml
```

This attaches the specified Cinder volume to the compute instance. Nova automatically assigns the next available device path (e.g., `/dev/vdb`). No explicit device path is requested.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the attachment to an instance and volume deployed in the same InfraPipeline:

```yaml
spec:
  instanceId:
    valueFrom:
      kind: OpenStackInstance
      name: app-server
      fieldPath: status.outputs.instance_id
  volumeId:
    valueFrom:
      kind: OpenStackVolume
      name: app-data
      fieldPath: status.outputs.volume_id
```

The InfraPipeline resolves the dependency graph, deploys the instance and volume first, then provisions the attachment with the resolved IDs.

## Key Configuration

These are the most important decisions when configuring an OpenStack volume attachment. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Instance and volume references** -- The `instanceId` and `volumeId` fields are both required. Each accepts a literal UUID or a ValueFromRef reference to another Cloud Resource. Using ValueFromRef creates an explicit dependency edge in InfraCharts, ensuring the instance and volume exist before the attachment is created.

**Device path** -- The `device` field optionally specifies where the volume appears inside the instance (e.g., `/dev/vdb`, `/dev/vdc`). If omitted, Nova selects the next available device. Specify a device path when your application expects the disk at a known location.

**Immutability** -- All fields on this resource are ForceNew. Changing `instanceId`, `volumeId`, `device`, or `region` recreates the attachment (detach + reattach). Plan attachment configuration carefully before deploying.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackInstance** | `instanceId` | `status.outputs.instance_id` |
| **OpenStackVolume** | `volumeId` | `status.outputs.volume_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Resource identifier for the attachment (`{instance_id}/{volume_id}`) | Audit logs, lifecycle tracking |
| `instance_id` | UUID of the compute instance the volume is attached to | Operational dashboards |
| `volume_id` | UUID of the Cinder volume that was attached | Operational dashboards |
| `device` | Device path where the volume appears in the instance | Application configuration, mount scripts |
| `region` | OpenStack region where the attachment was created | Multi-region deployment coordination |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard attachment** -- Connects a volume to an instance with auto-assigned device path. Nova picks the next available device. Suitable for most data volume use cases where device path is not significant. Start from the **Standard** preset.

## Works With

- [**OpenStack Instance**](/cloud-catalog/openstack-instance) -- the compute instance that the volume attaches to
- [**OpenStack Volume**](/cloud-catalog/openstack-volume) -- the Cinder block storage volume being attached