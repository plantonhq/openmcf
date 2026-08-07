# Pulumi Module to Deploy AwsRedshiftCluster

This Pulumi Go program deploys an Amazon Redshift provisioned cluster using the
Planton API and module: the cluster itself plus its folded satellites (subnet
group, parameter group, audit logging, cross-region snapshot copy). Security
groups, IAM roles, KMS keys, and Elastic IPs attach by reference -- the module
never creates them.

## Requirements

- Planton CLI built locally
- Valid AWS credential provided via the CLI stack input (not in `spec`)

## CLI commands

Preview:

```shell
planton pulumi preview \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .
```

Update (apply):

```shell
planton pulumi update \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir . \
  --yes
```

Destroy:

```shell
planton pulumi destroy \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir . \
  --yes
```

## Resources Created

1. **`redshift.Cluster`** — the core data warehouse cluster
2. **`redshift.SubnetGroup`** (conditional) — from `subnetIds`; or bring an
   existing group via `clusterSubnetGroupName`
3. **`redshift.ParameterGroup`** (conditional) — from inline `parameters`; or
   bring an existing group via `clusterParameterGroupName`
4. **`redshift.Logging`** (conditional) — when `logging` is configured
5. **`redshift.SnapshotCopy`** (conditional) — when `snapshotCopy` is configured

## Outputs

| Key | Description |
|-----|-------------|
| `cluster_identifier` | Unique identifier of the Redshift cluster |
| `cluster_arn` | ARN of the cluster |
| `cluster_namespace_arn` | Namespace ARN for data sharing and the Data API |
| `endpoint` | Connection endpoint (address:port) |
| `dns_name` | DNS hostname (without port) |
| `database_name` | First database name |
| `port` | TCP port for connections |
| `subnet_group_name` | Subnet group name in use (managed or referenced) |
| `parameter_group_name` | Parameter group name in use (managed or referenced) |
| `master_password_secret_arn` | Secrets Manager secret ARN (managed password only) |

## Debugging

Optionally enable debugging by setting a binary in `Pulumi.yaml` and using the
`debug.sh` script.
