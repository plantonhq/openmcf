# AWS Redshift Cluster

Deploys an Amazon Redshift provisioned cluster -- a petabyte-scale columnar data warehouse for analytical (OLAP) queries on structured and semi-structured data -- with automatic subnet group creation, optional Secrets Manager password management, KMS encryption, audit logging, cross-region snapshot copy, and inline parameter group support. Security groups, IAM roles, KMS keys, and Elastic IPs compose by reference; warehouse ingress rules live on the referenced `AwsSecurityGroup` nodes.

## What Gets Created

When you deploy an AwsRedshiftCluster resource, Planton provisions:

- **Redshift Cluster** — a `redshift.Cluster` with the specified node type, node count, credentials, encryption settings, snapshot configuration, and optional Multi-AZ deployment or availability-zone relocation
- **Subnet Group** — a `redshift.SubnetGroup` created automatically when `subnetIds` are provided and `clusterSubnetGroupName` is not set, placing the cluster across the specified subnets
- **Parameter Group** — a `redshift.ParameterGroup` (family `redshift-1.0`) created when inline `parameters` are provided
- **Logging Configuration** — audit logging enabled when `logging` is specified, sending connection/user-activity/user logs to S3 or CloudWatch Logs
- **Snapshot Copy Configuration** — cross-region snapshot copy enabled when `snapshotCopy` is specified, for disaster recovery

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **At least two subnets** in different Availability Zones, or an existing Redshift subnet group name
- **Security group IDs** if the cluster should not use the VPC's default security group -- ingress rules (e.g. port 5439 from BI tooling) belong on the referenced `AwsSecurityGroup` nodes
- **A KMS key ARN** if enabling encryption with a customer-managed key or encrypting the managed password secret
- **IAM role ARNs** if the cluster needs to access S3, DynamoDB, Glue Data Catalog, or other AWS services

## Quick Start

Create a file `redshift-cluster.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRedshiftCluster
metadata:
  name: my-warehouse
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsRedshiftCluster.my-warehouse
spec:
  region: us-west-2
  nodeType: ra3.large
  subnetIds:
    - value: "<private-subnet-id-az1>"
    - value: "<private-subnet-id-az2>"
  masterUsername: admin
  manageMasterPassword: true
  skipFinalSnapshot: true
```

Deploy:

```shell
planton apply -f redshift-cluster.yaml
```

