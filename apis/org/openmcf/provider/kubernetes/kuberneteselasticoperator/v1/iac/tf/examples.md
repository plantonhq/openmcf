# Kubernetes Elastic Operator - Terraform Examples

This document provides practical Terraform examples for deploying the ECK operator.

## Example 1: Basic Deployment

```hcl
module "eck_operator" {
  source = "./path/to/kubernetes-elastic-operator/iac/tf"

  metadata = {
    name = "eck-operator"
    id   = "eck-op-prod"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "elastic-system"
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
  }
}

output "eck_namespace" {
  value = module.eck_operator.namespace
}
```

## Example 2: High-Availability Production

```hcl
module "eck_operator_ha" {
  source = "./path/to/kubernetes-elastic-operator/iac/tf"

  metadata = {
    name = "eck-operator"
    id   = "eck-op-prod"
    org  = "platform"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "production-gke-cluster"
    }
    namespace = "elastic-system"
    container = {
      resources = {
        requests = {
          cpu    = "200m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "2000m"
          memory = "2Gi"
        }
      }
    }
  }
}
```

## Example 3: Development Environment

```hcl
module "eck_operator_dev" {
  source = "./path/to/kubernetes-elastic-operator/iac/tf"

  metadata = {
    name = "eck-operator-dev"
    id   = "eck-op-dev"
    env  = "development"
  }

  spec = {
    target_cluster = {
      cluster_name = "dev-gke-cluster"
    }
    namespace = "elastic-system"
    container = {
      resources = {
        requests = {
          cpu    = "25m"
          memory = "64Mi"
        }
        limits = {
          cpu    = "500m"
          memory = "512Mi"
        }
      }
    }
  }
}
```

## Example 4: Multi-Environment with Terraform Workspaces

**main.tf**:
```hcl
locals {
  env_resources = {
    dev = {
      cpu_request = "25m"
      cpu_limit   = "500m"
      mem_request = "64Mi"
      mem_limit   = "512Mi"
    }
    staging = {
      cpu_request = "100m"
      cpu_limit   = "1000m"
      mem_request = "256Mi"
      mem_limit   = "1Gi"
    }
    production = {
      cpu_request = "200m"
      cpu_limit   = "2000m"
      mem_request = "512Mi"
      mem_limit   = "2Gi"
    }
  }
  current_env = terraform.workspace
  resources   = local.env_resources[local.current_env]
}

module "eck_operator" {
  source = "./path/to/kubernetes-elastic-operator/iac/tf"

  metadata = {
    name = "eck-operator"
    id   = "eck-op-${local.current_env}"
    env  = local.current_env
  }

  spec = {
    target_cluster = {
      cluster_name = "${local.current_env}-gke-cluster"
    }
    namespace = "elastic-system"
    container = {
      resources = {
        requests = {
          cpu    = local.resources.cpu_request
          memory = local.resources.mem_request
        }
        limits = {
          cpu    = local.resources.cpu_limit
          memory = local.resources.mem_limit
        }
      }
    }
  }
}
```

**Usage**:
```bash
# Deploy to dev
terraform workspace select dev
terraform apply

# Deploy to production
terraform workspace select production
terraform apply
```

## Verification

```bash
# Check Terraform outputs
terraform output

# Verify Kubernetes resources
kubectl get pods -n elastic-system
kubectl get crds | grep elastic
```

## Cleanup

```bash
terraform destroy
```

