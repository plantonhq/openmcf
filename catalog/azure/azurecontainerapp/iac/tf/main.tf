# The Container App: a continuously running containerized service inside
# a Container App Environment. Name, resource group, and environment are
# ForceNew; everything in the template creates a new revision on change.
resource "azurerm_container_app" "main" {
  name                         = var.spec.container_app_name
  resource_group_name          = var.spec.resource_group
  container_app_environment_id = var.spec.container_app_environment_id

  # Single = one active revision (new replaces old); Multiple = several
  # active at once with traffic splitting (blue-green / canary).
  revision_mode = local.revision_mode

  # Omitted runs on the environment's serverless Consumption profile.
  workload_profile_name  = var.spec.workload_profile_name
  max_inactive_revisions = var.spec.max_inactive_revisions

  # The template is the revision unit. Replica bounds and the scaler
  # dials carry documented defaults resolved in variables.tf (the
  # platform never materializes proto defaults).
  template {
    min_replicas                     = var.spec.min_replicas
    max_replicas                     = var.spec.max_replicas
    cooldown_period_in_seconds       = var.spec.cooldown_period_in_seconds
    polling_interval_in_seconds      = var.spec.polling_interval_in_seconds
    revision_suffix                  = var.spec.revision_suffix
    termination_grace_period_seconds = var.spec.termination_grace_period_seconds

    dynamic "container" {
      for_each = var.spec.containers
      content {
        name   = container.value.name
        image  = container.value.image
        cpu    = container.value.cpu
        memory = container.value.memory
        # The provider expects command/args as attributes; empty lists
        # are sent as null so the revision hash stays stable.
        command = length(container.value.command) > 0 ? container.value.command : null
        args    = length(container.value.args) > 0 ? container.value.args : null

        dynamic "env" {
          for_each = container.value.env
          content {
            name        = env.value.name
            value       = env.value.value
            secret_name = env.value.secret_name
          }
        }

        # Health probes. The per-type contracts (success threshold is
        # readiness-only, per-type failure ceilings) are spec-enforced;
        # each type gets Azure's own initial-delay default when unset
        # (1 for liveness, 0 for readiness/startup).
        dynamic "liveness_probe" {
          for_each = container.value.liveness_probe != null ? [container.value.liveness_probe] : []
          content {
            transport               = local.probe_transport_map[liveness_probe.value.transport]
            port                    = liveness_probe.value.port
            path                    = liveness_probe.value.path
            host                    = liveness_probe.value.host
            initial_delay           = coalesce(liveness_probe.value.initial_delay_in_seconds, 1)
            interval_seconds        = liveness_probe.value.interval_seconds
            timeout                 = liveness_probe.value.timeout_seconds
            failure_count_threshold = liveness_probe.value.failure_count_threshold

            dynamic "header" {
              for_each = liveness_probe.value.headers
              content {
                name  = header.value.name
                value = header.value.value
              }
            }
          }
        }

        dynamic "readiness_probe" {
          for_each = container.value.readiness_probe != null ? [container.value.readiness_probe] : []
          content {
            transport               = local.probe_transport_map[readiness_probe.value.transport]
            port                    = readiness_probe.value.port
            path                    = readiness_probe.value.path
            host                    = readiness_probe.value.host
            initial_delay           = coalesce(readiness_probe.value.initial_delay_in_seconds, 0)
            interval_seconds        = readiness_probe.value.interval_seconds
            timeout                 = readiness_probe.value.timeout_seconds
            failure_count_threshold = readiness_probe.value.failure_count_threshold
            success_count_threshold = readiness_probe.value.success_count_threshold

            dynamic "header" {
              for_each = readiness_probe.value.headers
              content {
                name  = header.value.name
                value = header.value.value
              }
            }
          }
        }

        dynamic "startup_probe" {
          for_each = container.value.startup_probe != null ? [container.value.startup_probe] : []
          content {
            transport               = local.probe_transport_map[startup_probe.value.transport]
            port                    = startup_probe.value.port
            path                    = startup_probe.value.path
            host                    = startup_probe.value.host
            initial_delay           = coalesce(startup_probe.value.initial_delay_in_seconds, 0)
            interval_seconds        = startup_probe.value.interval_seconds
            timeout                 = startup_probe.value.timeout_seconds
            failure_count_threshold = startup_probe.value.failure_count_threshold

            dynamic "header" {
              for_each = startup_probe.value.headers
              content {
                name  = header.value.name
                value = header.value.value
              }
            }
          }
        }

        dynamic "volume_mounts" {
          for_each = container.value.volume_mounts
          content {
            name     = volume_mounts.value.name
            path     = volume_mounts.value.path
            sub_path = volume_mounts.value.sub_path
          }
        }
      }
    }

    # Init containers run to completion before main containers start;
    # cpu/memory omitted inherit the app's overall allocation.
    dynamic "init_container" {
      for_each = var.spec.init_containers
      content {
        name    = init_container.value.name
        image   = init_container.value.image
        cpu     = init_container.value.cpu
        memory  = init_container.value.memory
        command = length(init_container.value.command) > 0 ? init_container.value.command : null
        args    = length(init_container.value.args) > 0 ? init_container.value.args : null

        dynamic "env" {
          for_each = init_container.value.env
          content {
            name        = env.value.name
            value       = env.value.value
            secret_name = env.value.secret_name
          }
        }

        dynamic "volume_mounts" {
          for_each = init_container.value.volume_mounts
          content {
            name     = volume_mounts.value.name
            path     = volume_mounts.value.path
            sub_path = volume_mounts.value.sub_path
          }
        }
      }
    }

    # Volumes: EmptyDir scratch, share-backed AzureFile/NfsAzureFile via
    # an environment storage registration, or Secret mounts of the app's
    # secrets. Unspecified storage_type deploys EmptyDir.
    dynamic "volume" {
      for_each = var.spec.volumes
      content {
        name          = volume.value.name
        storage_type  = volume.value.storage_type != null ? local.volume_storage_type_map[volume.value.storage_type] : "EmptyDir"
        storage_name  = volume.value.storage_name
        mount_options = volume.value.mount_options
      }
    }

    dynamic "http_scale_rule" {
      for_each = var.spec.http_scale_rules
      content {
        name                = http_scale_rule.value.name
        concurrent_requests = http_scale_rule.value.concurrent_requests

        dynamic "authentication" {
          for_each = http_scale_rule.value.authentication
          content {
            secret_name       = authentication.value.secret_name
            trigger_parameter = authentication.value.trigger_parameter
          }
        }
      }
    }

    dynamic "tcp_scale_rule" {
      for_each = var.spec.tcp_scale_rules
      content {
        name                = tcp_scale_rule.value.name
        concurrent_requests = tcp_scale_rule.value.concurrent_requests

        dynamic "authentication" {
          for_each = tcp_scale_rule.value.authentication
          content {
            secret_name       = authentication.value.secret_name
            trigger_parameter = authentication.value.trigger_parameter
          }
        }
      }
    }

    dynamic "azure_queue_scale_rule" {
      for_each = var.spec.azure_queue_scale_rules
      content {
        name         = azure_queue_scale_rule.value.name
        queue_name   = azure_queue_scale_rule.value.queue_name
        queue_length = azure_queue_scale_rule.value.queue_length

        dynamic "authentication" {
          for_each = azure_queue_scale_rule.value.authentication
          content {
            secret_name       = authentication.value.secret_name
            trigger_parameter = authentication.value.trigger_parameter
          }
        }
      }
    }

    dynamic "custom_scale_rule" {
      for_each = var.spec.custom_scale_rules
      content {
        name             = custom_scale_rule.value.name
        custom_rule_type = custom_scale_rule.value.custom_rule_type
        metadata         = custom_scale_rule.value.metadata
        # Workload identity for the scaler instead of connection-string
        # secrets ("System" or a user-assigned identity ARM id; foreign-key
        # references arrive pre-resolved to the literal id).
        identity_id = custom_scale_rule.value.identity_id

        dynamic "authentication" {
          for_each = custom_scale_rule.value.authentication
          content {
            secret_name       = authentication.value.secret_name
            trigger_parameter = authentication.value.trigger_parameter
          }
        }
      }
    }
  }

  # App secrets: plain values or Key Vault references read through a
  # managed identity (the pairing is spec-enforced).
  dynamic "secret" {
    for_each = var.spec.secrets
    content {
      name                = secret.value.name
      value               = secret.value.value
      key_vault_secret_id = secret.value.key_vault_secret_id
      identity            = secret.value.identity
    }
  }

  # Private registry credentials: exactly one auth mode (spec-enforced).
  dynamic "registry" {
    for_each = var.spec.registries
    content {
      server               = registry.value.server
      username             = registry.value.username
      password_secret_name = registry.value.password_secret_name
      identity             = registry.value.identity
    }
  }

  # Ingress: without it the app is only reachable inside the environment
  # via service discovery.
  dynamic "ingress" {
    for_each = var.spec.ingress != null ? [var.spec.ingress] : []
    content {
      external_enabled           = ingress.value.external_enabled
      target_port                = ingress.value.target_port
      exposed_port               = ingress.value.exposed_port
      transport                  = ingress.value.transport != null ? local.ingress_transport_map[ingress.value.transport] : "auto"
      allow_insecure_connections = ingress.value.allow_insecure_connections
      # Sent only when chosen: unset leaves Azure's default behavior
      # (no client certificate requirement).
      client_certificate_mode = ingress.value.client_certificate_mode != null ? local.client_certificate_mode_map[ingress.value.client_certificate_mode] : null

      dynamic "traffic_weight" {
        for_each = ingress.value.traffic_weight
        content {
          latest_revision = traffic_weight.value.latest_revision
          revision_suffix = traffic_weight.value.revision_suffix
          percentage      = traffic_weight.value.percentage
          label           = traffic_weight.value.label
        }
      }

      dynamic "ip_security_restriction" {
        for_each = ingress.value.ip_security_restrictions
        content {
          name             = ip_security_restriction.value.name
          action           = local.ip_restriction_action_map[ip_security_restriction.value.action]
          ip_address_range = ip_security_restriction.value.ip_address_range
          description      = ip_security_restriction.value.description
        }
      }

      # CORS for browser-based clients.
      dynamic "cors" {
        for_each = ingress.value.cors != null ? [ingress.value.cors] : []
        content {
          allowed_origins           = cors.value.allowed_origins
          allowed_headers           = length(cors.value.allowed_headers) > 0 ? cors.value.allowed_headers : null
          allowed_methods           = length(cors.value.allowed_methods) > 0 ? cors.value.allowed_methods : null
          exposed_headers           = length(cors.value.exposed_headers) > 0 ? cors.value.exposed_headers : null
          max_age_in_seconds        = cors.value.max_age_in_seconds
          allow_credentials_enabled = cors.value.allow_credentials_enabled
        }
      }
    }
  }

  # Dapr sidecar; components are registered on the environment and scoped
  # to this app_id.
  dynamic "dapr" {
    for_each = var.spec.dapr != null ? [var.spec.dapr] : []
    content {
      app_id       = dapr.value.app_id
      app_port     = dapr.value.app_port
      app_protocol = dapr.value.app_protocol != null ? local.dapr_protocol_map[dapr.value.app_protocol] : "http"
    }
  }

  # Managed identity. The spec's CEL guarantees identity ids are present
  # exactly when the type includes UserAssigned.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = identity.value.user_assigned_identity_ids
    }
  }

  tags = local.final_tags
}
