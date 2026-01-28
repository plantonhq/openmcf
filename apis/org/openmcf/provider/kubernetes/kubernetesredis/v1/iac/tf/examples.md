# Redis Kubernetes Terraform Module - Examples

This document provides examples demonstrating various Redis configurations using the Bitnami Redis Helm chart.

## Prerequisites

- Terraform >= 1.0
- Kubernetes cluster configured
- **Bitnami Redis Helm chart** (deployed via this module)

---

## Example 1: Minimal Development Redis

```hcl
module "redis_dev" {
  source = "../../tf"

  metadata = {
    name = "dev-redis"
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
      enabled  = false
      hostname = ""
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

## Example 2: Production Redis with Persistence

```hcl
module "redis_production" {
  source = "../../tf"

  metadata = {
    name = "prod-redis"
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
      replicas = 1
      persistence_enabled = true
      disk_size = "10Gi"
      resources = {
        requests = {
          cpu    = "100m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "2000m"
          memory = "2Gi"
        }
      }
    }
    
    ingress = {
      enabled  = true
      hostname = "redis-prod.example.com"
    }
  }
}
```

**Key points**:
- Persistent storage enabled with 10Gi disk
- Production-grade resources
- External access via LoadBalancer with DNS
- Hostname managed by external-dns
- Suitable for production caching and session storage

---

## Example 3: High-Availability Redis Cluster

```hcl
module "redis_ha" {
  source = "../../tf"

  metadata = {
    name = "ha-redis"
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
      disk_size = "20Gi"
      resources = {
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "4000m"
          memory = "4Gi"
        }
      }
    }
    
    ingress = {
      enabled  = true
      hostname = "redis-ha.example.com"
    }
  }
}

# Output the connection details
output "redis_namespace" {
  value = module.redis_ha.namespace
}

output "redis_endpoint" {
  value = module.redis_ha.kube_endpoint
}

output "redis_password_secret" {
  value = module.redis_ha.password_secret_name
}
```

**Key points**:
- High availability with 3 replicas
- Large persistent storage (20Gi)
- High resource allocation
- Optimized for high-traffic applications

---

## Example 4: Minimal Resource Redis

```hcl
module "redis_minimal" {
  source = "../../tf"

  metadata = {
    name = "minimal-redis"
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
      enabled  = false
      hostname = ""
    }
  }
}
```

**Key points**:
- Minimal resources for small workloads
- No persistence
- Cost-optimized configuration
- Suitable for testing or low-traffic caching

---

## Example 5: Redis with External Access

```hcl
module "redis_public" {
  source = "../../tf"

  metadata = {
    name = "public-redis"
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
      disk_size = "5Gi"
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
    
    ingress = {
      enabled  = true
      hostname = "redis.example.com"
    }
  }
}

