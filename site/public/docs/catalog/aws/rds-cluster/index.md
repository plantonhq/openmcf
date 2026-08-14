---
title: "RDS Cluster"
description: "RDS Cluster deployment documentation"
icon: "package"
order: 100
componentName: "awsrdscluster"
---

# AWS RDS Cluster

Deploys an RDS database cluster in any of its three shapes — Aurora (an instance fleet on shared cluster storage, including Serverless v2), legacy Aurora Serverless v1, or a Multi-AZ community mysql/postgres cluster — with AWS Secrets Manager integration for managed master passwords, storage encryption, continuous backups with backtrack and fast clones, S3 XtraBackup migration, Kerberos through AWS Managed Microsoft AD, feature-scoped IAM role associations, custom reader endpoints, Database Activity Streams, Aurora Global Database membership, the RDS Data API, and CloudWatch log exports. The cluster integrates with Planton's Provider Connections for AWS credential management and supports ValueFromRef wiring to subnets, security groups, KMS keys, and IAM roles.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Aurora DB Cluster** -- a managed database cluster in the specified AWS region running Aurora MySQL or Aurora PostgreSQL, placed in a DB subnet group spanning at least two Availability Zones
- **DB Subnet Group** -- created from the provided `subnetIds`; skipped when an existing `dbSubnetGroupName` is specified instead
- **Cluster Instances** -- one provider resource per `instances[]` entry (Aurora topology), each keyed by name so scaling readers never touches the cluster; Serverless v1 and Multi-AZ clusters have AWS-owned compute instead
- **Cluster Parameter Group** -- created only when `parameters` entries are configured; applies engine-specific tuning parameters
- **Secrets Manager Secret** -- created only when `manageMasterUserPassword` is `true` (recommended); stores the master password with automatic rotation, optionally encrypted with a customer-managed KMS key
- **CloudWatch Log Streams** -- created only when `enabledCloudwatchLogsExports` entries are specified; exports engine-specific logs (audit, error, slowquery for MySQL; postgresql, upgrade for PostgreSQL)
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or cross-account trust authentication modes.

### AWS Account

- **At least two subnets** in distinct Availability Zones within the target VPC. Private subnets are recommended for production. Provide subnet IDs directly or reference an AwsVpc Cloud Resource via ValueFromRef. Alternatively, provide an existing `dbSubnetGroupName`.
- **A security group** (optional) to attach to the cluster for network access control. Provide security group IDs directly or reference an AwsSecurityGroup Cloud Resource — database ingress rules belong on the referenced security-group node.
- **A KMS key** (optional) for encrypting cluster storage and the managed master password beyond the default AWS-managed key. Provide the ARN directly or reference an AwsKmsKey Cloud Resource.

## Deploy

### Console

