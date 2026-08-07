# Terraform Module to Deploy AwsRedshiftServerlessWorkgroup

This module provisions an Amazon Redshift Serverless workgroup -- the compute
plane of the serverless warehouse -- aligned with the Planton API. The
namespace it serves, the subnets it places compute in, and the security groups
on its endpoint all attach by reference. The module never creates a resource
that deserves to be its own node.

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
- `provider.tf` — provider setup (`hashicorp/aws >= 6.0.0`;
  price_performance_target and track_name ride the v6 line)
- `locals.tf` — naming basis and identity tags
- `workgroup.tf` — the `aws_redshiftserverless_workgroup` resource
- `outputs.tf` — outputs matching `AwsRedshiftServerlessWorkgroupStackOutputs`

## Outputs

| Name | Description |
|------|-------------|
| workgroup_name | The handle the serverless APIs and credentials API address |
| workgroup_id | Unique identifier AWS assigns to the workgroup |
| arn | ARN of the workgroup |
| endpoint_address | DNS hostname SQL clients connect to |
| port | TCP port for connections |
