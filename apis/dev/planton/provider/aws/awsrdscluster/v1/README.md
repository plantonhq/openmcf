# AwsRdsCluster

An RDS DB cluster: an Aurora MySQL/PostgreSQL cluster (provisioned or Serverless v2), an Aurora Serverless v1 cluster, or a Multi-AZ RDS cluster for the community mysql/postgres engines.

The cluster is the shared-storage brain -- endpoints, credentials, backups, encryption, and engine lifecycle. The compute that serves queries is the folded `instances` list: each entry materializes as its own DB instance inside the cluster (a writer plus any readers), managed per-name so scaling readers never touches the cluster. Subnets, security groups, KMS keys, and IAM roles compose by reference.

## Spec highlights

- **Cluster shapes** -- Aurora provisioned (`instances` with provisioned classes), Aurora Serverless v2 (`serverlessV2Scaling` + `db.serverless` instances, scale-to-zero with `minCapacity: 0`), legacy Aurora Serverless v1 (`engineMode: serverless` + `serverlessV1Scaling`), and Multi-AZ RDS clusters (`engine: mysql|postgres` + `dbClusterInstanceClass` + storage sizing).
- **Credentials** -- `manageMasterUserPassword` (recommended) keeps the master password in Secrets Manager, generated and rotated by AWS, with the secret's ARN exported; or supply `masterPassword` directly (sensitive).
- **Data protection** -- storage encryption (create-time one-way door), continuous backups with `backupRetentionPeriod`, deletion protection, final-snapshot enforcement, `deleteAutomatedBackups: false` retention, Aurora MySQL backtrack.
- **Restore shapes** -- create from a snapshot (`snapshotIdentifier`) or from another cluster's continuous backup (`restoreToPointInTime`, including Aurora copy-on-write fast clones).
- **Integrations** -- IAM database authentication, engine `iamRoles` (S3 import/export, Lambda, ML), the Data API (`enableHttpEndpoint`), CloudWatch log exports, Performance Insights, Enhanced Monitoring, Database Insights, Aurora Global Database membership, and local/global write forwarding.
- **Parameters** -- inline `parameters` (a module-managed cluster parameter group with the family derived from the pinned engine version) or an existing `dbClusterParameterGroupName`.

## Stack outputs

`cluster_identifier`, `arn`, `cluster_resource_id`, `endpoint` (writer), `reader_endpoint`, `port`, `hosted_zone_id`, `engine_version_actual`, `master_user_secret_arn`, `db_subnet_group_name`, `db_cluster_parameter_group_name`, `instance_endpoints`.

## How it works

Planton provisions via the Pulumi or Terraform module in `iac/`, both implementing the same contract at full parity. The API contract is protobuf-based (`spec.proto`); stack execution is orchestrated using `AwsRdsClusterStackInput` (provider credentials + IaC info).

## References

- Aurora user guide: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/CHAP_AuroraOverview.html
- Multi-AZ DB clusters: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/multi-az-db-clusters-concepts.html
- Aurora Serverless v2: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/aurora-serverless-v2.html
