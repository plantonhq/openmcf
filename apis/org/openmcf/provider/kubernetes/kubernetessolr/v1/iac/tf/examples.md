# Terraform Examples for KubernetesSolr

This document provides practical examples for deploying Apache Solr on Kubernetes using Terraform. Each example demonstrates different configuration scenarios and can be adapted to your specific requirements.

## Prerequisites

Before running these examples:

1. **Kubernetes Cluster**: Access to a Kubernetes cluster with sufficient resources
2. **Solr Operator**: The Apache Solr Operator must be installed in the cluster
3. **cert-manager** (for ingress): Required if enabling external access
4. **Istio** (for ingress): Required for Gateway API support
5. **Terraform**: Version 1.0 or later

## Example 1: Basic Solr Deployment with Default Settings

This example demonstrates the most basic configuration for deploying a Solr instance. It configures a single Solr pod with default resource allocations and storage size.

### Configuration (`basic-solr.tfvars`)

```hcl
metadata = {
  name = "solr-instance-basic"
}

spec = {
  target_cluster = {
    name = "my-gke-cluster"
  }
  namespace        = "solr-instance-basic"
  create_namespace = true

  solr_container = {
    replicas = 1
    image = {
      repo = "solr"
      tag  = "8.7.0"
    }
    resources = {
      requests = {
        cpu    = "50m"
        memory = "256Mi"
      }
      limits = {
        cpu    = "1"
        memory = "1Gi"
      }
    }
    disk_size = "1Gi"
  }

  config = {}

  zookeeper_container = {
    replicas = 1
    resources = {
      requests = {
        cpu    = "50m"
        memory = "256Mi"
      }
      limits = {
        cpu    = "1"
        memory = "1Gi"
      }
    }
    disk_size = "1Gi"
  }

  ingress = {
    is_enabled = false
    dns_domain = ""
  }
}
```

### Deploy

```bash
cd iac/tf
terraform init
terraform plan -var-file="basic-solr.tfvars"
terraform apply -var-file="basic-solr.tfvars"
```

### Access Solr

```bash
# Port-forward to access Solr UI locally
kubectl port-forward -n solr-instance-basic service/solr-instance-basic-solrcloud-common 8080:8080

# Open browser to http://localhost:8080
```

---

## Example 2: Solr Deployment with Custom JVM and Persistent Storage

This example configures a Solr deployment with custom JVM memory settings, larger persistent volumes, and increased replica count for high availability.

### Configuration (`custom-jvm-solr.tfvars`)

```hcl
metadata = {
  name = "solr-instance-custom"
  id   = "solr-prod"
  org  = "myorg"
  env  = "production"
}

spec = {
  target_cluster = {
    name = "my-gke-cluster"
  }
  namespace        = "solr-instance-custom"
  create_namespace = true

  solr_container = {
    replicas = 3
    image = {
      repo = "solr"
      tag  = "8.7.0"
    }
    resources = {
      requests = {
        cpu    = "100m"
        memory = "512Mi"
      }
      limits = {
        cpu    = "2"
        memory = "2Gi"
      }
    }
    disk_size = "5Gi"
  }

  config = {
    java_mem                  = "-Xms2g -Xmx4g"
    opts                      = "-Dsolr.autoSoftCommit.maxTime=5000"
    garbage_collection_tuning = "-XX:SurvivorRatio=6 -XX:MaxTenuringThreshold=10"
  }

  zookeeper_container = {
    replicas = 3
    resources = {
      requests = {
        cpu    = "100m"
        memory = "512Mi"
      }
      limits = {
        cpu    = "2"
        memory = "2Gi"
      }
    }
    disk_size = "5Gi"
  }

  ingress = {
    is_enabled = false
    dns_domain = ""
  }
}
```

### Deploy

```bash
cd iac/tf
terraform init
terraform plan -var-file="custom-jvm-solr.tfvars"
terraform apply -var-file="custom-jvm-solr.tfvars"
```

### Verify Deployment

```bash
# Check pod status
kubectl get pods -n solr-prod

# Expected output: 3 Solr pods and 3 Zookeeper pods running
```

---

## Example 3: Solr Deployment with Ingress Enabled

This example demonstrates how to deploy Solr with external ingress enabled. This allows access to Solr through a public hostname with TLS encryption.

### Configuration (`ingress-solr.tfvars`)

