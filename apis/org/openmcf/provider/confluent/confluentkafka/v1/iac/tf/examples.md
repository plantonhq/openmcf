# Confluent Kafka Terraform Examples

This document provides complete Terraform examples for deploying Confluent Cloud Kafka clusters using the ConfluentKafka module.

## Table of Contents

- [Basic Development Cluster](#basic-development-cluster)
- [Standard Production Cluster](#standard-production-cluster)
- [Enterprise Cluster with Private Networking](#enterprise-cluster-with-private-networking)
- [Dedicated Cluster with Provisioned Capacity](#dedicated-cluster-with-provisioned-capacity)
- [Complete Example with All Features](#complete-example-with-all-features)

---

## Basic Development Cluster

A minimal configuration suitable for development and testing environments.

### main.tf

```hcl
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    confluent = {
      source  = "confluentinc/confluent"
      version = "~> 2.0"
    }
  }
}

module "dev_kafka" {
  source = "github.com/plantonhq/openmcf//apis/org/openmcf/provider/confluent/confluentkafka/v1/iac/tf"

  metadata = {
    name = "dev-kafka-cluster"
    labels = {
      environment = "development"
      team        = "platform"
    }
  }

  spec = {
    cloud          = "AWS"
    region         = "us-east-2"
    availability   = "SINGLE_ZONE"
    environment_id = "env-dev-abc123"
    cluster_type   = "BASIC"
  }

  confluent_api_key    = var.confluent_api_key
  confluent_api_secret = var.confluent_api_secret
}

output "bootstrap_endpoint" {
  value = module.dev_kafka.bootstrap_endpoint
}

output "cluster_id" {
  value = module.dev_kafka.id
}
```

### variables.tf

```hcl
variable "confluent_api_key" {
  description = "Confluent Cloud API Key"
  type        = string
  sensitive   = true
}

variable "confluent_api_secret" {
  description = "Confluent Cloud API Secret"
  type        = string
  sensitive   = true
}
```

### terraform.tfvars.example

```hcl
# Rename this file to terraform.tfvars and fill in your values
# DO NOT commit terraform.tfvars to version control

confluent_api_key    = "your-confluent-api-key"
confluent_api_secret = "your-confluent-api-secret"
```

### Usage

```bash
# Initialize Terraform
terraform init

# Plan the deployment
terraform plan

# Apply the configuration
terraform apply

# Get outputs
terraform output bootstrap_endpoint
```

---

## Standard Production Cluster

A production-ready cluster with multi-zone high availability and elastic scaling.

### main.tf

```hcl
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    confluent = {
      source  = "confluentinc/confluent"
      version = "~> 2.0"
    }
  }

  backend "s3" {
    bucket = "my-terraform-state"
    key    = "confluent/prod-kafka/terraform.tfstate"
    region = "us-east-1"
  }
}

module "prod_kafka" {
  source = "github.com/plantonhq/openmcf//apis/org/openmcf/provider/confluent/confluentkafka/v1/iac/tf"

  metadata = {
    name = "prod-orders-kafka"
    labels = {
      environment  = "production"
      team         = "platform"
      cost-center  = "engineering"
      criticality  = "high"
    }
  }

  spec = {
    cloud          = "GCP"
    region         = "us-central1"
    availability   = "MULTI_ZONE"
    environment_id = "env-prod-xyz789"
    cluster_type   = "STANDARD"
    display_name   = "Production Orders Kafka Cluster"
  }

  confluent_api_key    = var.confluent_api_key
  confluent_api_secret = var.confluent_api_secret
}

# Export all outputs for integration with other infrastructure
output "cluster_id" {
  description = "Kafka cluster ID"
  value       = module.prod_kafka.id
}

output "bootstrap_endpoint" {
  description = "Bootstrap endpoint for Kafka clients"
  value       = module.prod_kafka.bootstrap_endpoint
}

output "rest_endpoint" {
  description = "REST API endpoint"
  value       = module.prod_kafka.rest_endpoint
}

output "crn" {
  description = "Confluent Resource Name"
  value       = module.prod_kafka.crn
}
```

### variables.tf

```hcl
variable "confluent_api_key" {
  description = "Confluent Cloud API Key"
  type        = string
  sensitive   = true
}

variable "confluent_api_secret" {
  description = "Confluent Cloud API Secret"
  type        = string
  sensitive   = true
}
```

---

## Enterprise Cluster with Private Networking

An enterprise-grade cluster with AWS PrivateLink for secure, private connectivity.

### main.tf

```hcl
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    confluent = {
      source  = "confluentinc/confluent"
      version = "~> 2.0"
    }
  }
}

# First, create or reference a Confluent Cloud network
# This example assumes the network already exists
# If creating new, see: https://registry.terraform.io/providers/confluentinc/confluent/latest/docs/resources/confluent_network

module "enterprise_kafka" {
  source = "github.com/plantonhq/openmcf//apis/org/openmcf/provider/confluent/confluentkafka/v1/iac/tf"

  metadata = {
    name = "enterprise-secure-kafka"
    labels = {
      environment   = "production"
      security-zone = "restricted"
      compliance    = "pci-dss"
    }
  }

  spec = {
    cloud          = "AWS"
    region         = "us-west-2"
    availability   = "MULTI_ZONE"
    environment_id = var.confluent_environment_id
    cluster_type   = "ENTERPRISE"
    display_name   = "Enterprise Secure Kafka Cluster"
    
    network_config = {
      network_id = var.confluent_network_id
    }
  }

  confluent_api_key    = var.confluent_api_key
  confluent_api_secret = var.confluent_api_secret
}

# Outputs
output "bootstrap_endpoint" {
  description = "Private bootstrap endpoint (accessible via PrivateLink)"
  value       = module.enterprise_kafka.bootstrap_endpoint
}

output "cluster_id" {
  value = module.enterprise_kafka.id
}
```

### variables.tf

```hcl
variable "confluent_api_key" {
  description = "Confluent Cloud API Key"
  type        = string
  sensitive   = true
}

variable "confluent_api_secret" {
  description = "Confluent Cloud API Secret"
  type        = string
  sensitive   = true
}

variable "confluent_environment_id" {
  description = "Confluent Cloud Environment ID"
  type        = string
}

variable "confluent_network_id" {
  description = "Confluent Cloud Network ID (pre-created with PrivateLink)"
  type        = string
}
```

### terraform.tfvars.example

```hcl
confluent_environment_id = "env-prod-xyz789"
confluent_network_id     = "n-abc123"

# Set via environment variables:
# export TF_VAR_confluent_api_key="your-key"
# export TF_VAR_confluent_api_secret="your-secret"
```

---

## Dedicated Cluster with Provisioned Capacity

A dedicated, single-tenant cluster with 4 CKU for high-throughput workloads.

### main.tf

```hcl
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    confluent = {
      source  = "confluentinc/confluent"
      version = "~> 2.0"
    }
  }
}

module "dedicated_kafka" {
  source = "github.com/plantonhq/openmcf//apis/org/openmcf/provider/confluent/confluentkafka/v1/iac/tf"

  metadata = {
    name = "dedicated-high-throughput"
    labels = {
      environment = "production"
      workload    = "high-throughput"
      criticality = "tier-1"
      team        = "data-platform"
    }
  }

  spec = {
    cloud          = "AZURE"
    region         = "eastus"
    availability   = "MULTI_ZONE"
    environment_id = var.confluent_environment_id
    cluster_type   = "DEDICATED"
    display_name   = "High-Throughput Dedicated Kafka"
    
    dedicated_config = {
      cku = 4  # 4 Confluent Kafka Units
    }
    
    network_config = {
      network_id = var.confluent_network_id
    }
  }

  confluent_api_key    = var.confluent_api_key
  confluent_api_secret = var.confluent_api_secret
}

# Outputs
output "cluster_info" {
  description = "Complete cluster information"
  value = {
    id                 = module.dedicated_kafka.id
    bootstrap_endpoint = module.dedicated_kafka.bootstrap_endpoint
    rest_endpoint      = module.dedicated_kafka.rest_endpoint
    crn                = module.dedicated_kafka.crn
  }
}
```

### variables.tf

```hcl
variable "confluent_api_key" {
  description = "Confluent Cloud API Key"
  type        = string
  sensitive   = true
}

variable "confluent_api_secret" {
  description = "Confluent Cloud API Secret"
  type        = string
  sensitive   = true
}

variable "confluent_environment_id" {
  description = "Confluent Cloud Environment ID"
  type        = string
}

variable "confluent_network_id" {
  description = "Confluent Cloud Network ID for Azure Private Link"
  type        = string
}
```

---

## Complete Example with All Features

A comprehensive example demonstrating all available configuration options.

### Directory Structure

```
.
├── main.tf
├── variables.tf
├── outputs.tf
├── terraform.tfvars  (gitignored)
└── README.md
```

### main.tf

```hcl
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    confluent = {
      source  = "confluentinc/confluent"
      version = "~> 2.0"
    }
  }

  backend "s3" {
    bucket = "my-terraform-state"
    key    = "confluent/complete-example/terraform.tfstate"
    region = "us-east-1"
  }
}

# Development Cluster (Basic)
module "dev_kafka" {
  source = "github.com/plantonhq/openmcf//apis/org/openmcf/provider/confluent/confluentkafka/v1/iac/tf"

  metadata = {
    name = "dev-kafka"
    labels = {
      environment = "development"
      team        = "platform"
    }
  }

  spec = {
    cloud          = "AWS"
    region         = "us-east-2"
    availability   = "SINGLE_ZONE"
    environment_id = var.dev_environment_id
    cluster_type   = "BASIC"
  }

  confluent_api_key    = var.confluent_api_key
  confluent_api_secret = var.confluent_api_secret
}

# Staging Cluster (Standard)
module "staging_kafka" {
  source = "github.com/plantonhq/openmcf//apis/org/openmcf/provider/confluent/confluentkafka/v1/iac/tf"

  metadata = {
    name = "staging-kafka"
    labels = {
      environment = "staging"
      team        = "platform"
    }
  }

  spec = {
    cloud          = "GCP"
    region         = "us-central1"
    availability   = "MULTI_ZONE"
    environment_id = var.staging_environment_id
    cluster_type   = "STANDARD"
    display_name   = "Staging Kafka Cluster"
  }

  confluent_api_key    = var.confluent_api_key
  confluent_api_secret = var.confluent_api_secret
}

# Production Cluster (Dedicated with Private Networking)
module "prod_kafka" {
  source = "github.com/plantonhq/openmcf//apis/org/openmcf/provider/confluent/confluentkafka/v1/iac/tf"

  metadata = {
    name = "prod-kafka"
    labels = {
      environment  = "production"
      team         = "platform"
      criticality  = "high"
      compliance   = "sox"
    }
  }

  spec = {
    cloud          = "AWS"
    region         = "us-east-1"
    availability   = "MULTI_ZONE"
    environment_id = var.prod_environment_id
    cluster_type   = "DEDICATED"
    display_name   = "Production Kafka Cluster"
    
    dedicated_config = {
      cku = 4
    }
    
    network_config = {
      network_id = var.prod_network_id
    }
  }

  confluent_api_key    = var.confluent_api_key
  confluent_api_secret = var.confluent_api_secret
}
```

### variables.tf

```hcl
variable "confluent_api_key" {
  description = "Confluent Cloud API Key"
  type        = string
  sensitive   = true
}

variable "confluent_api_secret" {
  description = "Confluent Cloud API Secret"
  type        = string
  sensitive   = true
}

variable "dev_environment_id" {
  description = "Confluent Cloud Environment ID for development"
  type        = string
}

variable "staging_environment_id" {
  description = "Confluent Cloud Environment ID for staging"
  type        = string
}

variable "prod_environment_id" {
  description = "Confluent Cloud Environment ID for production"
  type        = string
}

variable "prod_network_id" {
  description = "Confluent Cloud Network ID for production private networking"
  type        = string
}
```

### outputs.tf

```hcl
# Development Outputs
output "dev_cluster_id" {
  value = module.dev_kafka.id
}

output "dev_bootstrap_endpoint" {
  value = module.dev_kafka.bootstrap_endpoint
}

# Staging Outputs
output "staging_cluster_id" {
  value = module.staging_kafka.id
}

output "staging_bootstrap_endpoint" {
  value = module.staging_kafka.bootstrap_endpoint
}

# Production Outputs
output "prod_cluster_id" {
  value       = module.prod_kafka.id
  description = "Production Kafka cluster ID"
}

output "prod_bootstrap_endpoint" {
  value       = module.prod_kafka.bootstrap_endpoint
  description = "Production bootstrap endpoint (private)"
  sensitive   = false
}

output "prod_rest_endpoint" {
  value = module.prod_kafka.rest_endpoint
}

output "prod_crn" {
  value = module.prod_kafka.crn
}

# Summary output
output "cluster_summary" {
  value = {
    development = {
      id       = module.dev_kafka.id
      endpoint = module.dev_kafka.bootstrap_endpoint
    }
    staging = {
      id       = module.staging_kafka.id
      endpoint = module.staging_kafka.bootstrap_endpoint
    }
    production = {
      id       = module.prod_kafka.id
      endpoint = module.prod_kafka.bootstrap_endpoint
    }
  }
}
```

### terraform.tfvars.example

```hcl
# Copy to terraform.tfvars and fill in your values
# DO NOT commit terraform.tfvars to version control

dev_environment_id     = "env-dev-abc123"
staging_environment_id = "env-staging-def456"
prod_environment_id    = "env-prod-ghi789"
prod_network_id        = "n-prod-private-123"

# Set credentials via environment variables:
# export TF_VAR_confluent_api_key="your-key"
# export TF_VAR_confluent_api_secret="your-secret"
```

---

## Best Practices

### 1. State Management

Always use remote state for production deployments:

```hcl
terraform {
  backend "s3" {
    bucket = "my-terraform-state"
    key    = "confluent/kafka/terraform.tfstate"
    region = "us-east-1"
    encrypt = true
    dynamodb_table = "terraform-locks"
  }
}
```

### 2. Credential Management

Use environment variables or secret management tools:

```bash
# Option 1: Environment variables
export TF_VAR_confluent_api_key="your-key"
export TF_VAR_confluent_api_secret="your-secret"

# Option 2: AWS Systems Manager Parameter Store
data "aws_ssm_parameter" "confluent_key" {
  name = "/confluent/api_key"
}
```

### 3. Workspace Management

Use Terraform workspaces for environment separation:

```bash
# Create and select workspace
terraform workspace new production
terraform workspace select production
terraform apply
```

### 4. Module Versioning

Pin module versions in production:

```hcl
module "prod_kafka" {
  source = "github.com/plantonhq/openmcf//apis/org/openmcf/provider/confluent/confluentkafka/v1/iac/tf?ref=v1.0.0"
  # ... config ...
}
```

---

## Troubleshooting

### Issue: Terraform doesn't find the module

**Solution**: Ensure you're using the correct source path:

```hcl
# For local development
source = "../../../../../../confluent/confluentkafka/v1/iac/tf"

# For production (from Git)
source = "github.com/plantonhq/openmcf//apis/org/openmcf/provider/confluent/confluentkafka/v1/iac/tf"
```

### Issue: Validation errors

**Solution**: Check that all required fields are provided and valid:

```bash
terraform validate
terraform plan
```

---

## Next Steps

- Review the [Terraform module README](./README.md) for detailed documentation
- Check the [YAML examples](../../examples.md) for API-level examples
- Explore [Confluent Cloud Terraform provider docs](https://registry.terraform.io/providers/confluentinc/confluent/latest/docs)

