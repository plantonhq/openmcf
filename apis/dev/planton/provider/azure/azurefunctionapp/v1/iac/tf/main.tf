# The Linux Function App: event-driven compute on an App Service Plan.
# Name, region, and resource group are ForceNew; moving between Dynamic
# (Consumption) and other tiers also forces recreation.
resource "azurerm_linux_function_app" "main" {
  name                = var.spec.function_app_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  service_plan_id     = var.spec.service_plan_id

  # The runtime-state storage binding: exactly one of account name or
  # Key Vault secret ID applies (spec-enforced); the access key and
  # managed identity are mutually exclusive authentication methods.
  storage_account_name          = var.spec.storage_account_name
  storage_account_access_key    = var.spec.storage_account_access_key
  storage_uses_managed_identity = var.spec.storage_uses_managed_identity ? true : null
  storage_key_vault_secret_id   = var.spec.storage_key_vault_secret_id

  functions_extension_version = var.spec.functions_extension_version

  # The Consumption-plan cost circuit breaker (GB-seconds per day; 0 =
  # unlimited).
  daily_memory_time_quota = var.spec.daily_memory_time_quota

  enabled                       = var.spec.enabled
  https_only                    = var.spec.https_only
  public_network_access_enabled = var.spec.public_network_access_enabled
  builtin_logging_enabled       = var.spec.builtin_logging_enabled

  # VNet integration. Image pulls and backup/restore traffic only ride
  # the VNet when their dedicated toggles say so (spec-gated to require
  # the subnet).
  virtual_network_subnet_id              = var.spec.virtual_network_subnet_id
  vnet_image_pull_enabled                = var.spec.vnet_image_pull_enabled
  virtual_network_backup_restore_enabled = var.spec.virtual_network_backup_restore_enabled

  key_vault_reference_identity_id    = var.spec.key_vault_reference_identity_id
  client_certificate_enabled         = var.spec.client_certificate_enabled
  client_certificate_mode            = local.client_certificate_mode
  client_certificate_exclusion_paths = var.spec.client_certificate_exclusion_paths
  content_share_force_disabled       = var.spec.content_share_force_disabled

  # Basic-auth publishing toggles: disabling both closes the classic
  # credential-based deployment paths (Web Deploy + FTP) and forces
  # identity-based deployment.
  webdeploy_publish_basic_authentication_enabled = var.spec.webdeploy_publish_basic_authentication_enabled
  ftp_publish_basic_authentication_enabled       = var.spec.ftp_publish_basic_authentication_enabled

  zip_deploy_file = var.spec.zip_deploy_file

  app_settings = var.spec.app_settings

  tags = local.final_tags

  # ---------------------------------------------------------------------------
  # Site Config
  # ---------------------------------------------------------------------------
  site_config {
    always_on                         = var.spec.site_config.always_on
    app_command_line                  = var.spec.site_config.app_command_line
    api_management_api_id             = var.spec.site_config.api_management_api_id
    api_definition_url                = var.spec.site_config.api_definition_url
    default_documents                 = length(var.spec.site_config.default_documents) > 0 ? var.spec.site_config.default_documents : null
    health_check_path                 = var.spec.site_config.health_check_path
    health_check_eviction_time_in_min = var.spec.site_config.health_check_eviction_time_in_min
    minimum_tls_version               = local.minimum_tls_version
    scm_minimum_tls_version           = local.scm_minimum_tls_version
    # Absent accepts Azure's platform default cipher set.
    minimum_tls_cipher_suite = var.spec.site_config.minimum_tls_cipher_suite

    # Serverless scaling dials (Consumption / Elastic Premium).
    app_scale_limit           = var.spec.site_config.app_scale_limit
    elastic_instance_minimum  = var.spec.site_config.elastic_instance_minimum
    pre_warmed_instance_count = var.spec.site_config.pre_warmed_instance_count

    worker_count                                   = var.spec.site_config.worker_count
    http2_enabled                                  = var.spec.site_config.http2_enabled
    websockets_enabled                             = var.spec.site_config.websockets_enabled
    use_32_bit_worker                              = var.spec.site_config.use_32_bit_worker
    vnet_route_all_enabled                         = var.spec.site_config.vnet_route_all_enabled
    ftps_state                                     = local.ftps_state
    load_balancing_mode                            = local.load_balancing_mode
    managed_pipeline_mode                          = local.managed_pipeline_mode
    remote_debugging_enabled                       = var.spec.site_config.remote_debugging_enabled
    runtime_scale_monitoring_enabled               = var.spec.site_config.runtime_scale_monitoring_enabled
    ip_restriction_default_action                  = local.ip_restriction_default_action
    scm_use_main_ip_restriction                    = var.spec.site_config.scm_use_main_ip_restriction
    scm_ip_restriction_default_action              = local.scm_ip_restriction_default_action
    container_registry_use_managed_identity        = var.spec.site_config.container_registry_use_managed_identity
    container_registry_managed_identity_client_id  = var.spec.site_config.container_registry_managed_identity_client_id
    application_insights_key                       = var.spec.site_config.application_insights_key
    application_insights_connection_string         = var.spec.application_insights_connection_string

    # --- Application Stack ---
    dynamic "application_stack" {
      for_each = var.spec.site_config.application_stack != null ? [var.spec.site_config.application_stack] : []
      content {
        dotnet_version = application_stack.value.dotnet_version
        # Only meaningful alongside dotnet_version -- the provider
        # rejects it when another runtime is selected.
        use_dotnet_isolated_runtime = application_stack.value.dotnet_version != null ? application_stack.value.use_dotnet_isolated_runtime : null
        node_version                = application_stack.value.node_version
        python_version              = application_stack.value.python_version
        java_version                = application_stack.value.java_version
        powershell_core_version     = application_stack.value.powershell_core_version
        use_custom_runtime          = application_stack.value.use_custom_runtime

        dynamic "docker" {
          for_each = application_stack.value.docker != null ? [application_stack.value.docker] : []
          content {
            registry_url      = docker.value.registry_url
            image_name        = docker.value.image_name
            image_tag         = docker.value.image_tag
            registry_username = docker.value.registry_username
            registry_password = docker.value.registry_password
          }
        }
      }
    }

    # --- CORS ---
    dynamic "cors" {
      for_each = var.spec.site_config.cors != null ? [var.spec.site_config.cors] : []
      content {
        allowed_origins     = cors.value.allowed_origins
        support_credentials = cors.value.support_credentials
      }
    }

    # --- IP Restrictions (main site) ---
    dynamic "ip_restriction" {
      for_each = var.spec.site_config.ip_restrictions
      content {
        name                      = ip_restriction.value.name
        priority                  = ip_restriction.value.priority
        action                    = local.ip_restriction_action_map[coalesce(ip_restriction.value.action, "ALLOW")]
        ip_address                = ip_restriction.value.ip_address
        service_tag               = ip_restriction.value.service_tag
        virtual_network_subnet_id = ip_restriction.value.virtual_network_subnet_id
        description               = ip_restriction.value.description

        dynamic "headers" {
          for_each = ip_restriction.value.headers != null ? [ip_restriction.value.headers] : []
          content {
            x_forwarded_for  = headers.value.x_forwarded_for
            x_forwarded_host = headers.value.x_forwarded_host
            # Each x_azure_fdid entry is a StringValueOrRef in the spec
            # (referencing AzureFrontDoorProfile.resource_guid by
            # default); the tfvars converter flattens the list to the
            # resolved GUID literals.
            x_azure_fdid      = headers.value.x_azure_fdid
            x_fd_health_probe = headers.value.x_fd_health_probe
          }
        }
      }
    }

    # --- SCM IP Restrictions ---
    dynamic "scm_ip_restriction" {
      for_each = var.spec.site_config.scm_ip_restrictions
      content {
        name                      = scm_ip_restriction.value.name
        priority                  = scm_ip_restriction.value.priority
        action                    = local.ip_restriction_action_map[coalesce(scm_ip_restriction.value.action, "ALLOW")]
        ip_address                = scm_ip_restriction.value.ip_address
        service_tag               = scm_ip_restriction.value.service_tag
        virtual_network_subnet_id = scm_ip_restriction.value.virtual_network_subnet_id
        description               = scm_ip_restriction.value.description

        dynamic "headers" {
          for_each = scm_ip_restriction.value.headers != null ? [scm_ip_restriction.value.headers] : []
          content {
            x_forwarded_for   = headers.value.x_forwarded_for
            x_forwarded_host  = headers.value.x_forwarded_host
            x_azure_fdid      = headers.value.x_azure_fdid
            x_fd_health_probe = headers.value.x_fd_health_probe
          }
        }
      }
    }

    # --- App Service Logs ---
    dynamic "app_service_logs" {
      for_each = var.spec.site_config.app_service_logs != null ? [var.spec.site_config.app_service_logs] : []
      content {
        disk_quota_mb         = app_service_logs.value.disk_quota_mb
        retention_period_days = app_service_logs.value.retention_period_days
      }
    }
  }

  # ---------------------------------------------------------------------------
  # Backup (Standard tier or above -- Azure rejects it elsewhere at apply time)
  # ---------------------------------------------------------------------------
  dynamic "backup" {
    for_each = var.spec.backup != null ? [var.spec.backup] : []
    content {
      name                = backup.value.name
      storage_account_url = backup.value.storage_account_url
      enabled             = backup.value.enabled

      schedule {
        frequency_interval       = backup.value.schedule.frequency_interval
        frequency_unit           = local.backup_frequency_unit_map[backup.value.schedule.frequency_unit]
        keep_at_least_one_backup = backup.value.schedule.keep_at_least_one_backup
        retention_period_days    = backup.value.schedule.retention_period_days
        start_time               = backup.value.schedule.start_time
      }
    }
  }

  # ---------------------------------------------------------------------------
  # Sticky settings (pinned to the production slot during swaps)
  # ---------------------------------------------------------------------------
  dynamic "sticky_settings" {
    for_each = var.spec.sticky_settings != null ? [var.spec.sticky_settings] : []
    content {
      app_setting_names       = length(sticky_settings.value.app_setting_names) > 0 ? sticky_settings.value.app_setting_names : null
      connection_string_names = length(sticky_settings.value.connection_string_names) > 0 ? sticky_settings.value.connection_string_names : null
    }
  }

  # ---------------------------------------------------------------------------
  # Connection Strings
  # ---------------------------------------------------------------------------
  dynamic "connection_string" {
    for_each = var.spec.connection_strings
    content {
      name  = connection_string.value.name
      type  = local.connection_string_type_map[connection_string.value.type]
      value = connection_string.value.value
    }
  }

  # ---------------------------------------------------------------------------
  # Storage Account Mounts
  # ---------------------------------------------------------------------------
  dynamic "storage_account" {
    for_each = var.spec.storage_mounts
    content {
      name         = storage_account.value.name
      type         = local.storage_mount_type_map[storage_account.value.type]
      account_name = storage_account.value.account_name
      share_name   = storage_account.value.share_name
      access_key   = storage_account.value.access_key
      mount_path   = storage_account.value.mount_path
    }
  }

  # ---------------------------------------------------------------------------
  # Identity
  # ---------------------------------------------------------------------------
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  # ---------------------------------------------------------------------------
  # Easy Auth v2 (App Service built-in authentication)
  # ---------------------------------------------------------------------------
  # Provider client secrets are referenced by APP SETTING NAME -- the
  # secret values live in app_settings (or Key Vault references), never
  # in this block.
  dynamic "auth_settings_v2" {
    for_each = var.spec.auth_settings_v2 != null ? [var.spec.auth_settings_v2] : []
    content {
      auth_enabled                            = auth_settings_v2.value.auth_enabled
      runtime_version                         = auth_settings_v2.value.runtime_version
      config_file_path                        = auth_settings_v2.value.config_file_path
      require_authentication                  = auth_settings_v2.value.require_authentication
      unauthenticated_action                  = local.unauthenticated_action_map[coalesce(auth_settings_v2.value.unauthenticated_action, "REDIRECT_TO_LOGIN_PAGE")]
      default_provider                        = auth_settings_v2.value.default_provider
      excluded_paths                          = auth_settings_v2.value.excluded_paths
      require_https                           = auth_settings_v2.value.require_https
      http_route_api_prefix                   = auth_settings_v2.value.http_route_api_prefix
      forward_proxy_convention                = local.forward_proxy_convention_map[coalesce(auth_settings_v2.value.forward_proxy_convention, "FORWARD_PROXY_NO_PROXY")]
      forward_proxy_custom_host_header_name   = auth_settings_v2.value.forward_proxy_custom_host_header_name
      forward_proxy_custom_scheme_header_name = auth_settings_v2.value.forward_proxy_custom_scheme_header_name

      login {
        logout_endpoint                   = auth_settings_v2.value.login.logout_endpoint
        token_store_enabled               = auth_settings_v2.value.login.token_store_enabled
        token_refresh_extension_time      = auth_settings_v2.value.login.token_refresh_extension_time
        token_store_path                  = auth_settings_v2.value.login.token_store_path
        token_store_sas_setting_name      = auth_settings_v2.value.login.token_store_sas_setting_name
        preserve_url_fragments_for_logins = auth_settings_v2.value.login.preserve_url_fragments_for_logins
        allowed_external_redirect_urls    = auth_settings_v2.value.login.allowed_external_redirect_urls
        cookie_expiration_convention      = local.cookie_expiration_convention_map[coalesce(auth_settings_v2.value.login.cookie_expiration_convention, "FIXED_TIME")]
        cookie_expiration_time            = auth_settings_v2.value.login.cookie_expiration_time
        validate_nonce                    = auth_settings_v2.value.login.validate_nonce
        nonce_expiration_time             = auth_settings_v2.value.login.nonce_expiration_time
      }

      dynamic "apple_v2" {
        for_each = auth_settings_v2.value.apple_v2 != null ? [auth_settings_v2.value.apple_v2] : []
        content {
          client_id                  = apple_v2.value.client_id
          client_secret_setting_name = apple_v2.value.client_secret_setting_name
        }
      }

      dynamic "active_directory_v2" {
        for_each = auth_settings_v2.value.active_directory_v2 != null ? [auth_settings_v2.value.active_directory_v2] : []
        content {
          client_id                            = active_directory_v2.value.client_id
          tenant_auth_endpoint                 = active_directory_v2.value.tenant_auth_endpoint
          client_secret_setting_name           = active_directory_v2.value.client_secret_setting_name
          client_secret_certificate_thumbprint = active_directory_v2.value.client_secret_certificate_thumbprint
          login_parameters                     = active_directory_v2.value.login_parameters
          www_authentication_disabled          = active_directory_v2.value.www_authentication_disabled
          jwt_allowed_groups                   = active_directory_v2.value.jwt_allowed_groups
          jwt_allowed_client_applications      = active_directory_v2.value.jwt_allowed_client_applications
          allowed_groups                       = active_directory_v2.value.allowed_groups
          allowed_identities                   = active_directory_v2.value.allowed_identities
          allowed_applications                 = active_directory_v2.value.allowed_applications
          allowed_audiences                    = active_directory_v2.value.allowed_audiences
        }
      }

      dynamic "azure_static_web_app_v2" {
        for_each = auth_settings_v2.value.azure_static_web_app_v2 != null ? [auth_settings_v2.value.azure_static_web_app_v2] : []
        content {
          client_id = azure_static_web_app_v2.value.client_id
        }
      }

      dynamic "custom_oidc_v2" {
        for_each = auth_settings_v2.value.custom_oidc_v2
        content {
          name                          = custom_oidc_v2.value.name
          client_id                     = custom_oidc_v2.value.client_id
          openid_configuration_endpoint = custom_oidc_v2.value.openid_configuration_endpoint
          name_claim_type               = custom_oidc_v2.value.name_claim_type
          scopes                        = custom_oidc_v2.value.scopes
        }
      }

      dynamic "facebook_v2" {
        for_each = auth_settings_v2.value.facebook_v2 != null ? [auth_settings_v2.value.facebook_v2] : []
        content {
          app_id                  = facebook_v2.value.app_id
          app_secret_setting_name = facebook_v2.value.app_secret_setting_name
          graph_api_version       = facebook_v2.value.graph_api_version
          login_scopes            = facebook_v2.value.login_scopes
        }
      }

      dynamic "github_v2" {
        for_each = auth_settings_v2.value.github_v2 != null ? [auth_settings_v2.value.github_v2] : []
        content {
          client_id                  = github_v2.value.client_id
          client_secret_setting_name = github_v2.value.client_secret_setting_name
          login_scopes               = github_v2.value.login_scopes
        }
      }

      dynamic "google_v2" {
        for_each = auth_settings_v2.value.google_v2 != null ? [auth_settings_v2.value.google_v2] : []
        content {
          client_id                  = google_v2.value.client_id
          client_secret_setting_name = google_v2.value.client_secret_setting_name
          allowed_audiences          = google_v2.value.allowed_audiences
          login_scopes               = google_v2.value.login_scopes
        }
      }

      dynamic "microsoft_v2" {
        for_each = auth_settings_v2.value.microsoft_v2 != null ? [auth_settings_v2.value.microsoft_v2] : []
        content {
          client_id                  = microsoft_v2.value.client_id
          client_secret_setting_name = microsoft_v2.value.client_secret_setting_name
          allowed_audiences          = microsoft_v2.value.allowed_audiences
          login_scopes               = microsoft_v2.value.login_scopes
        }
      }

      dynamic "twitter_v2" {
        for_each = auth_settings_v2.value.twitter_v2 != null ? [auth_settings_v2.value.twitter_v2] : []
        content {
          consumer_key                 = twitter_v2.value.consumer_key
          consumer_secret_setting_name = twitter_v2.value.consumer_secret_setting_name
        }
      }
    }
  }
}
