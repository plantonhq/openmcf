---
title: "Block Volume"
description: "Block Volume deployment documentation"
icon: "package"
order: 100
componentName: "ociblockvolume"
---

# Block Volume on OCI

Deploys an Oracle Cloud Infrastructure Block Volume -- a high-performance, durable block storage device that can be attached to compute instances. Supports configurable performance tiers (VPUs/GB), automatic performance tuning via autotune policies, cross-region replication for disaster recovery, and optional backup policy assignment for scheduled backups. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Block Volume** -- a `core.Volume` in the specified compartment and availability domain with configurable size, performance tier, encryption, autotune policies, and cross-region replicas
- **Backup Policy Assignment** -- created only when `backupPolicyId` is set; a `core.VolumeBackupPolicyAssignment` linking the volume to an Oracle-defined or custom backup policy for scheduled backups
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the volume

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the volume in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- An availability domain name (e.g., `Uocm:US-ASHBURN-AD-1`). The volume and any attached compute instance must be in the same AD. Changing the AD forces recreation.
- For customer-managed encryption: a KMS master encryption key OCID. When omitted, Oracle-managed keys are used.
- For scheduled backups: an OCI backup policy OCID (Oracle-defined Gold/Silver/Bronze or custom).

## Deploy

### Console

Open the deployment store, find **Block Volume on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Balanced with Backup** preset in the [Presets](#presets) tab to pre-populate a cost-effective volume with scheduled backups and detach-based autotune.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciBlockVolume
metadata:
  name: data-volume
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  availabilityDomain: "Uocm:US-ASHBURN-AD-1"
  sizeInGbs: 100
  vpusPerGb: 10
```

```shell
planton apply -f block-volume.yaml
```

This creates a 100 GB volume at the Balanced performance tier (60 IOPS/GB, 480 KB/s/GB). No backup policy, encryption key, autotune policies, or cross-region replicas are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the volume to a compartment deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: workloads
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartment first, then provisions the block volume with the resolved compartment OCID.

## Key Configuration

These are the most important decisions when configuring a block volume. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Performance tier** -- Set `vpusPerGb` to control IOPS and throughput. `0` = Lower Cost (2 IOPS/GB), `10` = Balanced (60 IOPS/GB, the OCI default), `20` = Higher Performance (75 IOPS/GB), `30-120` = Ultra High Performance (in increments of 10). Higher tiers cost proportionally more. Use Balanced for general workloads; Higher Performance or Ultra High for databases and latency-sensitive applications.

**Autotune policies** -- Add entries to `autotunePolicies` for automatic performance adjustment. `detached_volume` drops VPUs to 0 (Lower Cost) when the volume is detached and restores the original tier on reattach -- saves cost on infrequently used volumes. `performance_based` dynamically adjusts VPUs up to `maxVpusPerGb` based on workload demand.

**Size** -- Set `sizeInGbs` between 50 and 32768 (50 GB to 32 TB). Must be specified explicitly to avoid OCI's 1 TB default. Size can be increased after creation but not decreased.

**Cross-region replicas** -- Add entries to `blockVolumeReplicas` for asynchronous cross-region replication. Each replica is placed in a target availability domain (can be in a different region) and optionally encrypted with a region-specific KMS key. Used for disaster recovery.

**Backup policy** -- Set `backupPolicyId` to an Oracle-defined policy (Gold = daily + weekly + monthly, Silver = weekly + monthly, Bronze = monthly) or a custom policy for scheduled backups. The policy assignment is created as a separate sub-resource.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `volume_id` | OCID of the block volume | Compute instance volume attachments, volume group membership |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Balanced with backup** -- A 100 GB Balanced-tier volume with a backup policy and detach-based autotune. Cost-effective for general workloads with automated protection. Start from the **Balanced with Backup** preset.

**High-performance encrypted** -- A 200 GB Higher Performance volume with customer-managed KMS encryption, performance-based autotune (up to 40 VPUs/GB), and a cross-region replica for disaster recovery. Designed for database and latency-sensitive workloads. Start from the **High-Performance Encrypted** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this block volume