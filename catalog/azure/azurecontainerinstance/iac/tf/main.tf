# Create the container group. Almost the whole resource is
# create-only: after create, Azure applies only identity and tag
# changes in place -- anything else replaces the group.
#
# The spec's volume UNION (azure_file XOR empty_dir XOR git_repo XOR
# secret, validated before apply) flattens back to the provider's flat
# volume block here. The spec's singular probe http_get becomes the
# provider's one-element list (the provider schema accepts a list but
# its own code keeps only the last entry -- the wire carries one).
#
# Send-only-when-set throughout: unset optional scalars ride the
# provider defaults (ip_address_type "Public", restart_policy
# "Always", sku "Standard", dns_name_label_reuse_policy "Unsecure",
# port protocol "TCP").
resource "azurerm_container_group" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region
  os_type             = var.spec.os_type

  restart_policy = (
    var.spec.restart_policy != null && var.spec.restart_policy != ""
    ? var.spec.restart_policy
    : null
  )
  sku = (
    var.spec.sku != null && var.spec.sku != ""
    ? var.spec.sku
    : null
  )
  priority = (
    var.spec.priority != null && var.spec.priority != ""
    ? var.spec.priority
    : null
  )
  ip_address_type = (
    var.spec.ip_address_type != null && var.spec.ip_address_type != ""
    ? var.spec.ip_address_type
    : null
  )
  dns_name_label = (
    var.spec.dns_name_label != null && var.spec.dns_name_label != ""
    ? var.spec.dns_name_label
    : null
  )
  dns_name_label_reuse_policy = (
    var.spec.dns_name_label_reuse_policy != null && var.spec.dns_name_label_reuse_policy != ""
    ? var.spec.dns_name_label_reuse_policy
    : null
  )

  # The provider takes a set with at most one subnet; the spec models
  # the honest singular. Azure serializes container-group operations
  # per subnet.
  subnet_ids = (
    var.spec.subnet_id != null && var.spec.subnet_id != ""
    ? [var.spec.subnet_id]
    : null
  )

  zones = length(var.spec.zones) > 0 ? var.spec.zones : null

  # Group-level exposed ports (attribute mode in this provider). Omit
  # to expose every container port -- the provider derives the group
  # ports. Protocol normalizes to the provider default "TCP" so the
  # rendered set is deterministic.
  exposed_port = length(var.spec.exposed_ports) > 0 ? [
    for exposed in var.spec.exposed_ports : {
      port     = exposed.port
      protocol = exposed.protocol != "" ? exposed.protocol : "TCP"
    }
  ] : null

  # Customer-managed-key encryption: a VERSIONED Key Vault key
  # identifier (ACI pins the exact version -- rotation does not
  # follow). BEHAVIOR: the provider applies
  # key_vault_user_assigned_identity_id at CREATE only; a later change
  # to it alone is silently never applied (the provider's update path
  # covers only identity and tags) -- treat it as create-only.
  key_vault_key_id = (
    var.spec.key_vault_key_id != null && var.spec.key_vault_key_id != ""
    ? var.spec.key_vault_key_id
    : null
  )
  key_vault_user_assigned_identity_id = (
    var.spec.key_vault_user_assigned_identity_id != null && var.spec.key_vault_user_assigned_identity_id != ""
    ? var.spec.key_vault_user_assigned_identity_id
    : null
  )

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  dynamic "image_registry_credential" {
    for_each = var.spec.image_registry_credentials
    content {
      server = image_registry_credential.value.server
      username = (
        image_registry_credential.value.username != ""
        ? image_registry_credential.value.username
        : null
      )
      # Azure never returns the password on reads; the provider echoes
      # it from configuration (state), so it must stay present here.
      password = (
        image_registry_credential.value.password != ""
        ? image_registry_credential.value.password
        : null
      )
      user_assigned_identity_id = (
        image_registry_credential.value.user_assigned_identity_id != ""
        ? image_registry_credential.value.user_assigned_identity_id
        : null
      )
    }
  }

  dynamic "init_container" {
    for_each = var.spec.init_containers
    content {
      name  = init_container.value.name
      image = init_container.value.image

      environment_variables = (
        length(init_container.value.environment_variables) > 0
        ? init_container.value.environment_variables
        : null
      )
      # Azure never returns secure values on reads; the provider echoes
      # them from state.
      secure_environment_variables = (
        length(init_container.value.secure_environment_variables) > 0
        ? init_container.value.secure_environment_variables
        : null
      )
      commands = length(init_container.value.commands) > 0 ? init_container.value.commands : null

      dynamic "volume" {
        for_each = init_container.value.volumes
        content {
          name       = volume.value.name
          mount_path = volume.value.mount_path
          read_only  = volume.value.read_only

          # The spec's volume union flattened to the provider's flat
          # block -- exactly one form is populated (validated upstream).
          share_name           = volume.value.azure_file != null ? volume.value.azure_file.share_name : null
          storage_account_name = volume.value.azure_file != null ? volume.value.azure_file.storage_account_name : null
          # Azure never returns the storage key on reads; the provider
          # echoes it from configuration.
          storage_account_key = volume.value.azure_file != null ? volume.value.azure_file.storage_account_key : null

          empty_dir = volume.value.empty_dir
          secret    = length(volume.value.secret) > 0 ? volume.value.secret : null

          dynamic "git_repo" {
            for_each = volume.value.git_repo != null ? [volume.value.git_repo] : []
            content {
              url       = git_repo.value.url
              directory = git_repo.value.directory != "" ? git_repo.value.directory : null
              revision  = git_repo.value.revision != "" ? git_repo.value.revision : null
            }
          }
        }
      }

      dynamic "security" {
        for_each = init_container.value.security != null ? [init_container.value.security] : []
        content {
          privilege_enabled = security.value.privilege_enabled
        }
      }
    }
  }

  dynamic "container" {
    for_each = var.spec.containers
    content {
      name   = container.value.name
      image  = container.value.image
      cpu    = container.value.cpu
      memory = container.value.memory

      # BEHAVIOR: the provider applies the limits at CREATE only -- a
      # later change to either alone is silently never applied (the
      # provider's update path covers only identity and tags).
      cpu_limit    = container.value.cpu_limit
      memory_limit = container.value.memory_limit

      dynamic "ports" {
        for_each = container.value.ports
        content {
          port     = ports.value.port
          protocol = ports.value.protocol != "" ? ports.value.protocol : null
        }
      }

      environment_variables = (
        length(container.value.environment_variables) > 0
        ? container.value.environment_variables
        : null
      )
      # Azure never returns secure values on reads; the provider echoes
      # them from state.
      secure_environment_variables = (
        length(container.value.secure_environment_variables) > 0
        ? container.value.secure_environment_variables
        : null
      )

      # Unset lets the image's own entrypoint run -- reads echo it back
      # (expected, not drift).
      commands = length(container.value.commands) > 0 ? container.value.commands : null

      dynamic "volume" {
        for_each = container.value.volumes
        content {
          name       = volume.value.name
          mount_path = volume.value.mount_path
          read_only  = volume.value.read_only

          share_name           = volume.value.azure_file != null ? volume.value.azure_file.share_name : null
          storage_account_name = volume.value.azure_file != null ? volume.value.azure_file.storage_account_name : null
          storage_account_key  = volume.value.azure_file != null ? volume.value.azure_file.storage_account_key : null

          empty_dir = volume.value.empty_dir
          secret    = length(volume.value.secret) > 0 ? volume.value.secret : null

          dynamic "git_repo" {
            for_each = volume.value.git_repo != null ? [volume.value.git_repo] : []
            content {
              url       = git_repo.value.url
              directory = git_repo.value.directory != "" ? git_repo.value.directory : null
              revision  = git_repo.value.revision != "" ? git_repo.value.revision : null
            }
          }
        }
      }

      dynamic "security" {
        for_each = container.value.security != null ? [container.value.security] : []
        content {
          privilege_enabled = security.value.privilege_enabled
        }
      }

      dynamic "liveness_probe" {
        for_each = container.value.liveness_probe != null ? [container.value.liveness_probe] : []
        content {
          exec = length(liveness_probe.value.exec) > 0 ? liveness_probe.value.exec : null

          # The spec's singular http_get as the provider's one-element
          # list (the provider keeps only one on the wire).
          dynamic "http_get" {
            for_each = liveness_probe.value.http_get != null ? [liveness_probe.value.http_get] : []
            content {
              path = http_get.value.path != "" ? http_get.value.path : null
              port = http_get.value.port != 0 ? http_get.value.port : null
              # Explicit-send "http" when unset: ARM materializes the
              # scheme on reads and the provider treats it as
              # replace-forcing, so an omitted scheme re-plans as a
              # destroy+create (live-proven by the idempotency gate).
              scheme       = http_get.value.scheme != "" ? http_get.value.scheme : "http"
              http_headers = length(http_get.value.http_headers) > 0 ? http_get.value.http_headers : null
            }
          }

          # Zero means unset -- the provider sends timings only when > 0.
          initial_delay_seconds = liveness_probe.value.initial_delay_seconds != 0 ? liveness_probe.value.initial_delay_seconds : null
          period_seconds        = liveness_probe.value.period_seconds != 0 ? liveness_probe.value.period_seconds : null
          failure_threshold     = liveness_probe.value.failure_threshold != 0 ? liveness_probe.value.failure_threshold : null
          success_threshold     = liveness_probe.value.success_threshold != 0 ? liveness_probe.value.success_threshold : null
          timeout_seconds       = liveness_probe.value.timeout_seconds != 0 ? liveness_probe.value.timeout_seconds : null
        }
      }

      dynamic "readiness_probe" {
        for_each = container.value.readiness_probe != null ? [container.value.readiness_probe] : []
        content {
          exec = length(readiness_probe.value.exec) > 0 ? readiness_probe.value.exec : null

          dynamic "http_get" {
            for_each = readiness_probe.value.http_get != null ? [readiness_probe.value.http_get] : []
            content {
              path = http_get.value.path != "" ? http_get.value.path : null
              port = http_get.value.port != 0 ? http_get.value.port : null
              # Explicit-send "http" when unset -- see the liveness
              # probe's scheme note (replace-forcing echo).
              scheme       = http_get.value.scheme != "" ? http_get.value.scheme : "http"
              http_headers = length(http_get.value.http_headers) > 0 ? http_get.value.http_headers : null
            }
          }

          initial_delay_seconds = readiness_probe.value.initial_delay_seconds != 0 ? readiness_probe.value.initial_delay_seconds : null
          period_seconds        = readiness_probe.value.period_seconds != 0 ? readiness_probe.value.period_seconds : null
          failure_threshold     = readiness_probe.value.failure_threshold != 0 ? readiness_probe.value.failure_threshold : null
          success_threshold     = readiness_probe.value.success_threshold != 0 ? readiness_probe.value.success_threshold : null
          timeout_seconds       = readiness_probe.value.timeout_seconds != 0 ? readiness_probe.value.timeout_seconds : null
        }
      }
    }
  }

  # The provider's diagnostics block has exactly one member -- this Log
  # Analytics form (the spec collapses the wrapper level).
  dynamic "diagnostics" {
    for_each = var.spec.diagnostics_log_analytics != null ? [var.spec.diagnostics_log_analytics] : []
    content {
      log_analytics {
        # The workspace CUSTOMER ID (GUID), not the ARM resource ID.
        workspace_id = diagnostics.value.workspace_id
        # Azure never returns the key on reads; the provider echoes it
        # from configuration.
        workspace_key = diagnostics.value.workspace_key
        log_type      = diagnostics.value.log_type != "" ? diagnostics.value.log_type : null
        # The provider only sends metadata alongside a log type
        # (validated upstream). Note the provider ALSO attaches an empty
        # metadata object whenever a log type is set, which ARM rejects
        # for ContainerInstanceLogs (LogAnalyticsMetadataNotAllowed,
        # live-proven) -- a spec CEL blocks that log type until the
        # provider stops sending the empty map.
        metadata = length(diagnostics.value.metadata) > 0 ? diagnostics.value.metadata : null
      }
    }
  }

  dynamic "dns_config" {
    for_each = var.spec.dns_config != null ? [var.spec.dns_config] : []
    content {
      nameservers    = dns_config.value.nameservers
      search_domains = length(dns_config.value.search_domains) > 0 ? dns_config.value.search_domains : null
      options        = length(dns_config.value.options) > 0 ? dns_config.value.options : null
    }
  }

  tags = local.final_tags
}
