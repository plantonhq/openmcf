# Harbor Kubernetes Terraform - Examples

## Namespace Management

All Harbor deployments require a namespace. Control namespace management with the `create_namespace` flag:

- **`create_namespace = true`** - Component creates and manages the namespace with appropriate labels
- **`create_namespace = false`** - Deploy into an existing namespace (must exist before deployment)

**Managed namespace example:**
```hcl
spec = {
  namespace        = "harbor-prod"
  create_namespace = true  # Component creates the namespace
  # ... rest of spec
}
```

**Existing namespace example:**
```hcl
spec = {
  namespace        = "container-registry"
  create_namespace = false  # Must exist before deployment
  # ... rest of spec
}
```

---

## Basic Example

```hcl
module "harbor_basic" {
  source = "../"

  harbor_kubernetes = {
    metadata = {
      name = "harbor-dev"
    }
    spec = {
      target_cluster = {
        cluster_name = "my-gke-cluster"
      }
      namespace        = "harbor-dev"
      create_namespace = true
      database = {
        is_external = false
      }
      cache = {
        is_external = false
      }
      storage = {
        type = "filesystem"
        filesystem = {
          disk_size = "100Gi"
        }
      }
    }
  }
}
```

## Production Example with External Dependencies

```hcl
module "harbor_production" {
  source = "../"

  harbor_kubernetes = {
    metadata = {
      name = "harbor-prod"
    }
    spec = {
      target_cluster = {
        cluster_name = "my-eks-cluster"
      }
      namespace        = "harbor-prod"
      create_namespace = true
      database = {
        is_external = true
        external_database = {
          host     = "postgres.example.com"
          port     = 5432
          username = "harbor"
          password = var.db_password
          use_ssl  = true
        }
      }
      cache = {
        is_external = true
        external_cache = {
          host     = "redis.example.com"
          port     = 6379
          password = var.redis_password
        }
      }
      storage = {
        type = "s3"
        s3 = {
          bucket     = "my-harbor-artifacts"
          region     = "us-west-2"
          access_key = var.aws_access_key
          secret_key = var.aws_secret_key
        }
      }
      ingress = {
        core = {
          enabled  = true
          hostname = "harbor.example.com"
        }
      }
    }
  }
}
```

## Note

These examples show the variable structure. The actual Terraform module would need to be
extended with Helm release resources to fully deploy Harbor. Refer to the Pulumi module
for a complete implementation reference.

