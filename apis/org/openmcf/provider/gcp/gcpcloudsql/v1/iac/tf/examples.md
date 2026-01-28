# GCP Cloud SQL Terraform Module - Examples

This document provides practical examples for deploying Google Cloud SQL instances using the Terraform module.

## StringValueOrRef Fields

The `project_id` and `network.vpc_id` fields now support the `StringValueOrRef` type, which allows either:
- **Literal values**: `{ value = "my-value" }`
- **References to other resources**: `{ value_from = { kind = "GcpProject", name = "my-project", field_path = "status.outputs.project_id" } }`

> **Note**: Currently, only literal values are supported. Reference resolution (`value_from`) will be implemented in a future release.

## Example 1: Basic MySQL Instance

Create a basic MySQL 8.0 instance for development:

```hcl
module "mysql_dev" {
  source = "../../tf"

  metadata = {
    name = "mysql-dev"
  }

  spec = {
    project_id       = { value = "my-dev-project" }
    region           = "us-central1"
    database_engine  = "MYSQL"
    database_version = "MYSQL_8_0"
    tier             = "db-f1-micro"
    storage_gb       = 10
    root_password    = "DevPassword123!"
  }
}

output "mysql_dev_connection" {
  value = module.mysql_dev.connection_name
}
```

## Example 2: PostgreSQL with Private IP

Production PostgreSQL instance with private IP connectivity:

```hcl
module "postgres_private" {
  source = "../../tf"

  metadata = {
    name = "postgres-app"
    org  = "acme"
    env  = "production"
  }

  spec = {
    project_id       = { value = "acme-prod" }
    region           = "us-central1"
    database_engine  = "POSTGRESQL"
    database_version = "POSTGRES_15"
    tier             = "db-n1-standard-2"
    storage_gb       = 20

    network = {
      vpc_id             = { value = "projects/acme-prod/global/networks/app-vpc" }
      private_ip_enabled = true
    }

    root_password = var.postgres_root_password
  }
}

output "postgres_private_ip" {
  value     = module.postgres_private.private_ip
  sensitive = true
}
```

## Example 3: MySQL High Availability Setup

Highly available MySQL instance with automated backups:

```hcl
module "mysql_ha" {
  source = "../../tf"

  metadata = {
    name = "mysql-ha-prod"
    org  = "acme"
    env  = "production"
  }

  spec = {
    project_id       = { value = "acme-prod" }
    region           = "us-east1"
    database_engine  = "MYSQL"
    database_version = "MYSQL_8_0"
    tier             = "db-n1-standard-4"
    storage_gb       = 100

    high_availability = {
      enabled = true
      zone    = "us-east1-b"
    }

    backup = {
      enabled        = true
      start_time     = "03:00"
      retention_days = 14
    }

    database_flags = {
      max_connections        = "250"
      innodb_buffer_pool_size = "536870912"
      slow_query_log         = "on"
    }

    root_password = var.mysql_root_password
  }
}
```

## Example 4: PostgreSQL with Authorized Networks

Public PostgreSQL instance with restricted network access:

```hcl
module "postgres_public" {
  source = "../../tf"

  metadata = {
    name = "postgres-analytics"
  }

  spec = {
    project_id       = { value = "analytics-project" }
    region           = "europe-west1"
    database_engine  = "POSTGRESQL"
    database_version = "POSTGRES_15"
    tier             = "db-n1-highmem-2"
    storage_gb       = 50

    network = {
      authorized_networks = [
        "203.0.113.0/24",  # Office network
        "198.51.100.0/24", # VPN network
      ]
    }

    backup = {
      enabled        = true
      start_time     = "02:00"
      retention_days = 7
    }

    root_password = var.postgres_root_password
  }
}
```

## Example 5: Multi-Environment Setup

Deploy instances across multiple environments:

