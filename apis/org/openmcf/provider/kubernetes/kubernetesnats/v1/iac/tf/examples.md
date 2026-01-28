# Terraform Examples for KubernetesNats

This document provides Terraform usage examples for deploying NATS on Kubernetes using the official NATS Helm chart.

## Example 1: Basic NATS Cluster with Default Settings

Deploy a simple NATS cluster with default replicas, resources, and JetStream enabled.

```hcl
module "nats_basic" {
  source = "./path/to/kubernetesnats/v1/iac/tf"

  metadata = {
    name = "nats-basic"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "nats-basic"
    create_namespace = true

    server_container = {
      replicas = 3
      resources = {
        limits = {
          cpu    = "1000m"
          memory = "2Gi"
        }
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
      }
      disk_size = "10Gi"
    }

    disable_jet_stream = false
    tls_enabled        = false
    disable_nats_box   = false
  }
}

# Output the connection details
output "nats_namespace" {
  value = module.nats_basic.namespace
}

output "nats_internal_url" {
  value = module.nats_basic.internal_client_url
}
```

## Example 2: NATS Cluster with Bearer Token Authentication

Set up a secure NATS cluster using bearer token authentication.

```hcl
module "nats_secure" {
  source = "./path/to/kubernetesnats/v1/iac/tf"

  metadata = {
    name = "nats-secure"
    org  = "my-organization"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "nats-secure"
    create_namespace = true

    server_container = {
      replicas = 5
      resources = {
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
      }
      disk_size = "20Gi"
    }

    auth = {
      enabled = true
      scheme  = "bearer_token"
    }

    tls_enabled      = true
    disable_jet_stream = false
  }
}

# Output connection and auth details
output "nats_namespace" {
  value = module.nats_secure.namespace
}

output "nats_internal_url" {
  value = module.nats_secure.internal_client_url
}

output "auth_secret_name" {
  value = "auth-nats"
}

output "auth_secret_key" {
  value = "token"
}
```

## Example 3: NATS Cluster with Basic Authentication

Deploy NATS with basic username/password authentication.

```hcl
module "nats_basic_auth" {
  source = "./path/to/kubernetesnats/v1/iac/tf"

  metadata = {
    name = "nats-basic-auth"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "nats-basic-auth"
    create_namespace = true

    server_container = {
      replicas = 3
      resources = {
        limits = {
          cpu    = "1000m"
          memory = "2Gi"
        }
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
      }
      disk_size = "10Gi"
    }

    auth = {
      enabled = true
      scheme  = "basic_auth"
    }

    tls_enabled = false
  }
}

# Output authentication details
output "nats_namespace" {
  value = module.nats_basic_auth.namespace
}

output "username" {
  value = "nats" # Default admin username
}

output "auth_secret_name" {
  value = "auth-nats"
}

output "password_key" {
  value = "password"
}
```

## Example 4: NATS Cluster with External Access via Ingress

Deploy NATS configured with ingress to allow external clients to access messaging services.

```hcl
module "nats_external" {
  source = "./path/to/kubernetesnats/v1/iac/tf"

  metadata = {
    name = "nats-external"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "nats-external"
    create_namespace = true

    server_container = {
      replicas = 3
      resources = {
        limits = {
          cpu    = "1000m"
          memory = "2Gi"
        }
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
      }
      disk_size = "10Gi"
    }

    ingress = {
      enabled  = true
      hostname = "nats.example.com"
    }

    auth = {
      enabled = true
      scheme  = "basic_auth"
    }

    tls_enabled = true
  }
}

# Output connection details
output "nats_namespace" {
  value = module.nats_external.namespace
}

output "external_hostname" {
  value = module.nats_external.external_hostname
}

output "internal_client_url" {
  value = module.nats_external.internal_client_url
}

output "external_client_url" {
  value = "nats://nats.example.com:4222"
}
```

## Example 5: Lightweight NATS Cluster Without JetStream

Set up a lightweight NATS messaging cluster without JetStream persistence.

