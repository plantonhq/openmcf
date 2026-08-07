# AWS RDS Cluster

Deploys an RDS DB cluster -- an Aurora MySQL/PostgreSQL cluster
(provisioned or Serverless v2 with scale-to-zero), a legacy Aurora
Serverless v1 cluster, or a Multi-AZ RDS cluster of the community
mysql/postgres engines -- with its writer and reader instances folded
into the same spec, an AWS-managed master password in Secrets Manager,
and every attachment (subnets, security groups, KMS keys, IAM roles)
composed by reference.

## What Gets Created

When you deploy an AwsRdsCluster resource, Planton provisions:

- **RDS cluster** — an `aws_rds_cluster` / `rds.Cluster` with the chosen
  engine and shape: provisioned Aurora, Aurora Serverless v2 bounds
  (including automatic pause at `minCapacity: 0`), legacy Serverless v1
  scaling, or a Multi-AZ RDS cluster sized by instance class and storage
- **Cluster instances** — one `aws_rds_cluster_instance` /
  `rds.ClusterInstance` per `instances` entry, keyed by name so adding or
  removing a reader is an in-place update; the lowest promotion tier
  becomes the writer
- **DB subnet group** — managed automatically from `subnetIds` (pure
  glue: a named list of subnets), or an existing group by name
- **Cluster parameter group** — managed automatically when inline
  `parameters` are provided, with the family derived from the pinned
  engine version

The cluster never modifies a resource it merely references: security
groups carry their own ingress rules, IAM roles own their policies, and
KMS keys govern their own rotation.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **Two subnets in distinct AZs** (`AwsSubnet`) or an existing DB subnet group.
- **A security group** (`AwsSecurityGroup`) allowing the database port from your application tier -- or omit to use the VPC default group.
- **A KMS key** (`AwsKmsKey`) only when replacing the AWS-managed keys for storage, Performance Insights, or the managed master-user secret.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRdsCluster
metadata:
  name: orders-db
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
  engine: aurora-postgresql
  manageMasterUserPassword: true
  storageEncrypted: true
  skipFinalSnapshot: true
  serverlessV2Scaling:
    minCapacity: 0
    maxCapacity: 4
  instances:
    - name: writer
      instanceClass: db.serverless
