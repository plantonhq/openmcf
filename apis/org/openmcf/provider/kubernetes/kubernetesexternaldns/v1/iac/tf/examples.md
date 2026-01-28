# Terraform Examples for KubernetesExternalDns

Complete Terraform configurations for deploying ExternalDNS on different cloud providers.

---

## Example 1: GKE with Cloud DNS

```hcl
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.23.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.11.0"
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

module "external_dns_gke" {
  source = "../"

  metadata = {
    name = "external-dns-prod"
  }

  spec = {
    namespace = {
      value = "external-dns"
    }
    create_namespace = true  # Creates namespace if it doesn't exist
    gke = {
      project_id = {
        value = "my-gcp-project"
      }
      dns_zone_id = {
        value = "my-cloud-dns-zone-id"
      }
    }
  }
}

output "namespace" {
  value = module.external_dns_gke.namespace
}

output "gke_service_account_email" {
  value = module.external_dns_gke.gke_service_account_email
}
```

---

## Example 2: EKS with Route53

```hcl
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.23.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.11.0"
    }
  }
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.cluster.token
}

provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority[0].data)
    token                  = data.aws_eks_cluster_auth.cluster.token
  }
}

data "aws_eks_cluster" "cluster" {
  name = "my-eks-cluster"
}

data "aws_eks_cluster_auth" "cluster" {
  name = "my-eks-cluster"
}

module "external_dns_eks" {
  source = "../"

  metadata = {
    name = "external-dns-eks-prod"
  }

  spec = {
    namespace = {
      value = "external-dns"
    }
    create_namespace = true  # Creates namespace if it doesn't exist
    eks = {
      route53_zone_id = {
        value = "Z1234567890ABC"
      }
      irsa_role_arn_override = "arn:aws:iam::123456789012:role/external-dns-role"
    }
  }
}

output "service_account_name" {
  value = module.external_dns_eks.service_account_name
}
```

---

## Example 3: AKS with Azure DNS

```hcl
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.23.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.11.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "kubernetes" {
  host                   = data.azurerm_kubernetes_cluster.cluster.kube_config.0.host
  client_certificate     = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.client_certificate)
  client_key             = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.client_key)
  cluster_ca_certificate = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.cluster_ca_certificate)
}

provider "helm" {
  kubernetes {
    host                   = data.azurerm_kubernetes_cluster.cluster.kube_config.0.host
    client_certificate     = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.client_certificate)
    client_key             = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.client_key)
    cluster_ca_certificate = base64decode(data.azurerm_kubernetes_cluster.cluster.kube_config.0.cluster_ca_certificate)
  }
}

data "azurerm_kubernetes_cluster" "cluster" {
  name                = "my-aks-cluster"
  resource_group_name = "my-resource-group"
}

module "external_dns_aks" {
  source = "../"

  metadata = {
    name = "external-dns-aks-prod"
  }

  spec = {
    namespace = {
      value = "external-dns"
    }
    create_namespace = true  # Creates namespace if it doesn't exist
    aks = {
      dns_zone_id = {
        value = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/my-rg/providers/Microsoft.Network/dnszones/example.com"
      }
      managed_identity_client_id = "12345678-1234-1234-1234-123456789012"
    }
  }
}
```

---

## Example 4: Cloudflare DNS

```hcl
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.23.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.11.0"
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

variable "cloudflare_api_token" {
  description = "Cloudflare API token"
  type        = string
  sensitive   = true
}

module "external_dns_cloudflare" {
  source = "../"

  metadata = {
    name = "external-dns-cloudflare"
  }

  spec = {
    namespace = {
      value = "external-dns"
    }
    create_namespace = true  # Creates namespace if it doesn't exist
    cloudflare = {
      api_token = var.cloudflare_api_token
      dns_zone_id = {
        value = "1234567890abcdef1234567890abcdef"
      }
      is_proxied  = true
    }
  }
}

output "cloudflare_secret_name" {
  value = module.external_dns_cloudflare.cloudflare_secret_name
}
```

