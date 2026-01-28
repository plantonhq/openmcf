# Kubernetes Ingress NGINX - Terraform Module Examples

This document provides Terraform-specific examples for deploying the NGINX Ingress Controller.

## Table of Contents

1. [Basic Usage](#example-1-basic-usage)
2. [GKE with Static IP](#example-2-gke-with-static-ip)
3. [GKE Internal Load Balancer](#example-3-gke-internal-load-balancer)
4. [EKS with Network Load Balancer](#example-4-eks-with-network-load-balancer)
5. [AKS with Static Public IP](#example-5-aks-with-static-public-ip)
6. [Multi-Environment Setup](#example-6-multi-environment-setup)

---

## Example 1: Basic Usage

Deploy NGINX Ingress Controller with external load balancer.

### main.tf

```hcl
module "ingress_nginx" {
  source = "path/to/kubernetesingressnginx/v1/iac/tf"

  metadata = {
    name = "basic-ingress"
  }

  spec = {
    namespace        = "ingress-nginx"
    create_namespace = true
    chart_version    = "4.11.1"
    internal         = false
  }
}

output "namespace" {
  value = module.ingress_nginx.namespace
}

output "service_name" {
  value = module.ingress_nginx.service_name
}
```

### Usage

```bash
terraform init
terraform plan
terraform apply
```

---

## Example 2: GKE with Static IP

Deploy on GKE with a reserved static IP address.

### main.tf

```hcl
module "gke_ingress" {
  source = "path/to/kubernetesingressnginx/v1/iac/tf"

  metadata = {
    name = "gke-ingress"
    env  = "production"
    labels = {
      cloud = "gcp"
      team  = "platform"
    }
  }

  spec = {
    namespace        = "ingress-nginx"
    create_namespace = true
    chart_version    = "4.11.1"
    internal         = false
    gke = {
      static_ip_name = "prod-ingress-static-ip"
    }
  }
}

output "ingress_namespace" {
  description = "Namespace where ingress controller is deployed"
  value       = module.gke_ingress.namespace
}

output "load_balancer_service" {
  description = "Service name for the load balancer"
  value       = module.gke_ingress.service_name
}
```

**Prerequisites:**

```bash
# Create static IP
gcloud compute addresses create prod-ingress-static-ip --global

# Verify creation
gcloud compute addresses describe prod-ingress-static-ip --global
```

---

## Example 3: GKE Internal Load Balancer

Deploy an internal load balancer on GKE with specific subnetwork.

### main.tf

```hcl
module "gke_internal_ingress" {
  source = "path/to/kubernetesingressnginx/v1/iac/tf"

  metadata = {
    name = "internal-ingress"
    env  = "production"
  }

  spec = {
    namespace        = "ingress-nginx"
    create_namespace = true
    chart_version    = "4.11.1"
    internal         = true
    gke = {
      subnetwork_self_link = "projects/my-project/regions/us-west1/subnetworks/private-subnet"
    }
  }
}
```

---

## Example 4: EKS with Network Load Balancer

Deploy on EKS with specific subnets and security groups.

### main.tf

```hcl
module "eks_ingress" {
  source = "path/to/kubernetesingressnginx/v1/iac/tf"

  metadata = {
    name = "eks-ingress"
    env  = "production"
    labels = {
      cloud = "aws"
      region = "us-west-2"
    }
  }

  spec = {
    namespace        = "ingress-nginx"
    create_namespace = true
    chart_version    = "4.11.1"
    internal         = false
    eks = {
      subnet_ids = [
        "subnet-public-1a",
        "subnet-public-1b",
        "subnet-public-1c"
      ]
      additional_security_group_ids = [
        "sg-web-access",
        "sg-api-access"
      ]
      irsa_role_arn_override = "arn:aws:iam::123456789012:role/ingress-nginx-role"
    }
  }
}

output "controller_service" {
  description = "Ingress controller service name"
  value       = module.eks_ingress.service_name
}
```

---

## Example 5: AKS with Static Public IP

Deploy on Azure AKS with a pre-created public IP.

### main.tf

```hcl
module "aks_ingress" {
  source = "path/to/kubernetesingressnginx/v1/iac/tf"

  metadata = {
    name = "aks-ingress"
    env  = "production"
    labels = {
      cloud = "azure"
    }
  }

  spec = {
    namespace        = "ingress-nginx"
    create_namespace = true
    chart_version    = "4.11.1"
    internal         = false
    aks = {
      managed_identity_client_id = "12345678-1234-1234-1234-123456789012"
      public_ip_name             = "prod-ingress-public-ip"
    }
  }
}
```

**Prerequisites:**

```bash
# Create public IP
az network public-ip create \
  --resource-group myResourceGroup \
  --name prod-ingress-public-ip \
  --sku Standard \
  --allocation-method Static

# Create managed identity
az identity create \
  --resource-group myResourceGroup \
  --name ingress-nginx-identity
```

---

## Example 6: Multi-Environment Setup

Use Terraform workspaces for multi-environment deployments.

### variables.tf

```hcl
variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "cluster_credential" {
  description = "Kubernetes cluster credential ID"
  type        = string
}

variable "is_internal" {
  description = "Deploy with internal load balancer"
  type        = bool
  default     = false
}
```

### terraform.tfvars (dev)

```hcl
environment        = "dev"
cluster_credential = "dev-cluster-credential"
is_internal        = false
```

### terraform.tfvars (prod)

```hcl
environment        = "prod"
cluster_credential = "prod-cluster-credential"
is_internal        = false
```

### main.tf

```hcl
module "ingress_nginx" {
  source = "path/to/kubernetesingressnginx/v1/iac/tf"

  metadata = {
    name = "ingress-${var.environment}"
    env  = var.environment
  }

  spec = {
    namespace        = "ingress-nginx"
    create_namespace = true
    chart_version    = "4.11.1"
    internal         = var.is_internal
  }
}

output "namespace" {
  value = module.ingress_nginx.namespace
}
```

### Usage

```bash
# Development
terraform workspace select dev
terraform apply -var-file=terraform.tfvars

# Production
terraform workspace select prod
terraform apply -var-file=terraform.tfvars
```

---

## Verifying Deployment

After applying:

```bash
# Check Helm release
helm list -n kubernetes-ingress-nginx

# Check controller pods
kubectl get pods -n kubernetes-ingress-nginx

# Get load balancer IP/hostname
kubectl get svc -n kubernetes-ingress-nginx kubernetes-ingress-nginx-controller

# View Terraform outputs
terraform output namespace
terraform output service_name
```

## Testing Ingress

Create a test ingress resource:

```hcl
resource "kubernetes_ingress_v1" "test" {
  metadata {
    name      = "test-ingress"
    namespace = "default"
  }

  spec {
    ingress_class_name = "nginx"
    
    rule {
      host = "test.example.com"
      http {
        path {
          path      = "/"
          path_type = "Prefix"
          backend {
            service {
              name = "test-service"
              port {
                number = 80
              }
            }
          }
        }
      }
    }
  }

  depends_on = [module.ingress_nginx]
}
```

## Common Terraform Commands

```bash
# Initialize
terraform init

# Format code
terraform fmt

# Validate configuration
terraform validate

# Plan changes
terraform plan

# Apply changes
terraform apply

# Show current state
terraform show

# List resources
terraform state list

# View outputs
terraform output

# Destroy resources
terraform destroy
```

## Troubleshooting

### Load Balancer Not Created

Check Terraform state and Kubernetes events:

```bash
terraform show
kubectl get events -n kubernetes-ingress-nginx --sort-by='.lastTimestamp'
```

### Helm Release Failed

View Helm release status:

```bash
helm list -n kubernetes-ingress-nginx
helm status kubernetes-ingress-nginx -n kubernetes-ingress-nginx
```

### Module Source Issues

Ensure module path is correct:

```hcl
# Relative path
source = "../../kubernetesingressnginx/v1/iac/tf"

# Absolute path
source = "/absolute/path/to/module"
```

## Additional Resources

- [Terraform Documentation](https://www.terraform.io/docs)
- [Terraform Helm Provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
- [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)

