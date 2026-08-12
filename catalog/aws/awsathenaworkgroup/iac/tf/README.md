# Terraform Module: AWS Athena Workgroup

Provisions an Amazon Athena workgroup using Terraform.

## Resources Created

- `aws_athena_workgroup` — The Athena workgroup carrying the full
  configuration surface: result storage (customer S3 or AWS-managed),
  per-query spend limits, engine selection (SQL or Spark), Identity
  Center integration, S3 Access Grants, and the three monitoring/logging
  arms (CloudWatch, managed, S3).

## Usage

```hcl
module "athena_workgroup" {
  source = "./path/to/module"

  metadata = {
    name = "analytics"
    org  = "my-org"
    env  = "production"
    id   = "awsathena-abc123"
  }

  spec = {
    region = "us-west-2"

    result_configuration = {
      output_location   = "s3://my-results-bucket/athena/"
      encryption_option = "SSE_S3"
    }

    enforce_workgroup_configuration     = true
    publish_cloudwatch_metrics_enabled  = true
    bytes_scanned_cutoff_per_query      = 10737418240
  }
}
```

Result storage is exactly one of customer S3 (`result_configuration`
with `output_location`) or AWS-managed (`managed_query_results`); the
spec validation enforces the exclusivity. The `managed_query_results`
KMS key is the one KMS field on this surface that accepts a full key
ARN only (its three siblings also accept `alias/...` forms).

## Inputs

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `metadata` | object | yes | Resource metadata (name, org, env, id) |
| `spec` | object | yes | AwsAthenaWorkgroup spec |

## Outputs

| Name | Description |
|------|-------------|
| `workgroup_arn` | ARN of the workgroup |
| `workgroup_name` | Name of the workgroup |
| `effective_engine_version` | The engine version Athena actually runs (resolved from AUTO) |
