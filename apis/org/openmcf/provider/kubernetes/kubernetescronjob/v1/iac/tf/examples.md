# Kubernetes CronJob - Terraform Examples

This document provides practical examples for deploying Kubernetes CronJobs using the **KubernetesCronJob** Terraform module. Each example demonstrates different configuration aspects including scheduling, environment variables, secrets, concurrency policies, and retry strategies.

> **Note:** These examples show how to use the Terraform module directly. The module expects variables to match the protobuf specification defined in the KubernetesCronJob API.

---

## 1. Minimal Configuration

A basic example that runs a simple container on a daily schedule at midnight with default concurrency and retry settings.

```hcl
module "daily_backup_cronjob" {
  source = "./path/to/kubernetes-cronjob-module"

  metadata = {
    name = "daily-backup"
    id   = "daily-backup-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    schedule = "0 0 * * *"
    
    image = {
      repo = "busybox"
      tag  = "latest"
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
  }
}
```

**Key Points:**
- **schedule**: Cron format - runs daily at midnight (00:00 UTC)
- **image**: Uses minimal busybox image for demonstration
- **resources**: Basic CPU/memory requests and limits for resource management

---

## 2. CronJob with Environment Variables

This example shows how to pass non-sensitive configuration values as environment variables.

```hcl
module "weekly_report_cronjob" {
  source = "./path/to/kubernetes-cronjob-module"

  metadata = {
    name = "weekly-report"
    id   = "weekly-report-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    schedule = "0 9 * * 1"  # Every Monday at 9:00 AM
    
    image = {
      repo = "my-org/reports-generator"
      tag  = "v1.0.0"
    }
    
    env = {
      variables = {
        REPORT_TYPE    = "weekly"
        S3_BUCKET_NAME = "my-report-bucket"
        OUTPUT_FORMAT  = "pdf"
      }
    }
    
    resources = {
      limits = {
        cpu    = "1"
        memory = "1Gi"
      }
      requests = {
        cpu    = "100m"
        memory = "256Mi"
      }
    }
  }
}
```

**Key Points:**
- **schedule**: `"0 9 * * 1"` runs every Monday at 9:00 AM
- **env.variables**: Map of environment variables for non-sensitive configuration

---

## 3. Using Secret Environment Variables

Sensitive data should be stored in secrets. This example shows two approaches for providing secrets.

### Option 1: Direct String Value (Development/Testing)

```hcl
module "db_maintenance_cronjob" {
  source = "./path/to/kubernetes-cronjob-module"

  metadata = {
    name = "db-maintenance"
    id   = "db-maintenance-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    schedule = "0 3 * * *"  # Daily at 3:00 AM
    
    image = {
      repo = "my-org/maintenance-tool"
      tag  = "stable"
    }
    
    env = {
      secrets = {
        DB_PASSWORD = {
          value = "password-from-secret-manager"
        }
        API_KEY = {
          value = "api-key-from-secret-manager"
        }
      }
      
      variables = {
        DB_HOST = "db-server.prod.svc.cluster.local"
        DB_NAME = "myappdb"
        DB_PORT = "5432"
      }
    }
    
    resources = {
      limits = {
        cpu    = "1"
        memory = "512Mi"
      }
      requests = {
        cpu    = "100m"
        memory = "128Mi"
      }
    }
  }
}
```

### Option 2: Kubernetes Secret Reference (Production)

```hcl
module "db_maintenance_cronjob" {
  source = "./path/to/kubernetes-cronjob-module"

  metadata = {
    name = "db-maintenance"
    id   = "db-maintenance-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    schedule = "0 3 * * *"  # Daily at 3:00 AM
    
    image = {
      repo = "my-org/maintenance-tool"
      tag  = "stable"
    }
    
    env = {
      secrets = {
        DB_PASSWORD = {
          secret_ref = {
            name = "postgres-credentials"  # Name of existing K8s Secret
            key  = "password"              # Key within the Secret
          }
        }
        API_KEY = {
          secret_ref = {
            name = "api-credentials"
            key  = "key"
          }
        }
      }
      
      variables = {
        DB_HOST = "db-server.prod.svc.cluster.local"
        DB_NAME = "myappdb"
        DB_PORT = "5432"
      }
    }
    
    resources = {
      limits = {
        cpu    = "1"
        memory = "512Mi"
      }
      requests = {
        cpu    = "100m"
        memory = "128Mi"
      }
    }
  }
}
```

