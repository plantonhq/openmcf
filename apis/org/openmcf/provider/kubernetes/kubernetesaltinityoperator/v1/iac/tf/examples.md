# KubernetesAltinityOperator Terraform Module Examples

This document provides practical examples for deploying the Altinity ClickHouse Operator using Terraform.

---

## Example 1: Basic Deployment

### Description

Deploy the operator with default resource allocations.

### Configuration Files

#### `variables.tfvars`

```hcl
metadata = {
  name = "kubernetes-altinity-operator-basic"
}

spec = {
  namespace        = "kubernetes-altinity-operator"  # Optional: defaults to "kubernetes-altinity-operator"
  create_namespace = true  # Create the namespace (default)
  
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
}
```

#### `main.tf`

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

module "kubernetes_altinity_operator" {
  source = "../"

  metadata = var.metadata
  spec     = var.spec
}

output "namespace" {
  value = module.kubernetes_altinity_operator.namespace
}
```

### Deploy

```bash
terraform init
terraform plan -var-file="variables.tfvars"
terraform apply -var-file="variables.tfvars"
```

---

## Example 2: Production Deployment with Custom Resources

### Description

Deploy the operator with increased resources for production workloads.

### Configuration

#### `production.tfvars`

```hcl
metadata = {
  name = "kubernetes-altinity-operator-prod"
  labels = {
    environment = "production"
    managed-by  = "terraform"
    team        = "data-platform"
  }
}