```hcl
metadata = {
  name = "solr-instance-ingress"
  id   = "solr-web"
}

spec = {
  target_cluster = {
    name = "my-gke-cluster"
  }
  namespace        = "solr-instance-ingress"
  create_namespace = true

  solr_container = {
    replicas = 2
    image = {
      repo = "solr"
      tag  = "8.7.0"
    }
    resources = {
      requests = {
        cpu    = "50m"
        memory = "256Mi"
      }
      limits = {
        cpu    = "1"
        memory = "1Gi"
      }
    }
    disk_size = "2Gi"
  }

  config = {}

  zookeeper_container = {
    replicas = 2
    resources = {
      requests = {
        cpu    = "50m"
        memory = "256Mi"
      }
      limits = {
        cpu    = "1"
        memory = "1Gi"
      }
    }
    disk_size = "2Gi"
  }

  ingress = {
    is_enabled = true
    dns_domain = "example.com"
  }
}
```

### Deploy

```bash
cd iac/tf
terraform init
terraform plan -var-file="ingress-solr.tfvars"
terraform apply -var-file="ingress-solr.tfvars"
```

### Access Solr via Ingress

```bash
# Wait for certificate to be issued
kubectl get certificate -n istio-ingress cert-solr-web

# Access Solr at the external hostname
# https://solr-web.example.com
```

### Ingress Architecture

When ingress is enabled, this module creates:

1. **TLS Certificate**: Automatically provisioned via cert-manager using the DNS domain as the cluster issuer
2. **Gateway**: Istio Gateway resource listening on ports 80 (HTTP) and 443 (HTTPS)
3. **HTTP Redirect**: HTTPRoute that redirects HTTP (port 80) to HTTPS (port 443)
4. **HTTPS Route**: HTTPRoute that routes HTTPS traffic to the Solr service
5. **Hostnames**:
   - External: `<resource-id>.<dns-domain>` (e.g., `solr-web.example.com`)
   - Internal: `<resource-id>-internal.<dns-domain>` (e.g., `solr-web-internal.example.com`)

---

## Example 4: Solr Deployment with Custom Garbage Collection (No Ingress)

This configuration deploys a Solr cluster with custom garbage collection tuning but without external ingress, relying on Kubernetes internal networking and port-forwarding for access.

### Configuration (`gc-tuning-solr.tfvars`)

```hcl
metadata = {
  name = "solr-instance-gc-tuned"
}

spec = {
  target_cluster = {
    name = "my-gke-cluster"
  }
  namespace        = "solr-instance-gc-tuned"
  create_namespace = true

  solr_container = {
    replicas = 1
    image = {
      repo = "solr"
      tag  = "8.7.0"
    }
    resources = {
      requests = {
        cpu    = "50m"
        memory = "256Mi"
      }
      limits = {
        cpu    = "1"
        memory = "1Gi"
      }
    }
    disk_size = "1Gi"
  }

  config = {
    garbage_collection_tuning = "-XX:SurvivorRatio=4 -XX:TargetSurvivorRatio=85 -XX:MaxTenuringThreshold=6"
  }

  zookeeper_container = {
    replicas = 1
    resources = {
      requests = {
        cpu    = "50m"
        memory = "256Mi"
      }
      limits = {
        cpu    = "1"
        memory = "1Gi"
      }
    }
    disk_size = "1Gi"
  }

  ingress = {
    is_enabled = false
    dns_domain = ""
  }
}
```

### Deploy

```bash
cd iac/tf
terraform init
terraform plan -var-file="gc-tuning-solr.tfvars"
terraform apply -var-file="gc-tuning-solr.tfvars"
```

---

## Example 5: Using Existing Namespace

This example shows how to deploy Solr into an existing namespace that's managed separately. This is useful when multiple components share a namespace or when namespace management is centralized.

### Configuration (`existing-namespace-solr.tfvars`)

```hcl
metadata = {
  name = "solr-shared-namespace"
}

spec = {
  target_cluster = {
    name = "my-gke-cluster"
  }
  namespace        = "shared-services"
  create_namespace = false  # Don't create namespace, use existing one

  solr_container = {
    replicas = 1
    image = {
      repo = "solr"
      tag  = "8.7.0"
    }
    resources = {
      requests = {
        cpu    = "50m"
        memory = "256Mi"
      }
      limits = {
        cpu    = "1"
        memory = "1Gi"
      }
    }
    disk_size = "1Gi"
  }

  config = {}

  zookeeper_container = {
    replicas = 1
    resources = {
      requests = {
        cpu    = "50m"
        memory = "256Mi"
      }
      limits = {
        cpu    = "1"
        memory = "1Gi"
      }
    }
    disk_size = "1Gi"
  }

  ingress = {
    is_enabled = false
    dns_domain = ""
  }
}
```

### Prerequisites

Ensure the namespace exists before deployment:

```bash
# Create the namespace if it doesn't exist
kubectl create namespace shared-services

# Or verify it exists
kubectl get namespace shared-services
```

### Deploy

```bash
cd iac/tf
terraform init
terraform plan -var-file="existing-namespace-solr.tfvars"
terraform apply -var-file="existing-namespace-solr.tfvars"
```

