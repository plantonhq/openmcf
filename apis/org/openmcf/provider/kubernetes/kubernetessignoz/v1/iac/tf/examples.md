# SigNoz Kubernetes Terraform - Example Configurations

## Example 1: Basic Configuration (Evaluation/Development)

This example demonstrates a minimal configuration for evaluating SigNoz with a single-node ClickHouse instance.

```hcl
module "signoz_dev" {
  source = "./path/to/module"

  metadata = {
    name = "signoz-dev"
    org  = "mycompany"
    env  = "development"
  }

  spec = {
    namespace        = "signoz-dev"
    create_namespace = true

    signoz_container = {
      replicas = 1
      resources = {
        requests = {
          cpu    = "200m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "1000m"
          memory = "2Gi"
        }
      }
    }

    otel_collector_container = {
      replicas = 2
      resources = {
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
    }

    database = {
      is_external = false
      managed_database = {
        container = {
          replicas               = 1
          persistence_enabled = true
          disk_size              = "20Gi"
          resources = {
            requests = {
              cpu    = "500m"
              memory = "1Gi"
            }
            limits = {
              cpu    = "2000m"
              memory = "4Gi"
            }
          }
        }
        cluster = {
          is_enabled = false
        }
        zookeeper = {
          is_enabled = false
        }
      }
    }
  }
}
```

---

## Example 2: Production Configuration (High Availability)

This example demonstrates a production-ready deployment with distributed ClickHouse cluster, sharding, replication, and Zookeeper coordination.

```hcl
module "signoz_production" {
  source = "./path/to/module"

  metadata = {
    name = "signoz-prod"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    namespace        = "signoz-prod"
    create_namespace = true

    signoz_container = {
      replicas = 2
      resources = {
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
    }

    otel_collector_container = {
      replicas = 3
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

    database = {
      is_external = false
      managed_database = {
        container = {
          replicas               = 1
          persistence_enabled = true
          disk_size              = "100Gi"
          resources = {
            requests = {
              cpu    = "1000m"
              memory = "4Gi"
            }
            limits = {
              cpu    = "4000m"
              memory = "16Gi"
            }
          }
        }
        cluster = {
          is_enabled    = true
          shard_count   = 2
          replica_count = 2
        }
        zookeeper = {
          is_enabled = true
          container = {
            replicas = 3
            disk_size = "10Gi"
            resources = {
              requests = {
                cpu    = "200m"
                memory = "512Mi"
              }
              limits = {
                cpu    = "1000m"
                memory = "1Gi"
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

## Example 3: External ClickHouse with Password from Kubernetes Secret (Recommended)

This example demonstrates connecting SigNoz to an existing external ClickHouse instance using a password stored in a Kubernetes Secret. This is the recommended approach for production deployments.

```hcl
# First, create a Kubernetes Secret with the ClickHouse password:
# kubectl create secret generic clickhouse-credentials \
#   --from-literal=password=your-secure-password \
#   --namespace=signoz-external

module "signoz_external_clickhouse" {
  source = "./path/to/module"

