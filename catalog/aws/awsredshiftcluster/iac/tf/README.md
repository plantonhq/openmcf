# Terraform Module to Deploy AwsRedshiftCluster

This module provisions an Amazon Redshift provisioned cluster aligned with the
Planton API: the cluster itself plus its folded satellites (subnet group,
parameter group, audit logging, cross-region snapshot copy). Security groups,
IAM roles, KMS keys, and Elastic IPs attach by reference -- the module never
creates them.

## CLI (local backend)

```shell
planton tofu init --manifest ../../e2e/manifest.yaml
planton tofu plan --manifest ../../e2e/manifest.yaml
planton tofu apply --manifest ../../e2e/manifest.yaml --auto-approve
planton tofu destroy --manifest ../../e2e/manifest.yaml --auto-approve
```

Credentials are passed via the stack input through the CLI, not in `spec`.

## Files

- `variables.tf` (generated; do not edit)
- `provider.tf` — provider setup (`hashicorp/aws >= 6.0.0`; audit logging and
  snapshot copy are standalone resources on the v6 line)
- `locals.tf` — computed locals and folding flags
- `redshift_cluster.tf` — the `aws_redshift_cluster` resource
- `subnet_group.tf` — Redshift subnet group when `subnet_ids` are provided
- `cluster_param_group.tf` — parameter group from inline `parameters`
- `logging.tf` — audit logging when `spec.logging` is set
- `snapshot_copy.tf` — cross-region snapshot copy when `spec.snapshot_copy` is set
- `outputs.tf` — outputs matching `AwsRedshiftClusterStackOutputs`

## Conditional Resources

- **Subnet group** — created from `spec.subnet_ids` (≥ 2 entries); or bring an
  existing group via `spec.cluster_subnet_group_name`
- **Parameter group** — created from inline `spec.parameters`; or bring an
  existing group via `spec.cluster_parameter_group_name`
- **Logging** — created when `spec.logging` is provided
- **Snapshot copy** — created when `spec.snapshot_copy` is provided

## Outputs

| Name | Description |
|------|-------------|
| cluster_identifier | Unique identifier of the Redshift cluster |
| cluster_arn | ARN of the cluster |
| cluster_namespace_arn | Namespace ARN for data sharing and the Data API |
| endpoint | Connection endpoint (address:port) |
| dns_name | DNS hostname (without port) |
| database_name | First database name |
| port | TCP port for connections |
| subnet_group_name | Subnet group name in use (managed or referenced) |
| parameter_group_name | Parameter group name in use (managed or referenced) |
| master_password_secret_arn | Secrets Manager secret ARN (managed password only) |
