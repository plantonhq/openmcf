# The Flex Consumption Function App: Azure's newest serverless Functions
# hosting model, on an FC1-SKU plan (the provider verifies the plan's SKU
# at create time and rejects anything else). Name, region, resource
# group, and service plan are ForceNew.
resource "azurerm_function_app_flex_consumption" "main" {
  name                = var.spec.function_app_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  service_plan_id     = var.spec.service_plan_id

  # The deployment storage binding. blobContainer is the type's single
  # legal value, so the spec omits it and the module sends the constant.
  storage_container_type     = "blobContainer"
  storage_container_endpoint = var.spec.storage_container_endpoint

  # The per-mode requiredness (key with connection-string auth, identity
  # id with user-assigned auth) is spec-enforced, mirroring the
  # provider's own create-time checks. Empty optionals become null so the
  # provider's not-empty validators never see a present-but-empty value.
  storage_authentication_type       = local.storage_authentication_type
  storage_access_key                = var.spec.storage_access_key != "" ? var.spec.storage_access_key : null
  storage_user_assigned_identity_id = var.spec.storage_user_assigned_identity_id != "" ? var.spec.storage_user_assigned_identity_id : null

  # The flat runtime declaration (Flex Consumption has no
  # application_stack block and no container form).
  runtime_name    = local.runtime_name
  runtime_version = var.spec.runtime_version

  # Scale: per-instance memory, the fan-out ceiling, and optional
  # per-instance HTTP concurrency (absent lets Azure pick the runtime's
  # default for the memory size).
  instance_memory_in_mb  = var.spec.instance_memory_in_mb
  maximum_instance_count = var.spec.maximum_instance_count
  http_concurrency       = var.spec.http_concurrency

  # Always-ready pools: the counts' sum must stay within
  # maximum_instance_count (Azure enforces at apply time). Azure
  # lower-cases pool names on save.
  dynamic "always_ready" {
    for_each = var.spec.always_ready
    content {
      name           = always_ready.value.name
      instance_count = always_ready.value.instance_count
    }
  }

  enabled                       = var.spec.enabled
  https_only                    = var.spec.https_only
  public_network_access_enabled = var.spec.public_network_access_enabled

  virtual_network_subnet_id = var.spec.virtual_network_subnet_id != "" ? var.spec.virtual_network_subnet_id : null

  client_certificate_enabled         = var.spec.client_certificate_enabled
  client_certificate_mode            = local.client_certificate_mode
  client_certificate_exclusion_paths = var.spec.client_certificate_exclusion_paths != "" ? var.spec.client_certificate_exclusion_paths : null

  # Basic-auth publishing toggle: disabling it closes the classic
  # credential-based deployment path and forces identity-based
  # deployment.
  webdeploy_publish_basic_authentication_enabled = var.spec.webdeploy_publish_basic_authentication_enabled

  zip_deploy_file = var.spec.zip_deploy_file != "" ? var.spec.zip_deploy_file : null

  app_settings = var.spec.app_settings

  tags = local.final_tags

  # ---------------------------------------------------------------------------
  # Site Config
  # ---------------------------------------------------------------------------
  site_config {
    api_management_api_id = var.spec.site_config.api_management_api_id != "" ? var.spec.site_config.api_management_api_id : null
    api_definition_url    = var.spec.site_config.api_definition_url != "" ? var.spec.site_config.api_definition_url : null
    app_command_line      = var.spec.site_config.app_command_line != "" ? var.spec.site_config.app_command_line : null

    # Both App Insights values travel as app settings on the wire; the
    # connection string lives on the parent spec (a typed reference to
    # AzureApplicationInsights), the classic key here.
    application_insights_key               = var.spec.site_config.application_insights_key != "" ? var.spec.site_config.application_insights_key : null
    application_insights_connection_string = var.spec.application_insights_connection_string != "" ? var.spec.application_insights_connection_string : null

    container_registry_use_managed_identity = var.spec.site_config.container_registry_use_managed_identity

    default_documents = length(var.spec.site_config.default_documents) > 0 ? var.spec.site_config.default_documents : null

    # Accepted by Azure but never read back on this hosting model
    # (always_ready is the flex-native warm-instance mechanism).
    elastic_instance_minimum = var.spec.site_config.elastic_instance_minimum

    http2_enabled                     = var.spec.site_config.http2_enabled
    websockets_enabled                = var.spec.site_config.websockets_enabled
    vnet_route_all_enabled            = var.spec.site_config.vnet_route_all_enabled
    load_balancing_mode               = local.load_balancing_mode
    managed_pipeline_mode             = local.managed_pipeline_mode
    remote_debugging_enabled          = var.spec.site_config.remote_debugging_enabled
    remote_debugging_version          = var.spec.site_config.remote_debugging_version != "" ? var.spec.site_config.remote_debugging_version : null
    runtime_scale_monitoring_enabled  = var.spec.site_config.runtime_scale_monitoring_enabled
    health_check_path                 = var.spec.site_config.health_check_path != "" ? var.spec.site_config.health_check_path : null
    health_check_eviction_time_in_min = var.spec.site_config.health_check_eviction_time_in_min
    worker_count                      = var.spec.site_config.worker_count
    minimum_tls_version               = local.minimum_tls_version
    scm_minimum_tls_version           = local.scm_minimum_tls_version

    ip_restriction_default_action     = local.ip_restriction_default_action
    scm_use_main_ip_restriction       = var.spec.site_config.scm_use_main_ip_restriction
    scm_ip_restriction_default_action = local.scm_ip_restriction_default_action

    # --- CORS ---
    dynamic "cors" {
      for_each = var.spec.site_config.cors != null ? [var.spec.site_config.cors] : []
      content {
        allowed_origins     = cors.value.allowed_origins
        support_credentials = cors.value.support_credentials
      }
    }

    # --- IP Restrictions (main site) ---
    # The exactly-one-of ip/service_tag/subnet contract fires on argument
    # PRESENCE, so empty optionals become null.
    dynamic "ip_restriction" {
      for_each = var.spec.site_config.ip_restrictions
      content {
        name                      = ip_restriction.value.name
        priority                  = ip_restriction.value.priority
        action                    = local.ip_restriction_action_map[coalesce(ip_restriction.value.action, "ALLOW")]
        ip_address                = ip_restriction.value.ip_address != "" ? ip_restriction.value.ip_address : null
        service_tag               = ip_restriction.value.service_tag != "" ? ip_restriction.value.service_tag : null
        virtual_network_subnet_id = ip_restriction.value.virtual_network_subnet_id != "" ? ip_restriction.value.virtual_network_subnet_id : null
        description               = ip_restriction.value.description != "" ? ip_restriction.value.description : null

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
        ip_address                = scm_ip_restriction.value.ip_address != "" ? scm_ip_restriction.value.ip_address : null
        service_tag               = scm_ip_restriction.value.service_tag != "" ? scm_ip_restriction.value.service_tag : null
        virtual_network_subnet_id = scm_ip_restriction.value.virtual_network_subnet_id != "" ? scm_ip_restriction.value.virtual_network_subnet_id : null
        description               = scm_ip_restriction.value.description != "" ? scm_ip_restriction.value.description : null

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
    # Azure applies this block on update operations only and never
    # returns it on read.
    dynamic "app_service_logs" {
      for_each = var.spec.site_config.app_service_logs != null ? [var.spec.site_config.app_service_logs] : []
      content {
        disk_quota_mb         = app_service_logs.value.disk_quota_mb
        retention_period_days = app_service_logs.value.retention_period_days
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
      config_file_path                        = auth_settings_v2.value.config_file_path != "" ? auth_settings_v2.value.config_file_path : null
      require_authentication                  = auth_settings_v2.value.require_authentication
      unauthenticated_action                  = local.unauthenticated_action_map[coalesce(auth_settings_v2.value.unauthenticated_action, "REDIRECT_TO_LOGIN_PAGE")]
      default_provider                        = auth_settings_v2.value.default_provider != "" ? auth_settings_v2.value.default_provider : null
      excluded_paths                          = auth_settings_v2.value.excluded_paths
      require_https                           = auth_settings_v2.value.require_https
      http_route_api_prefix                   = auth_settings_v2.value.http_route_api_prefix
      forward_proxy_convention                = local.forward_proxy_convention_map[coalesce(auth_settings_v2.value.forward_proxy_convention, "FORWARD_PROXY_NO_PROXY")]
      forward_proxy_custom_host_header_name   = auth_settings_v2.value.forward_proxy_custom_host_header_name != "" ? auth_settings_v2.value.forward_proxy_custom_host_header_name : null
      forward_proxy_custom_scheme_header_name = auth_settings_v2.value.forward_proxy_custom_scheme_header_name != "" ? auth_settings_v2.value.forward_proxy_custom_scheme_header_name : null

      login {
        logout_endpoint              = auth_settings_v2.value.login.logout_endpoint != "" ? auth_settings_v2.value.login.logout_endpoint : null
        token_store_enabled          = auth_settings_v2.value.login.token_store_enabled
        token_refresh_extension_time = auth_settings_v2.value.login.token_refresh_extension_time
        # The two token-store backings are ConflictsWith partners in the
        # provider, and the conflict fires on argument PRESENCE -- empty
        # optionals must become null.
        token_store_path                  = auth_settings_v2.value.login.token_store_path != "" ? auth_settings_v2.value.login.token_store_path : null
        token_store_sas_setting_name      = auth_settings_v2.value.login.token_store_sas_setting_name != "" ? auth_settings_v2.value.login.token_store_sas_setting_name : null
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
          client_id            = active_directory_v2.value.client_id
          tenant_auth_endpoint = active_directory_v2.value.tenant_auth_endpoint
          # The two credential forms are ConflictsWith partners in the
          # provider, firing on argument PRESENCE -- empty optionals must
          # become null.
          client_secret_setting_name           = active_directory_v2.value.client_secret_setting_name != "" ? active_directory_v2.value.client_secret_setting_name : null
          client_secret_certificate_thumbprint = active_directory_v2.value.client_secret_certificate_thumbprint != "" ? active_directory_v2.value.client_secret_certificate_thumbprint : null
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
          name_claim_type               = custom_oidc_v2.value.name_claim_type != "" ? custom_oidc_v2.value.name_claim_type : null
          scopes                        = custom_oidc_v2.value.scopes
        }
      }

      dynamic "facebook_v2" {
        for_each = auth_settings_v2.value.facebook_v2 != null ? [auth_settings_v2.value.facebook_v2] : []
        content {
          app_id                  = facebook_v2.value.app_id
          app_secret_setting_name = facebook_v2.value.app_secret_setting_name
          graph_api_version       = facebook_v2.value.graph_api_version != "" ? facebook_v2.value.graph_api_version : null
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
