# Prometheus Kubernetes Terraform Module - Examples

This document provides examples demonstrating various Prometheus configurations for Kubernetes deployments.

## Prerequisites

- Terraform >= 1.0
- Kubernetes cluster configured
- **kube-prometheus-stack Helm chart** (for actual Prometheus deployment)
- Prometheus Operator installed (typically deployed via kube-prometheus-stack)

---

## Example 1: Minimal Development Prometheus

```hcl
module "prometheus_dev" {
  source = "../../tf"

  metadata = {
    name = "dev-prometheus"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace = "my-namespace"
    create_namespace = true
    
    container = {
      replicas = 1
      persistence_enabled = false
      disk_size = ""
      resources = {
        requests = {
          cpu    = "50m"
          memory = "256Mi"
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
```

**Key points**:
- Single replica for development
- No persistence (in-memory only)
- Minimal resources
- No external access (ingress disabled)
- Suitable for local development and testing

---

## Example 2: Production Prometheus with Persistence

```hcl
module "prometheus_production" {
  source = "../../tf"

  metadata = {
    name = "prod-prometheus"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace = "my-namespace"
    create_namespace = true
    
    container = {
      replicas = 2
      persistence_enabled = true
      disk_size = "50Gi"
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
      dns_domain = "example.com"
    }
  }
}
```

**Key points**:
- High-availability with 2 replicas
- Persistent storage enabled with 50Gi disk
- Production-grade resources
- External access via ingress
- Suitable for production monitoring

---

## Example 3: Large-Scale Analytics Prometheus

```hcl
module "prometheus_analytics" {
  source = "../../tf"

  metadata = {
    name = "analytics-prometheus"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace = "my-namespace"
    create_namespace = true
    
    container = {
      replicas = 3
      persistence_enabled = true
      disk_size = "200Gi"
      resources = {
        requests = {
          cpu    = "2000m"
          memory = "8Gi"
        }
        limits = {
          cpu    = "8000m"
          memory = "32Gi"
        }
      }
    }
    
    ingress = {
      is_enabled = true
      dns_domain = "metrics.example.com"
    }
  }
}

# Output the connection details
output "prometheus_namespace" {
  value = module.prometheus_analytics.namespace
}

output "prometheus_endpoint" {
  value = module.prometheus_analytics.kube_endpoint
}
```

**Key points**:
- High resource allocation for large-scale monitoring
- Large persistent storage (200Gi)
- Multiple replicas for high availability
- Optimized for high cardinality metrics and long retention

---

## Example 4: Minimal Resource Prometheus

```hcl
module "prometheus_minimal" {
  source = "../../tf"

  metadata = {
    name = "minimal-prometheus"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace = "my-namespace"
    create_namespace = true
    
    container = {
      replicas = 1
      persistence_enabled = false
      disk_size = ""
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
      is_enabled = false
      dns_domain = ""
    }
  }
}
```

**Key points**:
- Minimal resources for small workloads
- No persistence
- Cost-optimized configuration
- Suitable for testing or low-traffic monitoring

---

## Example 5: Prometheus with Ingress and Custom Domain

```hcl
module "prometheus_public" {
  source = "../../tf"

  metadata = {
    name = "public-prometheus"
    org  = "my-org"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace = "my-namespace"
    create_namespace = true
    
    container = {
      replicas = 1
      persistence_enabled = true
      disk_size = "30Gi"
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
      dns_domain = "monitoring.example.com"
    }
  }
}

# Output the external URL
output "prometheus_url" {
  value = "https://${module.prometheus_public.external_hostname}"
}
```

**Key points**:
- External access via ingress
- Custom DNS domain for monitoring
- Persistent storage enabled
- Outputs for connecting applications

---

## Example 6: Multi-Environment Configuration

```hcl
variable "environment" {
  type    = string
  default = "dev"
}

variable "resource_tier" {
  type = map(object({
    cpu_request    = string
    cpu_limit      = string
    memory_request = string
    memory_limit   = string
    disk_size      = string
    replicas       = number
    persistence    = bool
  }))
  default = {
    dev = {
      cpu_request    = "50m"
      cpu_limit      = "500m"
      memory_request = "128Mi"
      memory_limit   = "512Mi"
      disk_size      = ""
      replicas       = 1
      persistence    = false
    }
    staging = {
      cpu_request    = "500m"
      cpu_limit      = "2000m"
      memory_request = "1Gi"
      memory_limit   = "4Gi"
      disk_size      = "20Gi"
      replicas       = 1
      persistence    = true
    }
    production = {
      cpu_request    = "2000m"
      cpu_limit      = "8000m"
      memory_request = "8Gi"
      memory_limit   = "32Gi"
      disk_size      = "100Gi"
      replicas       = 2
      persistence    = true
    }
  }
}

module "prometheus_env" {
  source = "../../tf"

  metadata = {
    name = "app-prometheus"
    env  = var.environment
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace = "my-namespace"
    create_namespace = true
    
    container = {
      replicas = var.resource_tier[var.environment].replicas
      persistence_enabled = var.resource_tier[var.environment].persistence
      disk_size = var.resource_tier[var.environment].disk_size
      resources = {
        requests = {
          cpu    = var.resource_tier[var.environment].cpu_request
          memory = var.resource_tier[var.environment].memory_request
        }
        limits = {
          cpu    = var.resource_tier[var.environment].cpu_limit
          memory = var.resource_tier[var.environment].memory_limit
        }
      }
    }
    
    ingress = {
      is_enabled = var.environment != "dev"
      dns_domain = var.environment != "dev" ? "prometheus-${var.environment}.example.com" : ""
    }
  }
}
```

