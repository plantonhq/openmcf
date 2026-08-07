# AWS Redshift Cluster

Deploys a managed Amazon Redshift data warehouse cluster with configurable node topology, availability posture (Multi-AZ or zone relocation), managed password storage in AWS Secrets Manager, KMS encryption, audit logging, cross-region snapshot copy, and VPC placement. The cluster integrates with Planton's Provider Connections for AWS credential management and supports ValueFromRef wiring to subnets, security groups, IAM roles, KMS keys, and Elastic IPs.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Redshift Cluster** -- a columnar data warehouse cluster running the Redshift engine with the specified node type and count, placed in the configured subnets with configurable public or private access
- **Subnet Group** -- created from the provided `subnetIds` spanning at least two Availability Zones; skipped when an existing `clusterSubnetGroupName` is specified instead
- **Parameter Group** -- created only when inline `parameters` entries are configured; applies Redshift-specific tuning parameters (family `redshift-1.0` unless `parameterGroupFamily` selects another)
- **Logging Configuration** -- created only when `logging` is configured; delivers audit logs (connection, user activity, user DDL) to S3 or CloudWatch Logs
- **Snapshot Copy Configuration** -- created only when `snapshotCopy` names a destination region; automatically copies the cluster's snapshots there for disaster recovery
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied automatically for tracking and governance

Security groups are composed, never created: the cluster attaches the referenced `securityGroupIds`, and the ingress rules that open the warehouse port belong on those first-class AwsSecurityGroup resources.

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or cross-account trust authentication modes.

### AWS Account

- **At least two subnets** in distinct Availability Zones. Private subnets are recommended for production. Reference AwsSubnet Cloud Resources via ValueFromRef or provide subnet IDs directly. Alternatively, provide an existing `clusterSubnetGroupName`.
- **Security groups** (optional) governing who reaches the warehouse port. Reference AwsSecurityGroup Cloud Resources or provide security group IDs; empty keeps the VPC's default group.
- **IAM roles** (optional) for COPY, UNLOAD, and Redshift Spectrum queries that access S3, DynamoDB, or the Glue Data Catalog. Provide role ARNs directly or reference AwsIamRole Cloud Resources.
- **A KMS key** (optional) for customer-managed encryption of cluster data and the managed master password. Provide the ARN directly or reference an AwsKmsKey Cloud Resource.
- **An Elastic IP** (optional) for a stable public address on a publicly accessible cluster. Reference an AwsElasticIp Cloud Resource's `public_ip` output or provide the IP address directly.

## Deploy

### Console

Open the deployment store, find **AWS Redshift Cluster**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Single Node Dev** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRedshiftCluster
metadata:
  name: analytics-warehouse
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  nodeType: ra3.xlplus
  numberOfNodes: 2
  subnetIds:
    - value: "subnet-0a1b2c3d4e5f00001"
    - value: "subnet-0a1b2c3d4e5f00002"
  masterUsername: admin
  manageMasterPassword: true
  encrypted: true
  skipFinalSnapshot: false
  finalSnapshotIdentifier: analytics-warehouse-final
```

```shell
planton apply -f redshift-cluster.yaml
```

This creates a two-node RA3 Redshift cluster with a managed master password (stored in Secrets Manager), encrypted storage using the default AWS-managed key, and a final-snapshot deletion contract. A Stack Job tracks the provisioning and streams progress in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the Redshift cluster to subnets, an IAM role, and a KMS key deployed in the same InfraPipeline:

```yaml
spec:
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: warehouse-subnet-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: warehouse-subnet-b
        fieldPath: status.outputs.subnet_id
  securityGroupIds:
    - valueFrom:
        kind: AwsSecurityGroup
        name: warehouse-sg
        fieldPath: status.outputs.security_group_id
  kmsKeyId:
    valueFrom:
      kind: AwsKmsKey
      name: warehouse-key
      fieldPath: status.outputs.key_arn
  iamRoles:
    - valueFrom:
        kind: AwsIamRole
        name: redshift-service-role
        fieldPath: status.outputs.role_arn
  defaultIamRoleArn:
    valueFrom:
      kind: AwsIamRole
      name: redshift-service-role
      fieldPath: status.outputs.role_arn
