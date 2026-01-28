# Terraform Examples for KubernetesNeo4j

This document provides comprehensive Terraform examples for deploying Neo4j Community Edition on Kubernetes.

## Namespace Management

The `create_namespace` flag controls whether the Neo4j component creates the namespace or expects it to already exist:

- **`create_namespace = true`** (recommended for new deployments): The component creates and manages the namespace
- **`create_namespace = false`**: The component expects the namespace to already exist (useful when namespace is managed separately or by another component)

---

## Example 1: Basic Neo4j Deployment

A minimal configuration for development and testing.

```hcl
module "neo4j_basic" {
  source = "./path/to/kubernetesneo4j/v1/iac/tf"

  metadata = {
    name = "basic-neo4j"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    container = {
      resources = {
        limits = {
          cpu    = "1000m"
          memory = "1Gi"
        }
        requests = {
          cpu    = "50m"
          memory = "100Mi"
        }
      }
      persistence_enabled = true
      disk_size          = "5Gi"
    }

    ingress = {
      enabled  = false
      hostname = ""
    }
  }
}

# Output connection details
output "neo4j_basic_namespace" {
  value = module.neo4j_basic.namespace
}

output "neo4j_basic_bolt_uri" {
  value = module.neo4j_basic.bolt_uri_kube_endpoint
}

output "neo4j_basic_http_uri" {
  value = module.neo4j_basic.http_uri_kube_endpoint
}
```

## Example 2: Neo4j with Custom Memory Configuration

Optimize memory allocation for better performance.

```hcl
module "neo4j_memory_optimized" {
  source = "./path/to/kubernetesneo4j/v1/iac/tf"

  metadata = {
    name = "neo4j-memory-optimized"
    env  = "development"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    container = {
      resources = {
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
        requests = {
          cpu    = "500m"
          memory = "2Gi"
        }
      }
      persistence_enabled = true
      disk_size          = "20Gi"
    }

    memory_config = {
      heap_max   = "1Gi"   # 25% of 4GB for heap
      page_cache = "2Gi"   # 50% of 4GB for page cache
    }

    ingress = {
      enabled  = false
      hostname = ""
    }
  }
}
```

## Example 3: Neo4j with External Access

Enable external access via LoadBalancer and external-DNS.

```hcl
module "neo4j_external" {
  source = "./path/to/kubernetesneo4j/v1/iac/tf"

  metadata = {
    name = "neo4j-external"
    org  = "my-company"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    container = {
      resources = {
        limits = {
          cpu    = "2000m"
          memory = "2Gi"
        }
        requests = {
          cpu    = "100m"
          memory = "512Mi"
        }
      }
      persistence_enabled = true
      disk_size          = "10Gi"
    }

    ingress = {
      enabled  = true
      hostname = "neo4j.example.com"
    }
  }
}

# Output external connection details
output "neo4j_external_hostname" {
  value = module.neo4j_external.external_hostname
}

output "neo4j_browser_url" {
  value = "http://${module.neo4j_external.external_hostname}:7474"
}

output "neo4j_bolt_external_url" {
  value = "bolt://${module.neo4j_external.external_hostname}:7687"
}
```

## Example 4: Production Neo4j Deployment

A production-ready configuration with optimal settings.

```hcl
module "neo4j_production" {
  source = "./path/to/kubernetesneo4j/v1/iac/tf"

  metadata = {
    name = "production-neo4j"
    id   = "prod-neo4j-001"
    org  = "acme-corp"
    env  = "production"
    labels = {
      team       = "data-engineering"
      component  = "graph-database"
      cost_center = "engineering"
    }
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    container = {
      resources = {
        limits = {
          cpu    = "4000m"
          memory = "8Gi"
        }
        requests = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
      persistence_enabled = true
      disk_size          = "100Gi"
    }

    memory_config = {
      heap_max   = "2Gi"   # 25% of 8GB
      page_cache = "4Gi"   # 50% of 8GB
    }

    ingress = {
      enabled  = true
      hostname = "neo4j.prod.example.com"
    }
  }
}

# Comprehensive outputs for production
output "neo4j_production_details" {
  description = "Complete connection details for production Neo4j"
  value = {
    namespace            = module.neo4j_production.namespace
    service             = module.neo4j_production.service
    bolt_uri_internal   = module.neo4j_production.bolt_uri_kube_endpoint
    http_uri_internal   = module.neo4j_production.http_uri_kube_endpoint
    external_hostname   = module.neo4j_production.external_hostname
    username            = module.neo4j_production.username
    password_secret     = module.neo4j_production.password_secret_name
    port_forward_cmd    = module.neo4j_production.port_forward_command
  }
}
```