To use:

```bash
terraform apply -var="cloudflare_api_token=your-token-here"
```

---

## Example 5: Custom Versions

```hcl
module "external_dns_custom" {
  source = "../"

  metadata = {
    name = "external-dns-custom"
    env  = "staging"
    org  = "platform-team"
  }

  spec = {
    namespace = {
      value = "dns-automation"
    }
    kubernetes_external_dns_version = "v0.14.0"
    helm_chart_version              = "1.14.0"

    gke = {
      project_id = {
        value = "my-gcp-project"
      }
      dns_zone_id = {
        value = "staging-dns-zone-id"
      }
    }
  }
}
```

---

## Example 6: Multiple Instances (Different Zones)

Deploy multiple ExternalDNS instances for different DNS zones:

```hcl
# Production domain
module "external_dns_prod" {
  source = "../"

  metadata = {
    name = "external-dns-prod-domain"
  }

  spec = {
    gke = {
      project_id = {
        value = "my-gcp-project"
      }
      dns_zone_id = {
        value = "prod-example-com-zone"
      }
    }
  }
}

# Staging domain
module "external_dns_staging" {
  source = "../"

  metadata = {
    name = "external-dns-staging-domain"
  }

  spec = {
    gke = {
      project_id = {
        value = "my-gcp-project"
      }
      dns_zone_id = {
        value = "staging-example-com-zone"
      }
    }
  }
}

output "prod_namespace" {
  value = module.external_dns_prod.namespace
}

output "staging_namespace" {
  value = module.external_dns_staging.namespace
}
```

---

## Example 7: Using Terraform Variables

Create a `terraform.tfvars` file:

```hcl
project_id  = "my-gcp-project"
dns_zone_id = "my-dns-zone-id"
environment = "production"
```

Main configuration:

```hcl
variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "dns_zone_id" {
  description = "Cloud DNS zone ID"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

module "external_dns" {
  source = "../"

  metadata = {
    name = "external-dns-${var.environment}"
    env  = var.environment
  }

  spec = {
    gke = {
      project_id = {
        value = var.project_id
      }
      dns_zone_id = {
        value = var.dns_zone_id
      }
    }
  }
}
```

---

## Running the Examples

### Initialize Terraform

```bash
terraform init
```

### Plan the deployment

```bash
terraform plan
```

### Apply the configuration

```bash
terraform apply
```

### Destroy resources

```bash
terraform destroy
```

---

## Testing ExternalDNS

After deployment, test by creating an Ingress with DNS annotation:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-app
  namespace: default
  annotations:
    external-dns.alpha.kubernetes.io/hostname: test.example.com
spec:
  rules:
  - host: test.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: test-app
            port:
              number: 80
```

Check if DNS record was created:

```bash
# For Route53
aws route53 list-resource-record-sets --hosted-zone-id $(terraform output -raw route53_zone_id)

# For Cloud DNS
gcloud dns record-sets list --zone=$(terraform output -raw dns_zone_id)

# For Cloudflare
curl -X GET "https://api.cloudflare.com/client/v4/zones/$(terraform output -raw dns_zone_id)/dns_records" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Troubleshooting

### Module Validation Errors

Ensure you provide only **one** provider configuration (gke, eks, aks, or cloudflare).

### Helm Release Fails

Check Helm status:

```bash
helm list -n $(terraform output -raw namespace)
```

### DNS Records Not Created

Check ExternalDNS logs:

```bash
kubectl logs -n $(terraform output -raw namespace) -l app.kubernetes.io/name=external-dns
```

---

## Best Practices

1. **Use remote state**: Store Terraform state remotely (S3, GCS, Azure Storage)
2. **Version control**: Commit your Terraform files to Git
3. **Use variables**: Avoid hardcoding values
4. **Sensitive data**: Use Terraform Cloud/Enterprise for secrets or environment variables
5. **Module versioning**: Pin module versions for reproducibility

---

For more information, see the [README](README.md).