Open the deployment store, find **AWS RDS Cluster**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. The [Presets](#presets) tab offers three starting configurations: **Aurora PostgreSQL**, **Aurora MySQL**, and **Aurora Serverless v2**.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRdsCluster
metadata:
  name: app-database
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  subnetIds:
    - value: "subnet-0a1b2c3d4e5f00001"
    - value: "subnet-0a1b2c3d4e5f00002"
  engine: aurora-postgresql
  engineVersion: "16.4"
  masterUsername: dbadmin
  manageMasterUserPassword: true
  storageEncrypted: true
  deletionProtection: true
  backupRetentionPeriod: 7
  skipFinalSnapshot: false
  finalSnapshotIdentifier: app-database-final
  instances:
    - name: writer
      instanceClass: db.r6g.large
      promotionTier: 0
    - name: reader-1
      instanceClass: db.r6g.large
      promotionTier: 1
```

```shell
planton apply -f rds-cluster.yaml
```

This creates an Aurora PostgreSQL cluster served by a writer and a reader, with the master password managed in Secrets Manager, encrypted storage, deletion protection, and 7-day continuous backups. A Stack Job tracks the provisioning and streams progress in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the Aurora cluster to a VPC and security group deployed in the same InfraPipeline:

```yaml
spec:
  subnetIds:
    - valueFrom:
        kind: AwsVpc
        name: production-vpc
        fieldPath: status.outputs.private_subnets.[0].id
    - valueFrom:
        kind: AwsVpc
        name: production-vpc
        fieldPath: status.outputs.private_subnets.[1].id
  securityGroupIds:
    - valueFrom:
        kind: AwsSecurityGroup
        name: database-sg
        fieldPath: status.outputs.security_group_id
  kmsKeyId:
    valueFrom:
      kind: AwsKmsKey
      name: database-key
      fieldPath: status.outputs.key_arn
```

The InfraPipeline resolves the dependency graph, deploys the VPC, security group, and KMS key first, then provisions the Aurora cluster with the resolved values.

## Key Configuration

These are the most important decisions when configuring an Aurora cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Topology** -- Three mutually exclusive cluster shapes share this spec: Aurora engines (`aurora-postgresql`, `aurora-mysql`) with an `instances[]` fleet; `engineMode: serverless` for legacy Aurora Serverless v1 (no instances, `serverlessV1Scaling` applies); or a community engine (`postgres`, `mysql`) with `dbClusterInstanceClass` + `allocatedStorageGb` + `iops` for a Multi-AZ RDS cluster (AWS manages one writer and two readable standbys). The engine and topology are immutable after creation.

**Instance fleet** -- Each `instances[]` entry is its own provider resource keyed by `name`: one writer (lowest `promotionTier`) plus readers. Use `db.serverless` as an instance class with `serverlessV2Scaling` ACU bounds for Serverless v2 — `minCapacity: 0` enables scale-to-zero auto-pause.

**Creation source** -- A new cluster configures everything; `snapshotIdentifier` restores a snapshot, `restoreToPointInTime` restores continuous backup (with `restoreType: copy-on-write` for an Aurora fast clone — prod-data staging that pays only for divergence), `s3Import` restores a Percona XtraBackup from S3 (the Aurora MySQL migration on-ramp), and `replicationSourceIdentifier` creates a cross-region replica (promote by clearing it).

**Password management** -- Set `manageMasterUserPassword: true` (recommended) to let AWS Secrets Manager store and automatically rotate the master password; its ARN is exported as `master_user_secret_arn`. A supplied `masterPassword` is mutually exclusive and must be an org-secret reference.

**Networking** -- Provide either `subnetIds` (at least two subnets in distinct AZs) or `dbSubnetGroupName` (an existing subnet group). Attach security groups via `securityGroupIds` — database ingress rules belong on the referenced AwsSecurityGroup node.

**Encryption and compliance** -- Set `storageEncrypted: true` (create-time only) and optionally `kmsKeyId` for a customer-managed key. Enable `deletionProtection`, keep `skipFinalSnapshot: false` with a `finalSnapshotIdentifier`, and enable `iamDatabaseAuthenticationEnabled` for IAM-based authentication. `backtrackWindowSeconds` (Aurora MySQL only) enables in-place rewind.

**Global and advanced** -- `globalClusterIdentifier` joins an Aurora Global Database (write forwarding available); `enableHttpEndpoint` turns on the Data API for connection-averse callers like Lambda; `parameters` manages inline cluster parameters (mutually exclusive with `dbClusterParameterGroupName`).

**Access and audit** -- `iamRoles` attaches feature-scoped engine roles (one association per entry, with `featureName` like `s3Export` or `Lambda`; the role's trust policy must allow `rds.amazonaws.com` to assume it — AWS validates that server-side at association time and rejects the deploy otherwise); `domain` + `domainIamRoleName` join an AWS Managed Microsoft AD for Kerberos authentication; `customEndpoints` carve stable READER/ANY DNS names over chosen instances; `activityStream` streams every database event to a KMS-encrypted Kinesis stream for compliance-grade auditing (Aurora only).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AwsSubnet** | `subnetIds` | `status.outputs.subnet_id` |
| **AwsSecurityGroup** (optional) | `securityGroupIds` | `status.outputs.security_group_id` |
| **AwsKmsKey** (optional) | `kmsKeyId`, `masterUserSecretKmsKeyId`, `performanceInsightsKmsKeyId` | `status.outputs.key_arn` |
| **AwsIamRole** (optional) | `iamRoles[].role`, `monitoringRoleArn`, `s3Import.ingestionRole`, per-instance `monitoringRoleArn` | `status.outputs.role_arn` |
| **AwsKmsKey** (optional, audit) | `activityStream.kmsKeyId` | `status.outputs.key_arn` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `endpoint` | Primary writer endpoint | Application connection strings, DNS CNAME records |
| `reader_endpoint` | Reader endpoint (round-robin across readers) | Read-heavy application connection strings |
| `port` | Port the cluster accepts connections on | Application connection configuration |
| `cluster_identifier` | AWS identifier of the DB cluster | Monitoring dashboards, CloudWatch alarms |
| `arn` | Amazon Resource Name | IAM policies, CloudWatch alarms, resource tagging |
| `cluster_resource_id` | Immutable cluster resource ID | Point-in-time restore sources, IAM auth policies |
| `hosted_zone_id` | The endpoints' Route 53 hosted zone | Alias records |
| `engine_version_actual` | The running engine version | Upgrade auditing |
| `master_user_secret_arn` | Secrets Manager ARN of the AWS-managed master secret | Application credential reads |
| `db_subnet_group_name` | DB subnet group name | Audit, related resource lookups |
| `db_cluster_parameter_group_name` | Cluster parameter group in use | Parameter auditing |
| `instance_endpoints` | Per-instance endpoints | Instance-pinned diagnostics |
| `custom_endpoints` | Name + DNS of each custom endpoint | Workload-scoped connection strings (e.g. analytics) |
| `activity_stream_kinesis_stream_name` | Kinesis stream receiving the Database Activity Stream | Audit/SIEM consumers, GuardDuty RDS Protection |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Aurora PostgreSQL production** -- Encrypted storage, managed master password in Secrets Manager, deletion protection, 7-day backup retention with final snapshot, and CloudWatch log export for PostgreSQL logs. Start from the **Aurora PostgreSQL** preset.

**Aurora MySQL production** -- Same resilience posture as PostgreSQL but with Aurora MySQL 8.0. Exports error and slow query logs to CloudWatch. Suitable for applications that require MySQL-specific features. Start from the **Aurora MySQL** preset.

**Aurora Serverless v2** -- Aurora PostgreSQL with a `db.serverless` writer, ACU bounds with scale-to-zero auto-pause, and the Data API enabled. Ideal for variable or connection-averse workloads. Start from the **Aurora Serverless v2** preset.

**MySQL migration via S3** -- Restore a Percona XtraBackup from S3 into a new Aurora MySQL cluster with a backtrack window for cutover safety. Start from the **MySQL to Aurora via S3 Import** preset.

## Works With

- [**AWS Subnet**](/cloud-catalog/aws-subnet) -- provides the subnets for the DB subnet group across multiple Availability Zones
- [**AWS Security Group**](/cloud-catalog/aws-security-group) -- provides network access control for the cluster endpoints
- [**AWS KMS Key**](/cloud-catalog/aws-kms-key) -- provides customer-managed keys for storage, the master secret, and Performance Insights
- [**AWS IAM Role**](/cloud-catalog/aws-iam-role) -- service-integration roles (S3 import/export, Lambda, ML) and the Enhanced Monitoring role
- [**AWS RDS Instance**](/cloud-catalog/aws-rds-instance) -- the standalone-instance sibling for Oracle, SQL Server, and single-instance community engines
