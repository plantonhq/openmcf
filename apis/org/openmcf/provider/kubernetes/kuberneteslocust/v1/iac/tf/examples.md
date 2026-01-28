# Terraform Examples for KubernetesLocust

This document provides Terraform usage examples for deploying Locust on Kubernetes.

## Example 1: Basic Locust Deployment

```hcl
module "locust_basic" {
  source = "./path/to/kuberneteslocust/v1/iac/tf"

  metadata = {
    name = "locust-basic"
  }

  spec = {
    target_cluster_name = "my-gke-cluster"
    namespace           = "locust-test"
    create_namespace    = true

    master_container = {
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      replicas = 1
    }

    worker_container = {
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      replicas = 2
    }

    load_test = {
      name            = "basic-load-test"
      main_py_content = <<-EOT
        from locust import HttpUser, task

        class MyUser(HttpUser):
            @task
            def my_task(self):
                self.client.get("/api/test")
      EOT
    }

    ingress = {
      enabled = false
    }
  }
}
```

## Example 2: Locust with Ingress and Custom Domain

```hcl
module "locust_with_ingress" {
  source = "./path/to/kuberneteslocust/v1/iac/tf"

  metadata = {
    name = "locust-ingress"
  }

  spec = {
    target_cluster_name = "production-gke-cluster"
    namespace           = "locust-prod"
    create_namespace    = true

    master_container = {
      resources = {
        requests = {
          cpu    = "200m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "2"
          memory = "2Gi"
        }
      }
      replicas = 1
    }

    worker_container = {
      resources = {
        requests = {
          cpu    = "200m"
          memory = "512Mi"
        }
        limits = {
          cpu    = "2"
          memory = "2Gi"
        }
      }
      replicas = 5
    }

    load_test = {
      name            = "custom-load-test"
      main_py_content = <<-EOT
        from locust import HttpUser, task

        class MyUser(HttpUser):
            @task
            def my_task(self):
                self.client.post("/api/test", json={"key": "value"})
      EOT

      lib_files_content = {
        "utils.py" = <<-EOT
          def helper_function():
              return "Helper"
        EOT
      }

      pip_packages = [
        "requests",
        "locust"
      ]
    }

    ingress = {
      enabled            = true
      ingress_class_name = "nginx"
      hosts = [
        {
          host  = "locust.mydomain.com"
          paths = ["/"]
        }
      ]
    }
  }
}
```

## Example 3: Locust with TLS

```hcl
module "locust_tls" {
  source = "./path/to/kuberneteslocust/v1/iac/tf"

  metadata = {
    name = "locust-tls"
  }

  spec = {
    target_cluster_name = "my-gke-cluster"
    namespace           = "locust-tls"
    create_namespace    = true

    master_container = {
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      replicas = 1
    }

    worker_container = {
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      replicas = 3
    }

    load_test = {
      name            = "tls-load-test"
      main_py_content = <<-EOT
        from locust import HttpUser, task

        class MyUser(HttpUser):
            @task
            def my_task(self):
                self.client.get("/secure-api/test")
      EOT
    }

    ingress = {
      enabled            = true
      ingress_class_name = "nginx"
      hosts = [
        {
          host  = "locust-tls.mydomain.com"
          paths = ["/"]
        }
      ]
      tls = [
        {
          secret_name = "locust-tls-cert"
          hosts       = ["locust-tls.mydomain.com"]
        }
      ]
    }
  }
}
```

## Example 4: Locust with External Libraries

```hcl
module "locust_external_lib" {
  source = "./path/to/kuberneteslocust/v1/iac/tf"

  metadata = {
    name = "locust-external-lib"
  }

  spec = {
    target_cluster_name = "dev-cluster"
    namespace           = "locust-dev"
    create_namespace    = true

    master_container = {
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      replicas = 1
    }

    worker_container = {
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      replicas = 2
    }

    load_test = {
      name            = "external-lib-load-test"
      main_py_content = <<-EOT
        from locust import HttpUser, task
        from utils import helper_function

        class MyUser(HttpUser):
            @task
            def my_task(self):
                result = helper_function()
                self.client.get(f"/api/test?result={result}")
      EOT

      lib_files_content = {
        "utils.py" = <<-EOT
          def helper_function():
              return "Hello from helper!"
        EOT
      }

      pip_packages = [
        "requests",
        "locust"
      ]
    }

    ingress = {
      enabled = false
    }
  }
}
```

## Example 5: Minimal Configuration

```hcl
module "locust_minimal" {
  source = "./path/to/kuberneteslocust/v1/iac/tf"

  metadata = {
    name = "locust-minimal"
  }

  spec = {
    target_cluster_name = "test-cluster"
    namespace           = "locust-minimal"
    create_namespace    = true

    master_container = {
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
      replicas = 1
    }

    worker_container = {
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
      replicas = 1
    }

    load_test = {
      name            = "minimal-test"
      main_py_content = <<-EOT
        from locust import HttpUser, task

        class MyUser(HttpUser):
            @task
            def my_task(self):
                self.client.get("/")
      EOT
    }

    ingress = {
      enabled = false
    }
  }
}
```

## Example 6: Using Existing Namespace

```hcl
module "locust_existing_namespace" {
  source = "./path/to/kuberneteslocust/v1/iac/tf"

  metadata = {
    name = "locust-prod"
  }

  spec = {
    target_cluster_name = "prod-gke-cluster"
    namespace           = "shared-load-testing"
    create_namespace    = false  # Use existing namespace

    master_container = {
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      replicas = 1
    }

    worker_container = {
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      replicas = 2
    }

    load_test = {
      name            = "existing-ns-test"
      main_py_content = <<-EOT
        from locust import HttpUser, task

        class MyUser(HttpUser):
            @task
            def my_task(self):
                self.client.get("/api/test")
      EOT
    }

    ingress = {
      enabled = false
    }
  }
}
```

## Output Values

The module provides the following outputs:

- Locust service endpoints
- Ingress configuration (if enabled)
- Namespace information
- Other resource identifiers

Refer to `outputs.tf` for the complete list of available outputs.

