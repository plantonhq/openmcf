---
title: "FSx for OpenZFS"
description: "FSx for OpenZFS deployment documentation"
icon: "package"
order: 100
componentName: "awsfsxopenzfsfilesystem"
---

# AWS FSx for OpenZFS

Deploys an Amazon FSx for OpenZFS file system with configurable NFS exports, ZSTD/LZ4 data compression, per-user/group quotas, provisioned IOPS, and optional Multi-AZ high availability with automatic failover. The component creates the file system and configures its root volume; child volumes are managed independently.

## What Gets Created

When you deploy an AwsFsxOpenzfsFileSystem resource, Planton provisions:

- **OpenZFS File System** — an `aws_fsx_openzfs_file_system` resource placed in the specified subnets with encryption at rest, tagged with Planton resource metadata
- **Root Volume** — configured inline with data compression, NFS export rules, record size tuning, and user/group quotas as specified
- **Disk IOPS Configuration** — created only when `diskIopsConfiguration` is specified; controls SSD IOPS in AUTOMATIC or USER_PROVISIONED mode
- **Multi-AZ Route Entries** — created only for MULTI_AZ_1 deployments; AWS manages routes in the specified route tables for seamless failover

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **At least one subnet** for SINGLE_AZ deployments, or **two subnets in different AZs** for MULTI_AZ_1
- **A security group** allowing NFS traffic: TCP 111 (portmapper), TCP 2049 (NFS), TCP 20001-20003 (NFS mount)
- **A KMS key ARN** if using customer-managed encryption at rest (optional — AWS-managed key used by default)
- **Route table IDs** if deploying MULTI_AZ_1 (for automatic failover routing)

## Quick Start

Create a file `openzfs.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsFsxOpenzfsFileSystem
metadata:
  name: my-openzfs
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsFsxOpenzfsFileSystem.my-openzfs
spec:
  region: us-east-1
  storageCapacityGib: 256
  throughputCapacity: 160
  subnetIds:
    - subnet-0123456789abcdef0
```

Deploy:

```shell
planton apply -f openzfs.yaml
```