## Example 5: Development Neo4j (Minimal Resources)

Lightweight deployment for local development.

```hcl
module "neo4j_dev" {
  source = "./path/to/kubernetesneo4j/v1/iac/tf"

  metadata = {
    name = "dev-neo4j"
    env  = "development"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    container = {
      resources = {
        limits = {
          cpu    = "500m"
          memory = "512Mi"
        }
        requests = {
          cpu    = "50m"
          memory = "128Mi"
        }
      }
      persistence_enabled = false  # No persistence for dev
      disk_size          = "1Gi"
    }

    ingress = {
      enabled  = false
      hostname = ""
    }
  }
}

# Simple dev outputs
output "neo4j_dev_port_forward" {
  description = "Run this command to access Neo4j locally"
  value       = module.neo4j_dev.port_forward_command
}
```

## Example 6: Neo4j with All Metadata

Complete example showing all metadata options.

```hcl
module "neo4j_complete" {
  source = "./path/to/kubernetesneo4j/v1/iac/tf"

  metadata = {
    name = "complete-neo4j"
    id   = "neo4j-abc123"
    org  = "acme-corporation"
    env  = "staging"
    labels = {
      application = "knowledge-graph"
      team        = "data-science"
      project     = "recommendation-engine"
    }
    tags = ["graph-db", "neo4j", "analytics"]
    version = {
      id      = "v1.0.0"
      message = "Initial production release"
    }
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    container = {
      resources = {
        limits = {
          cpu    = "3000m"
          memory = "6Gi"
        }
        requests = {
          cpu    = "1000m"
          memory = "3Gi"
        }
      }
      persistence_enabled = true
      disk_size          = "50Gi"
    }

    memory_config = {
      heap_max   = "1536m"  # ~25% of 6GB
      page_cache = "3Gi"    # ~50% of 6GB
    }

    ingress = {
      enabled  = true
      hostname = "neo4j.staging.example.com"
    }
  }
}
```

## Example 7: Multiple Neo4j Instances

Deploy multiple independent Neo4j instances in the same cluster.

```hcl
# Application database
module "neo4j_app" {
  source = "./path/to/kubernetesneo4j/v1/iac/tf"

  metadata = {
    name = "app-graph-db"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    container = {
      resources = {
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
        requests = {
          cpu    = "500m"
          memory = "2Gi"
        }
      }
      persistence_enabled = true
      disk_size          = "30Gi"
    }

    memory_config = {
      heap_max   = "1Gi"
      page_cache = "2Gi"
    }

    ingress = {
      enabled  = true
      hostname = "app-neo4j.example.com"
    }
  }
}

# Analytics database
module "neo4j_analytics" {
  source = "./path/to/kubernetesneo4j/v1/iac/tf"

  metadata = {
    name = "analytics-graph-db"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    container = {
      resources = {
        limits = {
          cpu    = "4000m"
          memory = "8Gi"
        }
        requests = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
      persistence_enabled = true
      disk_size          = "100Gi"
    }

    memory_config = {
      heap_max   = "2Gi"
      page_cache = "4Gi"
    }

    ingress = {
      enabled  = true
      hostname = "analytics-neo4j.example.com"
    }
  }
}

# Compare outputs
output "databases" {
  value = {
    application = {
      namespace = module.neo4j_app.namespace
      bolt_uri  = module.neo4j_app.bolt_uri_kube_endpoint
      hostname  = module.neo4j_app.external_hostname
    }
    analytics = {
      namespace = module.neo4j_analytics.namespace
      bolt_uri  = module.neo4j_analytics.bolt_uri_kube_endpoint
      hostname  = module.neo4j_analytics.external_hostname
    }
  }
}
```

