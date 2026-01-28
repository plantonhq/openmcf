# Terraform Examples for Argo CD on Kubernetes

This document provides Terraform configuration examples for deploying Argo CD on Kubernetes using the KubernetesArgocd module.

## Example 1: Basic Argo CD Deployment

```hcl
module "argocd_basic" {
  source = "path/to/kubernetesargocd/v1/iac/tf"

  metadata = {
    name = "dev-argocd"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "argocd-dev"
    container = {
      resources = {
        requests = {
          cpu    = "50m"
          memory = "100Mi"
        }
        limits = {
          cpu    = "1000m"
          memory = "1Gi"
        }
      }
    }
    ingress = {
      is_enabled = false
      dns_domain = ""
    }
  }
}

output "namespace" {
  value = module.argocd_basic.namespace
}

output "port_forward_command" {
  value = module.argocd_basic.port_forward_command
}
```

**Use Case:** Development environment with minimal resources and no external access.

---

## Example 2: Production Argo CD with Ingress

```hcl
module "argocd_production" {
  source = "path/to/kubernetesargocd/v1/iac/tf"

  metadata = {
    name = "prod-argocd"
    id   = "prod-argocd-01"
    org  = "acme-corp"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "prod-gke-cluster"
    }
    namespace = "argocd-prod"
    container = {
      resources = {
        requests = {
          cpu    = "200m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
    }
    ingress = {
      is_enabled = true
      dns_domain = "acme-corp.com"
    }
  }
}

output "namespace" {
  value = module.argocd_production.namespace
}

output "external_hostname" {
  value = module.argocd_production.external_hostname
}

output "internal_hostname" {
  value = module.argocd_production.internal_hostname
}

output "kube_endpoint" {
  value = module.argocd_production.kube_endpoint
}
```

**Use Case:** Production environment with higher resources and ingress for external access.

**Endpoints:**
- External: `https://argo-prod-argocd-01.acme-corp.com`
- Internal: `https://argo-prod-argocd-01-internal.acme-corp.com`

---

## Example 3: High-Performance Argo CD

```hcl
module "argocd_large_scale" {
  source = "path/to/kubernetesargocd/v1/iac/tf"

  metadata = {
    name = "platform-argocd"
    id   = "platform-argocd"
    labels = {
      "team"        = "platform"
      "cost-center" = "engineering"
    }
  }

  spec = {
    target_cluster = {
      cluster_name = "platform-gke-cluster"
    }
    namespace = "argocd-platform"
    container = {
      resources = {
        requests = {
          cpu    = "500m"
          memory = "2Gi"
        }
        limits = {
          cpu    = "4000m"
          memory = "8Gi"
        }
      }
    }
    ingress = {
      is_enabled = true
      dns_domain = "platform.example.com"
    }
  }
}
```

**Use Case:** Large-scale GitOps platform managing hundreds of applications across multiple clusters.

---

## Example 4: Multi-Environment Setup

```hcl
# Development Environment
module "argocd_dev" {
  source = "path/to/kubernetesargocd/v1/iac/tf"

  metadata = {
    name = "argocd"
    env  = "dev"
  }

  spec = {
    target_cluster = {
      cluster_name = "dev-gke-cluster"
    }
    namespace = "argocd-dev"
    container = {
      resources = {
        requests = {
          cpu    = "50m"
          memory = "128Mi"
        }
        limits = {
          cpu    = "500m"
          memory = "512Mi"
        }
      }
    }
    ingress = {
      is_enabled = true
      dns_domain = "dev.example.com"
    }
  }
}

# Staging Environment
module "argocd_staging" {
  source = "path/to/kubernetesargocd/v1/iac/tf"

  metadata = {
    name = "argocd"
    env  = "staging"
  }

  spec = {
    target_cluster = {
      cluster_name = "staging-gke-cluster"
    }
    namespace = "argocd-staging"
    container = {
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "1000m"
          memory = "1Gi"
        }
      }
    }
    ingress = {
      is_enabled = true
      dns_domain = "staging.example.com"
    }
  }
}

# Production Environment
module "argocd_prod" {
  source = "path/to/kubernetesargocd/v1/iac/tf"

  metadata = {
    name = "argocd"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "prod-gke-cluster"
    }
    namespace = "argocd-prod"
    container = {
      resources = {
        requests = {
          cpu    = "200m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
    }
    ingress = {
      is_enabled = true
      dns_domain = "example.com"
    }
  }
}

# Outputs for each environment
output "dev_hostname" {
  value = module.argocd_dev.external_hostname
}

output "staging_hostname" {
  value = module.argocd_staging.external_hostname
}

output "prod_hostname" {
  value = module.argocd_prod.external_hostname
}
```

**Use Case:** Consistent Argo CD deployments across development, staging, and production environments with appropriate resource sizing for each.

---

## Backend Configuration

For production use, configure a remote backend for Terraform state:

```hcl
terraform {
  backend "s3" {
    bucket         = "my-terraform-state"
    key            = "argocd/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-state-lock"
    encrypt        = true
  }
}
```

---

## Provider Requirements

Ensure you have the required providers configured:

```hcl
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}
```

---

## Common Variables Pattern

You can create a reusable variables file:

```hcl
# variables.tf
variable "environment" {
  description = "Environment name"
  type        = string
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
}

variable "dns_domain" {
  description = "DNS domain for ingress"
  type        = string
}

variable "enable_ingress" {
  description = "Enable ingress for external access"
  type        = bool
  default     = false
}

# main.tf
module "argocd" {
  source = "path/to/kubernetesargocd/v1/iac/tf"

  metadata = {
    name = "argocd-${var.environment}"
    env  = var.environment
  }

  spec = {
    target_cluster = {
      cluster_name = var.cluster_name
    }
    namespace = "argocd-${var.environment}"
    container = {
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "1000m"
          memory = "2Gi"
        }
      }
    }
    ingress = {
      is_enabled = var.enable_ingress
      dns_domain = var.dns_domain
    }
  }
}
```

---

## Tips and Best Practices

1. **Resource Sizing**: Start with conservative resource requests and adjust based on actual usage
2. **Ingress**: Enable ingress for production deployments with proper TLS configuration
3. **Labels**: Use comprehensive labeling for cost tracking and resource organization
4. **State Management**: Always use remote state backends for production deployments
5. **Security**: Configure OIDC/SAML for authentication instead of relying on the admin user
6. **Monitoring**: Integrate with Prometheus and Grafana for observability
7. **High Availability**: For production, consider enabling redis-ha and scaling controller/repo-server replicas