This creates a SINGLE_AZ_2 OpenZFS file system with 256 GiB SSD storage, 160 MB/s throughput, no compression, and default NFS settings.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region where the FSx OpenZFS file system will be created (e.g., `us-east-1`). | Required; non-empty |
| `throughputCapacity` | `int32` | Throughput in MB/s. Value sets per deployment generation: SINGLE_AZ_1: 64, 128, 256, 512, 1024, 2048, 3072, 4096. SINGLE_AZ_2/MULTI_AZ_1: 160, 320, 640, 1280, 2560, 3840, 5120, 7680, 10240. | Generation value set (CEL-enforced) |
| `subnetIds` | `StringValueOrRef[]` | Subnet IDs. Exactly 1 for single-AZ types (incl. HA), exactly 2 for MULTI_AZ_1. Can reference AwsSubnet via `valueFrom`. | Count matches deployment type |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `deploymentType` | `string` | `SINGLE_AZ_2` | `SINGLE_AZ_1`, `SINGLE_AZ_2`, `SINGLE_AZ_HA_1`, `SINGLE_AZ_HA_2`, or `MULTI_AZ_1`. ForceNew. |
| `storageCapacityGib` | `int32` | — | 64–524288 GiB. Required for SSD storage (unless restoring from `backupId`); forbidden for `INTELLIGENT_TIERING`. Grows in place. |
| `storageType` | `string` | `SSD` | `SSD` or `INTELLIGENT_TIERING` (MULTI_AZ_1 only, elastic capacity). ForceNew. |
| `readCacheConfiguration.sizingMode` | `string` | — | INTELLIGENT_TIERING read cache: `PROPORTIONAL_TO_THROUGHPUT_CAPACITY`, `USER_PROVISIONED`, or `NO_CACHE`. Required for INTELLIGENT_TIERING. |
| `readCacheConfiguration.sizeGib` | `int32` | — | Cache size when `USER_PROVISIONED`. |
| `backupId` | `string` | — | FSx backup to restore from. ForceNew; excludes `storageCapacityGib`. |
| `securityGroupIds` | `StringValueOrRef[]` | `[]` | Security group IDs (up to 50). Can reference AwsSecurityGroup via `valueFrom`. ForceNew. |
| `preferredSubnetId` | `StringValueOrRef` | — | Active file server subnet. Required for MULTI_AZ_1, invalid otherwise. ForceNew. |
| `endpointIpAddressRange` | `string` | — | CIDR range for endpoint floating IPs. MULTI_AZ_1 only. ForceNew. |
| `routeTableIds` | `StringValueOrRef[]` | `[]` | Route tables for failover routing (up to 50). MULTI_AZ_1 only. |
| `kmsKeyId` | `StringValueOrRef` | AWS-managed | Customer-managed KMS key ARN. Can reference AwsKmsKey via `valueFrom`. ForceNew. |
| `diskIopsConfiguration.mode` | `string` | `AUTOMATIC` | `AUTOMATIC` scales with storage. `USER_PROVISIONED` uses explicit IOPS. |
| `diskIopsConfiguration.iops` | `int32` | — | Total SSD IOPS. Only when mode is `USER_PROVISIONED`. Ceilings: 160000 (SINGLE_AZ_1), 400000 (SINGLE_AZ_2). |
| `rootVolumeConfiguration.dataCompressionType` | `string` | `NONE` | `NONE`, `ZSTD` (best ratio), or `LZ4` (fastest). |
| `rootVolumeConfiguration.nfsExports.clientConfigurations` | `object[]` | — | NFS client access rules. Each entry: `clients` (IP/CIDR/wildcard) + `options` (mount options). |
| `rootVolumeConfiguration.readOnly` | `bool` | `false` | Makes the root volume read-only. |
| `rootVolumeConfiguration.recordSizeKib` | `int32` | `128` | ZFS record size: 4, 8, 16, 32, 64, 128, 256, 512, or 1024 KiB. |
| `rootVolumeConfiguration.userAndGroupQuotas` | `object[]` | — | Per-user/group storage quotas. Each: `id`, `storageCapacityQuotaGib`, `type` (USER/GROUP). |
| `rootVolumeConfiguration.copyTagsToSnapshots` | `bool` | `false` | Propagate root volume tags to snapshots. |
| `automaticBackupRetentionDays` | `int32` | `0` | Days to retain automatic backups (0–90). 0 disables. |
| `dailyAutomaticBackupStartTime` | `string` | — | Backup window in HH:MM UTC format. |
| `copyTagsToBackups` | `bool` | `false` | Propagate file system tags to backups. |
| `copyTagsToVolumes` | `bool` | `false` | Propagate file system tags to volumes. |
| `skipFinalBackup` | `bool` | `true` | Skip final backup on deletion. |
| `finalBackupTags` | `map` | — | Tags applied to the final backup when one is taken. |
| `deleteOptions` | `string[]` | `[]` | `DELETE_CHILD_VOLUMES_AND_SNAPSHOTS` opts into cascading deletion. |
| `weeklyMaintenanceStartTime` | `string` | — | Maintenance window in d:HH:MM UTC (1=Mon, 7=Sun). |

## Examples

### Production Single-AZ with Compression

A SINGLE_AZ_2 file system with ZSTD compression, NFS exports open to the VPC, encryption, and daily backups:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsFsxOpenzfsFileSystem
metadata:
  name: app-storage
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsFsxOpenzfsFileSystem.app-storage
spec:
  region: us-east-1
  deploymentType: SINGLE_AZ_2
  storageCapacityGib: 1024
  throughputCapacity: 640
  subnetIds:
    - subnet-private-az1
  securityGroupIds:
    - sg-nfs-access
  kmsKeyId: arn:aws:kms:us-east-1:123456789012:key/my-key
  rootVolumeConfiguration:
    dataCompressionType: ZSTD
    nfsExports:
      clientConfigurations:
        - clients: "*"
          options:
            - rw
            - crossmnt
            - no_root_squash
  automaticBackupRetentionDays: 7
  dailyAutomaticBackupStartTime: "05:00"
  copyTagsToBackups: true