```

The InfraPipeline resolves the dependency graph, deploys the subnets, KMS key, and IAM role first, then provisions the Redshift cluster with the resolved values.

## Key Configuration

These are the most important decisions when configuring a Redshift cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Node type and topology** -- Choose between dc2 (legacy dense compute, node-local SSD only) and ra3 (managed storage with automatic S3 tiering). RA3 nodes are the right call for nearly all new clusters -- compute and storage scale independently, and RA3 is the only family supporting Multi-AZ and zone relocation. Set `numberOfNodes` to 1 for a single-node development cluster or 2+ for multi-node production with a dedicated leader node. Resizing later is in-place, never a replacement.

**Availability posture** -- At most one of two mutually exclusive mechanisms: `multiAz: true` runs compute in two Availability Zones with automatic failover behind a single endpoint (requires RA3 and a multi-node cluster); `availabilityZoneRelocationEnabled: true` lets AWS move the single cluster to a healthy zone during outages (requires RA3 and a port within 5431-5455 or 8191-8215).

**Creation source** -- Leave the restore fields empty for a brand-new cluster, or restore from an existing snapshot with `snapshotIdentifier` (by name, optionally disambiguated by `snapshotClusterIdentifier`) or `snapshotArn` (the cross-account/cross-region share shape). A restored cluster inherits the source's credentials and databases, so `masterUsername` stays empty; `ownerAccount` identifies a snapshot shared by another AWS account.

**Password management** -- Set `manageMasterPassword: true` (recommended) to let AWS Secrets Manager store and rotate the master password automatically; the secret's ARN is exported as `master_password_secret_arn`. Optionally provide `masterPasswordSecretKmsKeyId` to encrypt the secret with a customer-managed KMS key. When `manageMasterPassword` is `false`, provide `masterPassword` as an org-secret reference.

**Audit logging** -- Configure `logging` with `logDestinationType` set to `"s3"` or `"cloudwatch"`. CloudWatch delivery exports exactly the log types you list in `logExports` (`connectionlog`, `useractivitylog`, `userlog`); S3 delivery always writes all three to the designated bucket. The user activity log only produces data when the `enable_user_activity_logging` parameter is set.

**Snapshots and disaster recovery** -- Set `automatedSnapshotRetentionPeriod` (0-35 days; production typically wants 7+) and `manualSnapshotRetentionPeriod` (-1 retains indefinitely). Keep `skipFinalSnapshot: false` with a `finalSnapshotIdentifier` so deletion always leaves a recovery handle. Configure `snapshotCopy` with a `destinationRegion` to automatically copy snapshots to a second region -- KMS-encrypted clusters additionally need a `snapshotCopyGrantName` in the destination region.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AwsSubnet** | `subnetIds` | `status.outputs.subnet_id` |
| **AwsSecurityGroup** (optional) | `securityGroupIds` | `status.outputs.security_group_id` |
| **AwsElasticIp** (optional) | `elasticIp` | `status.outputs.public_ip` |
| **AwsKmsKey** (optional) | `kmsKeyId` | `status.outputs.key_arn` |
| **AwsKmsKey** (optional) | `masterPasswordSecretKmsKeyId` | `status.outputs.key_arn` |
| **AwsIamRole** (optional) | `iamRoles` | `status.outputs.role_arn` |
| **AwsIamRole** (optional) | `defaultIamRoleArn` | `status.outputs.role_arn` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_identifier` | Unique identifier of the Redshift cluster | Monitoring dashboards, operational scripts |
| `cluster_arn` | Amazon Resource Name of the cluster | IAM policies, CloudWatch alarms, resource tagging |
| `cluster_namespace_arn` | Namespace ARN of the cluster | Redshift data sharing, the Redshift Data API |
| `endpoint` | Connection endpoint in address:port format | Application connection strings, BI tool configuration |
| `dns_name` | DNS hostname of the cluster's leader node (without port) | DNS CNAME records, connection strings |
| `database_name` | Name of the first database in the cluster | Application connection configuration |
| `port` | TCP port for client connections | Application connection configuration, security group rules |
| `subnet_group_name` | Redshift subnet group in use (module-managed or referenced) | Audit, related resource lookups |
| `parameter_group_name` | Cluster parameter group in use (module-managed or referenced) | Parameter auditing |
| `master_password_secret_arn` | Secrets Manager secret ARN (when the password is AWS-managed) | Application secret retrieval, rotation monitoring |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Single-node development** -- An RA3 single-node cluster for development and prototyping with minimal cost. No Multi-AZ, no audit logging, managed master password. Start from the **Single Node Dev** preset.

**Multi-node production** -- RA3 nodes with managed master password, encrypted storage, 7-day automated snapshot retention, and deletion protection via final snapshot. Suitable for production analytics workloads. Start from the **Multi-Node Production** preset.

**Analytics workload** -- Multi-node RA3 cluster with IAM roles attached for Spectrum queries, enhanced VPC routing for network governance, CloudWatch audit logging, cross-region snapshot copy, and Multi-AZ. Designed for data lake analytics with S3 integration. Start from the **Analytics Workload** preset.

## Works With

- [**AWS Subnet**](/cloud-catalog/aws-subnet) -- provides the subnets the Redshift subnet group is built from
- [**AWS Security Group**](/cloud-catalog/aws-security-group) -- provides network access control for the cluster endpoint
- [**AWS Elastic IP**](/cloud-catalog/aws-elastic-ip) -- provides a stable public address for a publicly accessible cluster
- [**AWS KMS Key**](/cloud-catalog/aws-kms-key) -- provides a customer-managed key for cluster encryption and Secrets Manager password encryption
- [**AWS IAM Role**](/cloud-catalog/aws-iam-role) -- provides service roles for COPY, UNLOAD, and Redshift Spectrum access to S3 and Glue
