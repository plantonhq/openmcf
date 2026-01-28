# Kubernetes GitLab Terraform Examples

## Namespace Management

All examples demonstrate namespace management using the `create_namespace` variable:
- When `create_namespace = true`, the Terraform module creates the namespace automatically
- When `create_namespace = false`, the namespace must exist before applying Terraform

---

# Example 1: Basic GitLab Deployment with Namespace Creation

```hcl
module "gitlab_instance" {
  source = "./terraform/modules/kubernetes-gitlab"

  metadata = {
    name = "gitlab-instance"
    id   = "gitlab-inst-001"
    org  = "my-org"
    env  = "dev"
  }

  spec = {
    namespace        = "gitlab-instance"
    create_namespace = true

    container = {
      resources = {
        requests = {
          cpu    = "50m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "1"
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
  value = module.gitlab_instance.namespace
}

output "service_name" {
  value = module.gitlab_instance.service_name
}
```

---

# Example 2: GitLab with Ingress and Custom Hostname

```hcl
module "gitlab_production" {
  source = "./terraform/modules/kubernetes-gitlab"

  metadata = {
    name = "gitlab-production"
    id   = "gitlab-prod-001"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    namespace        = "gitlab-production"
    create_namespace = true

    container = {
      resources = {
        requests = {
          cpu    = "200m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "2"
          memory = "2Gi"
        }
      }
    }

    ingress = {
      is_enabled = true
      hostname   = "gitlab.example.com"
    }
  }
}

output "external_url" {
  value       = module.gitlab_production.external_url
  description = "External URL for accessing GitLab"
}
```

---

# Example 3: GitLab Using Existing Namespace

This example demonstrates deploying GitLab into a pre-existing namespace.
The namespace must be created before applying this Terraform configuration.

```hcl
# First, ensure the namespace exists
resource "kubernetes_namespace" "shared_services" {
  metadata {
    name = "shared-services"
    labels = {
      environment = "shared"
      managed_by  = "terraform"
    }
  }
}

module "gitlab_shared" {
  source = "./terraform/modules/kubernetes-gitlab"

  metadata = {
    name = "gitlab-shared"
    id   = "gitlab-shared-001"
    org  = "my-org"
    env  = "shared"
  }

  spec = {
    namespace        = "shared-services"
    create_namespace = false  # Use existing namespace

    container = {
      resources = {
        requests = {
          cpu    = "250m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
    }

    ingress = {
      is_enabled = false
      dns_domain = ""
    }
  }

  # Ensure namespace exists before creating GitLab resources
  depends_on = [kubernetes_namespace.shared_services]
}
```

---

# Example 4: GitLab with Custom Resources

```hcl
module "gitlab_custom" {
  source = "./terraform/modules/kubernetes-gitlab"

  metadata = {
    name = "gitlab-custom"
    id   = "gitlab-custom-001"
    org  = "my-org"
    env  = "staging"
  }

  spec = {
    namespace        = "gitlab-custom"
    create_namespace = true

    container = {
      resources = {
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "3"
          memory = "4Gi"
        }
      }
    }

    ingress = {
      is_enabled = false
      dns_domain = ""
    }
  }
}
```

---

# Example 5: Minimal GitLab Deployment

```hcl
module "minimal_gitlab" {
  source = "./terraform/modules/kubernetes-gitlab"

  metadata = {
    name = "minimal-gitlab"
  }

  spec = {
    namespace        = "minimal-gitlab"
    create_namespace = true

    container = {
      resources = {
        requests = {
          cpu    = "50m"
          memory = "100Mi"
        }
        limits = {
          cpu    = "1"
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
```

---

## Usage

1. Initialize Terraform:
   ```bash
   terraform init
   ```

2. Review the plan:
   ```bash
   terraform plan
   ```

3. Apply the configuration:
   ```bash
   terraform apply
   ```

4. Access GitLab:
   - If ingress is enabled, use the external URL from outputs
   - If ingress is disabled, use the port-forward command from outputs

## Important Notes

- When `create_namespace = false`, ensure the namespace exists before applying
- The module automatically handles namespace dependencies
- Resource labels are applied for tracking and organization
- Ingress configuration requires appropriate ingress controller in the cluster
