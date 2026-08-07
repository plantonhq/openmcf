# AliCloud NAS File System

Deploys an Alibaba Cloud Network Attached Storage (NAS) file system with a bundled access group and VPC mount target. NAS provides fully managed, elastic, shared file storage supporting NFS and SMB protocols. The file system, access group, and mount target are deployed as a single atomic unit because a file system without a mount target is unreachable from any VPC resource.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **NAS File System** -- an `alicloud_nas_file_system` resource with the chosen protocol, storage type, and optional encryption
- **Access Group** (conditional) -- an `alicloud_nas_access_group` with `alicloud_nas_access_rule` entries when `accessRules` are specified. Controls which IP ranges can mount and with what permissions
- **Mount Target** -- an `alicloud_nas_mount_target` in the specified VPC/VSwitch, associated with the custom access group or the default VPC group

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **An existing VPC and VSwitch** -- the mount target requires placement in a VSwitch within a VPC.
- **Protocol selection** -- `protocolType` (NFS or SMB) is immutable after creation. NFS for Linux/Unix workloads, SMB for Windows.
- **Storage type selection** -- `storageType` is immutable. For standard NAS: Performance (SSD), Capacity (HDD), or Premium. For extreme NAS: standard or advance.

## Deploy

### Console

Open the deployment store, find **AliCloud NAS File System**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including protocol, storage type, VPC, and access rules.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudNasFileSystem
metadata:
  name: app-shared-fs
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  protocolType: NFS
  storageType: Performance
  description: Shared file storage for application workloads
  vpcId:
    value: vpc-bp1234567890
  vswitchId:
    value: vsw-app-zone-a
  accessRules:
    - sourceCidrIp: "10.0.0.0/8"
      rwAccessType: RDWR
      userAccessType: root_squash
```

```shell
planton apply -f alicloud-nas.yaml
```

This creates a Performance NFS file system with a mount target and access rule restricting root squash. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of an application stack, use ValueFromRef to wire VPC and VSwitch dependencies:

```yaml
spec:
  region: cn-hangzhou
  protocolType: NFS
  storageType: Performance
  vpcId:
    valueFrom:
      kind: AliCloudVpc
      name: platform-vpc
      fieldPath: status.outputs.vpc_id
  vswitchId:
    valueFrom:
      kind: AliCloudVswitch
      name: app-vswitch-a
      fieldPath: status.outputs.vswitch_id
```

The InfraPipeline resolves the dependency graph and provisions VPC and VSwitch before the NAS file system.

## Key Configuration

These are the most important decisions when configuring a NAS file system. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**File system type** -- The `fileSystemType` field selects "standard" (default, auto-scaling capacity) or "extreme" (dedicated throughput, fixed capacity, requires `zoneId` and `capacity`).

**Protocol** -- The `protocolType` field is immutable. "NFS" for Linux/Unix. "SMB" for Windows.

**Storage type** -- The `storageType` field is immutable. "Performance" (SSD) for low latency, "Capacity" (HDD) for cost-effective warm storage, "Premium" for next-gen SSD.

**Access rules** -- When `accessRules` are specified, a custom access group is auto-created. Each rule controls which CIDR ranges can mount, with what read/write permissions, and with what user identity mapping (no_squash, root_squash, all_squash).

**Encryption** -- The `encryption` block enables at-rest encryption. Type 1 uses NAS-managed keys. Type 2 uses a customer-managed KMS key (create with AliCloudKmsKey).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVpc** | `vpcId` | `status.outputs.vpc_id` |
| **AliCloudVswitch** | `vswitchId` | `status.outputs.vswitch_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `file_system_id` | NAS file system ID (e.g., 1ca404a348) | Snapshot policies, lifecycle policies |
| `mount_target_domain` | Mount endpoint domain name | NFS/SMB mount commands from ECS instances and containers |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard NFS** -- A Performance NFS file system with default VPC access. The simplest starting point. Start from the **Standard NFS** preset.

**Production encrypted** -- A Performance NFS file system with NAS-managed encryption and restricted access rules with root squash. Start from the **Production Encrypted** preset.

## Works With

- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- the VPC the mount target is created in
- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- the VSwitch for mount target placement
- [**AliCloud KMS Key**](/cloud-catalog/ali-cloud-kms-key) -- customer-managed encryption key (when encryption type is 2)
- [**AliCloud ECS Instance**](/cloud-catalog/ali-cloud-ecs-instance) -- compute instances that mount the file system
- [**AliCloud Kubernetes Cluster**](/cloud-catalog/ali-cloud-kubernetes-cluster) -- Kubernetes pods that mount the file system via PV/PVC
