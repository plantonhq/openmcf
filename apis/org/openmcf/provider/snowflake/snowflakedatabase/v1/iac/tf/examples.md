# Snowflake Database Terraform Module Examples

This document provides practical examples of using the Snowflake Database Terraform module for various use cases.

## Basic Example

A minimal configuration for creating a production database:

```hcl
module "analytics_database" {
  source = "./iac/tf"

  metadata = {
    name = "analytics-db"
    id   = "snowdb-analytics-001"
    org  = "engineering"
    env  = "production"
  }

  spec = {
    name                        = "ANALYTICS_PROD"
    comment                     = "Production analytics database"
    is_transient                = false
    data_retention_time_in_days = 7

    user_task = {
      managed_initial_warehouse_size = "MEDIUM"
      minimum_trigger_interval_in_seconds = 60
      timeout_ms = 3600000
    }

    # Optional parameters with defaults
    catalog                     = ""
    default_ddl_collation       = ""
    drop_public_schema_on_creation = false
    enable_console_output       = true
    external_volume             = ""
    log_level                   = "INFO"
    max_data_extension_time_in_days = 0
    quoted_identifiers_ignore_case = false
    replace_invalid_characters  = false
    storage_serialization_policy = ""
    suspend_task_after_num_failures = 10
    task_auto_retry_attempts    = 3
    trace_level                 = "ON_EVENT"
  }

  snowflake_credential = {
    account  = var.snowflake_account
    region   = var.snowflake_region
    username = var.snowflake_username
    password = var.snowflake_password
  }
}
```

## Development Database (Cost-Optimized)

Creating a cost-optimized transient database for development environments:

```hcl
module "dev_database" {
  source = "./iac/tf"

  metadata = {
    name = "dev-database"
    id   = "snowdb-dev-001"
    org  = "engineering"
    env  = "development"
    labels = {
      team        = "data-engineering"
      cost-center = "development"
      temporary   = "true"
    }
  }

  spec = {
    name                        = "ANALYTICS_DEV"
    comment                     = "Development database - transient for cost savings"
    is_transient                = true  # Eliminates 7-day Fail-safe period
    data_retention_time_in_days = 1     # Minimal Time Travel retention

    user_task = {
      managed_initial_warehouse_size = "XSMALL"
      minimum_trigger_interval_in_seconds = 0
      timeout_ms = 0
    }

    catalog                     = ""
    default_ddl_collation       = ""
    drop_public_schema_on_creation = true  # Clean slate for dev
    enable_console_output       = true     # Helpful for debugging
    external_volume             = ""
    log_level                   = "DEBUG"  # Verbose logging for development
    max_data_extension_time_in_days = 0
    quoted_identifiers_ignore_case = false
    replace_invalid_characters  = false
    storage_serialization_policy = ""
    suspend_task_after_num_failures = 0
    task_auto_retry_attempts    = 0
    trace_level                 = "OFF"
  }

  snowflake_credential = {
    account  = var.snowflake_account
    region   = var.snowflake_region
    username = var.snowflake_username
    password = var.snowflake_password
  }
}
```

## CI/CD Pipeline Database

Ultra-minimal configuration for automated testing pipelines:

```hcl
module "ci_database" {
  source = "./iac/tf"

  metadata = {
    name = "ci-test-db"
    id   = "snowdb-ci-${var.pipeline_run_id}"
    org  = "engineering"
    env  = "ci"
    labels = {
      pipeline    = var.pipeline_name
      ephemeral   = "true"
      auto-delete = "true"
    }
  }

  spec = {
    name                        = "CI_TEST_${upper(var.pipeline_run_id)}"
    comment                     = "Ephemeral CI database - pipeline ${var.pipeline_name}"
    is_transient                = true
    data_retention_time_in_days = 0  # No retention needed for CI

    user_task = {
      managed_initial_warehouse_size = ""
      minimum_trigger_interval_in_seconds = 0
      timeout_ms = 0
    }

    catalog                     = ""
    default_ddl_collation       = ""
    drop_public_schema_on_creation = true
    enable_console_output       = false
    external_volume             = ""
    log_level                   = ""
    max_data_extension_time_in_days = 0
    quoted_identifiers_ignore_case = false
    replace_invalid_characters  = false
    storage_serialization_policy = ""
    suspend_task_after_num_failures = 0
    task_auto_retry_attempts    = 0
    trace_level                 = ""
  }

  snowflake_credential = {
    account  = var.snowflake_account
    region   = var.snowflake_region
    username = var.snowflake_ci_username
    password = var.snowflake_ci_password
  }
}
```

## Iceberg-Enabled Lakehouse Database

Database configured for Apache Iceberg tables with external storage:

