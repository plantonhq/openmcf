# AwsMwaaEnvironment — Terraform IaC Module

Terraform module for provisioning AWS MWAA (Managed Workflows for Apache Airflow) environments using the Planton `AwsMwaaEnvironmentSpec`.

## Overview

This module creates:
- An MWAA Environment (`aws_mwaa_environment`) with configurable Airflow version, S3 source, IAM execution role, VPC networking, encryption, sizing, logging, maintenance, and worker replacement strategy.

Network ingress is composed, never embedded: the environment attaches the referenced `security_group_ids` directly, and the rules MWAA needs (self-referencing all-traffic, HTTPS 443 ingress, egress) live on those first-class security group nodes.

## Usage

```hcl
module "mwaa" {
  source = "./path/to/this/module"

  provider_config = {
    region = "us-east-1"
  }

  metadata = {
    id   = "prod-data-pipelines"
    name = "prod-data-pipelines"
    org  = "myorg"
    env  = "production"
  }

  spec = {
    airflow_version   = "2.10.1"
    source_bucket_arn = "arn:aws:s3:::prod-airflow-dags"
    dag_s3_path       = "dags/"
    execution_role_arn = "arn:aws:iam::111122223333:role/mwaa-prod-role"

    subnet_ids = ["subnet-aaa111", "subnet-bbb222"]

    security_group_ids = ["sg-0mwaa001"]
    environment_class  = "mw1.medium"
    min_workers        = 2
    max_workers        = 10

    logging_configuration = {
      task_logs = {
        enabled   = true
        log_level = "INFO"
      }
      worker_logs = {
        enabled   = true
        log_level = "INFO"
      }
    }
  }
}
```

## Inputs

| Variable | Type | Required | Description |
|----------|------|----------|-------------|
| `provider_config` | object | yes | AWS region and optional credentials |
| `metadata` | object | yes | Resource ID, name, org, env |
| `spec` | object | yes | `AwsMwaaEnvironmentSpec` — see `variables.tf` for full type |

See `variables.tf` for the complete type definition of `spec`, including all optional fields and their defaults.

### Spec Fields

| Field | Type | Default | Description |
|---|---|---|---|
| `airflow_version` | string | `""` | Apache Airflow version |
| `airflow_configuration_options` | map(string) | `{}` | Airflow config overrides (`section.property` format) |
| `source_bucket_arn` | string | **required** | S3 bucket ARN for DAGs/plugins/requirements |
| `dag_s3_path` | string | **required** | Relative path to DAG folder in S3 |
| `plugins_s3_path` | string | `""` | Relative path to plugins.zip |
| `plugins_s3_object_version` | string | `""` | S3 object version for plugins.zip |
| `requirements_s3_path` | string | `""` | Relative path to requirements.txt |
| `requirements_s3_object_version` | string | `""` | S3 object version for requirements.txt |
| `startup_script_s3_path` | string | `""` | Relative path to startup script |
| `startup_script_s3_object_version` | string | `""` | S3 object version for startup script |
| `execution_role_arn` | string | **required** | IAM execution role ARN |
| `subnet_ids` | list(string) | **required** | 2 private subnet IDs in different AZs |
| `security_group_ids` | list(string) | **required** (≥1) | Security groups attached to the MWAA VPC endpoints; the MWAA ingress pattern lives on these groups |
| `kms_key_arn` | string | `""` | KMS key ARN for at-rest encryption |
| `environment_class` | string | `""` | Environment class (mw1.micro through mw1.2xlarge) |
| `min_workers` | number | `0` | Min Celery workers |
| `max_workers` | number | `0` | Max Celery workers |
| `min_webservers` | number | `0` | Min webservers |
| `max_webservers` | number | `0` | Max webservers |
| `schedulers` | number | `0` | Number of schedulers |
| `webserver_access_mode` | string | `"PRIVATE_ONLY"` | `PRIVATE_ONLY` or `PUBLIC_ONLY` |
| `endpoint_management` | string | `""` | `SERVICE` or `CUSTOMER` |
| `logging_configuration` | object | `null` | Per-module log config (5 modules) |
| `weekly_maintenance_window_start` | string | `""` | Maintenance window (`DAY:HH:MM` UTC) |
| `worker_replacement_strategy` | string | `""` | `FORCED` or `GRACEFUL` |

## Outputs

| Output | Description |
|--------|-------------|
| `environment_arn` | ARN of the MWAA environment |
| `environment_name` | Environment name |
| `webserver_url` | Airflow web UI URL |
| `airflow_version` | Effective Airflow version |
| `service_role_arn` | AWS service role ARN |
| `environment_class` | Effective environment class |
| `status` | Current environment status |
| `created_at` | Environment creation timestamp |
| `database_vpc_endpoint_service` | Endpoint service name for the metadata database (CUSTOMER endpoint management) |
| `webserver_vpc_endpoint_service` | Endpoint service name for the webserver (CUSTOMER endpoint management) |

## File Structure

| File | Purpose |
|------|---------|
| `provider.tf` | AWS provider configuration (hashicorp/aws >= 6.11.0) |
| `variables.tf` | Input variable definitions (generator-owned contract) |
| `locals.tf` | Environment name basis and resource-identity tags |
| `main.tf` | The MWAA environment resource |
| `outputs.tf` | Output definitions |

## How It Works

### Logging Configuration

The logging block uses Terraform `dynamic` blocks to conditionally configure each of the 5 log modules. Each module is only included when its configuration is non-null:

```hcl
dynamic "dag_processing_logs" {
  for_each = logging_configuration.value.dag_processing_logs != null ? [...] : []
  content { ... }
}
```

### Conditional Fields

Most optional fields use the `non-empty → set, empty → null` pattern:

```hcl
airflow_version = var.spec.airflow_version != "" ? var.spec.airflow_version : null
min_workers     = var.spec.min_workers > 0 ? var.spec.min_workers : null
```

This ensures Terraform does not send zero-value fields to the AWS API, allowing AWS defaults to take effect.

## Prerequisites

- Terraform 1.5+ / OpenTofu
- AWS provider >= 6.11.0
- AWS credentials (via provider config or ambient)

## Running

```bash
# Navigate to the Terraform module directory
cd catalog/aws/awsmwaaenvironment/iac/tf

# Initialize providers
terraform init

# Preview changes
terraform plan -var-file=terraform.tfvars

# Apply changes
terraform apply -var-file=terraform.tfvars

# View outputs
terraform output

# Destroy resources
terraform destroy -var-file=terraform.tfvars
```

## Related

- [Spec reference](../../README.md)
