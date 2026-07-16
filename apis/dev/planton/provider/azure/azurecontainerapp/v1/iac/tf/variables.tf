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
  description = "Azure Container App specification"
  type = object({
    # The Azure Resource Group name. References are resolved to a literal
    # name by the platform before the module runs. ForceNew.
    resource_group = string

    # The name of the Container App (max 32 lowercase alphanumerics,
    # hyphens, dots; no consecutive hyphens). ForceNew.
    container_app_name = string

    # The Container App Environment ARM ID. ForceNew.
    container_app_environment_id = string

    # The revision mode, as the spec enum's name string (SINGLE /
    # MULTIPLE). Absent deploys SINGLE.
    revision_mode = optional(string)

    # The environment workload profile to run on. Absent uses the
    # serverless Consumption profile.
    workload_profile_name = optional(string)

    # Maximum inactive revisions retained (0-100).
    max_inactive_revisions = optional(number)

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

    # Replica bounds and scaler dials (documented defaults applied here
    # because the platform never materializes proto defaults).
    min_replicas                     = optional(number, 0)
    max_replicas                     = optional(number, 10)
    cooldown_period_in_seconds       = optional(number, 300)
    polling_interval_in_seconds      = optional(number, 30)
    revision_suffix                  = optional(string)
    termination_grace_period_seconds = optional(number, 0)

    # Scale rules. Authentication trigger_parameter is optional on
    # HTTP/TCP rules and spec-required elsewhere.
    http_scale_rules = optional(list(object({
      name                = string
      concurrent_requests = string
      authentication = optional(list(object({
        secret_name       = string
        trigger_parameter = optional(string)
      })), [])
    })), [])
    tcp_scale_rules = optional(list(object({
      name                = string
      concurrent_requests = string
      authentication = optional(list(object({
        secret_name       = string
        trigger_parameter = optional(string)
      })), [])
    })), [])
    azure_queue_scale_rules = optional(list(object({
      name         = string
      queue_name   = string
      queue_length = number
      authentication = list(object({
        secret_name       = string
        trigger_parameter = optional(string)
      }))
    })), [])
    custom_scale_rules = optional(list(object({
      name             = string
      custom_rule_type = string
      metadata         = map(string)
      authentication = optional(list(object({
        secret_name       = string
        trigger_parameter = optional(string)
      })), [])
      identity_id = optional(string)
    })), [])

    # App secrets: plain value XOR key_vault_secret_id (+ identity),
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

    # Ingress. transport / client_certificate_mode / restriction actions
    # are the spec enums' name strings.
    ingress = optional(object({
      external_enabled           = optional(bool, false)
      target_port                = number
      exposed_port               = optional(number)
      transport                  = optional(string)
      allow_insecure_connections = optional(bool, false)
      client_certificate_mode    = optional(string)
      traffic_weight = list(object({
        latest_revision = optional(bool, false)
        revision_suffix = optional(string)
        percentage      = number
        label           = optional(string)
      }))
      ip_security_restrictions = optional(list(object({
        name             = string
        action           = string
        ip_address_range = string
        description      = optional(string)
      })), [])
      cors = optional(object({
        allowed_origins           = list(string)
        allowed_headers           = optional(list(string), [])
        allowed_methods           = optional(list(string), [])
        exposed_headers           = optional(list(string), [])
        max_age_in_seconds        = optional(number)
        allow_credentials_enabled = optional(bool, false)
      }))
    }))

    # Dapr sidecar. app_protocol is the spec enum's name string
    # (DAPR_HTTP / DAPR_GRPC); absent deploys http.
    dapr = optional(object({
      app_id       = string
      app_port     = optional(number)
      app_protocol = optional(string)
    }))

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