**Key points**:
- Environment-based resource allocation
- Single configuration for all environments
- Automatic resource scaling based on environment
- Conditional ingress based on environment

---

## Deployment Architecture

All deployments create the namespace foundation for Prometheus. The actual Prometheus deployment uses the **kube-prometheus-stack Helm chart** which provides:

- **Prometheus Operator**: Manages Prometheus instances using Kubernetes CRDs
- **Prometheus Server**: Time-series database and query engine
- **Grafana**: Visualization and dashboards
- **Alertmanager**: Alert routing and notification
- **ServiceMonitors**: Auto-discovery of scrape targets
- **PrometheusRules**: Recording and alerting rules
- **Exporters**: node-exporter, kube-state-metrics for comprehensive monitoring

---

## Resource Planning Guidelines

### Development/Testing
- **CPU**: 50m-500m
- **Memory**: 128Mi-512Mi
- **Disk**: N/A (no persistence)
- **Replicas**: 1

### Staging
- **CPU**: 500m-2000m
- **Memory**: 1Gi-4Gi
- **Disk**: 20Gi-50Gi
- **Replicas**: 1-2

### Production
- **CPU**: 2000m-8000m
- **Memory**: 8Gi-32Gi
- **Disk**: 50Gi-200Gi
- **Replicas**: 2-3

---

## Common Patterns

### Internal Monitoring (No Ingress)
- Set `ingress.is_enabled = false`
- Access via internal Kubernetes service
- More secure for internal-only monitoring

### External Monitoring (With Ingress)
- Set `ingress.is_enabled = true`
- Provide `ingress.dns_domain` for DNS
- External hostname created automatically
- Secure with authentication/authorization

### Persistence Strategy
- **Development**: Disable persistence for faster iterations
- **Production**: Enable persistence for data retention
- Set appropriate disk size based on retention and cardinality

---

## Outputs

The module provides these outputs:

- `namespace`: Kubernetes namespace where Prometheus is deployed
- `service`: Name of the Prometheus service (for internal connections)
- `port_forward_command`: Command for local port-forwarding
- `kube_endpoint`: Internal cluster endpoint (FQDN)
- `external_hostname`: External hostname when ingress is enabled
- `internal_hostname`: Internal hostname for private access

Example usage:

```hcl
output "prometheus_access_command" {
  value = module.prometheus_production.port_forward_command
}

output "prometheus_external_url" {
  value = "https://${module.prometheus_production.external_hostname}"
}
```

---

## Security Considerations

1. **Authentication**: Use OAuth2 proxy or built-in Prometheus authentication
2. **Authorization**: Configure RBAC for ServiceMonitor creation
3. **Network Policies**: Restrict access to Prometheus endpoints
4. **TLS**: Enable TLS for ingress endpoints
5. **Secrets Management**: Store sensitive configuration in Kubernetes secrets

---

## Troubleshooting

### Deployment Issues
- Ensure Prometheus Operator is installed
- Check operator logs: `kubectl logs -n prometheus-operator -l app=prometheus-operator`
- Verify CRD exists: `kubectl get crd prometheuses.monitoring.coreos.com`

### Storage Issues
- Verify StorageClass exists for persistent volumes
- Check PVC status: `kubectl get pvc -n <namespace>`
- Review disk size matches spec

### Resource Issues
- Monitor resource usage: `kubectl top pods -n <namespace>`
- Adjust requests/limits based on actual usage
- Consider vertical pod autoscaling

### Ingress Issues
- Verify ingress controller is installed
- Check ingress resource: `kubectl get ingress -n <namespace>`
- Verify DNS records point to ingress load balancer

---

## Additional Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Prometheus Operator Documentation](https://prometheus-operator.dev/)
- [kube-prometheus-stack Helm Chart](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)
- [Kubernetes Monitoring Best Practices](https://kubernetes.io/docs/tasks/debug-application-cluster/resource-usage-monitoring/)

