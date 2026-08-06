# AWS RDS Instance

Deploys a single RDS DB instance -- postgres, mysql, mariadb, oracle, or
sqlserver -- with Multi-AZ standby failover, storage autoscaling, an
AWS-managed master password in Secrets Manager, read-replica and
point-in-time-restore create shapes, and Blue/Green near-zero-downtime
updates, with every attachment (subnets, security groups, KMS keys, the
monitoring role) composed by reference.

## What Gets Created

When you deploy an AwsRdsInstance resource, Planton provisions:

- **DB instance** — an `aws_db_instance` / `rds.Instance` with the chosen
  engine, class, and storage (gp3 IOPS/throughput tuning, autoscaling
  ceiling, dedicated log volume), optionally Multi-AZ, publicly
  addressable, or joined to an Active Directory
- **DB subnet group** — managed automatically from `subnetIds` (pure
  glue: a named list of subnets), or an existing group by name

A read replica (`replicateSourceDb`), a snapshot restore
(`snapshotIdentifier`), or a point-in-time restore
(`restoreToPointInTime`) inherits engine, storage, and credentials from
its source -- the spec validation keeps those fields empty so drift is
impossible.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **Two subnets in distinct AZs** (`AwsSubnet`) or an existing DB subnet group -- AWS requires two AZs even for a single-AZ instance.
- **A security group** (`AwsSecurityGroup`) allowing the database port from your application tier -- or omit to use the VPC default group.
- **A KMS key** (`AwsKmsKey`) only when replacing the AWS-managed keys.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRdsInstance
metadata:
  name: billing-db
spec:
  region: us-west-2
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: platform-private-1
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: platform-private-2
        fieldPath: status.outputs.subnet_id
  securityGroupIds:
    - valueFrom:
        kind: AwsSecurityGroup
        name: database-sg
        fieldPath: status.outputs.security_group_id
  engine: postgres
  instanceClass: db.m6g.large
  allocatedStorageGb: 50
  storageType: gp3
  storageEncrypted: true
  manageMasterUserPassword: true
  multiAz: true
  skipFinalSnapshot: true
```

```shell
planton apply -f rds-instance.yaml
```

This creates a Multi-AZ PostgreSQL instance on gp3 storage with its
master password managed in Secrets Manager.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the referenced subnets/SGs/keys. | Required; non-empty |
| `instanceClass` | `string` | The compute size, e.g. `db.t4g.micro`, `db.m6g.large`. | Required; `db.` prefix |
| `engine` | `string` | The database engine. Required for a new instance; inherited on replicas/restores. | Conditional |
| `allocatedStorageGb` | `int` | Storage in GiB. Required for a new instance; inherited on replicas/restores. | Conditional |
| networking | — | At least two `subnetIds` (distinct AZs) or an existing `dbSubnetGroupName`. | Enforced |

### Storage Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `storageType` | `string` | AWS default | `gp3` (modern default), `gp2`, `io1`, `io2`, `standard`. |
| `iops` / `storageThroughput` | `int` | type baseline | Raise gp3 above 3000 IOPS / 125 MiB/s; required for io1/io2. |
| `maxAllocatedStorageGb` | `int` | 0 (off) | Storage autoscaling ceiling -- insurance against disk-full outages. |
| `dedicatedLogVolume` | `bool` | `false` | A separate EBS volume for logs. |
| `storageEncrypted` + `kmsKeyId` | — | recommended on | Create-time one-way door; reference an `AwsKmsKey` `key_arn` output. |

### Credentials

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `manageMasterUserPassword` | `bool` | recommended `true` | AWS generates, stores, and rotates the master password in Secrets Manager; ARN exported as `master_user_secret_arn`. |
| `password` | `string` (sensitive) | — | Direct password; mutually exclusive with the managed strategy. |
| `username` | `string` | — | Required for a new instance (AWS has no default); create-only. Replicas and restores inherit it. |
| `iamDatabaseAuthenticationEnabled` | `bool` | `false` | Short-lived IAM auth tokens instead of passwords. |

### Availability, Replicas, and Restores

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `multiAz` | `bool` | `false` | Synchronous standby in a second AZ with automatic failover. |
| `availabilityZone` | `string` | AWS-placed | Single-AZ pin; mutually exclusive with `multiAz`. |
| `replicateSourceDb` | `string` | — | Make this a read replica (identifier same-region, ARN cross-region); `replicaMode: mounted` for Oracle DR. |
| `snapshotIdentifier` / `restoreToPointInTime` | — | — | Create from a snapshot or another instance's continuous backup. |
| `blueGreenUpdateEnabled` | `bool` | `false` | Updates run on a synchronized green copy and switch over in under a minute. |

### Data Protection and Operations

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `backupRetentionPeriod` | `int` | engine default | 0-35 days; 0 disables automated backups and PITR. |
| `skipFinalSnapshot` / `finalSnapshotIdentifier` | — | safe | A final-snapshot name is required unless skipping is explicit. |
| `deletionProtection` / `deleteAutomatedBackups` | — | — | The same two-step-delete and backup-retention posture as the cluster kind. |
| `backupWindow` / `maintenanceWindow` | `string` | AWS-assigned | UTC windows; must not overlap. |
| `autoMinorVersionUpgrade` | `bool` | `true` | Tri-state; disable only for manual patch control. |
| `allowMajorVersionUpgrade` / `applyImmediately` | `bool` | `false` | Major-version guard; immediate-vs-window application. |

### Observability and Enterprise

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `performanceInsightsEnabled` + retention/KMS | — | off | Per-query telemetry; free at 7-day retention. |
| `monitoringInterval` + `monitoringRoleArn` | — | off | Enhanced Monitoring through a referenced `AwsIamRole`. |
| `databaseInsightsMode` | `string` | `standard` | `advanced` for fleet-level analysis. |
| `enabledCloudwatchLogsExports` | `list` | `[]` | Engine-specific log types. |
| `activeDirectory` | `object` | — | AWS-managed (`domain` + role) or self-managed (fqdn/OU/secret ARN/DNS IPs) AD join. |
| `licenseModel` / `characterSetName` / `ncharCharacterSetName` / `timezone` | `string` | engine default | Oracle / SQL Server create-time knobs. |
| `parameterGroupName` / `optionGroupName` | `string` | engine default | Engine-configuration attachments by name. |

## Stack Outputs

| Output | Description |
| --- | --- |
| `instance_identifier` | The instance identifier. |
| `arn` | The instance ARN. |
| `resource_id` | The immutable resource ID -- keys PITR, IAM auth policies, CloudWatch. |
| `endpoint` | The connection endpoint (`address:port`). |
| `address` | The bare hostname. |
| `port` | The listening port. |
| `hosted_zone_id` | For Route53 alias records. |
| `engine_version_actual` | The resolved running version. |
| `master_user_secret_arn` | The Secrets Manager ARN of the managed master password (managed strategy only). |
| `db_subnet_group_name` | The subnet group in use. |

## Related Resources

- [AwsSubnet](/docs/catalog/aws/awssubnet) — the private subnets behind the subnet group
- [AwsSecurityGroup](/docs/catalog/aws/awssecuritygroup) — database ingress rules live here
- [AwsKmsKey](/docs/catalog/aws/awskmskey) — customer-managed encryption keys
- [AwsIamRole](/docs/catalog/aws/awsiamrole) — the Enhanced Monitoring role
- [AwsRdsCluster](/docs/catalog/aws/awsrdscluster) — Aurora and Multi-AZ cluster shapes
