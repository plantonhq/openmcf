# AWS Elastic File System

Deploys an AWS Elastic File System with mount targets across specified subnets (optionally with static IPv4/IPv6 addressing), lifecycle policies for cost-optimized storage tiering, automatic backups, an optional IAM resource policy, replication-overwrite protection, and cross-region/cross-AZ replication for disaster recovery. The component bundles everything needed to make the file system mountable immediately after deployment. Application entry points are separate [AwsEfsAccessPoint](/docs/catalog/aws/awsefsaccesspoint) resources that reference the file system.

## What Gets Created

When you deploy an AwsElasticFileSystem resource, Planton provisions:

- **EFS File System** — an `efs.FileSystem` resource with the configured encryption, performance mode, throughput mode, lifecycle policies, and replication-overwrite protection
- **Mount Targets** — one `efs.MountTarget` per declared mount target, placing an elastic network interface in each Availability Zone for NFS client access on TCP port 2049
- **Backup Policy** — an `efs.BackupPolicy` enabling automatic daily backups via AWS Backup, created only when `backupEnabled` is `true`
- **File System Policy** — an `efs.FileSystemPolicy` attaching an IAM resource policy to the file system, created only when `policy` is provided
- **Replication Configuration** — an `efs.ReplicationConfiguration` replicating the file system to another region and/or AZ, created only when `replication` is provided

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **At least one subnet** where mount targets will be created (one subnet per AZ for multi-AZ availability)
- **A security group** allowing inbound NFS traffic (TCP port 2049) from the clients that will mount the file system (when omitted, AWS attaches the VPC's default security group)
- **A KMS key ARN** if using customer-managed encryption (otherwise EFS uses the AWS-managed `aws/elasticfilesystem` key)

## Quick Start

Create a file `efs.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsElasticFileSystem
metadata:
  name: my-efs
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsElasticFileSystem.my-efs
spec:
  region: us-east-1
  mountTargets:
    - subnetId:
        value: subnet-0a1b2c3d4e5f00001
    - subnetId:
        value: subnet-0a1b2c3d4e5f00002
  securityGroupIds:
    - value: sg-0a1b2c3d4e5f00001
```

Deploy:

```shell
planton apply -f efs.yaml
```

This creates an unencrypted, bursting-throughput EFS file system with mount targets in two subnets.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region where the file system will be created (e.g., `us-east-1`, `us-west-2`). | Required; non-empty |
| `mountTargets` | `object[]` | Mount targets — one per subnet, at most one per Availability Zone. | Minimum 1 item required |
| `mountTargets[].subnetId` | `StringValueOrRef` | Subnet for the mount target. Can reference an AwsSubnet resource via `valueFrom`. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `encrypted` | `bool` | `false` | Enable encryption at rest. ForceNew. Recommended: `true`. |
| `kmsKeyId` | `StringValueOrRef` | AWS-managed key | Customer-managed KMS key ARN for encryption. Requires `encrypted` to be `true`. ForceNew. Can reference AwsKmsKey via `valueFrom`. |
| `performanceMode` | `string` | `generalPurpose` | File system performance mode. ForceNew. Valid values: `generalPurpose`, `maxIO`. |
| `throughputMode` | `string` | `bursting` | Throughput mode. Valid values: `bursting`, `provisioned`, `elastic`. Elastic requires generalPurpose. |
| `provisionedThroughputInMibps` | `double` | — | Fixed throughput in MiB/s. Required when `throughputMode` is `provisioned`. AWS accepts 1.0–3414.0 (generalPurpose), 1.0–1024.0 (maxIO). |
| `availabilityZoneName` | `string` | — | AZ name for One Zone storage (e.g., `us-east-1a`). ForceNew. When set, only one mount target is allowed. |
| `transitionToIa` | `string` | — | Transition to Infrequent Access after period. Valid values: `AFTER_1_DAY`, `AFTER_7_DAYS`, `AFTER_14_DAYS`, `AFTER_30_DAYS`, `AFTER_60_DAYS`, `AFTER_90_DAYS`, `AFTER_180_DAYS`, `AFTER_270_DAYS`, `AFTER_365_DAYS`. |
| `transitionToArchive` | `string` | — | Transition to Archive after period. Requires `transitionToIa` to be set. Same valid values as `transitionToIa`. |
| `transitionToPrimaryStorageClass` | `string` | — | Move files back to Standard on access. Only valid value: `AFTER_1_ACCESS`. |
| `backupEnabled` | `bool` | `false` | Enable automatic daily backups via AWS Backup. |
| `replicationOverwriteProtection` | `string` | `ENABLED` | `ENABLED` or `DISABLED`. Must be `DISABLED` before this file system can be a replication destination. |
| `mountTargets[].ipAddress` | `string` | auto-assigned | Static IPv4 address from the subnet's CIDR. ForceNew. |
| `mountTargets[].ipAddressType` | `string` | `IPV4_ONLY` | Address family: `IPV4_ONLY`, `IPV6_ONLY`, or `DUAL_STACK`. ForceNew. |
| `mountTargets[].ipv6Address` | `string` | auto-assigned | Static IPv6 address. Requires an IPv6-capable `ipAddressType`. ForceNew. |
| `securityGroupIds` | `StringValueOrRef[]` | VPC default SG | Security groups applied to all mount targets. Must allow inbound TCP 2049. Can reference AwsSecurityGroup via `valueFrom`. |
| `policy` | `object` | — | IAM resource policy for the file system as a JSON object structure. |
| `bypassPolicyLockoutSafetyCheck` | `bool` | `false` | Skip AWS's policy-lockout safety check. Only set when deliberately deploying a policy that denies the deploying principal future policy updates. Requires `policy`. |
| `replication.destinationRegion` | `string` | — | Region for the replica. At least one of region / AZ is required when `replication` is set. |
| `replication.destinationAvailabilityZoneName` | `string` | — | AZ for a One Zone replica — the cheaper DR shape. |
| `replication.destinationKmsKeyId` | `StringValueOrRef` | AWS-managed key | KMS key for the replica (replicas are always encrypted). |
| `replication.destinationFileSystemId` | `StringValueOrRef` | AWS creates one | Replicate into an existing file system (its overwrite protection must be `DISABLED`). |

## Examples

### Encrypted Multi-AZ File System

Production-ready file system with encryption and mount targets across two Availability Zones:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsElasticFileSystem
metadata:
  name: prod-efs
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsElasticFileSystem.prod-efs
spec:
  region: us-east-1
  encrypted: true
  mountTargets:
    - subnetId:
        value: subnet-private-az1
    - subnetId:
        value: subnet-private-az2
  securityGroupIds:
    - value: sg-nfs-access
```

### One Zone Development with Lifecycle Policies

Cost-optimized single-AZ file system with automatic tiering to Infrequent Access and Archive storage:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsElasticFileSystem
metadata:
  name: dev-efs
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsElasticFileSystem.dev-efs
spec:
  region: us-east-1
  availabilityZoneName: us-east-1a
  throughputMode: elastic
  transitionToIa: AFTER_30_DAYS
  transitionToArchive: AFTER_90_DAYS
  transitionToPrimaryStorageClass: AFTER_1_ACCESS
  mountTargets:
    - subnetId:
        value: subnet-dev-az1
  securityGroupIds:
    - value: sg-nfs-dev
```

### Cross-Region Disaster Recovery

File system replicated to another region; the replica stays read-only and in sync automatically:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsElasticFileSystem
metadata:
  name: dr-efs
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsElasticFileSystem.dr-efs
spec:
  region: us-east-1
  encrypted: true
  backupEnabled: true
  throughputMode: elastic
  mountTargets:
    - subnetId:
        value: subnet-private-az1
    - subnetId:
        value: subnet-private-az2
  securityGroupIds:
    - value: sg-nfs-access
  replication:
    destinationRegion: us-west-2
```

### Using Foreign Key References

Reference Planton-managed VPC subnets, security groups, and KMS keys instead of hardcoding IDs:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsElasticFileSystem
metadata:
  name: ref-efs
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsElasticFileSystem.ref-efs
spec:
  region: us-east-1
  encrypted: true
  kmsKeyId:
    valueFrom:
      kind: AwsKmsKey
      name: efs-key
      fieldPath: status.outputs.key_arn
  mountTargets:
    - subnetId:
        valueFrom:
          kind: AwsSubnet
          name: my-private-subnet-a
          fieldPath: status.outputs.subnet_id
    - subnetId:
        valueFrom:
          kind: AwsSubnet
          name: my-private-subnet-b
          fieldPath: status.outputs.subnet_id
  securityGroupIds:
    - valueFrom:
        kind: AwsSecurityGroup
        name: nfs-sg
        fieldPath: status.outputs.security_group_id
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `file_system_id` | `string` | File system ID (e.g., `fs-0123456789abcdef0`). Used by EKS PersistentVolumes, ECS task definitions, and AwsEfsAccessPoint references. |
| `file_system_arn` | `string` | Amazon Resource Name of the file system for IAM policies |
| `dns_name` | `string` | Regional DNS name for NFS mounting (e.g., `fs-xxx.efs.us-east-1.amazonaws.com`) |
| `mount_target_ids` | `map<string, string>` | Map of subnet ID to mount target ID |
| `mount_target_ips` | `map<string, string>` | Map of subnet ID to mount target IPv4 address (empty for IPV6_ONLY targets) |
| `mount_target_ipv6_addresses` | `map<string, string>` | Map of subnet ID to mount target IPv6 address (IPV6_ONLY / DUAL_STACK targets) |
| `mount_target_dns_names` | `map<string, string>` | Map of subnet ID to AZ-specific mount target DNS name |
| `replication_destination_file_system_id` | `string` | Replica's file system ID; empty when replication is not configured |

## Related Components

- [AwsEfsAccessPoint](/docs/catalog/aws/awsefsaccesspoint) — application-specific entry points enforcing POSIX identity and root directory
- [AwsVpc](/docs/catalog/aws/awsvpc) — provides the subnets for mount target placement
- [AwsSecurityGroup](/docs/catalog/aws/awssecuritygroup) — controls NFS access (TCP port 2049) to the file system
- [AwsKmsKey](/docs/catalog/aws/awskmskey) — provides customer-managed encryption keys
- [AwsEksCluster](/docs/catalog/aws/awsekscluster) — consumes EFS via the EFS CSI driver for PersistentVolumes
- [AwsLambda](/docs/catalog/aws/awslambda) — mounts EFS via AwsEfsAccessPoint ARNs for serverless file access