spec = {
  namespace        = "kubernetes-altinity-operator-prod"
  create_namespace = true  # Create the namespace
  
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
```

### Deploy

```bash
terraform init
terraform plan -var-file="production.tfvars"
terraform apply -var-file="production.tfvars"
```

---

## Example 3: Development Environment

### Description

Minimal resource configuration for development clusters.

### Configuration

#### `dev.tfvars`

```hcl
metadata = {
  name = "kubernetes-altinity-operator-dev"
  env  = "development"
}

spec = {
  namespace        = "kubernetes-altinity-operator-dev"
  create_namespace = true  # Create the namespace
  
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
}
```

### Deploy

```bash
terraform init
terraform plan -var-file="dev.tfvars"
terraform apply -var-file="dev.tfvars"
```

---

## Example 4: Multi-Environment Setup

### Description

Deploy the operator to multiple environments using workspaces.

### Directory Structure

```
terraform/
├── main.tf
├── variables.tf
├── dev.tfvars
├── staging.tfvars
└── prod.tfvars
```

### Commands

```bash
# Initialize
terraform init

# Development
terraform workspace new dev
terraform apply -var-file="dev.tfvars"

# Staging
terraform workspace new staging
terraform apply -var-file="staging.tfvars"

# Production
terraform workspace new prod
terraform apply -var-file="prod.tfvars"
```

---

## Example 5: With Remote State

### Description

Use S3 backend for state management in production.

### Configuration

#### `backend.tf`

```hcl
terraform {
  backend "s3" {
    bucket         = "my-terraform-state"
    key            = "kubernetes-altinity-operator/terraform.tfstate"
    region         = "us-west-2"
    encrypt        = true
    dynamodb_table = "terraform-state-lock"
  }
}
```

#### Deploy

```bash
terraform init -backend-config="backend.tf"
terraform apply -var-file="production.tfvars"
```

---

## Example 6: With Outputs for Integration

### Description

Export operator information for use by other modules.

### Configuration

#### `outputs.tf`

```hcl
output "operator_namespace" {
  description = "Namespace where Altinity operator is deployed"
  value       = module.kubernetes_altinity_operator.namespace
}

output "operator_ready" {
  description = "Indicates operator deployment completion"
  value       = true
}
```

#### Using in Another Module

```hcl
module "clickhouse_cluster" {
  source = "./clickhouse-cluster"

  operator_namespace = module.kubernetes_altinity_operator.namespace
  
  # Wait for operator to be ready
  depends_on = [module.kubernetes_altinity_operator]
}
```

---

## Verification Commands

After applying any example, verify the deployment:

```bash
# Check namespace
kubectl get namespace kubernetes-altinity-operator

# Check operator pod
kubectl get pods -n kubernetes-altinity-operator

# Check Helm release
helm list -n kubernetes-altinity-operator

# Check operator logs
kubectl logs -n kubernetes-altinity-operator -l app.kubernetes.io/name=altinity-clickhouse-operator -f

# Verify CRDs
kubectl get crds | grep clickhouse
```

Expected CRDs:
- `clickhouseinstallations.clickhouse.altinity.com`
- `clickhouseoperatorconfigurations.clickhouse.altinity.com`

---

## Updating Resources

To update operator resources:

1. Modify the tfvars file
2. Run terraform plan to preview changes
3. Apply the changes

```bash
terraform plan -var-file="production.tfvars"
terraform apply -var-file="production.tfvars"
```

---

## Destroying Resources

To remove the operator:

```bash
terraform destroy -var-file="production.tfvars"
```

**Warning**: This will remove the operator but not the CRDs or any ClickHouse clusters managed by it.

---

## Namespace Management

The operator can either create a new namespace or use an existing one, controlled by the `create_namespace` flag in the spec.

### Create New Namespace (Default)

By default (`create_namespace = true`), the operator creates a new namespace:

```hcl
spec = {
  namespace        = "kubernetes-altinity-operator"
  create_namespace = true  # Explicitly create namespace (default)
  
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
}
```

### Use Existing Namespace

If the namespace already exists (created externally or by another Terraform resource), set `create_namespace = false`:

```hcl
spec = {
  namespace        = "kubernetes-altinity-operator"
  create_namespace = false  # Do not create, use existing namespace
  
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
}
```

### When to Use `create_namespace = false`

Set `create_namespace = false` in these scenarios:

- **External Namespace Management**: Namespace created by another Terraform module or external process
- **Pre-configured Namespace**: Namespace has specific ResourceQuotas, LimitRanges, or NetworkPolicies
- **Shared Namespace**: Multiple resources share the same namespace
- **Security Policies**: Namespace creation is restricted by organizational policies
- **Multi-tenant Environments**: Centralized namespace provisioning system

### Example: Using with Existing Namespace

First, create the namespace separately:

```hcl
resource "kubernetes_namespace" "altinity_operator" {
  metadata {
    name = "kubernetes-altinity-operator"
    
    labels = {
      "pod-security.kubernetes.io/enforce" = "restricted"
      "managed-by"                          = "terraform"
    }
  }
}

resource "kubernetes_resource_quota" "altinity_operator" {
  metadata {
    name      = "altinity-operator-quota"
    namespace = kubernetes_namespace.altinity_operator.metadata[0].name
  }
  
  spec {
    hard = {
      "requests.cpu"    = "10"
      "requests.memory" = "20Gi"
      "limits.cpu"      = "20"
      "limits.memory"   = "40Gi"
    }
  }
}
```

Then deploy the operator using the existing namespace:

```hcl
module "kubernetes_altinity_operator" {
  source = "../"

  metadata = {
    name = "kubernetes-altinity-operator"
  }
  
  spec = {
    namespace        = kubernetes_namespace.altinity_operator.metadata[0].name
    create_namespace = false  # Use namespace created above
    
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
  }
  
  depends_on = [kubernetes_resource_quota.altinity_operator]
}
```

### Example: Using with Data Source

If the namespace exists outside of Terraform:

```hcl
data "kubernetes_namespace" "existing" {
  metadata {
    name = "kubernetes-altinity-operator"
  }
}

module "kubernetes_altinity_operator" {
  source = "../"

  metadata = {
    name = "kubernetes-altinity-operator"
  }
  
  spec = {
    namespace        = data.kubernetes_namespace.existing.metadata[0].name
    create_namespace = false  # Namespace already exists
    
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
  }
}
```

---

## Advanced Configuration

### Custom Helm Values

To pass additional Helm values, modify `main.tf`:

```hcl
resource "helm_release" "kubernetes_altinity_operator" {
  # ... existing configuration ...

  # Additional custom values
  set {
    name  = "operator.watchNamespaces"
    value = "clickhouse-prod,clickhouse-dev"
  }

  set {
    name  = "operator.logLevel"
    value = "debug"
  }
}
```

### Using with Terragrunt

#### `terragrunt.hcl`

```hcl
terraform {
  source = "git::https://github.com/myorg/terraform-modules//kubernetes-altinity-operator"
}

inputs = {
  metadata = {
    name = "kubernetes-altinity-operator-${get_env("ENVIRONMENT")}"
  }
  
  spec = {
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
  }
}
```

---

## Troubleshooting

### Issue: Helm Provider Authentication Failed

**Solution**: Ensure kubectl is configured correctly:
```bash
kubectl cluster-info
kubectl auth can-i create namespace
```

### Issue: Insufficient Resources

**Solution**: Increase resource limits in tfvars:
```hcl
spec = {
  container = {
    resources = {
      limits = {
        cpu    = "2000m"
        memory = "2Gi"
      }
    }
  }
}
```

### Issue: CRDs Already Exist

**Solution**: If upgrading, the CRDs might already exist. This is normal and can be ignored.

---

## Next Steps

After successfully deploying the operator:

1. Deploy ClickHouse clusters using ClickHouseInstallation CRDs
2. Configure monitoring and alerting for the operator
3. Set up backup and disaster recovery procedures
4. Review operator logs regularly for any issues