This creates a single-node ra3.large Redshift cluster across two subnets, encrypted at rest (the AWS default this component preserves), with the admin password generated, stored, and rotated by AWS Secrets Manager.

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `region` | `string` | AWS region where the Redshift cluster will be created. Example: `us-west-2`, `eu-west-1`. |
| `nodeType` | `string` | Compute and storage capacity of each node: `ra3.large`, `ra3.xlplus`, `ra3.4xlarge`, `ra3.16xlarge` (managed storage -- recommended), or the legacy dense-compute `dc2` family. |
| `subnetIds` | `StringValueOrRef[]` | Subnet IDs for automatic Redshift subnet group creation. Provide at least two in distinct AZs. Can reference `AwsSubnet` outputs via `valueFrom`. Not required when `clusterSubnetGroupName` is set. |
| `masterUsername` | `string` | Admin username. Required for a new cluster; only snapshot restores leave it empty and inherit the source's credentials. Create-time only. |
| `finalSnapshotIdentifier` | `string` | Identifier for the final snapshot created on cluster deletion. Required when `skipFinalSnapshot` is `false`. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `numberOfNodes` | `int32` | `1` | Cluster size, 1-128. `1` = single-node (leader+compute combined); `2+` = multi-node (dedicated leader + N compute nodes). Resize is in-place. |
| `clusterVersion` | `string` | AWS default | Redshift engine version family (`"1.0"` -- the only family). Engine patches ride `maintenanceTrackName` and `allowVersionUpgrade`. |
| `databaseName` | `string` | AWS default (`dev`) | Name of the first database created in the cluster. 1-64 lowercase alphanumeric/underscore characters. |
| `masterPassword` | `string` | — | Admin password (8-64 chars, mixed case + digit), stored in IaC state. Mutually exclusive with `manageMasterPassword` -- prefer the managed path. |
| `manageMasterPassword` | `bool` | `false` (recommended: `true`) | When `true`, AWS Secrets Manager generates, rotates, and stores the admin password; the secret's ARN is exported. Mutually exclusive with `masterPassword`. |
| `masterPasswordSecretKmsKeyId` | `StringValueOrRef` | — | KMS key to encrypt the Secrets Manager secret holding the managed password. Only used when `manageMasterPassword` is `true`. Can reference `AwsKmsKey` via `valueFrom`. |
| `port` | `int32` | `5439` | TCP port for client connections. Valid range: 1115-65535 (5431-5455 or 8191-8215 when AZ relocation is enabled). |
| `clusterSubnetGroupName` | `StringValueOrRef` | — | Name of an existing Redshift subnet group. When set, `subnetIds` is not required. |
| `securityGroupIds` | `StringValueOrRef[]` | VPC default SG | Security groups attached to the cluster. Warehouse ingress rules live on the referenced `AwsSecurityGroup` nodes. Can reference `AwsSecurityGroup` via `valueFrom`. |
| `availabilityZone` | `string` | AWS placement | Pin the cluster to one AZ. Changing it on a live cluster requires `availabilityZoneRelocationEnabled`. |
| `availabilityZoneRelocationEnabled` | `bool` | `false` | Allow the cluster to be relocated to another AZ during outages or on demand. Requires RA3 node types; mutually exclusive with `multiAz`. |
| `publiclyAccessible` | `bool` | `false` | When `true`, the cluster gets a public IP and is reachable from outside the VPC. |
| `elasticIp` | `StringValueOrRef` | — | Static public IPv4 ADDRESS for the leader node (the IP itself, not an allocation ID). Requires `publiclyAccessible`. Can reference `AwsElasticIp` `public_ip` via `valueFrom`. |
| `enhancedVpcRouting` | `bool` | `false` | Forces all COPY/UNLOAD traffic through the VPC, enabling VPC flow logs and endpoint policies. |
| `multiAz` | `bool` | `false` | Multi-AZ deployment with automatic failover to a standby. Requires RA3 node types and a multi-node cluster; mutually exclusive with `availabilityZoneRelocationEnabled`. |
| `encrypted` | `bool` | `true` | At-rest encryption (the AWS default). Uses the AWS-managed Redshift service key unless `kmsKeyId` is specified. |
| `kmsKeyId` | `StringValueOrRef` | — | Customer-managed KMS key for cluster storage encryption. Requires `encrypted: true`. Can reference `AwsKmsKey` via `valueFrom`. |
| `iamRoles` | `StringValueOrRef[]` | `[]` | IAM roles the cluster assumes for COPY/UNLOAD/Spectrum access to S3, DynamoDB, Glue, etc. Maximum 10 roles. Can reference `AwsIamRole` via `valueFrom`. |
| `defaultIamRoleArn` | `StringValueOrRef` | — | IAM role used by default when SQL commands do not specify a role. Must also be present in `iamRoles`. Can reference `AwsIamRole` via `valueFrom`. |
| `automatedSnapshotRetentionPeriod` | `int32` | `1` | Days to retain automated snapshots, 0-35. `0` disables automated snapshots. |
| `manualSnapshotRetentionPeriod` | `int32` | AWS default (indefinite) | Days manual snapshots are retained: 1-3653, or `-1` for indefinite. |
| `skipFinalSnapshot` | `bool` | `false` | When `true`, no final snapshot is created on deletion. Set to `true` only for ephemeral dev/test clusters. |
| `snapshotIdentifier` | `string` | — | Restore from an existing snapshot by NAME at create time. Mutually exclusive with `snapshotArn`. |
| `snapshotArn` | `string` | — | Restore from an existing snapshot by ARN at create time (cross-account/cross-region shares). Mutually exclusive with `snapshotIdentifier`. |
| `snapshotClusterIdentifier` | `string` | — | The cluster the source snapshot was taken from -- disambiguates shared snapshot NAMEs. Only meaningful with `snapshotIdentifier`. |
| `ownerAccount` | `string` | — | The AWS account that owns the source snapshot, for cross-account restores. Only meaningful alongside a restore source. |
| `preferredMaintenanceWindow` | `string` | AWS assigned | Weekly UTC maintenance window. Format: `ddd:hh:mi-ddd:hh:mi` (e.g., `sat:03:00-sat:04:00`). |
| `allowVersionUpgrade` | `bool` | `true` | Permits AWS to apply engine version upgrades during the maintenance window. |
| `maintenanceTrackName` | `string` | AWS default (`current`) | Cluster maintenance track: `current` (latest approved release), `trailing` (one release behind), or a named track inherited from a snapshot restore. |
| `applyImmediately` | `bool` | `false` | When `true`, modifications apply immediately instead of during the next maintenance window. |
| `logging` | `object` | — | Audit logging configuration. See sub-fields below. |
| `logging.logDestinationType` | `string` | — | Where audit logs are delivered. `"s3"` or `"cloudwatch"`. Required when `logging` is set. |
| `logging.s3BucketName` | `string` | — | S3 bucket for log delivery. Required when `logDestinationType` is `"s3"`. |
| `logging.s3KeyPrefix` | `string` | — | Prefix for log objects in the S3 bucket. |
| `logging.logExports` | `string[]` | `[]` | Log types to export: `connectionlog`, `useractivitylog`, `userlog`. Required when `logDestinationType` is `"cloudwatch"`. |
| `snapshotCopy` | `object` | — | Cross-region snapshot copy for disaster recovery. See sub-fields below. |
| `snapshotCopy.destinationRegion` | `string` | — | The region snapshots are copied to. Required when `snapshotCopy` is set. |
| `snapshotCopy.retentionPeriod` | `int32` | AWS default (`7`) | Days copied automated snapshots are retained in the destination region, 1-35. |
| `snapshotCopy.manualSnapshotRetentionPeriod` | `int32` | AWS default (indefinite) | Days copied manual snapshots are retained: 1-3653, or `-1` for indefinite. |
| `snapshotCopy.snapshotCopyGrantName` | `string` | — | Snapshot copy grant for KMS-encrypted clusters (required by AWS in that case). |
| `clusterParameterGroupName` | `string` | — | Name of an existing parameter group to associate. Mutually exclusive with inline `parameters`. |
| `parameters` | `AwsRedshiftClusterParameter[]` | `[]` | Inline parameters managed as a dedicated group. Each entry has `name` and `value`. Common parameters: `require_ssl`, `enable_user_activity_logging`, `max_concurrency_scaling_clusters`. |
| `parameterGroupFamily` | `string` | `redshift-1.0` | Parameter-group family for the managed group: `redshift-1.0` (accepted everywhere) or `redshift-2.0` (the Redshift patch 2.0 generation). Only meaningful alongside `parameters`. |

