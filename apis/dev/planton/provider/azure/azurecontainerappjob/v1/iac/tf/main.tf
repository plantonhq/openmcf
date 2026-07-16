# The Container App Job: a run-to-completion workload (batch, scheduled
# tasks, queue workers) inside a Container App Environment. Name, region,
# resource group, environment, and the trigger choice are all ForceNew.
resource "azurerm_container_app_job" "main" {
  name                         = var.spec.job_name
  location                     = var.spec.region
  resource_group_name          = var.spec.resource_group
  container_app_environment_id = var.spec.container_app_environment_id

  # Hard per-replica deadline; a replica killed by it counts as failed
  # (and retries per replica_retry_limit).
  replica_timeout_in_seconds = var.spec.replica_timeout_in_seconds
  replica_retry_limit        = var.spec.replica_retry_limit

  # Omitted runs on the environment's serverless Consumption profile.
  workload_profile_name = var.spec.workload_profile_name

  # The template each execution runs: containers, init containers, and
  # volumes (no scale block -- executions are the scaling unit).
  template {
    dynamic "container" {
      for_each = var.spec.containers
      content {
        name    = container.value.name
        image   = container.value.image
        cpu     = container.value.cpu
        memory  = container.value.memory
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

    dynamic "volume" {
      for_each = var.spec.volumes
      content {
        name          = volume.value.name
        storage_type  = volume.value.storage_type != null ? local.volume_storage_type_map[volume.value.storage_type] : "EmptyDir"
        storage_name  = volume.value.storage_name
        mount_options = volume.value.mount_options
      }
    }
  }

  # Exactly one trigger (spec-enforced); switching types is ForceNew.
  dynamic "manual_trigger_config" {
    for_each = var.spec.manual_trigger != null ? [var.spec.manual_trigger] : []
    content {
      parallelism              = manual_trigger_config.value.parallelism
      replica_completion_count = manual_trigger_config.value.replica_completion_count
    }
  }

  dynamic "schedule_trigger_config" {
    for_each = var.spec.schedule_trigger != null ? [var.spec.schedule_trigger] : []
    content {
      cron_expression          = schedule_trigger_config.value.cron_expression
      parallelism              = schedule_trigger_config.value.parallelism
      replica_completion_count = schedule_trigger_config.value.replica_completion_count
    }
  }

  dynamic "event_trigger_config" {
    for_each = var.spec.event_trigger != null ? [var.spec.event_trigger] : []
    content {
      parallelism              = event_trigger_config.value.parallelism
      replica_completion_count = event_trigger_config.value.replica_completion_count

      # The scale contract: how event pressure fans out into executions.
      dynamic "scale" {
        for_each = event_trigger_config.value.scale != null ? [event_trigger_config.value.scale] : []
        content {
          max_executions              = scale.value.max_executions
          min_executions              = scale.value.min_executions
          polling_interval_in_seconds = scale.value.polling_interval_in_seconds

          dynamic "rules" {
            for_each = scale.value.rules
            content {
              name             = rules.value.name
              custom_rule_type = rules.value.custom_rule_type
              metadata         = rules.value.metadata
              # Workload identity for the scaler instead of
              # connection-string secrets (foreign-key references arrive
              # pre-resolved to the literal id).
              identity_id = rules.value.identity_id

              dynamic "authentication" {
                for_each = rules.value.authentication
                content {
                  secret_name       = authentication.value.secret_name
                  trigger_parameter = authentication.value.trigger_parameter
                }
              }
            }
          }
        }
      }
    }
  }

  # Job secrets: plain values or Key Vault references read through a
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
