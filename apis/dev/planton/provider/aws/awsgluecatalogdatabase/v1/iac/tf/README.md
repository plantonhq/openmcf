# Terraform Module: AWS Glue Catalog Database

Provisions an AWS Glue Data Catalog database using Terraform.

## Resources Created

- `aws_glue_catalog_database` — The Glue Data Catalog database

## Usage

```hcl
module "glue_catalog_database" {
  source = "./path/to/module"

  metadata = {
    name = "analytics"
    org  = "my-org"
    env  = "production"
    id   = "awsglue-abc123"
  }

  spec = {
    region       = "us-west-2"
    description  = "Analytics data catalog for BI and ML pipelines"
    location_uri = "s3://prod-data-lake/databases/analytics/"
    parameters   = { classification = "parquet" }
    create_table_default_permissions = [
      { permissions = ["ALL"], principal = "IAM_ALLOWED_PRINCIPALS" }
    ]
  }
}
```

The module also carries the resource-link (`target_database`) and federated
(`federated_database`) creation shapes; a database is exactly one shape and
the spec validation enforces the exclusivity.

## Inputs

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `metadata` | object | yes | Resource metadata (name, org, env, id) |
| `spec` | object | yes | AwsGlueCatalogDatabase spec |

## Outputs

| Name | Description |
|------|-------------|
| `database_name` | Name of the Glue Data Catalog database |
| `database_arn` | ARN of the database |
| `catalog_id` | ID of the Glue Data Catalog (AWS Account ID) |
