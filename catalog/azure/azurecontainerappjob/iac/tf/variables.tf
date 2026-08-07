variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Azure Container App Job specification"
  type = object({
    # The Azure region (must match the environment's). ForceNew.
    region = string

    # The Azure Resource Group name. References are resolved to a literal
    # name by the platform before the module runs. ForceNew.
    resource_group = string

    # The job name (max 32 lowercase alphanumerics/hyphens, no
    # consecutive hyphens). ForceNew.
    job_name = string

    # The Container App Environment ARM ID. ForceNew.
    container_app_environment_id = string

    # Hard per-replica deadline in seconds; a replica killed by it counts
    # as failed.
    replica_timeout_in_seconds = number

    # Retries per failed replica before the execution is marked failed.
    replica_retry_limit = optional(number)

    # The environment workload profile to run on. Absent uses the
    # serverless Consumption profile.
    workload_profile_name = optional(string)

    # Main containers (at least one). Probe transports are the spec
    # enum's name strings (TCP_SOCKET / HTTP_GET / HTTPS_GET); the
    # per-type threshold contracts are spec-enforced.
    containers = list(object({
      name    = string
      image   = string
      cpu     = number
      memory  = string
      command = optional(list(string), [])
      args    = optional(list(string), [])
      env = optional(list(object({
        name        = string
        value       = optional(string)
        secret_name = optional(string)
      })), [])
      liveness_probe = optional(object({
        transport                = string
        port                     = number
        path                     = optional(string)
        host                     = optional(string)
        initial_delay_in_seconds = optional(number)
        interval_seconds         = optional(number, 10)
        timeout_seconds          = optional(number, 1)
        failure_count_threshold  = optional(number, 3)
        success_count_threshold  = optional(number)
        headers = optional(list(object({
          name  = string
          value = string
        })), [])
      }))
      readiness_probe = optional(object({
        transport                = string
        port                     = number
        path                     = optional(string)
        host                     = optional(string)
        initial_delay_in_seconds = optional(number)
        interval_seconds         = optional(number, 10)
        timeout_seconds          = optional(number, 1)
        failure_count_threshold  = optional(number, 3)
        success_count_threshold  = optional(number, 3)
        headers = optional(list(object({
          name  = string
          value = string
        })), [])
      }))
      startup_probe = optional(object({
        transport                = string
        port                     = number
        path                     = optional(string)
        host                     = optional(string)
        initial_delay_in_seconds = optional(number)
        interval_seconds         = optional(number, 10)
        timeout_seconds          = optional(number, 1)
        failure_count_threshold  = optional(number, 3)
        success_count_threshold  = optional(number)
        headers = optional(list(object({
          name  = string
          value = string
        })), [])
      }))
      volume_mounts = optional(list(object({
        name     = string
        path     = string
        sub_path = optional(string)
      })), [])
    }))

    # Init containers (run to completion before main containers; no
    # probes; cpu/memory optional).
    init_containers = optional(list(object({
      name    = string
      image   = string
      cpu     = optional(number)
      memory  = optional(string)
      command = optional(list(string), [])
      args    = optional(list(string), [])
      env = optional(list(object({
        name        = string
        value       = optional(string)
        secret_name = optional(string)
      })), [])
      volume_mounts = optional(list(object({
        name     = string
        path     = string
        sub_path = optional(string)
      })), [])
    })), [])

    # Volumes. storage_type is the spec enum's name string (EMPTY_DIR /
    # AZURE_FILE / NFS_AZURE_FILE / SECRET); storage_name pairs with the
    # share-backed types (spec-enforced).
    volumes = optional(list(object({
      name          = string
      storage_type  = optional(string)
      storage_name  = optional(string)
      mount_options = optional(string)
    })), [])

    # Exactly one trigger (spec-enforced); switching types is ForceNew.
    # Parallelism / completion counts carry documented defaults applied
    # here because the platform never materializes proto defaults.
    manual_trigger = optional(object({
      parallelism              = optional(number, 1)
      replica_completion_count = optional(number, 1)
    }))
    schedule_trigger = optional(object({
      cron_expression          = string
      parallelism              = optional(number, 1)
      replica_completion_count = optional(number, 1)
    }))
    event_trigger = optional(object({
      parallelism              = optional(number, 1)
      replica_completion_count = optional(number, 1)
      scale = optional(object({
        max_executions              = optional(number, 100)
        min_executions              = optional(number, 0)
        polling_interval_in_seconds = optional(number, 30)
        rules = optional(list(object({
          name             = string
          custom_rule_type = string
          metadata         = map(string)
          authentication = optional(list(object({
            secret_name       = string
            trigger_parameter = string
          })), [])
          identity_id = optional(string)
        })), [])
      }))
    }))

    # Job secrets: plain value XOR key_vault_secret_id (+ identity),
    # spec-enforced.
    secrets = optional(list(object({
      name                = string
      value               = optional(string)
      key_vault_secret_id = optional(string)
      identity            = optional(string)
    })), [])

    # Private registry credentials: managed identity XOR username +
    # password_secret_name, spec-enforced.
    registries = optional(list(object({
      server               = string
      username             = optional(string)
      password_secret_name = optional(string)
      identity             = optional(string)
    })), [])

    # Managed identity: type is the spec enum's name string
    # (SYSTEM_ASSIGNED / USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED).
    identity = optional(object({
      type                       = string
      user_assigned_identity_ids = optional(list(string), [])
    }))

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
