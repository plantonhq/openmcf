# AwsNeptuneCluster

An Amazon Neptune graph database cluster -- property graphs (Apache TinkerPop Gremlin, openCypher) and RDF (SPARQL) over shared cluster storage, provisioned or Serverless.

The cluster is the shared-storage brain -- endpoints, backups, encryption, and engine lifecycle. The compute that serves queries is the folded `instances` list: each entry materializes as its own DB instance inside the cluster (a writer plus any readers), managed per-name so scaling readers never touches the cluster. Subnets, security groups, KMS keys, and IAM roles compose by reference.

Neptune has no master username or password -- access is network reachability plus (optionally) IAM database authentication with SigV4-signed requests.

## Spec highlights

- **Cluster shapes** -- provisioned (`instances` with provisioned classes) or Neptune Serverless (`serverlessV2Scaling` NCU bounds 1-128 + `db.serverless` instances).
- **Authentication** -- `iamDatabaseAuthenticationEnabled` maps IAM identities to database access; engine `iamRoles` power the S3 bulk loader and Neptune ML.
- **Data protection** -- storage encryption (create-time one-way door), continuous backups with `backupRetentionPeriod`, deletion protection, and final-snapshot enforcement.
- **Restore and replication shapes** -- create from a snapshot (`snapshotIdentifier`), as a replica (`replicationSourceIdentifier`), or join a Neptune global database (`globalClusterIdentifier`).
- **Observability** -- `audit` and `slowquery` CloudWatch log exports (paired with their cluster parameters).
- **Parameters** -- inline `parameters` (a module-managed cluster parameter group with the family derived from the pinned engine version) or an existing `neptuneClusterParameterGroupName`; `neptuneInstanceParameterGroupName` accompanies major version upgrades.

## Stack outputs

`cluster_identifier`, `arn`, `cluster_resource_id`, `endpoint` (writer), `reader_endpoint`, `port`, `hosted_zone_id`, `engine_version_actual`, `neptune_subnet_group_name`, `neptune_cluster_parameter_group_name`, `instance_endpoints`.

## How it works

Planton provisions via the Pulumi or Terraform module in `iac/`, both implementing the same contract at full parity. The API contract is protobuf-based (`spec.proto`); stack execution is orchestrated using `AwsNeptuneClusterStackInput` (provider credentials + IaC info).

## References

- Neptune user guide: https://docs.aws.amazon.com/neptune/latest/userguide/intro.html
- Neptune Serverless: https://docs.aws.amazon.com/neptune/latest/userguide/neptune-serverless.html
- IAM database authentication: https://docs.aws.amazon.com/neptune/latest/userguide/iam-auth.html
