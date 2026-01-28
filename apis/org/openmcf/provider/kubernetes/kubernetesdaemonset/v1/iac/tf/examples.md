# KubernetesDaemonSet Terraform Examples

## Basic Log Collector

Deploy a Fluentd log collector on all nodes:

```hcl
module "fluentd" {
  source = "./path/to/module"

  metadata = {
    name = "fluentd"
  }

  spec = {
    namespace        = "logging"
    create_namespace = true

    container = {
      app = {
        image = {
          repo = "fluent/fluentd-kubernetes-daemonset"
          tag  = "v1.16-debian-elasticsearch8"
        }

        resources = {
          limits = {
            cpu    = "500m"
            memory = "512Mi"
          }
          requests = {
            cpu    = "100m"
            memory = "200Mi"
          }
        }

        env = {
          variables = {
            FLUENT_ELASTICSEARCH_HOST = "elasticsearch.logging.svc.cluster.local"
            FLUENT_ELASTICSEARCH_PORT = "9200"
          }
        }

        volume_mounts = [
          {
            name       = "varlog"
            mount_path = "/var/log"
            read_only  = true
            host_path = {
              path = "/var/log"
            }
          },
          {
            name       = "containers"
            mount_path = "/var/lib/docker/containers"
            read_only  = true
            host_path = {
              path = "/var/lib/docker/containers"
            }
          }
        ]

        ports = [
          {
            name             = "metrics"
            container_port   = 24231
            network_protocol = "TCP"
          }
        ]
      }
    }
  }
}
```

## Node Monitoring with Prometheus

Deploy Node Exporter on all nodes including masters:

```hcl
module "node_exporter" {
  source = "./path/to/module"

  metadata = {
    name = "node-exporter"
  }

  spec = {
    namespace        = "monitoring"
    create_namespace = true

    container = {
      app = {
        image = {
          repo = "prom/node-exporter"
          tag  = "v1.7.0"
        }

        resources = {
          limits = {
            cpu    = "250m"
            memory = "180Mi"
          }
          requests = {
            cpu    = "100m"
            memory = "100Mi"
          }
        }

        args = [
          "--path.procfs=/host/proc",
          "--path.sysfs=/host/sys",
          "--path.rootfs=/host/root",
          "--collector.filesystem.mount-points-exclude=^/(dev|proc|sys|var/lib/docker/.+)($|/)"
        ]

        volume_mounts = [
          {
            name       = "proc"
            mount_path = "/host/proc"
            read_only  = true
            host_path = {
              path = "/proc"
            }
          },
          {
            name       = "sys"
            mount_path = "/host/sys"
            read_only  = true
            host_path = {
              path = "/sys"
            }
          },
          {
            name       = "root"
            mount_path = "/host/root"
            read_only  = true
            host_path = {
              path = "/"
            }
          }
        ]

        ports = [
          {
            name             = "metrics"
            container_port   = 9100
            network_protocol = "TCP"
          }
        ]
      }
    }

    tolerations = [
      {
        key      = "node-role.kubernetes.io/master"
        operator = "Exists"
        effect   = "NoSchedule"
      },
      {
        key      = "node-role.kubernetes.io/control-plane"
        operator = "Exists"
        effect   = "NoSchedule"
      }
    ]
  }
}
```

## With Environment Secrets (Direct String Values)

Deploy with secrets using direct string values for development:

```hcl
module "logging_agent" {
  source = "./path/to/module"

  metadata = {
    name = "logging-agent"
  }

  spec = {
    namespace        = "logging"
    create_namespace = true

    container = {
      app = {
        image = {
          repo = "fluent/fluent-bit"
          tag  = "2.2.0"
        }

        resources = {
          limits = {
            cpu    = "200m"
            memory = "256Mi"
          }
          requests = {
            cpu    = "50m"
            memory = "64Mi"
          }
        }

        env = {
          variables = {
            LOG_LEVEL   = "info"
            OUTPUT_HOST = "elasticsearch.logging.svc.cluster.local"
          }
          secrets = {
            # Direct string value (for development/testing)
            API_KEY = {
              value = "my-dev-api-key"
            }
          }
        }

        volume_mounts = [
          {
            name       = "varlog"
            mount_path = "/var/log"
            read_only  = true
            host_path = {
              path = "/var/log"
            }
          }
        ]
      }
    }

    tolerations = [
      {
        operator = "Exists"
      }
    ]
  }
}
```

## With Kubernetes Secret References (Production)

Deploy with secrets referencing existing Kubernetes Secrets:

