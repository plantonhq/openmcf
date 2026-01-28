# KubernetesPerconaMongoOperator Terraform Module Examples

This document provides practical examples for deploying the Percona Operator for MongoDB using Terraform.

---

## Example 1: Basic Deployment

### Description

Deploy the operator with default resource allocations.

### Configuration Files

#### `variables.tfvars`

```hcl
metadata = {
  name = "percona-operator-basic"
}

spec = {
  namespace        = "percona-operator"  # Optional: defaults to "percona-operator"
  create_namespace = true
  
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

module "percona_operator" {
  source = "../"

  metadata = var.metadata
  spec     = var.spec
}

output "namespace" {
  value = module.percona_operator.namespace
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
  name = "percona-operator-prod"
  labels = {
    environment = "production"
    managed-by  = "terraform"
    team        = "data-platform"
  }
}

spec = {
  namespace        = "percona-operator-prod"
  create_namespace = true
  
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
  name = "percona-operator-dev"
  env  = "development"
}

spec = {
  namespace        = "percona-operator-dev"
  create_namespace = true
  
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

## Example 5: Using Existing Namespace

### Description

Deploy the operator to an existing namespace managed separately (e.g., via separate Terraform resources or GitOps).

### Prerequisites

Create the namespace first:

```bash
kubectl create namespace shared-operators
```

Or define it in Terraform:

```hcl
resource "kubernetes_namespace" "shared_operators" {
  metadata {
    name = "shared-operators"
    labels = {
      managed-by = "terraform"
      purpose    = "shared-services"
    }
  }
}
```

### Configuration

#### `existing-namespace.tfvars`

```hcl
metadata = {
  name = "percona-operator-shared"
}

spec = {
  namespace        = "shared-operators"
  create_namespace = false  # Use existing namespace
  
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

### Deploy

```bash
terraform init
terraform plan -var-file="existing-namespace.tfvars"
terraform apply -var-file="existing-namespace.tfvars"
```

**Important**: Ensure the namespace `shared-operators` exists before applying.

---

## Example 6: With Remote State

### Description

Use S3 backend for state management in production.

### Configuration

#### `backend.tf`

```hcl
terraform {
  backend "s3" {
    bucket         = "my-terraform-state"
    key            = "percona-operator/terraform.tfstate"
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

## Example 7: With Outputs for Integration

### Description

Export operator information for use by other modules.

### Configuration

#### `outputs.tf`

```hcl
output "operator_namespace" {
  description = "Namespace where Percona operator is deployed"
  value       = module.percona_operator.namespace
}

output "operator_ready" {
  description = "Indicates operator deployment completion"
  value       = true
}
```

#### Using in Another Module

```hcl
module "mongodb_cluster" {
  source = "./mongodb-cluster"

  operator_namespace = module.percona_operator.namespace
  
  # Wait for operator to be ready
  depends_on = [module.percona_operator]
}
```

---

## Verification Commands

After applying any example, verify the deployment:

```bash
# Check namespace
kubectl get namespace percona-operator

# Check operator pod
kubectl get pods -n percona-operator

# Check Helm release
helm list -n percona-operator

# Check operator logs
kubectl logs -n percona-operator -l app.kubernetes.io/name=kubernetes-percona-mongo-operator -f

# Verify CRDs
kubectl get crds | grep percona
```

Expected CRDs:
- `perconaservermongodbs.psmdb.percona.com`
- `perconaservermongodbbackups.psmdb.percona.com`
- `perconaservermongodbrestores.psmdb.percona.com`

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

**Warning**: This will remove the operator but not the CRDs or any MongoDB clusters managed by it.

---

## Advanced Configuration

### Custom Helm Values

To pass additional Helm values, modify `main.tf`:

```hcl
resource "helm_release" "percona_operator" {
  # ... existing configuration ...

  # Additional custom values
  set {
    name  = "watchNamespace"
    value = "mongodb-prod,mongodb-dev"
  }

  set {
    name  = "logLevel"
    value = "debug"
  }
}
```

### Using with Terragrunt

#### `terragrunt.hcl`

```hcl
terraform {
  source = "git::https://github.com/myorg/terraform-modules//percona-operator"
}

inputs = {
  metadata = {
    name = "percona-operator-${get_env("ENVIRONMENT")}"
  }
  
  spec = {
    create_namespace = true
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

1. Deploy MongoDB clusters using PerconaServerMongoDB CRDs
2. Configure monitoring and alerting for the operator
3. Set up backup and disaster recovery procedures
4. Review operator logs regularly for any issues