```

```shell
planton apply -f rds-cluster.yaml
```

This creates an Aurora PostgreSQL Serverless v2 cluster that scales with
demand, pauses to zero compute cost when idle, and keeps its master
password in Secrets Manager.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the referenced subnets/SGs/keys. | Required; non-empty |
| `engine` | `string` | `aurora-mysql`, `aurora-postgresql` (Aurora), or `mysql`, `postgres` (Multi-AZ RDS cluster). Create-only. | Required; one of the four |
| networking | — | At least two `subnetIds` (distinct AZs) or an existing `dbSubnetGroupName`. | Enforced |

### Cluster Shape Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `instances` | `list` | `[]` | The writer/reader instances (name, instanceClass incl. `db.serverless`, promotionTier, AZ pin, per-instance parameter group, PI/monitoring overrides). Empty only for Serverless v1 and Multi-AZ RDS clusters. |
| `serverlessV2Scaling` | `object` | — | ACU bounds for `db.serverless` instances; `minCapacity: 0` enables automatic pause (`secondsUntilAutoPause` tunes the idle window). |
| `engineMode` | `string` | `provisioned` | `serverless` selects legacy Aurora Serverless v1 (with `serverlessV1Scaling`). |
| `dbClusterInstanceClass` + `allocatedStorageGb` + `iops` | — | — | The Multi-AZ RDS cluster trio (community engines only); AWS manages one writer + two readers internally. |
| `storageType` | `string` | engine default | Aurora: `aurora-iopt1` (I/O-Optimized). Multi-AZ RDS: `io1`, `io2`, `gp3`. |

### Credentials and Encryption

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `manageMasterUserPassword` | `bool` | recommended `true` | AWS generates, stores, and rotates the master password in Secrets Manager; ARN exported as `master_user_secret_arn`. |
| `masterPassword` | `string` (sensitive) | — | Direct password; mutually exclusive with the managed strategy. |
| `masterUsername` | `string` | — | Required for a new cluster (AWS has no default); create-only. Restores, replicas, and global-database secondaries inherit it. |
| `storageEncrypted` | `bool` | recommended `true` | Create-time one-way door. |
| `kmsKeyId` / `masterUserSecretKmsKeyId` / `performanceInsightsKmsKeyId` | `string \| valueFrom` | AWS-managed keys | Reference `AwsKmsKey` `key_arn` outputs. |

### Data Protection

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `backupRetentionPeriod` | `int` | 1 | Days of continuous backup (1-35) -- bounds point-in-time recovery. |
| `skipFinalSnapshot` / `finalSnapshotIdentifier` | — | safe | A final-snapshot name is required unless skipping is explicit. |
| `deletionProtection` | `bool` | `false` | Deleting becomes a deliberate two-step. |
| `deleteAutomatedBackups` | `bool` | `true` | `false` retains backups after deletion. |
| `backtrackWindowSeconds` | `int` | 0 | Aurora MySQL in-place rewind (up to 72h). Create-time decision. |
| `snapshotIdentifier` / `restoreToPointInTime` | — | — | Create the cluster from a snapshot or another cluster's continuous backup (incl. copy-on-write fast clones). |

### Integrations

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `iamDatabaseAuthenticationEnabled` | `bool` | `false` | Short-lived IAM auth tokens instead of passwords. |
| `iamRoles` | `list \| valueFrom` | `[]` | Roles the ENGINE assumes (S3 import/export, Lambda, ML). Reference `AwsIamRole` `role_arn` outputs. |
| `enableHttpEndpoint` | `bool` | `false` | The Data API: SQL over HTTPS -- the Lambda-native access path. |
| `enabledCloudwatchLogsExports` | `list` | `[]` | Engine-family-validated log types. |
| `performanceInsightsEnabled` + retention/KMS | — | off | Per-query telemetry; free at 7-day retention. |
| `monitoringInterval` + `monitoringRoleArn` | — | off | Enhanced Monitoring (OS metrics) through a referenced role. |
| `globalClusterIdentifier` + write forwarding | — | — | Aurora Global Database membership; local/global write forwarding. |
| `parameters` / `dbClusterParameterGroupName` | — | engine default | Inline parameters (module-managed group) or an existing group. |

## Stack Outputs

| Output | Description |
| --- | --- |
| `cluster_identifier` | The cluster identifier. |
| `arn` | The cluster ARN. |
| `cluster_resource_id` | The immutable resource ID -- survives renames; keys PITR and CloudWatch. |
| `endpoint` | The writer endpoint. |
| `reader_endpoint` | Load-balances across reader instances. |
| `port` | The listening port. |
| `hosted_zone_id` | For Route53 alias records to the endpoints. |
| `engine_version_actual` | The resolved running version. |
| `master_user_secret_arn` | The Secrets Manager ARN of the managed master password (managed strategy only). |
| `db_subnet_group_name` | The subnet group in use. |
| `db_cluster_parameter_group_name` | The parameter group in use. |
| `instance_endpoints` | Per-instance endpoints of the folded instances, in spec order. |

## Related Resources

- [AwsSubnet](/docs/catalog/aws/awssubnet) — the private subnets the cluster spans
- [AwsSecurityGroup](/docs/catalog/aws/awssecuritygroup) — database ingress rules live here
- [AwsKmsKey](/docs/catalog/aws/awskmskey) — customer-managed encryption keys
- [AwsIamRole](/docs/catalog/aws/awsiamrole) — engine integration roles and the Enhanced Monitoring role
- [AwsRdsInstance](/docs/catalog/aws/awsrdsinstance) — the single-node alternative