```hcl
module "nats_minimal" {
  source = "./path/to/kubernetesnats/v1/iac/tf"

  metadata = {
    name = "nats-minimal"
    env  = "development"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "nats-minimal"
    create_namespace = true

    server_container = {
      replicas = 1
      resources = {
        limits = {
          cpu    = "500m"
          memory = "512Mi"
        }
        requests = {
          cpu    = "100m"
          memory = "128Mi"
        }
      }
      disk_size = "1Gi"
    }

    disable_jet_stream = true
    tls_enabled        = false
    disable_nats_box   = true
  }
}

# Output minimal connection details
output "nats_namespace" {
  value = module.nats_minimal.namespace
}

output "nats_internal_url" {
  value = module.nats_minimal.internal_client_url
}
```

## Example 6: High Availability NATS with Basic Auth and No-Auth User

Deploy a highly available NATS cluster with both authenticated and unauthenticated access.

```hcl
module "nats_ha" {
  source = "./path/to/kubernetesnats/v1/iac/tf"

  metadata = {
    name = "nats-ha-metrics"
    org  = "acme-corp"
    env  = "production"
    labels = {
      team    = "platform"
      project = "messaging"
    }
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "nats-ha-metrics"
    create_namespace = true

    server_container = {
      replicas = 7
      resources = {
        limits = {
          cpu    = "4000m"
          memory = "8Gi"
        }
        requests = {
          cpu    = "1000m"
          memory = "2Gi"
        }
      }
      disk_size = "50Gi"
    }

    auth = {
      enabled = true
      scheme  = "basic_auth"
      no_auth_user = {
        enabled = true
        publish_subjects = [
          "telemetry.>",
          "metrics.>",
        ]
      }
    }

    tls_enabled        = true
    disable_jet_stream = false

    ingress = {
      enabled  = true
      hostname = "nats-ha.example.com"
    }
  }
}

# Output comprehensive connection details
output "nats_namespace" {
  value = module.nats_ha.namespace
}

output "external_hostname" {
  value = module.nats_ha.external_hostname
}

output "internal_client_url" {
  value = module.nats_ha.internal_client_url
}

output "metrics_endpoint" {
  value = "http://nats-prom.${module.nats_ha.namespace}.svc.cluster.local:7777/metrics"
}
```

## Example 7: Complete Production Setup with All Features

A production-ready configuration with all security features enabled.

```hcl
module "nats_production" {
  source = "./path/to/kubernetesnats/v1/iac/tf"

  metadata = {
    name = "nats-prod"
    id   = "nats-prod-001"
    org  = "acme-corp"
    env  = "production"
    labels = {
      tier        = "platform"
      component   = "messaging"
      managed_by  = "terraform"
      cost_center = "engineering"
    }
  }

  spec = {
    target_cluster = {
      cluster_name = "prod-gke-cluster"
    }
    namespace        = "nats-prod"
    create_namespace = true

    server_container = {
      replicas = 5
      resources = {
        limits = {
          cpu    = "3000m"
          memory = "6Gi"
        }
        requests = {
          cpu    = "750m"
          memory = "1.5Gi"
        }
      }
      disk_size = "100Gi"
    }

    auth = {
      enabled = true
      scheme  = "bearer_token"
    }

    tls_enabled        = true
    disable_jet_stream = false
    disable_nats_box   = false

    ingress = {
      enabled  = true
      hostname = "nats.prod.example.com"
    }
  }
}

# Comprehensive outputs
output "nats_production_namespace" {
  description = "The namespace where NATS is deployed"
  value       = module.nats_production.namespace
}

output "nats_production_internal_url" {
  description = "Internal cluster URL for NATS clients"
  value       = module.nats_production.internal_client_url
}

output "nats_production_external_hostname" {
  description = "External hostname for NATS access"
  value       = module.nats_production.external_hostname
}
```

## Example 8: NATS with Pre-existing Namespace

Deploy NATS into an existing namespace that you manage separately.