```hcl
locals {
  environments = {
    dev = {
      tier       = "db-f1-micro"
      storage_gb = 10
      backup_retention = 3
    }
    staging = {
      tier       = "db-n1-standard-1"
      storage_gb = 20
      backup_retention = 7
    }
    production = {
      tier       = "db-n1-standard-4"
      storage_gb = 100
      backup_retention = 30
    }
  }
}

module "postgres_instances" {
  source   = "../../tf"
  for_each = local.environments

  metadata = {
    name = "postgres-${each.key}"
    env  = each.key
  }

  spec = {
    project_id       = { value = "my-project" }
    region           = "us-central1"
    database_engine  = "POSTGRESQL"
    database_version = "POSTGRES_15"
    tier             = each.value.tier
    storage_gb       = each.value.storage_gb

    backup = {
      enabled        = true
      start_time     = "03:00"
      retention_days = each.value.backup_retention
    }

    high_availability = {
      enabled = each.key == "production"
      zone    = each.key == "production" ? "us-central1-b" : null
    }

    root_password = var.postgres_passwords[each.key]
  }
}
```

## Example 6: PostgreSQL with Custom Database Flags

Optimize PostgreSQL for analytics workloads:

```hcl
module "postgres_analytics" {
  source = "../../tf"

  metadata = {
    name = "postgres-analytics"
  }

  spec = {
    project_id       = { value = "analytics-project" }
    region           = "us-central1"
    database_engine  = "POSTGRESQL"
    database_version = "POSTGRES_15"
    tier             = "db-n1-highmem-8"
    storage_gb       = 500

    database_flags = {
      max_connections           = "200"
      shared_buffers            = "2097152"  # 2GB
      effective_cache_size      = "6291456"  # 6GB
      work_mem                  = "10485"    # 10MB
      maintenance_work_mem      = "524288"   # 512MB
      checkpoint_completion_target = "0.9"
      wal_buffers               = "2048"
      default_statistics_target = "100"
      random_page_cost          = "1.1"
      effective_io_concurrency  = "200"
    }

    network = {
      vpc_id             = { value = "projects/analytics-project/global/networks/analytics-vpc" }
      private_ip_enabled = true
    }

    backup = {
      enabled        = true
      start_time     = "02:00"
      retention_days = 14
    }

    root_password = var.postgres_root_password
  }
}
```

## Example 7: Using with Variables File

Create a `terraform.tfvars` file:

```hcl
# terraform.tfvars
metadata = {
  name = "my-database"
  org  = "myorg"
  env  = "production"
}

spec = {
  project_id       = { value = "my-gcp-project" }
  region           = "us-central1"
  database_engine  = "MYSQL"
  database_version = "MYSQL_8_0"
  tier             = "db-n1-standard-2"
  storage_gb       = 50

  network = {
    vpc_id             = { value = "projects/my-gcp-project/global/networks/prod-vpc" }
    private_ip_enabled = true
  }

  high_availability = {
    enabled = true
    zone    = "us-central1-b"
  }

  backup = {
    enabled        = true
    start_time     = "03:00"
    retention_days = 30
  }

  database_flags = {
    max_connections = "200"
  }

  root_password = "ProductionPassword123!"
}
```

Then use it:

```shell
terraform apply -var-file="terraform.tfvars"
```

## Deployment Workflow

### Step 1: Initialize
```shell
terraform init
```

### Step 2: Validate
```shell
terraform validate
```

### Step 3: Plan
```shell
terraform plan -out=tfplan
```

### Step 4: Apply
```shell
terraform apply tfplan
```

### Step 5: Verify Outputs
```shell
terraform output
```

## Managing Secrets

### Using Environment Variables

```shell
export TF_VAR_mysql_root_password="SecurePassword123!"
terraform apply
```

### Using Google Secret Manager

```hcl
data "google_secret_manager_secret_version" "db_password" {
  secret = "mysql-root-password"
}

module "mysql_instance" {
  source = "../../tf"

  # ... other configuration ...

  spec = {
    # ... other spec fields ...
    root_password = data.google_secret_manager_secret_version.db_password.secret_data
  }
}
```

## Clean Up

To destroy all resources:

```shell
terraform destroy -auto-approve
```

To destroy specific module:

```shell
terraform destroy -target=module.mysql_dev
```

## Notes

- Always use strong, randomly generated passwords for production instances
- Store sensitive outputs (like passwords and IPs) securely
- Use remote state backend for team collaboration
- Enable deletion protection for production instances
- Review GCP quotas before creating multiple instances
- Consider using Cloud SQL Proxy for secure connections