**Key Points:**
- **env.secrets.value**: Direct string value for development/testing scenarios
- **env.secrets.secret_ref**: Reference to an existing Kubernetes Secret (recommended for production)
- **env.variables**: Non-sensitive configuration values
- The module only creates a secret when there are direct string values to store

---

## 4. Concurrency Control & Retry Configuration

Control how CronJobs behave under load and when failures occur.

```hcl
module "heavy_lift_cronjob" {
  source = "./path/to/kubernetes-cronjob-module"

  metadata = {
    name = "heavy-lift-job"
    id   = "heavy-lift-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    schedule = "*/30 * * * *"  # Every 30 minutes
    
    # Prevent concurrent executions
    concurrency_policy = "Forbid"
    
    # Retry logic
    backoff_limit = 3
    
    # Job history retention
    successful_jobs_history_limit = 2
    failed_jobs_history_limit      = 2
    
    image = {
      repo = "my-org/heavy-lift"
      tag  = "v2.1"
    }
    
    resources = {
      limits = {
        cpu    = "2"
        memory = "2Gi"
      }
      requests = {
        cpu    = "200m"
        memory = "256Mi"
      }
    }
  }
}
```

**Key Points:**
- **concurrency_policy**: `"Forbid"` prevents new runs if previous is still running
  - Other options: `"Allow"` (default), `"Replace"` (kills old run, starts new)
- **backoff_limit**: Number of retries before marking job as failed (default: 6)
- **successful_jobs_history_limit**: Keep last N successful job pods (default: 3)
- **failed_jobs_history_limit**: Keep last N failed job pods (default: 1)

---

## 5. Advanced Configuration with Custom Commands

For complex scheduling and execution requirements, including custom commands and startup deadlines.

```hcl
module "end_of_month_report" {
  source = "./path/to/kubernetes-cronjob-module"

  metadata = {
    name = "end-of-month-report"
    id   = "eom-report-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    # Run at 23:30 on the last day of each month
    schedule = "30 23 28-31 * *"
    
    # If schedule is missed for over 2 hours, don't run retroactively
    starting_deadline_seconds = 7200
    
    # Temporarily suspend scheduling (set to false to resume)
    suspend = false
    
    # Replace running job if new schedule arrives
    concurrency_policy = "Replace"
    
    backoff_limit = 1
    
    successful_jobs_history_limit = 5
    failed_jobs_history_limit      = 5
    
    image = {
      repo = "my-org/report-service"
      tag  = "monthly-latest"
    }
    
    env = {
      variables = {
        REPORT_TZ               = "UTC"
        SEND_EMAIL_NOTIFICATIONS = "true"
        REPORT_FORMAT           = "excel"
      }
    }
    
    resources = {
      limits = {
        cpu    = "2"
        memory = "2Gi"
      }
      requests = {
        cpu    = "500m"
        memory = "512Mi"
      }
    }
  }
}
```

**Key Points:**
- **schedule**: `"30 23 28-31 * *"` runs at 23:30 on days 28-31 (catches last day of month)
- **starting_deadline_seconds**: Skip job if it can't start within this time window
- **suspend**: Set to `true` to pause scheduling; `false` to resume
- **concurrency_policy**: `"Replace"` cancels running job to start new one

---

## 6. gRPC Service Invocation Example

Example CronJob that periodically invokes a gRPC service using grpcurl.

```hcl
module "grpc_invoker_cronjob" {
  source = "./path/to/kubernetes-cronjob-module"

  metadata = {
    name = "grpc-invoker"
    id   = "grpc-invoker-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-namespace"
    create_namespace = true
    
    schedule = "0 0 * * *"  # Daily at midnight
    
    concurrency_policy = "Forbid"
    
    restart_policy = "Never"
    
    image = {
      repo = "fullstorydev/grpcurl"
      tag  = "latest"
    }
    
    # Custom command and arguments
    command = ["grpcurl"]
    args = [
      "-plaintext",
      "my-grpc-service:50051",
      "my.package.Service/Method"
    ]
    
    env = {
      variables = {
        GRPC_TIMEOUT = "30s"
      }
    }
    
    resources = {
      limits = {
        cpu    = "200m"
        memory = "256Mi"
      }
      requests = {
        cpu    = "50m"
        memory = "100Mi"
      }
    }
    
    backoff_limit = 6
    
    successful_jobs_history_limit = 3
    failed_jobs_history_limit      = 1
  }
}
```