  metadata = {
    name = "signoz-external"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    namespace        = "signoz-external"
    create_namespace = true

    signoz_container = {
      replicas = 2
      resources = {
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
    }

    otel_collector_container = {
      replicas = 3
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

    database = {
      is_external = true
      external_database = {
        host         = "clickhouse.database.svc.cluster.local"
        http_port    = 8123
        tcp_port     = 9000
        cluster_name = "cluster"
        is_secure    = false
        username     = "signoz"
        # Reference an existing Kubernetes Secret for the password
        password = {
          secret_ref = {
            name = "clickhouse-credentials"
            key  = "password"
          }
        }
      }
    }
  }
}
```

---

## Example 3b: External ClickHouse with Plain String Password

This example demonstrates connecting SigNoz to an external ClickHouse instance using a plain string password. This approach is suitable for development/testing but not recommended for production.

```hcl
module "signoz_external_clickhouse_string" {
  source = "./path/to/module"

  metadata = {
    name = "signoz-external"
    org  = "mycompany"
    env  = "development"
  }

  spec = {
    namespace        = "signoz-external"
    create_namespace = true

    signoz_container = {
      replicas = 2
      resources = {
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
    }

    otel_collector_container = {
      replicas = 3
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

    database = {
      is_external = true
      external_database = {
        host         = "clickhouse.database.svc.cluster.local"
        http_port    = 8123
        tcp_port     = 9000
        cluster_name = "cluster"
        is_secure    = false
        username     = "signoz"
        # Plain string password (not recommended for production)
        password = {
          string_value = var.clickhouse_password
        }
      }
    }
  }
}
```

---

## Example 4: With Ingress Configuration

This example demonstrates configuring ingress for both SigNoz UI and OpenTelemetry Collector endpoints.

```hcl
module "signoz_with_ingress" {
  source = "./path/to/module"

  metadata = {
    name = "signoz-ingress"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    namespace        = "signoz-ingress"
    create_namespace = true

    signoz_container = {
      replicas = 2
      resources = {
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
    }

    otel_collector_container = {
      replicas = 3
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

    database = {
      is_external = false
      managed_database = {
        container = {
          replicas               = 1
          persistence_enabled = true
          disk_size              = "50Gi"
          resources = {
            requests = {
              cpu    = "500m"
              memory = "2Gi"
            }
            limits = {
              cpu    = "2000m"
              memory = "8Gi"
            }
          }
        }
        cluster = {
          is_enabled = false
        }
        zookeeper = {
          is_enabled = false
        }
      }
    }

    ingress = {
      ui = {
        enabled  = true
        hostname = "signoz.example.com"
      }
      otel_collector = {
        enabled  = true
        hostname = "signoz-ingest.example.com"
      }
    }
  }
}
```

---

## Example 5: Custom Container Images

This example demonstrates using custom container images for different SigNoz components.

```hcl
module "signoz_custom_images" {
  source = "./path/to/module"

  metadata = {
    name = "signoz-custom"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    target_cluster = {
      name = "my-gke-cluster"
    }
    namespace = "signoz-custom"

    signoz_container = {
      replicas = 2
      resources = {
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
      image = {
        repo = "signoz/signoz"
        tag  = "0.40.0"
      }
    }

    otel_collector_container = {
      replicas = 2
      resources = {
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
      image = {
        repo = "signoz/signoz-otel-collector"
        tag  = "0.88.11"
      }
    }

    database = {
      is_external = false
      managed_database = {
        container = {
          replicas               = 1
          persistence_enabled = true
          disk_size              = "50Gi"
          resources = {
            requests = {
              cpu    = "500m"
              memory = "2Gi"
            }
            limits = {
              cpu    = "2000m"
              memory = "8Gi"
            }
          }
          image = {
            repo = "clickhouse/clickhouse-server"
            tag  = "23.11-alpine"
          }
        }
        cluster = {
          is_enabled = false
        }
        zookeeper = {
          is_enabled = false
        }
      }
    }
  }
}
```

---

## Example 6: Custom Helm Values

This example demonstrates providing additional customization through Helm chart values.

```hcl
module "signoz_custom_helm" {
  source = "./path/to/module"

  metadata = {
    name = "signoz-custom"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    target_cluster = {
      name = "my-gke-cluster"
    }
    namespace = "signoz-custom"

    signoz_container = {
      replicas = 2
      resources = {
        requests = {
          cpu    = "500m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2000m"
          memory = "4Gi"
        }
      }
    }

    otel_collector_container = {
      replicas = 3
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

    database = {
      is_external = false
      managed_database = {
        container = {
          replicas               = 1
          persistence_enabled = true
          disk_size              = "100Gi"
          resources = {
            requests = {
              cpu    = "1000m"
              memory = "4Gi"
            }
            limits = {
              cpu    = "4000m"
              memory = "16Gi"
            }
          }
        }
        cluster = {
          is_enabled = false
        }
        zookeeper = {
          is_enabled = false
        }
      }
    }

    helm_values = {
      # Configure data retention for traces (in hours)
      "clickhouse.ttl.traces" = "720" # 30 days
      # Configure data retention for metrics (in hours)
      "clickhouse.ttl.metrics" = "2160" # 90 days
      # Configure data retention for logs (in hours)
      "clickhouse.ttl.logs" = "168" # 7 days
      # Enable email alerting
      "signoz.env.SMTP_HOST"     = "smtp.gmail.com"
      "signoz.env.SMTP_PORT"     = "587"
      "signoz.env.SMTP_USER"     = "alerts@example.com"
      "signoz.env.SMTP_PASSWORD" = var.smtp_password
    }
  }
}
```

---

## Example 7: High-Volume Ingestion

This example is optimized for environments with high telemetry data volumes.

```hcl
module "signoz_high_volume" {
  source = "./path/to/module"

  metadata = {
    name = "signoz-high-volume"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    namespace        = "signoz-high-volume"
    create_namespace = true

    signoz_container = {
      replicas = 3
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

    otel_collector_container = {
      replicas = 6
      resources = {
        requests = {
          cpu    = "2000m"
          memory = "4Gi"
        }
        limits = {
          cpu    = "8000m"
          memory = "16Gi"
        }
      }
    }

    database = {
      is_external = false
      managed_database = {
        container = {
          replicas               = 1
          persistence_enabled = true
          disk_size              = "500Gi"
          resources = {
            requests = {
              cpu    = "4000m"
              memory = "16Gi"
            }
            limits = {
              cpu    = "16000m"
              memory = "64Gi"
            }
          }
        }
        cluster = {
          is_enabled    = true
          shard_count   = 3
          replica_count = 2
        }
        zookeeper = {
          is_enabled = true
          container = {
            replicas = 5
            disk_size = "20Gi"
            resources = {
              requests = {
                cpu    = "500m"
                memory = "1Gi"
              }
              limits = {
                cpu    = "2000m"
                memory = "4Gi"
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

## Accessing Outputs

After applying the module, access outputs using:

```hcl
output "signoz_ui_endpoint" {
  description = "SigNoz UI internal endpoint"
  value       = module.signoz.kube_endpoint
}

output "port_forward_cmd" {
  description = "Command to access SigNoz UI locally"
  value       = module.signoz.port_forward_command
}

output "otel_grpc_endpoint" {
  description = "OpenTelemetry Collector gRPC endpoint"
  value       = module.signoz.otel_collector_grpc_endpoint
}

output "external_ui_hostname" {
  description = "External hostname for SigNoz UI (if ingress enabled)"
  value       = module.signoz.external_hostname
}
```

---

## Important Notes

### Storage Configuration

Ensure your Kubernetes cluster has a default StorageClass:

```bash
kubectl get storageclass
```

If you need to specify a custom StorageClass, add it through Helm values:

```hcl
helm_values = {
  "global.storageClass" = "gp2-csi"
}
```

### Ingress Annotations

For nginx-ingress, you may need separate ingress resources for gRPC and HTTP OTel Collector endpoints:

- gRPC requires: `nginx.ingress.kubernetes.io/backend-protocol: "GRPC"`
- HTTP uses standard configuration

Configure through Helm values if needed.

### External ClickHouse Prerequisites

When using external ClickHouse:
- Ensure network connectivity from Kubernetes cluster to ClickHouse
- A Zookeeper instance must be available for distributed cluster support
- ClickHouse must have a distributed cluster configured (default name: "cluster")
- Provided credentials must have database creation permissions

### Resource Planning

Monitor actual resource usage and adjust allocations based on:
- Number of instrumented applications
- Telemetry data volume (traces/sec, metrics/sec, logs/sec)
- Query patterns and concurrency
- Configured data retention periods