---

## Namespace Management

The Solr Kubernetes Terraform module provides flexible namespace management through the `create_namespace` variable:

### create_namespace = true (default)

When set to `true`, the module:
- Creates the namespace with proper labels
- Manages the namespace lifecycle
- Automatically destroys the namespace when running `terraform destroy`

**Use when:**
- Deploying a dedicated Solr instance with its own namespace
- You want Terraform to fully manage the namespace lifecycle
- No other components share the namespace

### create_namespace = false

When set to `false`, the module:
- Uses an existing namespace
- Does not create or manage the namespace
- Requires the namespace to exist before deployment

**Use when:**
- Multiple components share a namespace
- Namespace is managed by a separate Terraform module or process
- Organization policies require centralized namespace management
- Using a GitOps approach with namespace management in a separate layer

**Important**: When using `create_namespace = false`, ensure the namespace exists before running `terraform apply`, otherwise the deployment will fail.

---

## Common Operations

### Scaling Solr Replicas

Update the `solr_container.replicas` value in your `.tfvars` file and re-apply:

```hcl
spec = {
  target_cluster = {
    name = "my-gke-cluster"
  }
  namespace        = "your-namespace"
  create_namespace = true

  solr_container = {
    replicas = 5  # Scale from 3 to 5
    # ... rest of config
  }
}
```

```bash
terraform apply -var-file="your-config.tfvars"
```

### Updating Solr Image Version

Modify the `image.tag` in your configuration:

```hcl
spec = {
  target_cluster = {
    name = "my-gke-cluster"
  }
  namespace        = "your-namespace"
  create_namespace = true

  solr_container = {
    image = {
      repo = "solr"
      tag  = "8.11.0"  # Upgrade from 8.7.0
    }
    # ... rest of config
  }
}
```

```bash
terraform apply -var-file="your-config.tfvars"
```

**Note**: Test version upgrades in a non-production environment first.

### Destroying Resources

To remove all Solr resources:

```bash
terraform destroy -var-file="your-config.tfvars"
```

**Warning**: This will delete all persistent volumes and data. Back up your Solr indices before destroying.

---

## Terraform Outputs

After applying any of these configurations, Terraform provides useful outputs:

```bash
# View all outputs
terraform output

# View specific output
terraform output namespace
terraform output kube_service_name
terraform output kube_port_forward_command
```

### Available Outputs

- `namespace`: Kubernetes namespace where Solr is deployed
- `kube_service_name`: Internal Kubernetes service name for Solr
- `kube_service_fqdn`: Fully qualified domain name for internal service access
- `kube_port_forward_command`: Ready-to-use kubectl port-forward command
- `ingress_external_hostname`: External hostname (if ingress enabled)
- `ingress_internal_hostname`: Internal hostname (if ingress enabled)

---

## Troubleshooting

### Pods Not Starting

Check pod status and logs:

```bash
kubectl get pods -n <namespace>
kubectl describe pod -n <namespace> <pod-name>
kubectl logs -n <namespace> <pod-name>
```

### Ingress Not Working

Verify ingress resources:

```bash
# Check certificate status
kubectl get certificate -n istio-ingress

# Check gateway status
kubectl get gateway -n istio-ingress

# Check HTTPRoute status
kubectl get httproute -n <namespace>
```

### Insufficient Resources

If pods are pending due to insufficient resources, reduce resource requests or add more nodes to your cluster.

---

## Best Practices

1. **Production Deployments**:
   - Use at least 3 replicas for both Solr and Zookeeper
   - Allocate sufficient memory (2Gi+ for Solr in production)
   - Use persistent storage with appropriate IOPS

2. **Resource Planning**:
   - Monitor actual resource usage and adjust requests/limits accordingly
   - Use custom JVM settings to optimize memory usage
   - Tune garbage collection based on workload patterns

3. **Security**:
   - Enable ingress with TLS for external access
   - Use network policies to restrict traffic
   - Regularly update Solr image to latest stable version

4. **High Availability**:
   - Deploy odd number of Zookeeper replicas (3, 5, or 7)
   - Use pod anti-affinity to spread replicas across nodes
   - Configure pod disruption budgets for zero-downtime updates

5. **Backup and Recovery**:
   - Implement regular backup strategy for Solr indices
   - Test restore procedures before production deployment
   - Consider using snapshots for persistent volumes

---

## Additional Resources

- [Apache Solr Documentation](https://solr.apache.org/guide/)
- [Solr Operator Documentation](https://apache.github.io/solr-operator/)
- [Terraform Kubernetes Provider](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs)
- [Component README](../README.md)
- [Research Documentation](../docs/README.md)

---

## Support

For issues specific to this Terraform module, refer to the main component documentation. For general Solr questions, consult the Apache Solr community resources.
