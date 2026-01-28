# Kubernetes Deployment - Terraform Examples

This document provides comprehensive examples for deploying microservices to Kubernetes using the **KubernetesDeployment** Terraform module. Each example demonstrates different configuration patterns including environment variables, secrets, scaling, ingress, and more.

> **Note:** These examples show how to use the Terraform module directly. The module expects variables to match the protobuf specification defined in the KubernetesDeployment API.

---

## 1. Minimal Configuration

A basic example deploying a containerized application with default settings.

```hcl
module "minimal_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "minimal-example"
    id   = "minimal-example-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "minimal-example"
    create_namespace = true
    version = "main"
    
    container = {
      app = {
        image = {
          repo = "nginx"
          tag  = "latest"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 80
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = false
          }
        ]
        
        resources = {
          requests = {
            cpu    = "100m"
            memory = "128Mi"
          }
          limits = {
            cpu    = "500m"
            memory = "512Mi"
          }
        }
        
        env = {
          variables = {}
          secrets   = {}
        }
      }
    }
  }
}
```

**Key Points:**
- Minimal configuration with single container port
- No external ingress exposure
- Default resource allocation

---

## 2. Microservice with Environment Variables

Pass non-sensitive configuration values as environment variables.

```hcl
module "env_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "env-example"
    id   = "env-example-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "env-example"
    create_namespace = true
    version = "main"
    
    container = {
      app = {
        image = {
          repo = "my-org/my-app"
          tag  = "1.0.0"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 8080
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = false
          }
        ]
        
        env = {
          variables = {
            LOG_LEVEL          = "debug"
            FEATURE_X_ENABLED  = "true"
            MAX_CONNECTIONS    = "100"
          }
          secrets = {}
        }
        
        resources = {
          requests = {
            cpu    = "200m"
            memory = "256Mi"
          }
          limits = {
            cpu    = "1"
            memory = "1Gi"
          }
        }
      }
    }
  }
}
```

**Key Points:**
- `env.variables` for non-sensitive configuration
- Environment variables accessible in the container

---

## 3. Using Secrets for Sensitive Data (Direct String Values)

Store sensitive data directly. A Kubernetes Secret is automatically created.
This approach is suitable for development and testing.

```hcl
module "db_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "db-credentials-example"
    id   = "db-app-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "db-credentials-example"
    create_namespace = true
    version = "main"
    
    container = {
      app = {
        image = {
          repo = "my-org/database-connector"
          tag  = "stable"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 9090
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = true
          }
        ]
        
        env = {
          variables = {
            DB_HOST = "db.prod.svc.cluster.local"
            DB_PORT = "5432"
          }
          
          # Secrets with direct string values
          secrets = {
            DB_PASSWORD = {
              value = var.db_password
            }
            API_KEY = {
              value = var.api_key
            }
          }
        }
        
        resources = {
          requests = {
            cpu    = "100m"
            memory = "200Mi"
          }
          limits = {
            cpu    = "500m"
            memory = "1Gi"
          }
        }
      }
    }
  }
}

# Variables for sensitive data
variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "api_key" {
  description = "API key for external service"
  type        = string
  sensitive   = true
}
```

**Key Points:**
- `value` for direct secret values
- Secrets stored in a Kubernetes Secret resource created by the module
- Use Terraform variables with `sensitive = true` to protect values

---

## 3b. Using Secrets with Kubernetes Secret References

Reference existing Kubernetes Secrets for production deployments.
This is the recommended approach to avoid storing sensitive values in configuration files.

```hcl
module "db_microservice_secret_ref" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "db-credentials-example"
    id   = "db-app-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "db-credentials-example"
    create_namespace = true
    version = "main"
    
    container = {
      app = {
        image = {
          repo = "my-org/database-connector"
          tag  = "stable"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 9090
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = true
          }
        ]
        
        env = {
          variables = {
            DB_HOST = "db.prod.svc.cluster.local"
            DB_PORT = "5432"
          }
          
          # Reference existing Kubernetes Secrets
          secrets = {
            DB_PASSWORD = {
              secret_ref = {
                name = "postgres-credentials"
                key  = "password"
              }
            }
            API_KEY = {
              secret_ref = {
                name = "external-api-secrets"
                key  = "api-key"
              }
            }
          }
        }
        
        resources = {
          requests = {
            cpu    = "100m"
            memory = "200Mi"
          }
          limits = {
            cpu    = "500m"
            memory = "1Gi"
          }
        }
      }
    }
  }
}
```

