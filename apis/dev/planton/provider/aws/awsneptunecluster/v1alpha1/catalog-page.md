# AWS Neptune Cluster

Deploys an Amazon Neptune graph database cluster -- Gremlin, openCypher,
and SPARQL over shared cluster storage, provisioned or Serverless --
with its writer and reader instances folded into the same spec, IAM
database authentication as the credential mechanism, and every
attachment (subnets, security groups, KMS keys, IAM roles) composed by
reference.

## What Gets Created

When you deploy an AwsNeptuneCluster resource, Planton provisions:

- **Neptune cluster** — an `aws_neptune_cluster` / `neptune.Cluster`
  with the chosen shape: provisioned instances or Neptune Serverless NCU
  bounds
- **Cluster instances** — one `aws_neptune_cluster_instance` /
  `neptune.ClusterInstance` per `instances` entry, keyed by name so
  adding or removing a reader is an in-place update; the lowest
  promotion tier becomes the writer
- **Neptune subnet group** — managed automatically from `subnetIds`
  (pure glue: a named list of subnets), or an existing group by name
- **Cluster parameter group** — managed automatically when inline
  `parameters` are provided, with the family derived from the pinned
  engine version

The cluster never modifies a resource it merely references: security
groups carry their own ingress rules, IAM roles own their policies, and
KMS keys govern their own rotation.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **Two subnets in distinct AZs** (`AwsSubnet`) or an existing Neptune subnet group.
- **A security group** (`AwsSecurityGroup`) allowing port 8182 from your application tier -- or omit to use the VPC default group.
- **A KMS key** (`AwsKmsKey`) only when replacing the AWS-managed storage encryption key.
- **IAM roles** (`AwsIamRole`) only for the S3 bulk loader or Neptune ML integrations.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsNeptuneCluster
metadata:
  name: knowledge-graph
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
  storageEncrypted: true
  iamDatabaseAuthenticationEnabled: true
  skipFinalSnapshot: true
  serverlessV2Scaling:
    minCapacity: 1
    maxCapacity: 8
  instances:
    - name: writer
      instanceClass: db.serverless
```

```shell
planton apply -f neptune.yaml
```

This creates a Neptune Serverless cluster (the current AWS default
engine version) that scales between 1 and 8 NCUs with demand and
authenticates callers through IAM.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match the referenced subnets/SGs/keys. | Required; non-empty |
| networking | — | At least two `subnetIds` (distinct AZs) or an existing `neptuneSubnetGroupName`. | Enforced |
| `instances` | `list` | The writer/reader instances. Required unless the cluster starts headless (a restore, a replica, or a global-cluster member). | Enforced |

### Cluster Shape Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `instances` | `list` | `[]` | Per-name instances (name, instanceClass incl. `db.serverless`, promotionTier, AZ pin, instance parameter group, maintenance window). |
| `serverlessV2Scaling` | `object` | — | NCU bounds for `db.serverless` instances (1-128 on both ends). |
| `engineVersion` | `string` | AWS default | Pin for deliberate upgrades; empty never goes stale. |
| `storageType` | `string` | `standard` | `iopt1` enables I/O-Optimized storage (engine 1.3+). |
| `port` | `int` | 8182 | Create-only. |

### Authentication and Encryption

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `iamDatabaseAuthenticationEnabled` | `bool` | `false` | SigV4-signed requests from IAM identities -- Neptune's only credential mechanism (no master password exists). |
| `iamRoles` | `list \| valueFrom` | `[]` | Roles the ENGINE assumes (S3 bulk loader, Neptune ML). Reference `AwsIamRole` `role_arn` outputs. |
| `storageEncrypted` | `bool` | recommended `true` | Create-time one-way door. |
| `kmsKeyId` | `string \| valueFrom` | AWS-managed key | Reference an `AwsKmsKey` `key_arn` output. |

### Data Protection

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `backupRetentionPeriod` | `int` | 1 | Days of continuous backup (1-35) -- bounds point-in-time recovery. |
| `skipFinalSnapshot` / `finalSnapshotIdentifier` | — | safe | A final-snapshot name is required unless skipping is explicit. |
| `deletionProtection` | `bool` | `false` | Deleting becomes a deliberate two-step. |
| `copyTagsToSnapshot` | `bool` | `false` | Cluster tags flow onto snapshots. |
| `snapshotIdentifier` | `string` | — | Create the cluster from a snapshot. |
| `replicationSourceIdentifier` | `string` | — | Make this cluster a read replica of the source ARN. |
| `globalClusterIdentifier` | `string` | — | Join a Neptune global database. |

### Observability and Parameters

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `enabledCloudwatchLogsExports` | `list` | `[]` | `audit` and/or `slowquery` -- each also needs its matching cluster parameter before Neptune emits anything. |
| `parameters` / `neptuneClusterParameterGroupName` | — | engine default | Inline parameters (module-managed group) or an existing group. |
| `neptuneInstanceParameterGroupName` | `string` | — | Applied to instances during a major version upgrade (required with `allowMajorVersionUpgrade`). |

## Stack Outputs

| Output | Description |
| --- | --- |
| `cluster_identifier` | The cluster identifier. |
| `arn` | The cluster ARN. |
| `cluster_resource_id` | The immutable resource ID -- keys CloudWatch dimensions and IAM auth policies. |
| `endpoint` | The writer endpoint. |
| `reader_endpoint` | Load-balances read-only queries across reader instances. |
| `port` | The listening port. |
| `hosted_zone_id` | For Route53 alias records to the endpoints. |
| `engine_version_actual` | The resolved running version. |
| `neptune_subnet_group_name` | The subnet group in use. |
| `neptune_cluster_parameter_group_name` | The parameter group in use. |
| `instance_endpoints` | Per-instance endpoints of the folded instances, in spec order. |

## Related Resources

- [AwsSubnet](/docs/catalog/aws/awssubnet) — the private subnets the cluster spans
- [AwsSecurityGroup](/docs/catalog/aws/awssecuritygroup) — database ingress rules live here
- [AwsKmsKey](/docs/catalog/aws/awskmskey) — customer-managed encryption keys
- [AwsIamRole](/docs/catalog/aws/awsiamrole) — bulk loader and ML integration roles
- [AwsDocumentDb](/docs/catalog/aws/awsdocumentdb) — the document-database sibling with the same cluster anatomy