```hcl
module "nats_existing_ns" {
  source = "./path/to/kubernetesnats/v1/iac/tf"

  metadata = {
    name = "nats-existing-ns"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "shared-messaging"
    create_namespace = false # Don't create the namespace

    server_container = {
      replicas = 3
      resources = {
        limits = {
          cpu    = "1000m"
          memory = "2Gi"
        }
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
      }
      disk_size = "10Gi"
    }

    auth = {
      enabled = true
      scheme  = "basic_auth"
    }

    tls_enabled = true
  }
}

# Output connection details
output "nats_namespace" {
  description = "The namespace where NATS is deployed (pre-existing)"
  value       = module.nats_existing_ns.namespace
}

output "nats_internal_url" {
  description = "Internal cluster URL for NATS clients"
  value       = module.nats_existing_ns.internal_client_url
}
```

**Use Case:**

* Suitable when the namespace already exists with specific ResourceQuotas, LimitRanges, or NetworkPolicies
* Useful in GitOps workflows where namespace creation is handled separately
* Allows sharing a namespace across multiple messaging components

## Retrieving Authentication Credentials

### For Bearer Token Authentication:

```bash
# Get the bearer token from the secret
kubectl get secret auth-nats -n <namespace> -o jsonpath='{.data.token}' | base64 -d
```

### For Basic Authentication:

```bash
# Get username
kubectl get secret auth-nats -n <namespace> -o jsonpath='{.data.user}' | base64 -d

# Get password
kubectl get secret auth-nats -n <namespace> -o jsonpath='{.data.password}' | base64 -d
```

### Using Terraform Data Sources:

```hcl
data "kubernetes_secret_v1" "nats_auth" {
  metadata {
    name      = "auth-nats"
    namespace = module.nats_basic_auth.namespace
  }

  depends_on = [module.nats_basic_auth]
}

# Output the credentials (sensitive)
output "nats_username" {
  value     = data.kubernetes_secret_v1.nats_auth.data["user"]
  sensitive = true
}

output "nats_password" {
  value     = data.kubernetes_secret_v1.nats_auth.data["password"]
  sensitive = true
}
```

## Connecting to NATS

### From within the Kubernetes cluster:

```bash
# Use the internal client URL
nats://nats-basic-nats.nats-basic.svc.cluster.local:4222
```

### Using kubectl port-forward:

```bash
# Port-forward to NATS service
kubectl port-forward -n <namespace> service/<name>-nats 4222:4222

# Then connect to localhost:4222
nats://localhost:4222
```

### Using external hostname (if ingress is enabled):

```bash
# Connect using the external hostname
nats://<external-hostname>:4222
```

## Example 9: NATS with NACK Controller and JetStream Streams

Deploy NATS with the NACK (NATS Controllers for Kubernetes) controller for declarative JetStream stream management.

```hcl
module "nats_with_streams" {
  source = "./path/to/kubernetesnats/v1/iac/tf"

  metadata = {
    name = "nats-with-streams"
    org  = "my-org"
    env  = "development"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "nats-streams"
    create_namespace = true

    server_container = {
      replicas = 3
      resources = {
        limits = {
          cpu    = "1000m"
          memory = "2Gi"
        }
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
      }
      disk_size = "10Gi"
    }

    auth = {
      enabled = true
      scheme  = "basic_auth"
    }

    # Enable NACK controller for declarative stream management
    nack_controller = {
      enabled             = true
      enable_control_loop = true  # Required for KeyValue/ObjectStore support
    }

    # Define JetStream streams
    streams = [
      {
        name     = "orders"
        subjects = ["orders.*", "orders.>"]
        storage  = "file"
        replicas = 3
        retention = "limits"
        max_age  = "7d"
        max_bytes = 1073741824  # 1GB
        consumers = [
          {
            durable_name   = "orders-processor"
            deliver_policy = "all"
            ack_policy     = "explicit"
            max_ack_pending = 1000
            ack_wait       = "30s"
          }
        ]
      },
      {
        name     = "events"
        subjects = ["events.>"]
        storage  = "memory"
        replicas = 1
        retention = "interest"
      }
    ]
  }
}

# Output stream information
output "nats_namespace" {
  value = module.nats_with_streams.namespace
}

output "nats_internal_url" {
  value = module.nats_with_streams.internal_client_url
}

output "nack_enabled" {
  value = module.nats_with_streams.nack_controller_enabled
}

output "streams_created" {
  value = module.nats_with_streams.streams_created
}
```