**Key Points:**
- `secret_ref` references an existing Kubernetes Secret by name and key
- The referenced Secret must exist in the cluster before deployment
- Avoids storing sensitive values in Terraform state or configuration
- Recommended for production environments

---

## 4. Microservice with Sidecar Containers

Deploy a microservice with logging or monitoring sidecars.

```hcl
module "sidecar_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "sidecar-example"
    id   = "sidecar-app-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "sidecar-example"
    create_namespace = true
    version = "v2"
    
    container = {
      app = {
        image = {
          repo = "my-org/main-service"
          tag  = "2.3.4"
        }
        
        ports = [
          {
            name             = "main-port"
            container_port   = 8080
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = false
          }
        ]
        
        resources = {
          requests = {
            cpu    = "200m"
            memory = "256Mi"
          }
          limits = {
            cpu    = "1"
            memory = "1Gi"
          }
        }
        
        env = {
          variables = {}
          secrets   = {}
        }
      }
      
      sidecars = [
        {
          name  = "logger"
          image = "my-org/log-agent:latest"
          
          ports = [
            {
              name           = "agent-port"
              container_port = 4000
              protocol       = "TCP"
            }
          ]
          
          resources = {
            requests = {
              cpu    = "50m"
              memory = "64Mi"
            }
            limits = {
              cpu    = "100m"
              memory = "128Mi"
            }
          }
          
          env = [
            {
              name  = "LOG_LEVEL"
              value = "info"
            }
          ]
        }
      ]
    }
  }
}
```

**Key Points:**
- Main application and sidecar containers in same pod
- Each container has independent resource allocation
- Sidecars share network namespace with main container

---

## 5. Enabling Ingress with Istio

Expose your microservice externally using Istio Gateway.

```hcl
module "ingress_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "ingress-example"
    id   = "api-service-prod"
    org  = "my-org"
    env  = "production"
    labels = {
      team = "platform"
    }
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "ingress-example"
    create_namespace = true
    version = "main"
    
    container = {
      app = {
        image = {
          repo = "my-org/web-api"
          tag  = "v1.1"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 8080
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = true  # Enable ingress for this port
          }
        ]
        
        resources = {
          requests = {
            cpu    = "100m"
            memory = "128Mi"
          }
          limits = {
            cpu    = "1"
            memory = "1Gi"
          }
        }
        
        env = {
          variables = {}
          secrets   = {}
        }
      }
    }
    
    # Enable ingress
    ingress = {
      is_enabled = true
      dns_domain = "example.org"
    }
  }
}

# Outputs
output "external_url" {
  description = "External URL for accessing the microservice"
  value       = "https://api-service-prod.example.org"
}

output "internal_url" {
  description = "Internal URL for cluster-internal access"
  value       = module.ingress_microservice.service_fqdn
}
```

**Key Points:**
- `ingress.is_enabled = true` enables external access
- Service accessible at `<resource_id>.<dns_domain>`
- Automatically configures Istio Gateway and VirtualService
- TLS certificate provisioned automatically

---

## 6. Horizontal Pod Autoscaling (HPA)

Configure autoscaling based on CPU and memory utilization.

```hcl
module "hpa_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "hpa-example"
    id   = "scalable-service-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "hpa-example"
    create_namespace = true
    version = "v3.0"
    
    container = {
      app = {
        image = {
          repo = "my-org/hpa-service"
          tag  = "stable"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 3000
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = false
          }
        ]
        
        resources = {
          requests = {
            cpu    = "250m"
            memory = "256Mi"
          }
          limits = {
            cpu    = "2"
            memory = "2Gi"
          }
        }
        
        env = {
          variables = {}
          secrets   = {}
        }
      }
    }
    
    # Availability and autoscaling configuration
    availability = {
      min_replicas = 2
      
      horizontal_pod_autoscaling = {
        is_enabled                      = true
        target_cpu_utilization_percent  = 70.0
        target_memory_utilization       = "1Gi"
      }
    }
  }
}
```