**Key Points:**
- **command**: Override container's default ENTRYPOINT
- **args**: Override container's default CMD
- **restart_policy**: Typically `"Never"` for CronJobs (default)
  - Other option: `"OnFailure"` to retry within same pod

---

## 7. ConfigMap and Volume Mounts

Deploy a CronJob with inline ConfigMaps and volume mounts for configuration files.

```hcl
module "backup_cronjob" {
  source = "./path/to/kubernetes-cronjob-module"

  metadata = {
    name = "db-backup"
    id   = "db-backup-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "database"
    create_namespace = true
    
    schedule = "0 2 * * *"  # Daily at 2:00 AM
    
    concurrency_policy = "Forbid"
    
    image = {
      repo = "postgres"
      tag  = "15"
    }
    
    # Create ConfigMaps with inline content
    config_maps = {
      "backup-script" = <<-EOT
        #!/bin/bash
        echo "Starting backup at $(date)"
        pg_dump -h $DB_HOST -U $DB_USER $DB_NAME > /backup/dump.sql
        gzip /backup/dump.sql
        echo "Backup completed"
      EOT
    }
    
    # Mount ConfigMaps and create temporary storage
    volume_mounts = [
      {
        name       = "backup-script"
        mount_path = "/scripts/backup.sh"
        config_map = {
          name         = "backup-script"
          key          = "backup-script"
          path         = "backup.sh"
          default_mode = 493  # 0755 - executable
        }
      },
      {
        name       = "backup-data"
        mount_path = "/backup"
        empty_dir = {
          size_limit = "1Gi"
        }
      }
    ]
    
    command = ["/bin/bash", "/scripts/backup.sh"]
    
    env = {
      variables = {
        DB_HOST = "postgres.database.svc"
        DB_NAME = "myapp"
      }
      secrets = {
        DB_USER = {
          value = "admin"
        }
      }
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
    
    successful_jobs_history_limit = 3
    failed_jobs_history_limit     = 1
  }
}
```

**Key Points:**
- **config_maps**: Map of ConfigMap names to content (created by the module)
- **volume_mounts**: List of volumes to mount into the container
- **config_map volume**: Mount a ConfigMap key as a file with custom permissions
- **empty_dir volume**: Temporary storage for the backup (deleted when pod terminates)
- **default_mode = 493**: Decimal for 0755 permissions (executable script)

---

## 8. Multiple Volume Types

Example showing various volume mount types for complex configurations.

```hcl
module "data_processor_cronjob" {
  source = "./path/to/kubernetes-cronjob-module"

  metadata = {
    name = "data-processor"
    id   = "data-processor-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "batch-jobs"
    create_namespace = true
    
    schedule = "0 */6 * * *"  # Every 6 hours
    
    image = {
      repo = "python"
      tag  = "3.11-slim"
    }
    
    # Multiple ConfigMaps for different purposes
    config_maps = {
      "processor-script" = <<-EOT
        #!/bin/bash
        source /config/settings.env
        python /scripts/process.py --config /config/processor.yaml
      EOT
      
      "processor-config" = <<-EOT
        input_bucket: s3://data-input
        output_bucket: s3://data-output
        batch_size: 1000
        compression: gzip
      EOT
      
      "settings-env" = <<-EOT
        export AWS_REGION=us-west-2
        export LOG_LEVEL=info
      EOT
    }
    
    volume_mounts = [
      # Executable script from ConfigMap
      {
        name       = "processor-script"
        mount_path = "/scripts/run.sh"
        config_map = {
          name         = "processor-script"
          key          = "processor-script"
          path         = "run.sh"
          default_mode = 493  # 0755
        }
      },
      # YAML config file from ConfigMap
      {
        name       = "processor-config"
        mount_path = "/config/processor.yaml"
        config_map = {
          name = "processor-config"
          key  = "processor-config"
          path = "processor.yaml"
        }
      },
      # Environment file from ConfigMap
      {
        name       = "settings-env"
        mount_path = "/config/settings.env"
        config_map = {
          name = "settings-env"
          key  = "settings-env"
          path = "settings.env"
        }
      },
      # Memory-backed temp directory (fast I/O)
      {
        name       = "temp-work"
        mount_path = "/tmp/work"
        empty_dir = {
          medium     = "Memory"
          size_limit = "512Mi"
        }
      },
      # TLS certificates from existing Secret
      {
        name       = "tls-certs"
        mount_path = "/certs"
        read_only  = true
        secret = {
          name = "api-tls-certs"
        }
      }
    ]
    
    command = ["/bin/bash", "/scripts/run.sh"]
    
    resources = {
      limits = {
        cpu    = "2"
        memory = "4Gi"
      }
      requests = {
        cpu    = "500m"
        memory = "1Gi"
      }
    }
  }
}
```