## Retrieving the Neo4j Password

After deployment, retrieve the auto-generated password:

### Using kubectl

```bash
# Get the secret name from Terraform output
SECRET_NAME=$(terraform output -raw password_secret_name)
NAMESPACE=$(terraform output -raw namespace)

# Retrieve the password
kubectl get secret $SECRET_NAME -n $NAMESPACE \
  -o jsonpath='{.data.neo4j-password}' | base64 -d && echo
```

### Using Terraform Data Source

```hcl
data "kubernetes_secret_v1" "neo4j_password" {
  metadata {
    name      = module.neo4j_production.password_secret_name
    namespace = module.neo4j_production.namespace
  }

  depends_on = [module.neo4j_production]
}

output "neo4j_password" {
  description = "Neo4j admin password (sensitive)"
  value       = data.kubernetes_secret_v1.neo4j_password.data["neo4j-password"]
  sensitive   = true
}

# Access the password
# terraform output -raw neo4j_password
```

## Connecting to Neo4j

### From Within Kubernetes Cluster

```cypher
// Using Bolt protocol
bolt://<service_fqdn>:7687
// Username: neo4j
// Password: <from secret>
```

### Using Port Forwarding

```bash
# Get the port-forward command
terraform output -raw port_forward_command

# Execute it (example output)
kubectl port-forward -n production-neo4j service/production-neo4j-neo4j 7474:7474

# Access Neo4j Browser at http://localhost:7474
```

### From External Clients (if ingress enabled)

```bash
# Neo4j Browser
http://<external_hostname>:7474

# Bolt connection
bolt://<external_hostname>:7687
```

## Memory Configuration Guidelines

### Small Deployment (< 2GB)
```hcl
memory_config = {
  heap_max   = "512m"
  page_cache = "512m"
}
```

### Medium Deployment (4-8GB)
```hcl
memory_config = {
  heap_max   = "1Gi"
  page_cache = "2Gi"
}
```

### Large Deployment (16GB+)
```hcl
memory_config = {
  heap_max   = "4Gi"
  page_cache = "8Gi"
}
```

## Common Patterns

### With Terraform Workspace

```hcl
locals {
  environment = terraform.workspace
}

module "neo4j" {
  source = "./path/to/kubernetesneo4j/v1/iac/tf"

  metadata = {
    name = "${local.environment}-neo4j"
    env  = local.environment
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    container = {
      resources = local.environment == "production" ? {
        limits = {
          cpu    = "4000m"
          memory = "8Gi"
        }
        requests = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      } : {
        limits = {
          cpu    = "1000m"
          memory = "1Gi"
        }
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
      }
      persistence_enabled = local.environment == "production"
      disk_size          = local.environment == "production" ? "100Gi" : "10Gi"
    }

    ingress = {
      enabled  = local.environment == "production"
      hostname = local.environment == "production" ? "neo4j.prod.example.com" : ""
    }
  }
}
```

## Best Practices

1. **Always enable persistence for production** (`persistence_enabled = true`)
2. **Size disk appropriately** - account for data growth over time
3. **Configure memory limits** - use memory_config for optimal performance
4. **Use ingress for external access** - enable only when needed
5. **Secure credentials** - rotate default password after deployment
6. **Monitor resources** - ensure CPU and memory limits are appropriate
7. **Use labels** - tag resources for organization and cost tracking

## Notes

- Neo4j Community Edition is single-node only (no clustering)
- For high availability, consider Neo4j Enterprise Edition
- The Helm chart version is pinned to `2025.03.0` for stability
- Persistent volumes cannot be shrunk after creation
- External access requires a LoadBalancer-capable cluster