**Key Points:**
- `min_replicas` sets baseline pod count
- HPA scales up when CPU exceeds 70% utilization
- HPA scales based on memory usage reaching 1Gi
- Automatic scale-down when demand decreases

---

## 7. Production-Ready with Health Probes

Complete production configuration with liveness, readiness, and startup probes.

```hcl
module "production_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "production-api"
    id   = "prod-api-v1"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "production-api"
    create_namespace = true
    version = "v1.0"
    
    container = {
      app = {
        image = {
          repo = "my-org/production-service"
          tag  = "1.0.0"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 8080
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = true
          },
          {
            name             = "metrics"
            container_port   = 9090
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 9090
            is_ingress_port  = false
          }
        ]
        
        env = {
          variables = {
            ENVIRONMENT = "production"
            LOG_FORMAT  = "json"
          }
          secrets = {
            DATABASE_URL = {
              secret_ref = {
                name = "database-credentials"
                key  = "connection-string"
              }
            }
          }
        }
        
        resources = {
          requests = {
            cpu    = "500m"
            memory = "512Mi"
          }
          limits = {
            cpu    = "2"
            memory = "2Gi"
          }
        }
        
        # Health probes
        liveness_probe = {
          http_get = {
            path = "/healthz/live"
            port = 8080
          }
          initial_delay_seconds = 30
          period_seconds        = 10
          timeout_seconds       = 5
          failure_threshold     = 3
        }
        
        readiness_probe = {
          http_get = {
            path = "/healthz/ready"
            port = 8080
          }
          initial_delay_seconds = 10
          period_seconds        = 5
          timeout_seconds       = 3
          failure_threshold     = 3
        }
        
        startup_probe = {
          http_get = {
            path = "/healthz/startup"
            port = 8080
          }
          initial_delay_seconds = 0
          period_seconds        = 10
          timeout_seconds       = 3
          failure_threshold     = 30
        }
      }
    }
    
    availability = {
      min_replicas = 3
      
      horizontal_pod_autoscaling = {
        is_enabled                     = true
        target_cpu_utilization_percent = 75.0
      }
      
      deployment_strategy = {
        type = "RollingUpdate"
        rolling_update = {
          max_surge       = "25%"
          max_unavailable = "0"
        }
      }
      
      pod_disruption_budget = {
        is_enabled        = true
        min_available     = 2
        max_unavailable   = null
      }
    }
    
    ingress = {
      is_enabled = true
      dns_domain = "api.example.com"
    }
  }
}

# Note: DATABASE_URL references an existing Kubernetes Secret "database-credentials"
# which must be created separately (e.g., via ExternalSecrets, Vault, or manually)
```

**Key Points:**
- Comprehensive health probes for reliability
- Zero-downtime deployments with rolling update strategy
- Pod Disruption Budget prevents simultaneous pod termination
- Multiple replicas for high availability
- HPA for dynamic scaling

---

## 8. Private Container Registry

Use private container registry with authentication.

```hcl
module "private_registry_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "private-app"
    id   = "private-app-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "private-app"
    create_namespace = true
    version = "main"
    
    container = {
      app = {
        image = {
          repo             = "gcr.io/my-project/my-private-image"
          tag              = "v2.0.0"
          pull_secret_name = "gcr-pull-secret"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 8080
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = false
          }
        ]
        
        resources = {
          requests = {
            cpu    = "100m"
            memory = "128Mi"
          }
          limits = {
            cpu    = "1"
            memory = "512Mi"
          }
        }
        
        env = {
          variables = {}
          secrets   = {}
        }
      }
    }
  }

  # Docker config for private registry
  docker_config_json = jsonencode({
    auths = {
      "gcr.io" = {
        username = "_json_key"
        password = var.gcp_service_account_key
        email    = "service-account@project.iam.gserviceaccount.com"
        auth     = base64encode("_json_key:${var.gcp_service_account_key}")
      }
    }
  })
}

variable "gcp_service_account_key" {
  description = "GCP service account key for GCR access"
  type        = string
  sensitive   = true
}
```

