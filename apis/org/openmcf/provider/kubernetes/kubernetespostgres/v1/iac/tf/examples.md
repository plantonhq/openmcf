# PostgreSQL Kubernetes Terraform Module - Examples

This document provides examples demonstrating various PostgreSQL configurations using the Zalando PostgreSQL Operator.

## Prerequisites

The **Zalando PostgreSQL Operator** must be installed on your Kubernetes cluster before deploying these examples.

---

## Example 1: Minimal Development Database

```hcl
module "postgres_dev" {
  source = "../../tf"

  metadata = {
    name = "dev-postgres"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "dev-postgres"
    create_namespace = true
    container = {
      replicas = 1
      disk_size = "10Gi"
      resources = {
        requests = {
          cpu    = "250m"
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
- Minimal resources
- No external access (ingress disabled)
- Suitable for local development and testing

---

## Example 2: Production Database with Ingress

```hcl
module "postgres_production" {
  source = "../../tf"

  metadata = {
    name = "prod-postgres"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "prod-postgres"
    create_namespace = true
    container = {
      replicas = 1
      disk_size = "50Gi"
      resources = {
        requests = {
          cpu    = "500m"
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
      hostname = "postgres-prod.example.com"
    }
  }
}
```

**Key points**:
- Production-grade resources
- External access via LoadBalancer with DNS
- Hostname managed by external-dns
- Suitable for external database connections

---

## Example 3: High Resource Database for Analytics

```hcl
module "postgres_analytics" {
  source = "../../tf"

  metadata = {
    name = "analytics-db"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "analytics-db"
    create_namespace = true
    container = {
      replicas = 1
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
      enabled  = true
      hostname = "analytics-db.example.com"
    }
  }
}
```

**Key points**:
- High memory for analytical workloads
- Large disk allocation
- Optimized for complex queries and data processing

---

## Example 4: Minimal Resource Database

```hcl
module "postgres_minimal" {
  source = "../../tf"

  metadata = {
    name = "minimal-postgres"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "minimal-postgres"
    create_namespace = true
    container = {
      replicas = 1
      disk_size = "5Gi"
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
- Cost-optimized configuration
- Suitable for testing or low-traffic applications

---

## Example 5: Database with External Access

```hcl
module "postgres_public" {
  source = "../../tf"

  metadata = {
    name = "public-postgres"
    org  = "my-org"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "public-postgres"
    create_namespace = true
    container = {
      replicas = 1
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
      enabled  = true
      hostname = "postgres.example.com"
    }
  }
}

# Output the connection details
output "postgres_namespace" {
  value = module.postgres_public.namespace
}

output "postgres_service_name" {
  value = module.postgres_public.service_name
}
```

**Key points**:
- LoadBalancer service with DNS entry
- External hostname for remote access
- Port 5432 exposed via LoadBalancer
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
  }))
  default = {
    dev = {
      cpu_request    = "100m"
      cpu_limit      = "500m"
      memory_request = "128Mi"
      memory_limit   = "512Mi"
      disk_size      = "10Gi"
    }
    staging = {
      cpu_request    = "500m"
      cpu_limit      = "2000m"
      memory_request = "1Gi"
      memory_limit   = "4Gi"
      disk_size      = "50Gi"
    }
    production = {
      cpu_request    = "2000m"
      cpu_limit      = "8000m"
      memory_request = "8Gi"
      memory_limit   = "32Gi"
      disk_size      = "200Gi"
    }
  }
}

module "postgres_env" {
  source = "../../tf"

  metadata = {
    name = "app-postgres"
    env  = var.environment
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "app-postgres"
    create_namespace = true
    container = {
      replicas = 1
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
      hostname = var.environment != "dev" ? "postgres-${var.environment}.example.com" : ""
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

All deployments use the **Zalando PostgreSQL Operator** which:
- Creates PostgreSQL custom resources (`acid.zalan.do/v1`)
- Manages StatefulSets for database pods
- Provides automated backup and restore capabilities
- Handles connection pooling via PgBouncer
- Uses Patroni for high availability and failover
- Supports PostgreSQL 14 with customizable configurations

---

## Resource Planning Guidelines

### Development/Testing
- **CPU**: 50m-500m
- **Memory**: 128Mi-512Mi
- **Disk**: 5Gi-10Gi
- **Replicas**: 1

### Staging
- **CPU**: 500m-2000m
- **Memory**: 512Mi-4Gi
- **Disk**: 20Gi-50Gi
- **Replicas**: 1

### Production
- **CPU**: 1000m-8000m
- **Memory**: 2Gi-32Gi
- **Disk**: 50Gi-500Gi
- **Replicas**: 1 (or 3 for HA via operator configuration)

---

## Common Patterns

### Internal Service (No Ingress)
- Set `ingress.enabled = false`
- Access via internal Kubernetes service
- More secure for internal applications

### External Service (With Ingress)
- Set `ingress.enabled = true`
- Provide `ingress.hostname` for DNS
- LoadBalancer service created automatically
- External-DNS annotation for automatic DNS management

### Resource Limits
- Always set both requests and limits
- Limits should be 2-4x requests for burstable workloads
- Memory limits prevent OOM kills of other pods

---

## Outputs

The module provides these outputs:

- `namespace`: Kubernetes namespace where PostgreSQL is deployed
- `service_name`: Name of the PostgreSQL service (for internal connections)
- `service_port`: PostgreSQL service port (default: 5432)
- `load_balancer_ingress_hostname`: External hostname when ingress is enabled
- `load_balancer_ingress_ip`: External IP when ingress is enabled

Example usage:

```hcl
output "connection_string" {
  value = "postgresql://<username>@${module.postgres_production.service_name}.${module.postgres_production.namespace}.svc.cluster.local:5432/<database>"
}
```

---

## Security Considerations

1. **Credentials**: Managed by the Zalando operator, stored in Kubernetes secrets
2. **Network Policies**: Consider implementing network policies to restrict access
3. **TLS**: Operator supports TLS for client connections
4. **Backups**: Configure operator-level backup settings for WAL archiving
5. **External Access**: Use ingress only when necessary, prefer internal services

---

## Troubleshooting

### Deployment Issues
- Ensure Zalando operator is installed and running
- Check operator logs: `kubectl logs -n postgres-operator -l app.kubernetes.io/name=postgres-operator`
- Verify CRD exists: `kubectl get crd postgresqls.acid.zalan.do`

### Connection Issues
- Verify service exists: `kubectl get svc -n <namespace>`
- Check pod status: `kubectl get pods -n <namespace>`
- Review PostgreSQL logs: `kubectl logs -n <namespace> <pod-name>`

### Resource Issues
- Monitor resource usage: `kubectl top pods -n <namespace>`
- Adjust requests/limits based on actual usage
- Consider vertical pod autoscaling for dynamic sizing

---

## Additional Resources

- [Zalando PostgreSQL Operator Documentation](https://postgres-operator.readthedocs.io/)
- [PostgreSQL Best Practices](https://www.postgresql.org/docs/current/admin.html)
- [Kubernetes StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/)