# Access the Redis password
output "redis_password_command" {
  value = "kubectl get secret ${module.redis_public.password_secret_name} -n ${module.redis_public.namespace} -o jsonpath='{.data.${module.redis_public.password_secret_key}}' | base64 -d"
}
```

**Key points**:
- LoadBalancer service with DNS entry
- External hostname for remote access
- Port 6379 exposed via LoadBalancer
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
      cpu_request    = "100m"
      cpu_limit      = "2000m"
      memory_request = "512Mi"
      memory_limit   = "2Gi"
      disk_size      = "10Gi"
      replicas       = 1
      persistence    = true
    }
    production = {
      cpu_request    = "500m"
      cpu_limit      = "4000m"
      memory_request = "2Gi"
      memory_limit   = "8Gi"
      disk_size      = "50Gi"
      replicas       = 3
      persistence    = true
    }
  }
}

module "redis_env" {
  source = "../../tf"

  metadata = {
    name = "app-redis"
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
      enabled  = var.environment != "dev"
      hostname = var.environment != "dev" ? "redis-${var.environment}.example.com" : ""
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

All deployments use the **Bitnami Redis Helm chart** which provides:

- **Redis Server**: High-performance in-memory data store
- **Sentinel** (for HA configurations): Automatic failover and monitoring
- **Persistence Options**: RDB snapshots and AOF logs
- **Password Authentication**: Secure password stored in Kubernetes secrets
- **Resource Management**: CPU and memory limits/requests
- **Horizontal Scaling**: Multiple replicas for high availability

---

## Resource Planning Guidelines

### Development/Testing
- **CPU**: 50m-500m
- **Memory**: 128Mi-512Mi
- **Disk**: N/A (no persistence)
- **Replicas**: 1

### Staging
- **CPU**: 100m-2000m
- **Memory**: 512Mi-2Gi
- **Disk**: 5Gi-10Gi
- **Replicas**: 1

### Production
- **CPU**: 500m-4000m
- **Memory**: 2Gi-8Gi
- **Disk**: 10Gi-50Gi
- **Replicas**: 1-3

---

## Common Patterns

### Internal Service (No Ingress)
- Set `ingress.enabled = false`
- Access via internal Kubernetes service
- More secure for internal-only caching

### External Service (With Ingress)
- Set `ingress.enabled = true`
- Provide `ingress.hostname` for DNS
- LoadBalancer service created automatically
- External-DNS annotation for automatic DNS management

### Persistence Strategy
- **Development**: Disable persistence for faster iterations
- **Production**: Enable persistence for data durability
- Set appropriate disk size based on data volume

### Resource Limits
- Always set both requests and limits
- Limits should be 2-4x requests for burstable workloads
- Memory limits prevent OOM kills of other pods

---

## Outputs

The module provides these outputs:

- `namespace`: Kubernetes namespace where Redis is deployed
- `service`: Name of the Redis service (for internal connections)
- `port_forward_command`: Command for local port-forwarding
- `kube_endpoint`: Internal cluster endpoint (FQDN)
- `external_hostname`: External hostname when ingress is enabled
- `password_secret_name`: Name of the Kubernetes secret containing Redis password
- `password_secret_key`: Key in the secret containing the password
- `load_balancer_ingress_hostname`: LoadBalancer hostname when ingress is enabled
- `load_balancer_ingress_ip`: LoadBalancer IP when ingress is enabled

Example usage:

```hcl
output "redis_connection_string" {
  value = "redis://:PASSWORD@${module.redis_production.kube_endpoint}:6379"
  sensitive = true
}

output "get_password_command" {
  value = "kubectl get secret ${module.redis_production.password_secret_name} -n ${module.redis_production.namespace} -o jsonpath='{.data.${module.redis_production.password_secret_key}}' | base64 -d"
}
```

---

## Security Considerations

1. **Authentication**: Redis is deployed with password authentication enabled
2. **Password Management**: Password stored in Kubernetes secret, rotatable
3. **Network Policies**: Consider implementing network policies to restrict access
4. **TLS**: For external access, use TLS termination at the ingress level
5. **External Access**: Use ingress only when necessary, prefer internal services

---

## Troubleshooting

### Deployment Issues
- Check Helm release: `helm list -n <namespace>`
- Verify pod status: `kubectl get pods -n <namespace>`
- Review Redis logs: `kubectl logs -n <namespace> <pod-name>`

### Connection Issues
- Verify service exists: `kubectl get svc -n <namespace>`
- Check password: `kubectl get secret <secret-name> -n <namespace> -o yaml`
- Test connection: Use `redis-cli` from a pod in the cluster

### Resource Issues
- Monitor resource usage: `kubectl top pods -n <namespace>`
- Adjust requests/limits based on actual usage
- Consider vertical pod autoscaling

### Persistence Issues
- Verify PVC status: `kubectl get pvc -n <namespace>`
- Check StorageClass: `kubectl get storageclass`
- Review disk size matches spec

---

## Additional Resources

- [Redis Documentation](https://redis.io/docs/)
- [Bitnami Redis Helm Chart](https://github.com/bitnami/charts/tree/main/bitnami/redis)
- [Kubernetes Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
- [Valkey (Redis Fork)](https://valkey.io/) - Open-source alternative

