# KubernetesJob Terraform Examples

This document provides Terraform examples for deploying KubernetesJob resources.

## Basic Job

```hcl
module "hello_world_job" {
  source = "./path/to/module"

  metadata = {
    name = "hello-world"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-cluster"
    }
    namespace        = "default"
    create_namespace = false
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
        memory = "64Mi"
      }
    }
    command = ["echo", "Hello, World!"]
  }
}
```

## Database Migration Job

```hcl
module "db_migration" {
  source = "./path/to/module"

  metadata = {
    name = "db-migration"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    target_cluster = {
      cluster_name = "prod-cluster"
    }
    namespace        = "migrations"
    create_namespace = true
    image = {
      repo = "myregistry/migration-runner"
      tag  = "v2.0.0"
    }
    resources = {
      limits = {
        cpu    = "1000m"
        memory = "2Gi"
      }
      requests = {
        cpu    = "250m"
        memory = "512Mi"
      }
    }
    env = {
      variables = {
        DATABASE_HOST = {
          value = "postgres.production.svc.cluster.local"
        }
        DATABASE_NAME = {
          value = "myapp"
        }
      }
      secrets = {
        DATABASE_USER = {
          secret_ref = {
            name = "db-credentials"
            key  = "username"
          }
        }
        DATABASE_PASSWORD = {
          secret_ref = {
            name = "db-credentials"
            key  = "password"
          }
        }
      }
    }
    backoff_limit              = 3
    active_deadline_seconds    = 1800
    ttl_seconds_after_finished = 3600
    restart_policy             = "Never"
    command                    = ["python", "manage.py", "migrate"]
  }
}
```

## Parallel Processing Job

```hcl
module "parallel_processor" {
  source = "./path/to/module"

  metadata = {
    name = "parallel-processor"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-cluster"
    }
    namespace        = "data-processing"
    create_namespace = true
    image = {
      repo = "myregistry/data-processor"
      tag  = "v1.5.0"
    }
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
    parallelism   = 5
    completions   = 20
    backoff_limit = 6
    env = {
      variables = {
        INPUT_BUCKET = {
          value = "s3://my-bucket/input"
        }
        OUTPUT_BUCKET = {
          value = "s3://my-bucket/output"
        }
      }
    }
    command = ["python", "/app/process.py"]
  }
}
```

## Indexed Job for Partitioned Processing

```hcl
module "indexed_processor" {
  source = "./path/to/module"

  metadata = {
    name = "indexed-processor"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-cluster"
    }
    namespace        = "data-processing"
    create_namespace = true
    image = {
      repo = "myregistry/partition-processor"
      tag  = "v1.0.0"
    }
    resources = {
      limits = {
        cpu    = "1000m"
        memory = "2Gi"
      }
      requests = {
        cpu    = "250m"
        memory = "512Mi"
      }
    }
    completion_mode = "Indexed"
    completions     = 10
    parallelism     = 5
    backoff_limit   = 3
    env = {
      variables = {
        TOTAL_PARTITIONS = {
          value = "10"
        }
      }
    }
    command = ["/bin/sh", "-c", "echo Processing partition $JOB_COMPLETION_INDEX && python /app/process_partition.py --partition=$JOB_COMPLETION_INDEX"]
  }
}
```

## Job with ConfigMap Volume

```hcl
module "backup_job" {
  source = "./path/to/module"

  metadata = {
    name = "database-backup"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-cluster"
    }
    namespace        = "backups"
    create_namespace = true
    image = {
      repo = "postgres"
      tag  = "15"
    }
    resources = {
      limits = {
        cpu    = "500m"
        memory = "1Gi"
      }
      requests = {
        cpu    = "100m"
        memory = "256Mi"
      }
    }
    env = {
      variables = {
        PGHOST = {
          value = "postgres.production.svc.cluster.local"
        }
        PGDATABASE = {
          value = "myapp"
        }
      }
      secrets = {
        PGUSER = {
          secret_ref = {
            name = "db-credentials"
            key  = "username"
          }
        }
        PGPASSWORD = {
          secret_ref = {
            name = "db-credentials"
            key  = "password"
          }
        }
      }
    }
    backoff_limit              = 2
    active_deadline_seconds    = 7200
    ttl_seconds_after_finished = 604800
    config_maps = {
      "backup-script" = <<-EOT
        #!/bin/bash
        set -e
        BACKUP_FILE="/backup/backup-$(date +%Y%m%d-%H%M%S).sql.gz"
        echo "Creating backup: $BACKUP_FILE"
        pg_dump | gzip > "$BACKUP_FILE"
        echo "Backup complete: $BACKUP_FILE"
      EOT
    }
    volume_mounts = [
      {
        name       = "backup-script"
        mount_path = "/scripts/backup.sh"
        config_map = {
          name         = "backup-script"
          key          = "backup-script"
          path         = "backup.sh"
          default_mode = 493 # 0755
        }
      }
    ]
    command = ["/bin/bash", "/scripts/backup.sh"]
  }
}
```

## Job with Private Registry

```hcl
module "private_image_job" {
  source = "./path/to/module"

  metadata = {
    name = "private-job"
  }

  spec = {
    target_cluster = {
      cluster_name = "my-cluster"
    }
    namespace        = "jobs"
    create_namespace = true
    image = {
      repo = "private.registry.com/myapp/worker"
      tag  = "v1.0.0"
    }
    resources = {
      limits = {
        cpu    = "1000m"
        memory = "1Gi"
      }
      requests = {
        cpu    = "250m"
        memory = "256Mi"
      }
    }
    command = ["/app/run-job"]
  }

  docker_config_json = jsonencode({
    auths = {
      "private.registry.com" = {
        auth = base64encode("username:password")
      }
    }
  })
}
```