```hcl
module "datadog_agent" {
  source = "./path/to/module"

  metadata = {
    name = "datadog-agent"
  }

  spec = {
    namespace        = "monitoring"
    create_namespace = true

    container = {
      app = {
        image = {
          repo = "gcr.io/datadoghq/agent"
          tag  = "7.50.0"
        }

        resources = {
          limits = {
            cpu    = "500m"
            memory = "512Mi"
          }
          requests = {
            cpu    = "100m"
            memory = "256Mi"
          }
        }

        env = {
          variables = {
            DD_SITE         = "datadoghq.com"
            DD_LOGS_ENABLED = "true"
            DD_APM_ENABLED  = "true"
          }
          secrets = {
            # Reference existing Kubernetes Secret (production)
            DD_API_KEY = {
              secret_ref = {
                name = "datadog-credentials"
                key  = "api-key"
              }
            }
            DD_APP_KEY = {
              secret_ref = {
                name = "datadog-credentials"
                key  = "app-key"
              }
            }
          }
        }

        volume_mounts = [
          {
            name       = "dockersock"
            mount_path = "/var/run/docker.sock"
            host_path = {
              path = "/var/run/docker.sock"
            }
          },
          {
            name       = "procdir"
            mount_path = "/host/proc"
            read_only  = true
            host_path = {
              path = "/proc"
            }
          }
        ]

        security_context = {
          privileged = true
        }
      }
    }

    tolerations = [
      {
        operator = "Exists"
      }
    ]
  }
}
```

## With Mixed Secret Types

Deploy with both direct string values and Kubernetes Secret references:

```hcl
module "custom_collector" {
  source = "./path/to/module"

  metadata = {
    name = "custom-collector"
  }

  spec = {
    namespace        = "observability"
    create_namespace = true

    container = {
      app = {
        image = {
          repo = "myregistry/custom-collector"
          tag  = "v1.0.0"
        }

        resources = {
          limits = {
            cpu    = "300m"
            memory = "384Mi"
          }
          requests = {
            cpu    = "100m"
            memory = "128Mi"
          }
        }

        env = {
          variables = {
            COLLECTOR_MODE = "agent"
            CLUSTER_NAME   = "production"
          }
          secrets = {
            # Dev/test secrets can use direct values
            DEBUG_TOKEN = {
              value = "debug-only-token"
            }
            # Production secrets should reference external K8s Secrets
            DATABASE_PASSWORD = {
              secret_ref = {
                name = "collector-secrets"
                key  = "db-password"
              }
            }
            CLOUD_API_KEY = {
              secret_ref = {
                name = "cloud-credentials"
                key  = "api-key"
              }
            }
          }
        }
      }
    }

    tolerations = [
      {
        operator = "Exists"
      }
    ]
  }
}
```

## Privileged Network Tool with RBAC

Deploy a privileged network debugging tool with RBAC:

```hcl
module "network_debug" {
  source = "./path/to/module"

  metadata = {
    name = "network-debug"
  }

  spec = {
    namespace = "kube-system"

    create_service_account = true

    container = {
      app = {
        image = {
          repo = "nicolaka/netshoot"
          tag  = "latest"
        }

        resources = {
          limits = {
            cpu    = "200m"
            memory = "256Mi"
          }
          requests = {
            cpu    = "50m"
            memory = "64Mi"
          }
        }

        command = ["/bin/bash", "-c", "sleep infinity"]

        security_context = {
          privileged = true
          capabilities = {
            add = ["NET_ADMIN", "SYS_PTRACE"]
          }
        }
      }
    }

    node_selector = {
      "kubernetes.io/os" = "linux"
    }

    tolerations = [
      {
        operator = "Exists"
      }
    ]

    rbac = {
      cluster_rules = [
        {
          api_groups = [""]
          resources  = ["pods", "nodes"]
          verbs      = ["get", "list", "watch"]
        }
      ]
    }
  }
}
```

## With Update Strategy

Deploy with controlled rolling updates:

```hcl
module "kube_proxy_custom" {
  source = "./path/to/module"

  metadata = {
    name = "kube-proxy-custom"
  }

  spec = {
    namespace = "kube-system"

    container = {
      app = {
        image = {
          repo = "registry.k8s.io/kube-proxy"
          tag  = "v1.29.0"
        }

        resources = {
          limits = {
            cpu    = "500m"
            memory = "256Mi"
          }
          requests = {
            cpu    = "100m"
            memory = "128Mi"
          }
        }

        security_context = {
          privileged = true
        }
      }
    }

    update_strategy = {
      type = "RollingUpdate"
      rolling_update = {
        max_unavailable = "10%"
      }
    }

    min_ready_seconds = 30

    tolerations = [
      {
        operator = "Exists"
      }
    ]
  }
}
```