## Example 10: Production NATS with Multiple Streams and Consumers

A comprehensive production configuration with multiple streams, consumers, and external access.

```hcl
module "nats_production" {
  source = "./path/to/kubernetesnats/v1/iac/tf"

  metadata = {
    name = "nats-prod"
    org  = "acme-corp"
    env  = "production"
    labels = {
      team        = "platform"
      cost_center = "engineering"
    }
  }

  spec = {
    target_cluster = {
      cluster_name = "prod-gke-cluster"
    }
    namespace        = "nats-prod"
    create_namespace = true

    server_container = {
      replicas = 5
      resources = {
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
      }
      disk_size = "50Gi"
    }

    auth = {
      enabled = true
      scheme  = "basic_auth"
    }

    tls_enabled = true

    ingress = {
      enabled  = true
      hostname = "nats.prod.example.com"
    }

    # NACK controller with control-loop for reliable state enforcement
    nack_controller = {
      enabled             = true
      enable_control_loop = true
      helm_chart_version  = "0.31.1"
      app_version         = "0.21.1"
    }

    # Production streams configuration
    streams = [
      # High-throughput API events stream
      {
        name         = "api-events"
        subjects     = ["api.events.>"]
        storage      = "file"
        replicas     = 3
        retention    = "limits"
        max_age      = "24h"
        max_bytes    = 10737418240  # 10GB
        max_msgs     = 10000000
        discard      = "old"
        description  = "API event stream for real-time processing"
        consumers = [
          {
            durable_name    = "api-processor"
            deliver_policy  = "all"
            ack_policy      = "explicit"
            max_ack_pending = 5000
            max_deliver     = 5
            ack_wait        = "60s"
            replay_policy   = "instant"
            description     = "Main API event processor"
          },
          {
            durable_name   = "api-analytics"
            deliver_policy = "all"
            ack_policy     = "explicit"
            filter_subject = "api.events.orders.*"
            description    = "Analytics consumer for order events"
          }
        ]
      },
      # Work queue for background jobs
      {
        name        = "background-jobs"
        subjects    = ["jobs.>"]
        storage     = "file"
        replicas    = 3
        retention   = "workqueue"
        max_age     = "1h"
        description = "Work queue for background job processing"
        consumers = [
          {
            durable_name   = "job-worker"
            ack_policy     = "explicit"
            max_ack_pending = 100
            max_deliver    = 3
            ack_wait       = "5m"
          }
        ]
      },
      # Ephemeral notifications stream
      {
        name        = "notifications"
        subjects    = ["notify.>"]
        storage     = "memory"
        replicas    = 1
        retention   = "interest"
        max_age     = "5m"
        description = "Real-time notifications (ephemeral)"
      }
    ]

    nats_helm_chart_version = "2.12.3"
  }
}

# Comprehensive outputs
output "nats_production_namespace" {
  description = "The namespace where NATS is deployed"
  value       = module.nats_production.namespace
}

output "nats_production_internal_url" {
  description = "Internal cluster URL for NATS clients"
  value       = module.nats_production.internal_client_url
}

output "nats_production_external" {
  description = "External hostname for NATS access"
  value       = module.nats_production.external_hostname
}

output "nats_production_streams" {
  description = "List of JetStream streams created"
  value       = module.nats_production.streams_created
}

output "nack_controller_version" {
  description = "NACK controller version"
  value       = module.nats_production.nack_controller_version
}
```

## Example 11: Minimal NATS with Single Stream