**Key Points:**
- **Multiple ConfigMaps**: Different configs for different purposes
- **Memory-backed EmptyDir**: `medium = "Memory"` uses RAM for fast I/O
- **Secret volume**: Mount existing secrets (must be pre-created)
- **read_only**: Protect sensitive volumes from accidental writes

---

## 9. Docker Registry Authentication

When using private container registries, provide Docker credentials.

```hcl
module "private_registry_cronjob" {
  source = "./path/to/kubernetes-cronjob-module"

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
    
    namespace        = "my-namespace"
    create_namespace = true
    
    schedule = "0 */6 * * *"  # Every 6 hours
    
    image = {
      repo = "gcr.io/my-project/my-private-image"
      tag  = "v2.0.0"
    }
    
    resources = {
      limits = {
        cpu    = "1"
        memory = "512Mi"
      }
      requests = {
        cpu    = "100m"
        memory = "128Mi"
      }
    }
  }

  # Docker config for private registry authentication
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

# Variable for the service account key
variable "gcp_service_account_key" {
  description = "GCP service account key for accessing private container registry"
  type        = string
  sensitive   = true
}
```

**Key Points:**
- **docker_config_json**: Standard Docker config format for registry authentication
- The module creates a Kubernetes secret of type `kubernetes.io/dockerconfigjson`
- Automatically attached as imagePullSecret to the CronJob

---

## 10. Namespace Management

The `create_namespace` field controls whether the Terraform module creates a new namespace or references an existing one.

### Creating a New Namespace

Set `create_namespace = true` to automatically create the namespace with the CronJob:

```hcl
module "isolated_cronjob" {
  source = "./path/to/kubernetes-cronjob-module"

  metadata = {
    name = "isolated-job"
    id   = "isolated-job-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "my-cronjobs"
    create_namespace = true
    
    schedule = "0 2 * * *"
    
    image = {
      repo = "busybox"
      tag  = "latest"
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
  }
}
```

**Key Points:**
- **create_namespace = true**: The module creates the namespace `my-cronjobs` with appropriate labels
- **Use case**: Ideal when you want dedicated namespace isolation for this CronJob
- **Cleanup**: Destroying the module will also destroy the namespace and all its resources

### Using an Existing Namespace

Set `create_namespace = false` to deploy into a pre-existing namespace:

```hcl
module "shared_cronjob" {
  source = "./path/to/kubernetes-cronjob-module"

  metadata = {
    name = "shared-job"
    id   = "shared-job-prod"
    org  = "my-org"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-gke-cluster"
    }
    
    namespace        = "shared-batch-jobs"
    create_namespace = false
    
    schedule = "0 4 * * *"
    
    image = {
      repo = "busybox"
      tag  = "latest"
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
  }
}
```

**Key Points:**
- **create_namespace = false**: The module references the existing `shared-batch-jobs` namespace
- **Use case**: Ideal for multi-tenant scenarios where multiple CronJobs share a namespace
- **Important**: Ensure the namespace exists before applying, or Terraform will fail
- **Cleanup**: Destroying the module will only destroy the CronJob, not the shared namespace

### Terraform Implementation Details

When `create_namespace = true`:
- Uses `resource "kubernetes_namespace" "this"` with count = 1
- Namespace gets labels matching the CronJob metadata

When `create_namespace = false`:
- Uses `data "kubernetes_namespace" "existing"` to reference the namespace
- Namespace must exist before running `terraform apply`

### Best Practices