**Key Points:**
- `docker_config_json` provides registry authentication
- Automatically creates imagePullSecret
- Supports GCR, ECR, ACR, and other registries

---

## Module Outputs

The module provides useful outputs for integration:

```hcl
output "namespace" {
  description = "Kubernetes namespace where the microservice is deployed"
  value       = module.my_microservice.namespace
}

output "service_name" {
  description = "Name of the Kubernetes Service"
  value       = module.my_microservice.service_name
}

output "service_fqdn" {
  description = "Fully qualified domain name for cluster-internal access"
  value       = module.my_microservice.service_fqdn
}

output "external_hostname" {
  description = "External hostname (if ingress is enabled)"
  value       = module.my_microservice.external_hostname
}

output "internal_hostname" {
  description = "Internal hostname (if ingress is enabled)"
  value       = module.my_microservice.internal_hostname
}

output "port_forward_command" {
  description = "kubectl port-forward command for local access"
  value       = module.my_microservice.port_forward_command
}
```

---

## Best Practices

### 1. Resource Management

Always define resource requests and limits to ensure proper scheduling and QoS:

```hcl
resources = {
  requests = {
    cpu    = "100m"   # Guaranteed minimum
    memory = "128Mi"  # Guaranteed minimum
  }
  limits = {
    cpu    = "500m"   # Maximum allowed
    memory = "512Mi"  # Maximum allowed
  }
}
```

**QoS Classes:**
- **Guaranteed**: requests == limits (highest priority)
- **Burstable**: requests < limits (medium priority)
- **BestEffort**: no requests/limits (lowest priority)

### 2. Health Probes

Implement all three types of probes:

```hcl
# Startup probe - for slow-starting containers
startup_probe = {
  http_get = {
    path = "/healthz/startup"
    port = 8080
  }
  failure_threshold = 30  # Allow 5 minutes (30 * 10s)
  period_seconds    = 10
}

# Liveness probe - detects deadlocked containers
liveness_probe = {
  http_get = {
    path = "/healthz/live"
    port = 8080
  }
  initial_delay_seconds = 30
  period_seconds        = 10
  failure_threshold     = 3
}

# Readiness probe - controls traffic routing
readiness_probe = {
  http_get = {
    path = "/healthz/ready"
    port = 8080
  }
  initial_delay_seconds = 10
  period_seconds        = 5
  failure_threshold     = 3
}
```

### 3. Zero-Downtime Deployments

Configure proper deployment strategy:

```hcl
availability = {
  min_replicas = 3  # At least 3 for HA
  
  deployment_strategy = {
    type = "RollingUpdate"
    rolling_update = {
      max_surge       = "25%"  # Can scale up 25% during update
      max_unavailable = "0"     # Never reduce available pods
    }
  }
  
  pod_disruption_budget = {
    is_enabled    = true
    min_available = 2  # Always keep 2 pods running
  }
}
```

### 4. Environment-Specific Configuration

Use Terraform workspaces or variable files for different environments:

```hcl
# terraform.tfvars.prod
metadata = {
  name = "my-service"
  env  = "production"
}

spec = {
  availability = {
    min_replicas = 5
  }
}

# terraform.tfvars.dev
metadata = {
  name = "my-service"
  env  = "development"
}

spec = {
  availability = {
    min_replicas = 1
  }
}
```

### 5. Secret Management

Never commit secrets to version control. Use one of these approaches:

**Option 1: Direct values with environment variables (dev/test)**

```hcl
# Use environment variables
export TF_VAR_db_password="secret-value"

module "microservice" {
  # ...
  spec = {
    container = {
      app = {
        env = {
          secrets = {
            DB_PASSWORD = {
              value = var.db_password
            }
          }
        }
      }
    }
  }
}
```

**Option 2: Kubernetes Secret references (production - recommended)**

```hcl
# Reference existing Kubernetes Secrets
module "microservice" {
  # ...
  spec = {
    container = {
      app = {
        env = {
          secrets = {
            DB_PASSWORD = {
              secret_ref = {
                name = "postgres-credentials"
                key  = "password"
              }
            }
          }
        }
      }
    }
  }
}
```

**Option 3: Mixed approach**