## Examples

### Single-Node Development Cluster

A minimal single-node cluster for development and testing. Uses ra3.large (the smallest orderable class), skips the final snapshot, and retains automated snapshots for 1 day:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRedshiftCluster
metadata:
  name: dev-warehouse
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsRedshiftCluster.dev-warehouse
spec:
  region: us-west-2
  nodeType: ra3.large
  numberOfNodes: 1
  masterUsername: admin
  manageMasterPassword: true
  subnetIds:
    - value: "<private-subnet-id-az1>"
    - value: "<private-subnet-id-az2>"
  skipFinalSnapshot: true
  automatedSnapshotRetentionPeriod: 1
```

### Production Multi-Node with Encryption and Logging

A 2-node RA3 cluster with customer-managed KMS encryption, SSL enforcement, enhanced VPC routing, CloudWatch audit logging, and a 7-day snapshot retention policy:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRedshiftCluster
metadata:
  name: prod-warehouse
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsRedshiftCluster.prod-warehouse
spec:
  region: us-west-2
  nodeType: ra3.xlplus
  numberOfNodes: 2
  databaseName: analytics
  masterUsername: admin
  manageMasterPassword: true
  masterPasswordSecretKmsKeyId:
    value: "<kms-key-arn>"
  subnetIds:
    - value: "<private-subnet-id-az1>"
    - value: "<private-subnet-id-az2>"
  encrypted: true
  kmsKeyId:
    value: "<kms-key-arn>"
  enhancedVpcRouting: true
  automatedSnapshotRetentionPeriod: 7
  skipFinalSnapshot: false
  finalSnapshotIdentifier: prod-warehouse-final
  preferredMaintenanceWindow: "sat:03:00-sat:04:00"
  allowVersionUpgrade: true
  logging:
    logDestinationType: cloudwatch
    logExports:
      - connectionlog
      - useractivitylog
      - userlog
  parameters:
    - name: require_ssl
      value: "true"
    - name: enable_user_activity_logging
      value: "true"
  iamRoles:
    - value: "<redshift-s3-access-role-arn>"
```

