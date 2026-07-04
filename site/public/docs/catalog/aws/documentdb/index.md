---
title: "DocumentDB"
description: "DocumentDB deployment documentation"
icon: "package"
order: 100
componentName: "awsdocumentdb"
---

# AWS DocumentDB

Deploys an Amazon DocumentDB cluster -- a managed MongoDB-compatible
document database, provisioned or Serverless -- with its writer and
reader instances folded into the same spec, an AWS-managed master
password in Secrets Manager, and every attachment (subnets, security
groups, KMS keys) composed by reference.

## What Gets Created

When you deploy an AwsDocumentDb resource, Planton provisions:

- **DocumentDB cluster** — an `aws_docdb_cluster` / `docdb.Cluster` with
  the chosen shape: provisioned instances or DocumentDB Serverless DCU
  bounds
- **Cluster instances** — one `aws_docdb_cluster_instance` /
  `docdb.ClusterInstance` per `instances` entry, keyed by name so adding
  or removing a reader is an in-place update; the lowest promotion tier
  becomes the writer
- **DB subnet group** — managed automatically from `subnetIds` (pure
  glue: a named list of subnets), or an existing group by name
- **Cluster parameter group** — managed automatically when inline
  `parameters` are provided, with the family derived from the pinned
  engine version

The cluster never modifies a resource it merely references: security
groups carry their own ingress rules and KMS keys govern their own
rotation.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **Two subnets in distinct AZs** (`AwsSubnet`) or an existing DB subnet group.
- **A security group** (`AwsSecurityGroup`) allowing port 27017 from your application tier -- or omit to use the VPC default group.
- **A KMS key** (`AwsKmsKey`) only when replacing the AWS-managed keys for storage or Performance Insights.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsDocumentDb
metadata:
  name: orders-docdb
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
  masterUsername: docadmin
  manageMasterUserPassword: true
  storageEncrypted: true
  skipFinalSnapshot: true
  instances:
    - name: writer
      instanceClass: db.r6g.large
```

```shell
planton apply -f docdb.yaml
```

This creates a DocumentDB 5.0 cluster (the current AWS default version)
with one writer instance and its master password kept in Secrets
Manager.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the referenced subnets/SGs/keys. | Required; non-empty |
| networking | — | At least two `subnetIds` (distinct AZs) or an existing `dbSubnetGroupName`. | Enforced |
| `instances` | `list` | The writer/reader instances. Required unless the cluster starts headless (a restore or a global-cluster member). | Enforced |
| `masterUsername` | `string` | Required for a new cluster (AWS has no default); create-only. Restores and global-cluster members inherit it. | Enforced |

### Cluster Shape Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `instances` | `list` | `[]` | Per-name instances (name, instanceClass incl. `db.serverless`, promotionTier, AZ pin, maintenance window, Performance Insights, CA cert). |
| `serverlessV2Scaling` | `object` | — | DCU bounds for `db.serverless` instances (0.5-256, half-steps). Removing the block from a live cluster replaces it. |
| `engineVersion` | `string` | AWS default | Pin for deliberate upgrades; empty never goes stale. |
| `storageType` | `string` | `standard` | `iopt1` enables I/O-Optimized storage for I/O-heavy workloads. |
| `port` | `int` | 27017 | 1150-65535; create-only. |
| `networkType` | `string` | `IPV4` | `DUAL` for dual-stack IPv4+IPv6 (needs IPv6 subnets). |

### Credentials and Encryption

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `manageMasterUserPassword` | `bool` | recommended `true` | AWS generates, stores, and rotates the master password in Secrets Manager; ARN exported as `master_user_secret_arn`. |
| `masterPassword` | `string` (sensitive) | — | Direct password; mutually exclusive with the managed strategy. |
| `storageEncrypted` | `bool` | recommended `true` | Create-time one-way door. |
| `kmsKeyId` | `string \| valueFrom` | AWS-managed key | Reference an `AwsKmsKey` `key_arn` output. |

### Data Protection

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `backupRetentionPeriod` | `int` | 1 | Days of continuous backup (1-35) -- bounds point-in-time recovery. |
| `skipFinalSnapshot` / `finalSnapshotIdentifier` | — | safe | A final-snapshot name is required unless skipping is explicit. |
| `deletionProtection` | `bool` | `false` | Deleting becomes a deliberate two-step. |
| `snapshotIdentifier` / `restoreToPointInTime` | — | — | Create the cluster from a snapshot or another cluster's continuous backup (incl. copy-on-write fast clones). |
| `globalClusterIdentifier` | `string` | — | Join a DocumentDB global cluster. |

### Observability and Parameters

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `enabledCloudwatchLogsExports` | `list` | `[]` | `audit` and/or `profiler` -- each also needs its matching cluster parameter before DocumentDB emits anything. |
| per-instance `performanceInsightsEnabled` + KMS | — | off | Per-query telemetry; instance-scoped on DocumentDB. |
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

- [AwsSubnet](/docs/catalog/aws/subnet) — the private subnets the cluster spans
- [AwsSecurityGroup](/docs/catalog/aws/security-group) — database ingress rules live here
- [AwsKmsKey](/docs/catalog/aws/kms-key) — customer-managed encryption keys
- [AwsRdsCluster](/docs/catalog/aws/rds-cluster) — the relational sibling with the same cluster anatomy