- **Isolated workloads**: Use `create_namespace = true` for dedicated, single-purpose CronJobs
- **Shared namespaces**: Use `create_namespace = false` when multiple CronJobs need to coexist
- **GitOps workflows**: Set `create_namespace = false` and manage namespaces separately for better lifecycle control
- **Development environments**: Use `create_namespace = true` for quick iterations and easy cleanup
- **Production**: Consider managing namespaces separately with dedicated Terraform modules for better governance

---

## Common Cron Schedule Patterns

| Schedule         | Description                                    |
|------------------|------------------------------------------------|
| `"* * * * *"`    | Every minute                                   |
| `"*/5 * * * *"`  | Every 5 minutes                                |
| `"0 * * * *"`    | Every hour (on the hour)                       |
| `"0 0 * * *"`    | Daily at midnight                              |
| `"0 9 * * 1-5"`  | Weekdays at 9:00 AM                            |
| `"0 0 * * 0"`    | Every Sunday at midnight                       |
| `"0 0 1 * *"`    | First day of every month at midnight           |
| `"30 23 * * *"`  | Every day at 11:30 PM                          |
| `"0 */6 * * *"`  | Every 6 hours                                  |
| `"0 0 */2 * *"`  | Every 2 days at midnight                       |

---

## Output Values

The module provides these outputs:

```hcl
output "namespace" {
  description = "The Kubernetes namespace where the CronJob is deployed"
  value       = module.my_cronjob.namespace
}

output "cronjob_name" {
  description = "The name of the CronJob resource"
  value       = module.my_cronjob.cronjob_name
}

output "service_account_name" {
  description = "The service account used by the CronJob"
  value       = module.my_cronjob.service_account_name
}
```

---

## Best Practices

### 1. Resource Limits
Always set resource requests and limits to prevent resource starvation:

```hcl
resources = {
  limits = {
    cpu    = "1"      # Maximum CPU
    memory = "1Gi"    # Maximum memory
  }
  requests = {
    cpu    = "100m"   # Guaranteed CPU
    memory = "256Mi"  # Guaranteed memory
  }
}
```

### 2. Concurrency Control
Choose the right concurrency policy:
- **Forbid**: Prevent overlapping runs (recommended for most cases)
- **Allow**: Multiple instances can run simultaneously
- **Replace**: Kill old run and start new one

### 3. History Limits
Control how many completed jobs are retained:

```hcl
successful_jobs_history_limit = 3  # Keep last 3 successful runs
failed_jobs_history_limit      = 1  # Keep last 1 failed run
```

This helps with debugging while preventing cluster resource exhaustion.

### 4. Startup Deadlines
Prevent "thundering herd" of missed schedules:

```hcl
starting_deadline_seconds = 300  # Skip if can't start within 5 minutes
```

### 5. Restart Policy
For CronJobs, typically use `"Never"`:

```hcl
restart_policy = "Never"  # Don't restart failed containers
```

Use `"OnFailure"` only if you want the pod to retry internally.

---

## Troubleshooting

### Check CronJob Status
```bash
kubectl get cronjobs -n <namespace>
kubectl describe cronjob <cronjob-name> -n <namespace>
```

### View Job History
```bash
kubectl get jobs -n <namespace>
kubectl get pods -n <namespace> --show-all
```

### View Logs
```bash
kubectl logs -n <namespace> <pod-name>
kubectl logs -n <namespace> -l job-name=<job-name>
```

### Test Schedule Manually
```bash
kubectl create job --from=cronjob/<cronjob-name> <test-job-name> -n <namespace>
```

---

## Summary

- **Scheduling**: Use standard cron format in the `schedule` field
- **Resources**: Always define CPU/memory requests and limits
- **Environment**: Use `env.variables` for config, `env.secrets` for sensitive data
- **ConfigMaps**: Use `config_maps` to create inline configuration content
- **Volume Mounts**: Use `volume_mounts` to mount ConfigMaps, Secrets, HostPaths, EmptyDirs, and PVCs
- **Concurrency**: Set `concurrency_policy` to control overlapping runs
- **Retry**: Configure `backoff_limit` for failure handling
- **History**: Use history limits to manage completed job retention
- **Commands**: Override container behavior with `command` and `args`
- **Registry Auth**: Provide `docker_config_json` for private images

For more information, consult the [KubernetesCronJob API documentation](../../README.md) or review the [research documentation](../../docs/README.md) for production best practices.