```hcl
module "lakehouse_database" {
  source = "./iac/tf"

  metadata = {
    name = "lakehouse-prod"
    id   = "snowdb-lakehouse-001"
    org  = "data-platform"
    env  = "production"
    labels = {
      team          = "data-lake"
      storage-type  = "iceberg"
      external-data = "true"
    }
  }

  spec = {
    name                        = "LAKEHOUSE_PROD"
    comment                     = "Production lakehouse with Iceberg support"
    is_transient                = false
    data_retention_time_in_days = 14

    # Iceberg-specific configuration
    catalog                     = "iceberg_catalog"
    external_volume             = "s3_external_volume"
    storage_serialization_policy = "OPTIMIZED"

    user_task = {
      managed_initial_warehouse_size = "LARGE"
      minimum_trigger_interval_in_seconds = 300
      timeout_ms = 7200000
    }

    default_ddl_collation       = ""
    drop_public_schema_on_creation = false
    enable_console_output       = false
    log_level                   = "WARN"
    max_data_extension_time_in_days = 30
    quoted_identifiers_ignore_case = false
    replace_invalid_characters  = true
    suspend_task_after_num_failures = 5
    task_auto_retry_attempts    = 3
    trace_level                 = "ON_EVENT"
  }

  snowflake_credential = {
    account  = var.snowflake_account
    region   = var.snowflake_region
    username = var.snowflake_username
    password = var.snowflake_password
  }
}
```

## Multi-Environment Pattern

Managing multiple environments with shared configuration:

```hcl
locals {
  environments = {
    prod = {
      is_transient   = false
      retention_days = 30
      warehouse_size = "LARGE"
      log_level      = "WARN"
    }
    staging = {
      is_transient   = false
      retention_days = 7
      warehouse_size = "MEDIUM"
      log_level      = "INFO"
    }
    dev = {
      is_transient   = true
      retention_days = 1
      warehouse_size = "XSMALL"
      log_level      = "DEBUG"
    }
  }
}

module "databases" {
  source = "./iac/tf"
  
  for_each = local.environments

  metadata = {
    name = "analytics-${each.key}"
    id   = "snowdb-analytics-${each.key}"
    org  = "engineering"
    env  = each.key
    labels = {
      app     = "analytics"
      managed = "terraform"
    }
  }

  spec = {
    name                        = "ANALYTICS_${upper(each.key)}"
    comment                     = "${title(each.key)} analytics database"
    is_transient                = each.value.is_transient
    data_retention_time_in_days = each.value.retention_days

    user_task = {
      managed_initial_warehouse_size = each.value.warehouse_size
      minimum_trigger_interval_in_seconds = 60
      timeout_ms = 3600000
    }

    catalog                     = ""
    default_ddl_collation       = ""
    drop_public_schema_on_creation = false
    enable_console_output       = true
    external_volume             = ""
    log_level                   = each.value.log_level
    max_data_extension_time_in_days = 0
    quoted_identifiers_ignore_case = false
    replace_invalid_characters  = false
    storage_serialization_policy = ""
    suspend_task_after_num_failures = 10
    task_auto_retry_attempts    = 3
    trace_level                 = "ON_EVENT"
  }

  snowflake_credential = {
    account  = var.snowflake_account
    region   = var.snowflake_region
    username = var.snowflake_username
    password = var.snowflake_password
  }
}
```

## High-Compliance Database

Database with extended retention and strict auditing:

```hcl
module "compliance_database" {
  source = "./iac/tf"

  metadata = {
    name = "compliance-db"
    id   = "snowdb-compliance-001"
    org  = "finance"
    env  = "production"
    labels = {
      compliance    = "sox-hipaa"
      retention     = "extended"
      audit-enabled = "true"
    }
  }

  spec = {
    name                        = "COMPLIANCE_PROD"
    comment                     = "High-compliance database with extended retention"
    is_transient                = false
    data_retention_time_in_days = 90  # Maximum Time Travel retention

    user_task = {
      managed_initial_warehouse_size = "MEDIUM"
      minimum_trigger_interval_in_seconds = 300
      timeout_ms = 7200000
    }

    catalog                     = ""
    default_ddl_collation       = "en_US"
    drop_public_schema_on_creation = false
    enable_console_output       = false
    external_volume             = ""
    log_level                   = "INFO"
    max_data_extension_time_in_days = 60
    quoted_identifiers_ignore_case = false
    replace_invalid_characters  = false
    storage_serialization_policy = "COMPATIBLE"
    suspend_task_after_num_failures = 20
    task_auto_retry_attempts    = 5
    trace_level                 = "ALWAYS"
  }

  snowflake_credential = {
    account  = var.snowflake_account
    region   = var.snowflake_region
    username = var.snowflake_username
    password = var.snowflake_password
  }
}
```

## SaaS Multi-Tenant Pattern

Creating tenant-specific databases programmatically:

