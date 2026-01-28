# Terraform Examples for KubernetesMongodb

This document provides Terraform usage examples for deploying MongoDB on Kubernetes using the Percona Server for MongoDB Operator.

## Example 1: Basic MongoDB Deployment

```hcl
module "mongodb_basic" {
  source = "./path/to/kubernetesmongodb/v1/iac/tf"

  metadata = {
    name = "basic-mongodb"
  }

  spec = {
    target_cluster = {
      name = "my-gke-cluster"
    }
    namespace        = "basic-mongodb"
    create_namespace = true

    container = {
      replicas = 1

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

      persistence_enabled = false
      disk_size          = ""
    }

    ingress = {
      enabled  = false
      hostname = ""
    }
  }
}
```

## Example 2: MongoDB with Persistence Enabled

```hcl
module "mongodb_persistent" {
  source = "./path/to/kubernetesmongodb/v1/iac/tf"

  metadata = {
    name = "persistent-mongodb"
    org  = "my-organization"
    env  = "production"
  }

  spec = {
    target_cluster = {
      name = "my-gke-cluster"
    }
    namespace        = "persistent-mongodb"
    create_namespace = true

    container = {
      replicas = 3

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

      persistence_enabled = true
      disk_size          = "10Gi"
    }

    ingress = {
      enabled  = false
      hostname = ""
    }
  }
}
```

## Example 3: MongoDB with Ingress Enabled

```hcl
module "mongodb_with_ingress" {
  source = "./path/to/kubernetesmongodb/v1/iac/tf"

  metadata = {
    name = "ingress-mongodb"
  }

  spec = {
    target_cluster = {
      name = "my-gke-cluster"
    }
    namespace        = "ingress-mongodb"
    create_namespace = true

    container = {
      replicas = 1

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

      persistence_enabled = true
      disk_size          = "5Gi"
    }

    ingress = {
      enabled  = true
      hostname = "mongodb.example.com"
    }
  }
}
```

## Example 4: Production MongoDB Deployment

```hcl
module "mongodb_production" {
  source = "./path/to/kubernetesmongodb/v1/iac/tf"

  metadata = {
    name = "production-mongodb"
    id   = "prod-mdb-001"
    org  = "acme-corp"
    env  = "production"
    labels = {
      team    = "platform"
      project = "core-services"
    }
  }

  spec = {
    target_cluster = {
      name = "my-gke-cluster"
    }
    namespace        = "production-mongodb"
    create_namespace = true

    container = {
      replicas = 3

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

      persistence_enabled = true
      disk_size          = "50Gi"
    }

    ingress = {
      enabled  = true
      hostname = "mongodb.prod.example.com"
    }

    helm_values = {
      "backup.enabled" = "true"
      "monitoring.enabled" = "true"
    }
  }
}

# Output the MongoDB connection details
output "mongodb_namespace" {
  value = module.mongodb_production.namespace
}

output "mongodb_endpoint" {
  value = module.mongodb_production.kube_endpoint
}

output "mongodb_external_hostname" {
  value = module.mongodb_production.external_hostname
}

output "mongodb_username" {
  value = module.mongodb_production.username
}

output "mongodb_password_secret" {
  value = {
    name = module.mongodb_production.password_secret_name
    key  = module.mongodb_production.password_secret_key
  }
}
```

## Example 5: Development MongoDB with Minimal Resources

```hcl
module "mongodb_dev" {
  source = "./path/to/kubernetesmongodb/v1/iac/tf"

  metadata = {
    name = "dev-mongodb"
    env  = "development"
  }

  spec = {
    target_cluster = {
      name = "my-gke-cluster"
    }
    namespace        = "dev-mongodb"
    create_namespace = true

    container = {
      replicas = 1

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

      persistence_enabled = false
      disk_size          = ""
    }

    ingress = {
      enabled  = false
      hostname = ""
    }
  }
}
```

## Example 6: MongoDB with Custom Helm Values

```hcl
module "mongodb_custom" {
  source = "./path/to/kubernetesmongodb/v1/iac/tf"

  metadata = {
    name = "custom-mongodb"
  }

  spec = {
    target_cluster = {
      name = "my-gke-cluster"
    }
    namespace        = "custom-mongodb"
    create_namespace = true

    container = {
      replicas = 2

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

      persistence_enabled = true
      disk_size          = "20Gi"
    }

    ingress = {
      enabled  = false
      hostname = ""
    }

    helm_values = {
      "mongodbUsername" = "myuser"
      "mongodbDatabase" = "mydatabase"
    }
  }
}
```

## Output Values

The module provides the following outputs:

- `namespace` - The namespace in which MongoDB is deployed
- `service` - The base name of the MongoDB service
- `kube_endpoint` - Internal DNS name of the MongoDB service within the cluster
- `port_forward_command` - Command to port-forward traffic to the MongoDB service
- `external_hostname` - The external hostname for MongoDB if ingress is enabled
- `username` - The default MongoDB root username
- `password_secret_name` - Name of the Secret holding the MongoDB root password
- `password_secret_key` - Key within the Secret that contains the MongoDB root password

## Connecting to MongoDB

After deploying MongoDB, you can connect to it using the following methods:

### From within the Kubernetes cluster:

```bash
# Use the internal endpoint
mongodb://<username>:<password>@<kube_endpoint>:27017
```

### Using kubectl port-forward:

```bash
# Use the port-forward command from outputs
kubectl port-forward -n <namespace> service/<service-name> 8080:27017

# Then connect to localhost:8080
mongodb://<username>:<password>@localhost:8080
```

### Using external hostname (if ingress is enabled):

```bash
# Connect using the external hostname
mongodb://<username>:<password>@<external_hostname>:27017
```

## Retrieving the Password

The password is stored in a Kubernetes secret:

```bash
kubectl get secret <password_secret_name> -n <namespace> -o jsonpath='{.data.MONGODB_DATABASE_ADMIN_PASSWORD}' | base64 -d
```

## Example 7: Using Existing Namespace

This example demonstrates using an existing namespace instead of creating a new one:

```hcl
module "mongodb_existing_ns" {
  source = "./path/to/kubernetesmongodb/v1/iac/tf"

  metadata = {
    name = "existing-ns-mongodb"
  }

  spec = {
    target_cluster = {
      name = "my-gke-cluster"
    }
    namespace        = "existing-mongodb-namespace"
    create_namespace = false

    container = {
      replicas = 3

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

      persistence_enabled = true
      disk_size          = "20Gi"
    }

    ingress = {
      enabled  = false
      hostname = ""
    }
  }
}
```

**Note:** When `create_namespace = false`, the namespace must already exist before applying this module.

## Notes

- The deployment uses the Percona Server for MongoDB Operator for production-grade database management
- Replica sets are automatically configured based on the number of replicas
- For production deployments, it's recommended to use at least 3 replicas for high availability
- Persistence should be enabled for production workloads to ensure data durability
- The `unsafeFlags.replsetSize` option allows using fewer than 3 replicas for development/testing
- Use `create_namespace = false` when deploying to namespaces with pre-configured policies, quotas, or RBAC

