# KubernetesJob Terraform Module

This Terraform module deploys a KubernetesJob to a Kubernetes cluster.

## Overview

The module creates the following Kubernetes resources:

1. **Namespace** (optional) - Created if `create_namespace: true`
2. **Secret** (optional) - For environment secrets with direct values
3. **Image Pull Secret** (optional) - If Docker credentials are provided
4. **Job** - The main batch workload

The module never creates ServiceAccounts, RBAC objects, or ConfigMaps -- pods run as the ServiceAccount referenced in `spec.pod.service_account` (composed through KubernetesServiceAccount and KubernetesRbac resources), and configuration composes through KubernetesConfigMap resources. The Kubernetes permissions the IaC runner needs are declared in `../permissions.yaml`.

## Usage

### Basic Example

```hcl
module "kubernetes_job" {
  source = "./path/to/module"

  metadata = {
    name = "data-migration"
  }

  spec = {
    namespace        = "batch-jobs"
    create_namespace = true
    image = {
      repo = "myregistry/migration-runner"
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
    backoff_limit              = 3
    active_deadline_seconds    = 3600
    ttl_seconds_after_finished = 86400
    command = ["python", "/app/migrate.py"]
  }
}
```

### With Environment Variables

```hcl
module "kubernetes_job" {
  source = "./path/to/module"

  metadata = {
    name = "etl-job"
  }

  spec = {
    namespace        = "data-processing"
    create_namespace = true
    image = {
      repo = "myregistry/etl-runner"
      tag  = "v2.0.0"
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
    env = {
      variables = {
        INPUT_PATH = {
          value = "/data/input"
        }
        OUTPUT_PATH = {
          value = "/data/output"
        }
      }
      secrets = {
        DATABASE_PASSWORD = {
          secret_ref = {
            name = "db-credentials"
            key  = "password"
          }
        }
      }
    }
  }
}
```

### Parallel Job

```hcl
module "kubernetes_job" {
  source = "./path/to/module"

  metadata = {
    name = "parallel-processor"
  }

  spec = {
    namespace        = "batch"
    create_namespace = true
    image = {
      repo = "myregistry/processor"
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
    parallelism = 5
    completions = 20
    backoff_limit = 3
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| metadata | Resource metadata including name, org, env | object | yes |
| spec | Job specification | object | yes |
| docker_config_json | Docker credentials for private registries | string | no |

## Outputs

| Name | Description |
|------|-------------|
| namespace | The Kubernetes namespace |
| job_name | The name of the Job |
| selector_labels | Labels selecting the Job's pods |

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.0 |
| kubernetes | ~> 2.35 |

## Resources Created

- `kubernetes_namespace` (conditional)
- `kubernetes_secret_v1` (conditional -- env secrets and image pull)
- `kubernetes_job_v1`

## Notes

- Jobs run to completion and then stop
- Use `ttl_seconds_after_finished` for automatic cleanup
- Set `active_deadline_seconds` to prevent runaway jobs
- Use `backoff_limit` to control retry behavior
