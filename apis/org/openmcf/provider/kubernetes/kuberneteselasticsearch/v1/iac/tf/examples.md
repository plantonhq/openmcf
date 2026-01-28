# Terraform Examples for ElasticsearchKubernetes

This document provides examples of how to use the Terraform module for deploying Elasticsearch on Kubernetes.

---

## Example 1: Basic Elasticsearch and Kibana Deployment

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ElasticsearchKubernetes
metadata:
  name: logging-cluster
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace: logging
  elasticsearch:
    container:
      replicas: 1
      resources:
        requests:
          cpu: 500m
          memory: 1Gi
        limits:
          cpu: 1000m
          memory: 2Gi
      persistence_enabled: true
      disk_size: 10Gi
    ingress:
      enabled: false
      hostname: ""
  kibana:
    enabled: true
    container:
      replicas: 1
      resources:
        requests:
          cpu: 200m
          memory: 512Mi
        limits:
          cpu: 500m
          memory: 1Gi
    ingress:
      enabled: false
      hostname: ""
```

### Terraform Usage

```hcl
module "logging_cluster" {
  source = "../"

  metadata = {
    name = "logging-cluster"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "logging"
    create_namespace = true
    elasticsearch = {
      container = {
        replicas            = 1
        persistence_enabled = true
        disk_size          = "10Gi"
        resources = {
          requests = {
            cpu    = "500m"
            memory = "1Gi"
          }
          limits = {
            cpu    = "1000m"
            memory = "2Gi"
          }
        }
      }
      ingress = {
        enabled  = false
        hostname = ""
      }
    }
    kibana = {
      enabled = true
      container = {
        replicas = 1
        resources = {
          requests = {
            cpu    = "200m"
            memory = "512Mi"
          }
          limits = {
            cpu    = "500m"
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
}
```

---

## Example 2: Elasticsearch with Ingress Enabled

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ElasticsearchKubernetes
metadata:
  name: search-service
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace: search
  elasticsearch:
    container:
      replicas: 3
      resources:
        requests:
          cpu: 1
          memory: 2Gi
        limits:
          cpu: 2
          memory: 4Gi
      persistence_enabled: true
      disk_size: 50Gi
    ingress:
      enabled: true
      hostname: search.example.com
  kibana:
    enabled: true
    container:
      replicas: 1
      resources:
        requests:
          cpu: 200m
          memory: 512Mi
        limits:
          cpu: 500m
          memory: 1Gi
    ingress:
      enabled: true
      hostname: search-kibana.example.com
```

### Terraform Usage

```hcl
module "search_service" {
  source = "../"

  metadata = {
    name = "search-service"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "search"
    create_namespace = true
    elasticsearch = {
      container = {
        replicas            = 3
        persistence_enabled = true
        disk_size          = "50Gi"
        resources = {
          requests = {
            cpu    = "1"
            memory = "2Gi"
          }
          limits = {
            cpu    = "2"
            memory = "4Gi"
          }
        }
      }
      ingress = {
        enabled  = true
        hostname = "search.example.com"
      }
    }
    kibana = {
      enabled = true
      container = {
        replicas = 1
        resources = {
          requests = {
            cpu    = "200m"
            memory = "512Mi"
          }
          limits = {
            cpu    = "500m"
            memory = "1Gi"
          }
        }
      }
      ingress = {
        enabled  = true
        hostname = "search-kibana.example.com"
      }
    }
  }
}
```

---

## Example 3: Elasticsearch with Persistent Storage

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ElasticsearchKubernetes
metadata:
  name: persistent-cluster
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace: elasticsearch-prod
  elasticsearch:
    container:
      persistence_enabled: true
      disk_size: 100Gi
      replicas: 5
      resources:
        requests:
          cpu: 1
          memory: 2Gi
        limits:
          cpu: 2
          memory: 4Gi
    ingress:
      enabled: true
      hostname: elasticsearch.example.com
  kibana:
    enabled: true
    container:
      replicas: 2
      resources:
        requests:
          cpu: 200m
          memory: 512Mi
        limits:
          cpu: 500m
          memory: 1Gi
    ingress:
      enabled: true
      hostname: kibana.example.com
```

### Terraform Usage

```hcl
module "persistent_cluster" {
  source = "../"

  metadata = {
    name = "persistent-cluster"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "elasticsearch-prod"
    create_namespace = true
    elasticsearch = {
      container = {
        persistence_enabled = true
        disk_size          = "100Gi"
        replicas          = 5
        resources = {
          requests = {
            cpu    = "1"
            memory = "2Gi"
          }
          limits = {
            cpu    = "2"
            memory = "4Gi"
          }
        }
      }
      ingress = {
        enabled  = true
        hostname = "elasticsearch.example.com"
      }
    }
    kibana = {
      enabled = true
      container = {
        replicas = 2
        resources = {
          requests = {
            cpu    = "200m"
            memory = "512Mi"
          }
          limits = {
            cpu    = "500m"
            memory = "1Gi"
          }
        }
      }
      ingress = {
        enabled  = true
        hostname = "kibana.example.com"
      }
    }
  }
}
```

---

## Example 4: Minimal Elasticsearch Deployment

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ElasticsearchKubernetes
metadata:
  name: minimal-elasticsearch
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace: elasticsearch
  create_namespace: true
```

### Terraform Usage

```hcl
module "minimal_elasticsearch" {
  source = "../"

  metadata = {
    name = "minimal-elasticsearch"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "elasticsearch"
    create_namespace = true
    # All other fields will use defaults from the proto schema
  }
}
```

This minimal example provides only the required fields. All other fields use the default values configured in the proto schema.

---

## Example 5: Using Existing Namespace (create_namespace: false)

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: ElasticsearchKubernetes
metadata:
  name: shared-elasticsearch
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace: shared-services
  create_namespace: false
```

### Terraform Usage

```hcl
module "shared_elasticsearch" {
  source = "../"

  metadata = {
    name = "shared-elasticsearch"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    namespace        = "shared-services"
    create_namespace = false
    elasticsearch = {
      container = {
        replicas            = 1
        persistence_enabled = true
        disk_size          = "10Gi"
        resources = {
          requests = {
            cpu    = "500m"
            memory = "1Gi"
          }
          limits = {
            cpu    = "1000m"
            memory = "2Gi"
          }
        }
      }
      ingress = {
        enabled  = false
        hostname = ""
      }
    }
    kibana = {
      enabled = true
      container = {
        replicas = 1
        resources = {
          requests = {
            cpu    = "200m"
            memory = "512Mi"
          }
          limits = {
            cpu    = "500m"
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
}
```

**Note**: When `create_namespace: false`, ensure the namespace `shared-services` exists before applying this configuration.

---

## Running Terraform

To use any of these examples:

1. Save the Terraform code to a `.tf` file (e.g., `main.tf`)
2. Initialize Terraform:

```bash
terraform init
```

3. Review the planned changes:

```bash
terraform plan
```

4. Apply the configuration:

```bash
terraform apply
```

5. To destroy the resources:

```bash
terraform destroy
```

---

## Output Variables

After applying, you can access the following outputs:

```hcl
output "elasticsearch_endpoint" {
  value = module.search_service.elasticsearch_endpoint
}

output "kibana_endpoint" {
  value = module.search_service.kibana_endpoint
}

output "port_forward_command_elasticsearch" {
  value = module.search_service.port_forward_command_elasticsearch
}

output "port_forward_command_kibana" {
  value = module.search_service.port_forward_command_kibana
}
```

---

## Notes

- All examples require `target_cluster` and `namespace` fields - these are mandatory
- The `target_cluster.cluster_name` should reference an existing Kubernetes cluster
- Adjust resource limits and requests based on your workload requirements
- For production deployments, always enable persistence and use multiple replicas
- Ingress hostnames should be valid DNS names that resolve to your Kubernetes cluster's ingress controller
- The `namespace` field accepts a simple string value in Terraform (unlike YAML which uses StringValueOrRef)