A minimal setup with NACK for simple streaming use cases.

```hcl
module "nats_minimal_stream" {
  source = "./path/to/kubernetesnats/v1/iac/tf"

  metadata = {
    name = "nats-minimal"
    env  = "dev"
  }

  spec = {
    namespace        = "nats-dev"
    create_namespace = true

    server_container = {
      replicas = 1
      resources = {
        limits = {
          cpu    = "500m"
          memory = "512Mi"
        }
        requests = {
          cpu    = "100m"
          memory = "128Mi"
        }
      }
      disk_size = "5Gi"
    }

    # Enable NACK for stream management
    nack_controller = {
      enabled = true
    }

    # Single stream for development
    streams = [
      {
        name     = "dev-events"
        subjects = ["dev.>"]
        storage  = "memory"
        replicas = 1
        max_age  = "1h"
      }
    ]
  }
}

output "dev_nats_url" {
  value = module.nats_minimal_stream.internal_client_url
}

output "dev_streams" {
  value = module.nats_minimal_stream.streams_created
}
```

## Verifying NACK Stream Creation

After applying the Terraform configuration, you can verify the streams were created:

```bash
# Check NACK controller is running
kubectl get pods -n <namespace> -l app.kubernetes.io/name=nack

# List Stream custom resources
kubectl get streams -n <namespace>

# List Consumer custom resources
kubectl get consumers -n <namespace>

# Describe a specific stream
kubectl describe stream <stream-name> -n <namespace>

# Check stream status using nats-box
kubectl exec -it <nats-box-pod> -n <namespace> -- nats stream list
kubectl exec -it <nats-box-pod> -n <namespace> -- nats stream info <stream-name>
```

## NACK Controller Architecture

The deployment follows this order to avoid race conditions:

1. **NATS Helm Chart** - Deploys the NATS server with JetStream enabled
2. **NACK CRDs** - Installs JetStream CRDs (Stream, Consumer, KeyValue, ObjectStore)
3. **NACK Controller** - Deploys the controller that watches CRDs
4. **Stream/Consumer CRs** - Creates the actual stream and consumer resources

```
┌─────────────────────────────────────────────────────────────────┐
│                      Kubernetes Cluster                          │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │
│  │ NATS Server  │◄───│    NACK      │◄───│   Stream/    │       │
│  │ (JetStream)  │    │  Controller  │    │  Consumer    │       │
│  └──────────────┘    └──────────────┘    │    CRs       │       │
│         ▲                   │            └──────────────┘       │
│         │                   │                   ▲               │
│         │            ┌──────▼──────┐           │               │
│         └────────────│ Reconcile   │───────────┘               │
│                      │ JetStream   │                            │
│                      │ Resources   │                            │
│                      └─────────────┘                            │
└─────────────────────────────────────────────────────────────────┘
```

## Notes

- The deployment uses the official NATS Helm chart (version 2.12.3) for production-grade message streaming
- JetStream provides persistent streaming capabilities with file-based storage backed by PersistentVolumeClaims
- For production deployments, use an odd number of replicas (3, 5, 7) for optimal quorum in clustered mode
- TLS encryption is automatically configured with self-signed certificates when `tls_enabled` is true
- The `nats-box` utility pod is deployed by default for debugging and testing (can be disabled)
- Metrics are exposed at `http://nats-prom.<namespace>.svc.cluster.local:7777/metrics` for Prometheus scraping
- When using basic auth with no-auth user, the no-auth user can only publish to specified subjects
- The module automatically handles secret generation for authentication credentials
- **NACK Controller**: Enables declarative stream/consumer management via Kubernetes CRDs
- **Control-Loop Mode**: Required for KeyValue, ObjectStore, and Account support; provides more reliable state enforcement
- **CRD Installation**: CRDs are installed separately from the Helm chart for better control and to avoid preview/dry-run issues
- **Stream Replication**: Use odd numbers (1, 3, 5) for stream replicas to maintain quorum