```hcl
module "microservice" {
  # ...
  spec = {
    container = {
      app = {
        env = {
          secrets = {
            # Dev secret - direct value
            DEBUG_TOKEN = {
              value = "debug-only"
            }
            # Production secret - external reference
            DB_PASSWORD = {
              secret_ref = {
                name = "postgres-credentials"
                key  = "password"
              }
            }
          }
        }
      }
    }
  }
}
```

---

## Troubleshooting

### Check Deployment Status

```bash
# View deployment
kubectl get deployment <name> -n <namespace>

# View pods
kubectl get pods -n <namespace>

# Describe deployment
kubectl describe deployment <name> -n <namespace>

# View logs
kubectl logs -n <namespace> deployment/<name>
```

### Test Service Locally

```bash
# Port forward to local machine
kubectl port-forward -n <namespace> service/<service-name> 8080:80

# Test
curl http://localhost:8080
```

### Debug Pod Issues

```bash
# Get pod events
kubectl get events -n <namespace> --sort-by='.lastTimestamp'

# Check pod logs
kubectl logs -n <namespace> <pod-name>

# Execute into pod
kubectl exec -it -n <namespace> <pod-name> -- /bin/sh
```

---

## 9. Using an Existing Namespace

If the namespace already exists in the cluster (created by another process or team), you can skip namespace creation by setting `create_namespace = false`. This is useful when:
- Multiple deployments share the same namespace
- Namespaces are managed centrally by cluster administrators
- Using GitOps workflows where namespaces are managed separately

```hcl
module "existing_ns_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "existing-ns-example"
    id   = "existing-ns-example-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "shared-services"
    create_namespace = false  # Use existing namespace
    version = "main"
    
    container = {
      app = {
        image = {
          repo = "my-org/my-service"
          tag  = "1.0.0"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 8080
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = false
          }
        ]
        
        resources = {
          requests = {
            cpu    = "100m"
            memory = "128Mi"
          }
          limits = {
            cpu    = "500m"
            memory = "512Mi"
          }
        }
        
        env = {
          variables = {}
          secrets   = {}
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

**Key Points:**
- `create_namespace = false` tells the module to use the existing namespace without creating it
- The namespace "shared-services" must already exist in the cluster
- If the namespace doesn't exist, deployment will fail with a "namespace not found" error
- All resources (deployment, service, secrets) will still be created in the specified namespace

---

## 10. ConfigMaps and Volume Mounts

Deploy applications with configuration files using ConfigMaps and volume mounts.

```hcl
module "configmap_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "config-app"
    id   = "config-app-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "config-app"
    create_namespace = true
    version = "main"
    
    # ConfigMaps to create alongside the deployment
    config_maps = {
      "app-config" = <<-EOT
        database:
          host: postgres.default.svc
          port: 5432
        cache:
          host: redis.default.svc
          port: 6379
        logging:
          level: info
          format: json
      EOT
      
      "routes-config" = <<-EOT
        routes:
          - path: /api/v1
            upstream: http://backend:8080
          - path: /health
            upstream: http://localhost:8080/healthz
      EOT
    }
    
    container = {
      app = {
        image = {
          repo = "my-org/config-app"
          tag  = "v1.0.0"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 8080
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = false
          }
        ]
        
        resources = {
          requests = {
            cpu    = "100m"
            memory = "128Mi"
          }
          limits = {
            cpu    = "500m"
            memory = "512Mi"
          }
        }
        
        env = {
          variables = {
            CONFIG_PATH = "/etc/app/config.yaml"
          }
          secrets = {}
        }
        
        # Volume mounts to mount ConfigMaps as files
        volume_mounts = [
          {
            name       = "app-config-volume"
            mount_path = "/etc/app/config.yaml"
            config_map = {
              name = "app-config"
              key  = "app-config"  # Key in the ConfigMap
            }
          },
          {
            name       = "routes-config-volume"
            mount_path = "/etc/app/routes.yaml"
            config_map = {
              name = "routes-config"
              key  = "routes-config"
            }
          }
        ]
      }
    }
    
    ingress = {
      is_enabled = false
      dns_domain = ""
    }
  }
}
```

**Key Points:**
- `config_maps` creates ConfigMap resources with YAML/text content
- `volume_mounts` mounts ConfigMaps as files in the container
- Each ConfigMap key becomes a file when mounted
- Content is stored under a key with the same name as the ConfigMap

---

## 11. Command and Args Override

Override container entrypoint and arguments.

```hcl
module "custom_command_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "custom-cmd"
    id   = "custom-cmd-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "custom-cmd"
    create_namespace = true
    version = "main"
    
    container = {
      app = {
        image = {
          repo = "busybox"
          tag  = "latest"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 8080
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = false
          }
        ]
        
        resources = {
          requests = {
            cpu    = "50m"
            memory = "64Mi"
          }
          limits = {
            cpu    = "100m"
            memory = "128Mi"
          }
        }
        
        env = {
          variables = {}
          secrets   = {}
        }
        
        # Override container entrypoint
        command = ["/bin/sh", "-c"]
        
        # Override container arguments
        args = ["echo 'Starting application' && exec /app/server --config /etc/config.yaml"]
      }
    }
    
    ingress = {
      is_enabled = false
      dns_domain = ""
    }
  }
}
```

**Key Points:**
- `command` overrides the container image's ENTRYPOINT
- `args` overrides the container image's CMD
- Useful for custom startup scripts or debug containers

---

## 12. Multiple Volume Types

Combine different volume types for complex applications.

```hcl
module "multi_volume_microservice" {
  source = "./path/to/microservice-kubernetes-module"

