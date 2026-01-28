# Kubernetes Keycloak - Terraform Examples

This document provides Terraform-specific examples for deploying Keycloak on Kubernetes using the KubernetesKeycloak module.

## Prerequisites

- Terraform >= 1.0
- Kubernetes cluster (1.19+)
- kubectl configured
- Sufficient cluster resources

## Table of Contents

1. [Basic Keycloak Deployment](#example-1-basic-keycloak-deployment)
2. [Keycloak with Ingress](#example-2-keycloak-with-ingress)
3. [Minimal Development Setup](#example-3-minimal-development-setup)
4. [High Resource Allocation](#example-4-high-resource-allocation)
5. [Production High-Availability](#example-5-production-high-availability)

---

## Example 1: Basic Keycloak Deployment

Basic Keycloak deployment without ingress for internal testing.

```hcl
module "keycloak_basic" {
  source = "path/to/module"

  metadata = {
    name = "basic-keycloak"
  }

  spec = {
    target_cluster = {
      cluster_name = "dev-gke-cluster"
    }
    namespace = {
      value = "keycloak-basic"
    }
    create_namespace = true
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

output "keycloak_namespace" {
  value = module.keycloak_basic.namespace
}

output "port_forward_command" {
  value = module.keycloak_basic.port_forward_command
}
```

**Use Case:** Development and testing environments requiring basic identity management.

**Access Keycloak:**
```bash
# Get the port-forward command
terraform output port_forward_command

# Execute it
kubectl port-forward -n keycloak-basic-keycloak svc/keycloak-basic-keycloak 8080:8080

# Open browser
open http://localhost:8080
```

---

## Example 2: Keycloak with Ingress

Keycloak deployment with external access via ingress.

```hcl
module "keycloak_public" {
  source = "path/to/module"

  metadata = {
    name = "public-keycloak"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "prod-gke-cluster"
    }
    namespace = {
      value = "keycloak-public"
    }
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
    ingress = {
      is_enabled = true
      dns_domain = "auth.example.com"
    }
  }
}

output "keycloak_url" {
  value       = module.keycloak_public.external_hostname
  description = "Public Keycloak URL"
}

output "internal_url" {
  value       = module.keycloak_public.internal_hostname
  description = "Internal Keycloak URL"
}
```

**Use Case:** Production deployments requiring external authentication services.

**DNS Configuration:**
Ensure DNS record points to your ingress load balancer:
```bash
# Example for Google Cloud
gcloud dns record-sets create auth.example.com --type=A --ttl=300 --rrdatas=<LOAD_BALANCER_IP>
```

---

## Example 3: Minimal Development Setup

Minimal configuration for local development with reduced resources.

```hcl
module "keycloak_dev" {
  source = "path/to/module"

  metadata = {
    name = "dev-keycloak"
    env  = "development"
  }

  spec = {
    target_cluster = {
      cluster_name = "dev-gke-cluster"
    }
    namespace = {
      value = "keycloak-dev"
    }
    create_namespace = true
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
    ingress = {
      is_enabled = false
      dns_domain = ""
    }
  }
}
```

**Use Case:** Local development environments with limited resources.

**Note:** These minimal resources are not recommended for production use.

---

## Example 4: High Resource Allocation

Keycloak with increased resources for high-traffic scenarios.

```hcl
module "keycloak_high_resources" {
  source = "path/to/module"

  metadata = {
    name = "enterprise-keycloak"
    org  = "my-company"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "prod-gke-cluster"
    }
    namespace = {
      value = "keycloak-enterprise"
    }
    create_namespace = true
    container = {
      resources = {
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "4000m"
          memory = "8Gi"
        }
      }
    }
    ingress = {
      is_enabled = true
      dns_domain = "sso.company.com"
    }
  }
}

# Output all connection details
output "keycloak_endpoints" {
  value = {
    namespace       = module.keycloak_high_resources.namespace
    service         = module.keycloak_high_resources.service
    external_url    = module.keycloak_high_resources.external_hostname
    internal_url    = module.keycloak_high_resources.internal_hostname
    kube_endpoint   = module.keycloak_high_resources.kube_endpoint
  }
  description = "All Keycloak connection endpoints"
}
```

**Use Case:** Enterprise deployments handling thousands of concurrent users.

**Capacity Planning:**
- CPU: 500m - 4000m supports ~5000-10000 concurrent users
- Memory: 1Gi - 8Gi provides sufficient heap space
- Recommend horizontal scaling for higher loads

---

## Example 5: Production High-Availability

Enterprise-grade Keycloak deployment with monitoring and proper labeling.

```hcl
module "keycloak_production" {
  source = "path/to/module"

  metadata = {
    name = "prod-auth"
    id   = "keycloak-prod-001"
    org  = "acme-corp"
    env  = "production"
    labels = {
      "app"         = "keycloak"
      "tier"        = "authentication"
      "compliance"  = "soc2"
      "cost-center" = "engineering"
    }
  }

  spec = {
    target_cluster = {
      cluster_name = "prod-gke-cluster"
    }
    namespace = {
      value = "keycloak-prod-auth"
    }
    create_namespace = true
    container = {
      resources = {
        requests = {
          cpu    = "1000m"
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
      dns_domain = "auth.acmecorp.com"
    }
  }
}

# Production outputs
output "prod_keycloak_config" {
  value = {
    namespace         = module.keycloak_production.namespace
    service_name      = module.keycloak_production.service
    external_endpoint = module.keycloak_production.external_hostname
    internal_endpoint = module.keycloak_production.internal_hostname
    debug_command     = module.keycloak_production.port_forward_command
  }
  description = "Production Keycloak configuration"
}

# Resource for application to consume
resource "kubernetes_config_map" "keycloak_config" {
  metadata {
    name      = "keycloak-endpoints"
    namespace = "default"
  }

  data = {
    "KEYCLOAK_URL"      = module.keycloak_production.external_hostname
    "KEYCLOAK_INTERNAL" = module.keycloak_production.kube_endpoint
  }
}
```

**Use Case:** Mission-critical production deployments with strict SLAs.

**Production Checklist:**
- ✅ High resource allocations
- ✅ Ingress enabled with proper DNS
- ✅ Comprehensive labeling for cost tracking
- ✅ ConfigMap for service discovery
- ✅ Monitoring integration (add Prometheus/Grafana)
- ✅ Backup strategy (database snapshots)

---

## Common Patterns

### Accessing Keycloak from Applications

Connect applications running in the same cluster:

```hcl
resource "kubernetes_deployment" "app" {
  metadata {
    name      = "my-app"
    namespace = "default"
  }

  spec {
    selector {
      match_labels = {
        app = "my-app"
      }
    }

    template {
      metadata {
        labels = {
          app = "my-app"
        }
      }

      spec {
        container {
          name  = "app"
          image = "myapp:latest"

          env {
            name  = "KEYCLOAK_URL"
            value = module.keycloak_basic.kube_endpoint
          }

          env {
            name  = "KEYCLOAK_REALM"
            value = "myrealm"
          }
        }
      }
    }
  }
}
```

### Multi-Environment Setup

Deploy Keycloak across multiple environments:

```hcl
locals {
  environments = {
    dev = {
      cpu_request    = "50m"
      memory_request = "100Mi"
      cpu_limit      = "1000m"
      memory_limit   = "1Gi"
      ingress        = false
    }
    staging = {
      cpu_request    = "200m"
      memory_request = "512Mi"
      cpu_limit      = "2000m"
      memory_limit   = "2Gi"
      ingress        = true
    }
    prod = {
      cpu_request    = "1000m"
      memory_request = "2Gi"
      cpu_limit      = "4000m"
      memory_limit   = "8Gi"
      ingress        = true
    }
  }
}

module "keycloak" {
  for_each = local.environments
  source   = "path/to/module"

  metadata = {
    name = "${each.key}-keycloak"
    env  = each.key
  }

  spec = {
    target_cluster = {
      cluster_name = "${each.key}-gke-cluster"
    }
    namespace = {
      value = "keycloak-${each.key}"
    }
    create_namespace = true
    container = {
      resources = {
        requests = {
          cpu    = each.value.cpu_request
          memory = each.value.memory_request
        }
        limits = {
          cpu    = each.value.cpu_limit
          memory = each.value.memory_limit
        }
      }
    }
    ingress = {
      is_enabled = each.value.ingress
      dns_domain = each.value.ingress ? "auth-${each.key}.company.com" : ""
    }
  }
}
```

---

## Verification

After deployment, verify your Keycloak installation:

```bash
# Check namespace and pods
kubectl get pods -n keycloak-<name>

# Check service
kubectl get svc -n keycloak-<name>

# Port-forward for local access
kubectl port-forward -n keycloak-<name> svc/keycloak-<name> 8080:8080

# Test Keycloak admin console
curl http://localhost:8080/auth/admin/
```

---

## Troubleshooting

### Pods Not Starting

```bash
# Check pod status
kubectl describe pod -n keycloak-<name> <pod-name>

# Check events
kubectl get events -n keycloak-<name> --sort-by='.lastTimestamp'

# View logs
kubectl logs -n keycloak-<name> <pod-name>
```

### Database Connection Issues

If Keycloak can't connect to PostgreSQL:

```bash
# Check PostgreSQL pod
kubectl get pods -n keycloak-<name> -l app=postgresql

# View PostgreSQL logs
kubectl logs -n keycloak-<name> -l app=postgresql
```

### Ingress Not Working

```bash
# Check ingress resource
kubectl get ingress -n keycloak-<name>

# Describe ingress
kubectl describe ingress -n keycloak-<name>

# Verify DNS resolution
nslookup auth.example.com
```

---

## Best Practices

1. **Resource Planning**: Start with recommended resources and scale based on metrics
2. **Use Ingress**: Enable ingress for production deployments
3. **Database Backup**: Implement regular PostgreSQL backup strategy
4. **Monitoring**: Set up monitoring for Keycloak metrics
5. **Security**: Use HTTPS/TLS for all external endpoints
6. **Realm Configuration**: Create separate realms for different applications
7. **High Availability**: Deploy multiple Keycloak replicas for production
8. **Resource Limits**: Always set limits to prevent resource exhaustion

---

## Security Considerations

### Admin Credentials

After deployment, retrieve and secure admin credentials:

```bash
# Get admin password (if created by Helm chart)
kubectl get secret -n keycloak-<name> keycloak -o jsonpath='{.data.admin-password}' | base64 -d
```

### Network Policies

Implement network policies to restrict access:

```hcl
resource "kubernetes_network_policy" "keycloak_ingress" {
  metadata {
    name      = "keycloak-ingress-policy"
    namespace = module.keycloak_basic.namespace
  }

  spec {
    pod_selector {
      match_labels = {
        app = "keycloak"
      }
    }

    policy_types = ["Ingress"]

    ingress {
      from {
        namespace_selector {
          match_labels = {
            name = "default"
          }
        }
      }
    }
  }
}
```

---

## Additional Resources

- [Keycloak Official Documentation](https://www.keycloak.org/documentation)
- [Bitnami Keycloak Helm Chart](https://github.com/bitnami/charts/tree/main/bitnami/keycloak)
- [Terraform Kubernetes Provider](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs)
- [Module README](README.md)
- [Research Documentation](../docs/README.md)