### Analytics Workload with Multi-AZ and Spectrum IAM Roles

A 4-node ra3.4xlarge cluster for large-scale analytics. Multi-AZ provides automatic failover, concurrency scaling handles query bursts with up to 5 additional transient clusters, and two IAM roles are attached for S3 data loading and Redshift Spectrum external table queries:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRedshiftCluster
metadata:
  name: analytics-cluster
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsRedshiftCluster.analytics-cluster
spec:
  region: us-west-2
  nodeType: ra3.4xlarge
  numberOfNodes: 4
  databaseName: datalake
  masterUsername: admin
  manageMasterPassword: true
  masterPasswordSecretKmsKeyId:
    value: "<kms-key-arn>"
  subnetIds:
    - value: "<private-subnet-id-az1>"
    - value: "<private-subnet-id-az2>"
  encrypted: true
  kmsKeyId:
    value: "<kms-key-arn>"
  enhancedVpcRouting: true
  multiAz: true
  automatedSnapshotRetentionPeriod: 14
  skipFinalSnapshot: false
  finalSnapshotIdentifier: analytics-cluster-final
  preferredMaintenanceWindow: "sun:02:00-sun:04:00"
  allowVersionUpgrade: true
  applyImmediately: false
  logging:
    logDestinationType: cloudwatch
    logExports:
      - connectionlog
      - useractivitylog
      - userlog
  parameters:
    - name: require_ssl
      value: "true"
    - name: enable_user_activity_logging
      value: "true"
    - name: max_concurrency_scaling_clusters
      value: "5"
  iamRoles:
    - value: "<redshift-s3-access-role-arn>"
    - value: "<redshift-spectrum-role-arn>"
  defaultIamRoleArn:
    value: "<redshift-s3-access-role-arn>"
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `clusterIdentifier` | `string` | The unique identifier of the Redshift cluster |
| `clusterArn` | `string` | The Amazon Resource Name of the cluster, used for IAM policies and cross-service references |
| `clusterNamespaceArn` | `string` | The namespace ARN, used for Redshift data sharing and Serverless integration |
| `endpoint` | `string` | The connection endpoint in `address:port` format for SQL client connections |
| `dnsName` | `string` | The DNS hostname of the cluster (without port), for use in connection strings |
| `databaseName` | `string` | The name of the default database in the cluster |
| `port` | `int32` | The TCP port on which the cluster accepts connections |
| `subnetGroupName` | `string` | The name of the Redshift subnet group in use (managed or referenced) |
| `parameterGroupName` | `string` | The name of the parameter group in use (managed or referenced) |
| `masterPasswordSecretArn` | `string` | The ARN of the Secrets Manager secret containing the master password (only when `manageMasterPassword` is `true`) |

## Related Components

- [AwsVpc](/docs/catalog/aws/awsvpc) — provides subnets and VPC ID for cluster placement
- [AwsSecurityGroup](/docs/catalog/aws/awssecuritygroup) — controls network access to the cluster
- [AwsKmsKey](/docs/catalog/aws/awskmskey) — provides KMS keys for cluster encryption and Secrets Manager
- [AwsIamRole](/docs/catalog/aws/awsiamrole) — provides IAM roles for S3 access, Spectrum, and other AWS service integrations
- [AwsS3Bucket](/docs/catalog/aws/awss3bucket) — stores data for COPY/UNLOAD operations and audit logs (when using S3 log destination)