  metadata = {
    name = "multi-volume"
    id   = "multi-volume-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace = "multi-volume"
    create_namespace = true
    version = "main"
    
    config_maps = {
      "app-config" = <<-EOT
        server:
          port: 8080
          timeout: 30s
      EOT
    }
    
    container = {
      app = {
        image = {
          repo = "my-org/multi-volume-app"
          tag  = "v1.0.0"
        }
        
        ports = [
          {
            name             = "http"
            container_port   = 8080
            network_protocol = "TCP"
            app_protocol     = "http"
            service_port     = 80
            is_ingress_port  = false
          }
        ]
        
        resources = {
          requests = {
            cpu    = "200m"
            memory = "256Mi"
          }
          limits = {
            cpu    = "1"
            memory = "1Gi"
          }
        }
        
        env = {
          variables = {}
          secrets   = {}
        }
        
        volume_mounts = [
          # ConfigMap as file
          {
            name       = "config"
            mount_path = "/etc/app/config.yaml"
            config_map = {
              name = "app-config"
              key  = "app-config"
            }
          },
          # Secret volume
          {
            name       = "tls-certs"
            mount_path = "/etc/tls"
            read_only  = true
            secret = {
              name = "my-tls-secret"
            }
          },
          # EmptyDir for caching
          {
            name       = "cache"
            mount_path = "/tmp/cache"
            empty_dir = {
              medium     = "Memory"
              size_limit = "256Mi"
            }
          },
          # PVC for persistent data
          {
            name       = "data"
            mount_path = "/data"
            pvc = {
              claim_name = "data-pvc"
              read_only  = false
            }
          }
        ]
      }
    }
    
    ingress = {
      is_enabled = false
      dns_domain = ""
    }
  }
}
```

**Key Points:**
- Multiple volume types can be combined in one deployment
- `config_map`: Mount ConfigMaps as files
- `secret`: Mount Secrets as files (for TLS certs, credentials)
- `empty_dir`: Temporary storage, optionally memory-backed
- `pvc`: Persistent storage from PersistentVolumeClaim

---

## Summary

This Terraform module provides a comprehensive, production-ready solution for deploying microservices to Kubernetes:

- **Resource Management**: CPU/memory requests and limits with QoS guarantees
- **Health Monitoring**: Startup, liveness, and readiness probes
- **High Availability**: Multiple replicas, rolling updates, pod disruption budgets
- **Autoscaling**: Horizontal Pod Autoscaler based on CPU/memory
- **Networking**: Service discovery and optional Istio ingress
- **Security**: Secret management and private registry support
- **Observability**: Structured configuration for monitoring integration

For more information, see:
- [KubernetesDeployment API Documentation](../../README.md)
- [Research Documentation](../../docs/README.md) for production best practices
- [Pulumi Examples](../pulumi/examples.md) for alternative IaC patterns