```hcl
variable "tenants" {
  description = "List of tenant configurations"
  type = map(object({
    org_id         = string
    retention_days = number
  }))
}

module "tenant_databases" {
  source = "./iac/tf"
  
  for_each = var.tenants

  metadata = {
    name = "tenant-${each.key}"
    id   = "snowdb-tenant-${each.key}"
    org  = each.value.org_id
    env  = "production"
    labels = {
      tenant-id = each.key
      multi-tenant = "true"
    }
  }

  spec = {
    name                        = "TENANT_${upper(replace(each.key, "-", "_"))}"
    comment                     = "Database for tenant ${each.key}"
    is_transient                = false
    data_retention_time_in_days = each.value.retention_days

    user_task = {
      managed_initial_warehouse_size = "SMALL"
      minimum_trigger_interval_in_seconds = 120
      timeout_ms = 3600000
    }

    catalog                     = ""
    default_ddl_collation       = ""
    drop_public_schema_on_creation = false
    enable_console_output       = false
    external_volume             = ""
    log_level                   = "WARN"
    max_data_extension_time_in_days = 0
    quoted_identifiers_ignore_case = false
    replace_invalid_characters  = false
    storage_serialization_policy = ""
    suspend_task_after_num_failures = 10
    task_auto_retry_attempts    = 3
    trace_level                 = "ON_EVENT"
  }

  snowflake_credential = {
    account  = var.snowflake_account
    region   = var.snowflake_region
    username = var.snowflake_username
    password = var.snowflake_password
  }
}

# Example tenants configuration in terraform.tfvars:
# tenants = {
#   "customer-a" = {
#     org_id         = "org-customer-a"
#     retention_days = 30
#   }
#   "customer-b" = {
#     org_id         = "org-customer-b"
#     retention_days = 14
#   }
# }
```

## Variables File Example

Create a `terraform.tfvars` file to store your configuration:

```hcl
# terraform.tfvars

snowflake_account  = "xy12345.us-east-1"
snowflake_region   = "us-east-1"
snowflake_username = "terraform_service_account"

# IMPORTANT: Never commit passwords to version control
# Use environment variables instead:
# export TF_VAR_snowflake_password="your-password"
```

## Using with Terraform Cloud/Enterprise

When using Terraform Cloud or Enterprise, store credentials as sensitive variables:

```hcl
variable "snowflake_account" {
  description = "Snowflake account identifier"
  type        = string
}

variable "snowflake_region" {
  description = "Snowflake region"
  type        = string
}

variable "snowflake_username" {
  description = "Snowflake username"
  type        = string
}

variable "snowflake_password" {
  description = "Snowflake password"
  type        = string
  sensitive   = true
}

module "database" {
  source = "./iac/tf"

  # ... metadata and spec configuration ...

  snowflake_credential = {
    account  = var.snowflake_account
    region   = var.snowflake_region
    username = var.snowflake_username
    password = var.snowflake_password
  }
}
```

## Outputs Usage

Access module outputs for integration with other resources:

```hcl
module "database" {
  source = "./iac/tf"
  # ... configuration ...
}

output "database_id" {
  description = "Snowflake database ID"
  value       = module.database.id
}

output "database_name" {
  description = "Snowflake database name"
  value       = module.database.name
}

output "cost_indicators" {
  description = "Cost-related configuration"
  value = {
    is_transient   = module.database.is_transient
    retention_days = module.database.data_retention_time_in_days
  }
}

# Use outputs in other resources
resource "snowflake_schema" "example" {
  database = module.database.name
  name     = "PUBLIC"
}
```

## Cost Optimization Tips

1. **Use transient for non-production**: Development, staging, and CI databases should typically be transient
2. **Minimize retention**: Set `data_retention_time_in_days` to the minimum required
3. **Right-size warehouses**: Start with smaller warehouse sizes and scale up based on usage
4. **Tag for cost allocation**: Use metadata labels to track costs by team, project, or environment
5. **Consider lifecycle**: For ephemeral databases (CI/CD), ensure cleanup automation

## Best Practices

1. **Credential Management**: Use environment variables or secrets management tools
2. **Naming Conventions**: Use consistent, descriptive database names
3. **Documentation**: Add meaningful comments to spec
4. **Version Control**: Store Terraform configurations in git
5. **State Management**: Use remote state (S3, Terraform Cloud) for team collaboration
6. **Module Versioning**: Pin module versions in production
7. **Testing**: Test changes in dev/staging before applying to production

## Troubleshooting

### Authentication Errors

If you encounter authentication errors, verify:
- Account identifier includes region if required
- Username and password are correct
- User has `CREATE DATABASE` privilege

### Permission Issues

Grant necessary privileges:
```sql
GRANT CREATE DATABASE ON ACCOUNT TO ROLE terraform_role;
GRANT USAGE ON WAREHOUSE compute_wh TO ROLE terraform_role;
```

### State Conflicts

If multiple users manage the same infrastructure:
- Use remote state locking
- Coordinate changes through pull requests
- Consider Terraform Cloud for team collaboration

## Additional Resources

- [Terraform Best Practices](https://www.terraform.io/docs/cloud/guides/recommended-practices/index.html)
- [Snowflake Provider Documentation](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs)
- [Snowflake Database Parameters](https://docs.snowflake.com/en/sql-reference/parameters)




