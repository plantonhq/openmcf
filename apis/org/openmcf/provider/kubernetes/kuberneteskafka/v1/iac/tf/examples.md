# Kafka Kubernetes - Terraform Examples

This document provides Terraform-specific examples for deploying Apache Kafka on Kubernetes using the KafkaKubernetes module.

## Prerequisites

- Terraform >= 1.0
- Kubernetes cluster with Strimzi Operator installed
- kubectl configured
- Sufficient cluster resources

## Table of Contents

1. [Basic Kafka Cluster](#example-1-basic-kafka-cluster)
2. [Kafka with Schema Registry and UI](#example-2-kafka-with-schema-registry-and-ui)
3. [Minimal Development Setup](#example-3-minimal-development-setup)
4. [Custom Topic Configuration](#example-4-custom-topic-configuration)
5. [Schema Registry without Kafka UI](#example-5-schema-registry-without-kafka-ui)
6. [Production High-Availability Cluster](#example-6-production-high-availability-cluster)
7. [External Namespace Management](#example-7-external-namespace-management)

---

## Example 1: Basic Kafka Cluster

Basic Kafka deployment with minimal configuration for development/testing.

```hcl
module "kafka_basic" {
  source = "path/to/module"

  metadata = {
    name = "kafka-cluster-basic"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "kafka-basic"
    create_namespace = true

    kafka_topics = [
      {
        name       = "my-topic"
        partitions = 3
        replicas   = 2
      }
    ]

    broker_container = {
      replicas = 1
      resources = {
        requests = {
          cpu    = "100m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      disk_size = "20Gi"
    }

    zookeeper_container = {
      replicas = 3
      resources = {
        requests = {
          cpu    = "50m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "500m"
          memory = "512Mi"
        }
      }
      disk_size = "10Gi"
    }

    ingress = {
      enabled = false
    }

    is_deploy_kafka_ui = false
  }
}

output "kafka_namespace" {
  value = module.kafka_basic.namespace
}

output "kafka_bootstrap_servers" {
  value = module.kafka_basic.internal_bootstrap_endpoint
}
```

**Use Case:** Development and testing environments requiring basic Kafka functionality.

---

## Example 2: Kafka with Schema Registry and UI

Full-featured Kafka deployment with Schema Registry and Kafka UI for production use.

```hcl
module "kafka_full" {
  source = "path/to/module"

  metadata = {
    name = "kafka-cluster-with-schema-registry"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "kafka-prod"
    create_namespace = true

    kafka_topics = [
      {
        name       = "my-topic"
        partitions = 3
        replicas   = 3
      }
    ]

    broker_container = {
      replicas = 2
      resources = {
        requests = {
          cpu    = "200m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2"
          memory = "2Gi"
        }
      }
      disk_size = "50Gi"
    }

    zookeeper_container = {
      replicas = 3
      resources = {
        requests = {
          cpu    = "100m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      disk_size = "20Gi"
    }

    schema_registry_container = {
      is_enabled = true
      replicas   = 1
      resources = {
        requests = {
          cpu    = "50m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "500m"
          memory = "512Mi"
        }
      }
    }

    ingress = {
      enabled            = true
      ingress_class_name = "nginx"
      hosts = [
        {
          host = "kafka.mydomain.com"
          paths = ["/"]
        }
      ]
    }

    is_deploy_kafka_ui = true
  }
}

output "schema_registry_url" {
  value = module.kafka_full.schema_registry_endpoint
}

output "kafka_ui_url" {
  value = "https://kafka.mydomain.com"
}
```

**Use Case:** Production deployments requiring schema management and operational visibility.

---

## Example 3: Minimal Development Setup

Absolute minimal configuration for local development or testing.

```hcl
module "kafka_minimal" {
  source = "path/to/module"

  metadata = {
    name = "kafka-minimal"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "kafka-minimal"
    create_namespace = true

    broker_container = {
      replicas = 1
      resources = {
        requests = {
          cpu    = "50m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "500m"
          memory = "512Mi"
        }
      }
      disk_size = "10Gi"
    }

    zookeeper_container = {
      replicas = 1
      resources = {
        requests = {
          cpu    = "50m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "500m"
          memory = "512Mi"
        }
      }
      disk_size = "10Gi"
    }

    ingress = {
      enabled = false
    }

    is_deploy_kafka_ui = false
  }
}
```

**Use Case:** Local development, proof-of-concept, or resource-constrained environments.

**Note:** Single broker and ZooKeeper replica is not recommended for production.

---

## Example 4: Custom Topic Configuration

Kafka cluster with multiple topics having different retention and compaction policies.

```hcl
module "kafka_custom_topics" {
  source = "path/to/module"

  metadata = {
    name = "kafka-custom-topics"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "kafka-topics"
    create_namespace = true

    kafka_topics = [
      {
        name       = "logs"
        partitions = 5
        replicas   = 3
        config = {
          "retention.ms"    = "86400000"  # 24 hours
          "cleanup.policy"  = "delete"
        }
      },
      {
        name       = "metrics"
        partitions = 10
        replicas   = 2
        config = {
          "retention.ms"    = "3600000"   # 1 hour
          "cleanup.policy"  = "compact"
        }
      }
    ]

    broker_container = {
      replicas = 3
      resources = {
        requests = {
          cpu    = "500m"
          memory = "2Gi"
        }
        limits = {
          cpu    = "4"
          memory = "4Gi"
        }
      }
      disk_size = "100Gi"
    }

    zookeeper_container = {
      replicas = 3
      resources = {
        requests = {
          cpu    = "200m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2"
          memory = "2Gi"
        }
      }
      disk_size = "20Gi"
    }

    ingress = {
      enabled = true
    }

    is_deploy_kafka_ui = true
  }
}

# Output topic information
output "kafka_topics" {
  value = [
    for topic in module.kafka_custom_topics.kafka_topics : {
      name       = topic.name
      partitions = topic.partitions
      replicas   = topic.replicas
    }
  ]
}
```

**Use Case:** Multi-purpose clusters with different retention policies for various data types.

---

## Example 5: Schema Registry without Kafka UI

Kafka with Schema Registry but without the UI component for automated/API-only access.

```hcl
module "kafka_schema_only" {
  source = "path/to/module"

  metadata = {
    name = "kafka-with-schema-registry"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "kafka-schema"
    create_namespace = true

    kafka_topics = [
      {
        name       = "transactions"
        partitions = 3
        replicas   = 2
      }
    ]

    broker_container = {
      replicas = 2
      resources = {
        requests = {
          cpu    = "200m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2"
          memory = "2Gi"
        }
      }
      disk_size = "50Gi"
    }

    zookeeper_container = {
      replicas = 3
      resources = {
        requests = {
          cpu    = "100m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      disk_size = "20Gi"
    }

    schema_registry_container = {
      is_enabled = true
      replicas   = 2
      resources = {
        requests = {
          cpu    = "100m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
    }

    ingress = {
      enabled = true
    }

    is_deploy_kafka_ui = false
  }
}

# Access Schema Registry programmatically
output "schema_registry_internal_url" {
  value       = module.kafka_schema_only.schema_registry_endpoint
  description = "Internal URL for Schema Registry API access"
}
```

**Use Case:** Automated data pipelines requiring schema validation without manual cluster management.

---

## Example 6: Production High-Availability Cluster

Enterprise-grade Kafka cluster with full redundancy and monitoring.

```hcl
module "kafka_production" {
  source = "path/to/module"

  metadata = {
    name = "kafka-prod-ha"
    org  = "my-company"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-eks-cluster"
    }
    namespace        = "kafka-production"
    create_namespace = true

    kafka_topics = [
      {
        name       = "critical-events"
        partitions = 12
        replicas   = 3
        config = {
          "min.insync.replicas" = "2"
          "retention.ms"        = "604800000"  # 7 days
          "cleanup.policy"      = "delete"
        }
      },
      {
        name       = "user-activity"
        partitions = 20
        replicas   = 3
        config = {
          "min.insync.replicas" = "2"
          "retention.ms"        = "259200000"  # 3 days
        }
      }
    ]

    broker_container = {
      replicas = 5
      resources = {
        requests = {
          cpu    = "1"
          memory = "4Gi"
        }
        limits = {
          cpu    = "8"
          memory = "16Gi"
        }
      }
      disk_size = "500Gi"
    }

    zookeeper_container = {
      replicas = 5
      resources = {
        requests = {
          cpu    = "500m"
          memory = "2Gi"
        }
        limits = {
          cpu    = "4"
          memory = "8Gi"
        }
      }
      disk_size = "100Gi"
    }

    schema_registry_container = {
      is_enabled = true
      replicas   = 3
      resources = {
        requests = {
          cpu    = "200m"
          memory = "1Gi"
        }
        limits = {
          cpu    = "2"
          memory = "4Gi"
        }
      }
    }

    ingress = {
      enabled            = true
      ingress_class_name = "nginx"
      hosts = [
        {
          host = "kafka-prod.company.com"
          paths = ["/"]
        }
      ]
    }

    is_deploy_kafka_ui = true
  }
}

# Production outputs
output "prod_kafka_bootstrap" {
  value       = module.kafka_production.internal_bootstrap_endpoint
  description = "Production Kafka bootstrap servers"
  sensitive   = false
}

output "prod_schema_registry" {
  value       = module.kafka_production.schema_registry_endpoint
  description = "Production Schema Registry endpoint"
}

output "prod_admin_credentials" {
  value       = module.kafka_production.admin_username
  description = "Admin username for Kafka cluster"
}
```

**Use Case:** Mission-critical production workloads requiring high availability and data durability.

**Production Considerations:**
- 5+ broker replicas for resilience
- 5+ ZooKeeper nodes for quorum
- Multiple Schema Registry instances
- `min.insync.replicas` set for data safety
- Sufficient disk space for retention policies
- Monitoring via Kafka UI enabled

---

## Common Patterns

### Accessing Kafka from Applications

After deploying Kafka, connect from your applications:

```hcl
# In your application deployment
resource "kubernetes_deployment" "app" {
  # ... other config ...

  spec {
    template {
      spec {
        container {
          name  = "app"
          image = "myapp:latest"

          env {
            name  = "KAFKA_BOOTSTRAP_SERVERS"
            value = module.kafka_basic.internal_bootstrap_endpoint
          }

          env {
            name = "KAFKA_USERNAME"
            value_from {
              secret_key_ref {
                name = "${module.kafka_basic.admin_username}-secret"
                key  = "username"
              }
            }
          }

          env {
            name = "KAFKA_PASSWORD"
            value_from {
              secret_key_ref {
                name = "${module.kafka_basic.admin_username}-secret"
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

### Monitoring and Observability

Enable monitoring by accessing Kafka UI:

```hcl
# Port-forward to Kafka UI for local access
# kubectl port-forward -n kafka-namespace svc/kafka-ui 8080:80
```

---

## Verification

After deployment, verify your Kafka cluster:

```bash
# Check namespace and pods
kubectl get pods -n kafka-namespace

# Check Kafka cluster status
kubectl get kafka -n kafka-namespace

# Check topics
kubectl get kafkatopics -n kafka-namespace

# View admin credentials
kubectl get secret kafka-admin -n kafka-namespace -o jsonpath='{.data.username}' | base64 -d
```

---

## Troubleshooting

### Pods Not Starting

```bash
# Check pod status
kubectl describe pod -n kafka-namespace <pod-name>

# Check events
kubectl get events -n kafka-namespace --sort-by='.lastTimestamp'
```

### Topic Creation Issues

```bash
# Check topic status
kubectl describe kafkatopic <topic-name> -n kafka-namespace

# Check entity operator logs
kubectl logs -n kafka-namespace -l strimzi.io/name=<cluster-name>-entity-operator
```

### Storage Issues

Ensure your cluster has a default storage class:

```bash
kubectl get storageclass
```

---

## Best Practices

1. **Start Small**: Begin with minimal resources and scale based on metrics
2. **Use Odd Numbers**: Always use odd numbers for ZooKeeper replicas (3, 5, 7)
3. **Set Replica Count**: Match topic replicas to broker count (max)
4. **Monitor Disk**: Set appropriate retention policies to avoid disk exhaustion
5. **Enable UI**: Use Kafka UI for operational visibility in non-production
6. **Use Schema Registry**: For production data pipelines requiring schema evolution
7. **Set Resource Limits**: Prevent resource contention with appropriate limits
8. **Regular Backups**: Implement backup strategy for critical topics

---

## Example 7: External Namespace Management

Use an externally managed namespace (e.g., created separately via KubernetesNamespace component).

```hcl
module "kafka_external_ns" {
  source = "path/to/module"

  metadata = {
    name = "kafka-external-namespace"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "existing-kafka-namespace"
    create_namespace = false  # Namespace managed externally

    kafka_topics = [
      {
        name       = "my-topic"
        partitions = 3
        replicas   = 2
      }
    ]

    broker_container = {
      replicas = 1
      resources = {
        requests = {
          cpu    = "100m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      disk_size = "20Gi"
    }

    zookeeper_container = {
      replicas = 3
      resources = {
        requests = {
          cpu    = "50m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "500m"
          memory = "512Mi"
        }
      }
      disk_size = "10Gi"
    }

    ingress = {
      enabled = false
    }

    is_deploy_kafka_ui = false
  }
}
```

**Use Case:** Environments where namespace management is centralized or governed by policies.

**Prerequisites:**
- The namespace `existing-kafka-namespace` must already exist in the cluster
- Ensure the service account has appropriate RBAC permissions in that namespace

---

## Additional Resources

- [Strimzi Documentation](https://strimzi.io/documentation/)
- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [Terraform Kubernetes Provider](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs)
- [Module README](README.md)
- [Research Documentation](../docs/README.md)