```

### Multi-AZ High Availability

A MULTI_AZ_1 deployment with provisioned IOPS, user quotas, and two subnets across availability zones:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsFsxOpenzfsFileSystem
metadata:
  name: ha-nfs
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsFsxOpenzfsFileSystem.ha-nfs
spec:
  region: us-east-1
  deploymentType: MULTI_AZ_1
  storageCapacityGib: 2048
  throughputCapacity: 1280
  subnetIds:
    - subnet-private-az1
    - subnet-private-az2
  preferredSubnetId: subnet-private-az1
  routeTableIds:
    - rtb-private-az1
    - rtb-private-az2
  securityGroupIds:
    - sg-nfs-access
  kmsKeyId: arn:aws:kms:us-east-1:123456789012:key/prod-key
  diskIopsConfiguration:
    mode: USER_PROVISIONED
    iops: 100000
  rootVolumeConfiguration:
    dataCompressionType: ZSTD
    copyTagsToSnapshots: true
    nfsExports:
      clientConfigurations:
        - clients: "10.0.0.0/16"
          options:
            - rw
            - crossmnt
    userAndGroupQuotas:
      - id: 1000
        storageCapacityQuotaGib: 500
        type: USER
      - id: 100
        storageCapacityQuotaGib: 1000
        type: GROUP
  automaticBackupRetentionDays: 14
  dailyAutomaticBackupStartTime: "03:00"
  copyTagsToBackups: true
  copyTagsToVolumes: true
```

### Using Foreign Key References

Reference subnets, security groups, and KMS key from other Planton-managed resources:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsFsxOpenzfsFileSystem
metadata:
  name: ref-nfs
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsFsxOpenzfsFileSystem.ref-nfs
spec:
  region: us-east-1
  storageCapacityGib: 512
  throughputCapacity: 320
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: my-private-subnet-a
        fieldPath: status.outputs.subnet_id
  securityGroupIds:
    - valueFrom:
        kind: AwsSecurityGroup
        name: nfs-sg
        field: status.outputs.security_group_id
  kmsKeyId:
    valueFrom:
      kind: AwsKmsKey
      name: data-key
      field: status.outputs.key_arn
  rootVolumeConfiguration:
    dataCompressionType: LZ4
    nfsExports:
      clientConfigurations:
        - clients: "*"
          options:
            - rw
            - no_root_squash
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `file_system_id` | `string` | File system ID (e.g., `fs-0123456789abcdef0`). Primary identifier for EKS PVs, ECS volumes, and other integrations. |
| `file_system_arn` | `string` | ARN for IAM resource-level permissions and cross-account access. |
| `dns_name` | `string` | DNS name for NFS mount commands (e.g., `fs-xxx.fsx.us-east-1.amazonaws.com`). |
| `endpoint_ip_address` | `string` | Endpoint IP. For MULTI_AZ_1, this is the floating IP that follows the active file server. |
| `root_volume_id` | `string` | Root volume ID (e.g., `fsvol-xxx`). Use as `parentVolumeId` when creating child OpenZFS volumes. |
| `network_interface_ids` | `string[]` | ENI IDs created for the file system. 1 for SINGLE_AZ, 2 for MULTI_AZ. |
| `vpc_id` | `string` | VPC where the file system resides. |
| `owner_id` | `string` | AWS account ID of the file system owner. |

## Related Components

- [AwsVpc](/docs/catalog/aws/vpc) — provides the subnets for file system placement
- [AwsSecurityGroup](/docs/catalog/aws/security-group) — controls NFS traffic to/from the file system
- [AwsKmsKey](/docs/catalog/aws/kms-key) — provides customer-managed encryption keys
- [AwsElasticFileSystem](/docs/catalog/aws/elastic-file-system) — simpler serverless NFS alternative (EFS)
- [AwsFsxLustreFileSystem](/docs/catalog/aws/fsx-lustre-file-system) — HPC-optimized file system with S3 integration
