# Kubernetes Helm Release - Terraform Module Examples

This document provides Terraform-specific examples for deploying Helm charts using the KubernetesHelmRelease Terraform module.

## Table of Contents

1. [Basic Usage](#example-1-basic-usage)
2. [Helm Release with Custom Values](#example-2-helm-release-with-custom-values)
3. [WordPress with Ingress](#example-3-wordpress-with-ingress)
4. [Production PostgreSQL](#example-4-production-postgresql)
5. [Using Existing Namespace](#example-5-using-existing-namespace)
6. [Multi-Environment Setup](#example-6-multi-environment-setup)
7. [Using Terraform Variables](#example-7-using-terraform-variables)

---

## Prerequisites

- Terraform >= 1.0
- Kubernetes cluster access
- kubectl configured

## Module Setup

Add the module to your Terraform configuration:

```hcl
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.0"
    }
  }
}
```

---

## Example 1: Basic Usage

Deploy NGINX using the Bitnami Helm chart with default values.

### main.tf

```hcl
module "nginx_helm_release" {
  source = "path/to/kuberneteshelmrelease/v1/iac/tf"

  metadata = {
    name = "basic-nginx"
  }

  spec = {
    namespace = {
      value = "default"
    }
    create_namespace = true
    repo    = "https://charts.bitnami.com/bitnami"
    name    = "nginx"
    version = "15.14.0"
    values  = {}
  }
}

output "nginx_namespace" {
  value = module.nginx_helm_release.namespace
}
```

### Apply

```bash
terraform init
terraform plan
terraform apply
```

---

## Example 2: Helm Release with Custom Values

Deploy NGINX with custom configuration including replicas and resource limits.

### main.tf

```hcl
module "nginx_ha" {
  source = "path/to/kuberneteshelmrelease/v1/iac/tf"

  metadata = {
    name = "nginx-ha"
    labels = {
      environment = "production"
      team        = "platform"
    }
  }

  spec = {
    namespace = {
      value = "nginx"
    }
    create_namespace = true
    repo    = "https://charts.bitnami.com/bitnami"
    name    = "nginx"
    version = "15.14.0"
    values = {
      replicaCount = "3"
      "service.type" = "LoadBalancer"
      "resources.requests.memory" = "256Mi"
      "resources.requests.cpu" = "100m"
      "resources.limits.memory" = "512Mi"
      "resources.limits.cpu" = "500m"
    }
  }
}
```

**Note:** For nested Helm values, use dot notation in the key names.

---

## Example 3: WordPress with Ingress

Deploy WordPress with ingress enabled and persistent storage.

### main.tf

```hcl
module "wordpress" {
  source = "path/to/kuberneteshelmrelease/v1/iac/tf"

  metadata = {
    name = "wordpress-blog"
    labels = {
      app         = "wordpress"
      environment = "production"
    }
  }

  spec = {
    namespace = {
      value = "wordpress"
    }
    create_namespace = true
    repo    = "https://charts.bitnami.com/bitnami"
    name    = "wordpress"
    version = "19.0.0"
    values = {
      wordpressUsername = "admin"
      wordpressEmail    = "admin@example.com"
      "ingress.enabled"  = "true"
      "ingress.hostname" = "blog.example.com"
      "ingress.pathType" = "Prefix"
      "ingress.ingressClassName" = "nginx"
      "service.type" = "ClusterIP"
      "persistence.enabled" = "true"
      "persistence.size" = "10Gi"
    }
  }
}

output "wordpress_namespace" {
  description = "Namespace where WordPress is deployed"
  value       = module.wordpress.namespace
}

output "wordpress_url" {
  description = "WordPress URL"
  value       = "https://blog.example.com"
}
```

---

## Example 4: Production PostgreSQL

Deploy PostgreSQL with replication and monitoring.

### main.tf

```hcl
module "postgresql" {
  source = "path/to/kuberneteshelmrelease/v1/iac/tf"

  metadata = {
    name = "postgres-prod"
    labels = {
      environment = "production"
      app         = "postgresql"
      criticality = "high"
    }
  }

  spec = {
    namespace = {
      value = "postgres"
    }
    create_namespace = true
    repo    = "https://charts.bitnami.com/bitnami"
    name    = "postgresql"
    version = "14.3.0"
    values = {
      "global.postgresql.auth.username" = "appuser"
      "global.postgresql.auth.password" = "changeme"  # Use Terraform secrets in production
      "global.postgresql.auth.database" = "myapp"
      "primary.persistence.enabled" = "true"
      "primary.persistence.size" = "50Gi"
      "primary.persistence.storageClass" = "fast-ssd"
      "primary.resources.requests.memory" = "1Gi"
      "primary.resources.requests.cpu" = "500m"
      "primary.resources.limits.memory" = "2Gi"
      "primary.resources.limits.cpu" = "2000m"
      "readReplicas.replicaCount" = "2"
      "readReplicas.persistence.enabled" = "true"
      "readReplicas.persistence.size" = "50Gi"
      "metrics.enabled" = "true"
      "metrics.serviceMonitor.enabled" = "true"
    }
  }
}
```

---

## Example 5: Using Existing Namespace

Deploy a Helm release to an existing namespace without creating it.

### main.tf

```hcl
module "redis_existing_ns" {
  source = "path/to/kuberneteshelmrelease/v1/iac/tf"

  metadata = {
    name = "redis-cache"
  }

  spec = {
    namespace = {
      value = "existing-redis-namespace"
    }
    create_namespace = false  # Do not create namespace
    repo    = "https://charts.bitnami.com/bitnami"
    name    = "redis"
    version = "18.19.0"
    values = {}
  }
}
```

**Use Case:** Deploy to a namespace that already exists and is managed separately, such as a namespace created by another module or manually.

**Note:** Ensure the namespace exists before applying this configuration, otherwise the Helm release will fail.

---

## Example 6: Multi-Environment Setup

Use Terraform workspaces or separate configurations for different environments.

### variables.tf

```hcl
variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "helm_chart_version" {
  description = "Helm chart version"
  type        = string
}

variable "replica_count" {
  description = "Number of replicas"
  type        = string
}

variable "resource_limits" {
  description = "Resource limits for the application"
  type = object({
    cpu    = string
    memory = string
  })
}
```

### terraform.tfvars (dev)

```hcl
environment         = "dev"
helm_chart_version  = "3.2.0"
replica_count       = "1"
resource_limits = {
  cpu    = "500m"
  memory = "512Mi"
}
```

### terraform.tfvars (prod)

```hcl
environment         = "prod"
helm_chart_version  = "3.2.0"
replica_count       = "3"
resource_limits = {
  cpu    = "4000m"
  memory = "4Gi"
}
```

### main.tf

```hcl
module "myapp" {
  source = "path/to/kuberneteshelmrelease/v1/iac/tf"

  metadata = {
    name = "myapp-${var.environment}"
    labels = {
      environment = var.environment
      team        = "engineering"
    }
  }

  spec = {
    namespace = {
      value = "myapp-${var.environment}"
    }
    create_namespace = true
    repo    = "https://charts.company.com/stable"
    name    = "myapp"
    version = var.helm_chart_version
    values = {
      environment = var.environment
      replicaCount = var.replica_count
      "resources.limits.cpu" = var.resource_limits.cpu
      "resources.limits.memory" = var.resource_limits.memory
      "autoscaling.enabled" = var.environment == "prod" ? "true" : "false"
    }
  }
}
```

### Usage

```bash
# Development
terraform workspace select dev
terraform plan -var-file=terraform.tfvars
terraform apply -var-file=terraform.tfvars

# Production
terraform workspace select prod
terraform plan -var-file=terraform.tfvars
terraform apply -var-file=terraform.tfvars
```

---

## Example 7: Using Terraform Variables

Parameterize your Helm release configuration using Terraform variables.

### variables.tf

```hcl
variable "release_name" {
  description = "Name of the Helm release"
  type        = string
}

variable "chart_repository" {
  description = "Helm chart repository URL"
  type        = string
  default     = "https://charts.bitnami.com/bitnami"
}

variable "chart_name" {
  description = "Name of the Helm chart"
  type        = string
}

variable "chart_version" {
  description = "Version of the Helm chart"
  type        = string
}

variable "custom_values" {
  description = "Custom Helm values"
  type        = map(string)
  default     = {}
}

variable "enable_monitoring" {
  description = "Enable Prometheus monitoring"
  type        = bool
  default     = false
}
```

### main.tf

```hcl
locals {
  # Merge base values with custom values
  helm_values = merge(
    {
      "metrics.enabled" = var.enable_monitoring ? "true" : "false"
    },
    var.custom_values
  )
}

module "helm_release" {
  source = "path/to/kuberneteshelmrelease/v1/iac/tf"

  metadata = {
    name = var.release_name
  }

  spec = {
    namespace = {
      value = "default"
    }
    create_namespace = true
    repo    = var.chart_repository
    name    = var.chart_name
    version = var.chart_version
    values  = local.helm_values
  }
}
```

### terraform.tfvars

```hcl
release_name     = "redis-cache"
chart_name       = "redis"
chart_version    = "18.19.0"
enable_monitoring = true
custom_values = {
  "architecture" = "replication"
  "auth.enabled" = "true"
  "replica.replicaCount" = "3"
}
```

---

## Working with Complex Helm Values

### Nested Values

For deeply nested Helm values, use dot notation:

```hcl
spec = {
  repo    = "https://charts.bitnami.com/bitnami"
  name    = "postgresql"
  version = "14.3.0"
  values = {
    "global.postgresql.auth.username" = "myuser"
    "global.postgresql.auth.password" = "mypassword"
    "primary.persistence.enabled" = "true"
  }
}
```

### Array Values

For array values in Helm charts, use index notation:

```hcl
values = {
  "ingress.hosts.0" = "app.example.com"
  "ingress.hosts.1" = "app2.example.com"
  "extraEnvVars.0.name" = "ENV_VAR"
  "extraEnvVars.0.value" = "value"
}
```

### YAML Block Values

For complex YAML blocks, consider using a separate values file:

```hcl
# Create a local file with complex values
resource "local_file" "helm_values" {
  content = yamlencode({
    prometheus = {
      serverFiles = {
        "prometheus.yml" = {
          scrape_configs = [
            {
              job_name = "kubernetes-pods"
              kubernetes_sd_configs = [{
                role = "pod"
              }]
            }
          ]
        }
      }
    }
  })
  filename = "${path.module}/values.yaml"
}

# Reference in your module (requires custom implementation)
# This is a conceptual example
```

---

## Managing Secrets

### Using Terraform Sensitive Variables

```hcl
variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

module "app" {
  source = "path/to/kuberneteshelmrelease/v1/iac/tf"

  metadata = {
    name = "myapp"
  }

  spec = {
    repo    = "https://charts.company.com"
    name    = "myapp"
    version = "1.0.0"
    values = {
      "database.password" = var.db_password
    }
  }
}
```

### Using External Secrets

```hcl
# Fetch secret from AWS Secrets Manager
data "aws_secretsmanager_secret_version" "db_credentials" {
  secret_id = "prod/database/credentials"
}

locals {
  db_creds = jsondecode(data.aws_secretsmanager_secret_version.db_credentials.secret_string)
}

module "app" {
  source = "path/to/kuberneteshelmrelease/v1/iac/tf"

  metadata = {
    name = "myapp"
  }

  spec = {
    repo    = "https://charts.company.com"
    name    = "myapp"
    version = "1.0.0"
    values = {
      "database.host"     = local.db_creds.host
      "database.username" = local.db_creds.username
      "database.password" = local.db_creds.password
    }
  }
}
```

---

## Outputs and Dependencies

### Using Module Outputs

```hcl
module "database" {
  source = "path/to/kuberneteshelmrelease/v1/iac/tf"
  
  metadata = {
    name = "postgres"
  }
  
  spec = {
    repo    = "https://charts.bitnami.com/bitnami"
    name    = "postgresql"
    version = "14.3.0"
    values  = {}
  }
}

module "app" {
  source = "path/to/kuberneteshelmrelease/v1/iac/tf"
  
  metadata = {
    name = "myapp"
  }
  
  spec = {
    repo    = "https://charts.company.com"
    name    = "myapp"
    version = "1.0.0"
    values = {
      # Reference database namespace from output
      "database.namespace" = module.database.namespace
    }
  }
  
  # Ensure database is created first
  depends_on = [module.database]
}
```

---

## Best Practices

1. **Use Variables**: Parameterize chart versions and values for reusability
2. **Version Pinning**: Always specify exact chart versions in production
3. **Sensitive Data**: Mark sensitive variables and use external secret managers
4. **Workspaces**: Use Terraform workspaces for multi-environment deployments
5. **State Management**: Use remote state backends (S3, GCS, Azure Blob)
6. **Module Versioning**: Pin module versions when using from registries

## Common Terraform Commands

```bash
# Initialize Terraform
terraform init

# Format configuration
terraform fmt

# Validate configuration
terraform validate

# Plan deployment
terraform plan

# Apply changes
terraform apply

# Show current state
terraform show

# List resources
terraform state list

# Destroy resources
terraform destroy

# View outputs
terraform output

# Refresh state
terraform refresh
```

## Troubleshooting

### Helm Values Not Applied

If Helm values aren't being applied correctly:

1. Check dot notation syntax for nested values
2. Ensure all values are strings (Terraform requirement)
3. Verify chart supports the values (check chart's values.yaml)

### State Conflicts

If you encounter state lock conflicts:

```bash
# Force unlock (use with caution)
terraform force-unlock <lock-id>
```

### Module Path Issues

Ensure module source paths are correct:

```hcl
# Relative path
source = "../../kuberneteshelmrelease/v1/iac/tf"

# Absolute path
source = "/absolute/path/to/module"

# Registry module (if published)
source = "registry.terraform.io/company/kuberneteshelmrelease/kubernetes"
version = "1.0.0"
```

## Additional Resources

- [Terraform Documentation](https://www.terraform.io/docs)
- [Terraform Helm Provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs)
- [Terraform Kubernetes Provider](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs)
- [Helm Values Reference](https://helm.sh/docs/chart_template_guide/values_files/)

