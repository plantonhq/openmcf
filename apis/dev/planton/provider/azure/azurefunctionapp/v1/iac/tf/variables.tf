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
  description = "Azure Function App specification"
  type = object({
    # The Azure region where the Function App will be created. ForceNew.
    region = string

    # The Azure Resource Group name. References are resolved to a literal
    # name by the platform before the module runs. ForceNew.
    resource_group = string

    # The name of the Function App (globally unique -- it forms
    # {function_app_name}.azurewebsites.net). ForceNew.
    function_app_name = string

    # The App Service Plan resource ID. References are resolved before
    # the module runs.
    service_plan_id = string

    # The runtime-state storage binding: exactly one of the account name
    # or the Key Vault secret ID (spec-enforced). The access key is a
    # resolved secret literal by the time the module runs.
    storage_account_name          = optional(string)
    storage_account_access_key    = optional(string)
    storage_uses_managed_identity = optional(bool, false)
    storage_key_vault_secret_id   = optional(string)

    # The Azure Functions host version (e.g. "~4").
    functions_extension_version = optional(string, "~4")

    # Daily compute quota in GB-seconds (Consumption-plan cost circuit
    # breaker; 0 = unlimited).
    daily_memory_time_quota = optional(number, 0)

    # Site configuration including runtime, scaling, security, and networking.
    site_config = object({
      # Application stack (runtime selection). Exactly one runtime. The
      # registry password inside docker is a resolved secret literal.
      application_stack = optional(object({
        dotnet_version              = optional(string)
        use_dotnet_isolated_runtime = optional(bool, false)
        node_version                = optional(string)
        python_version              = optional(string)
        java_version                = optional(string)
        powershell_core_version     = optional(string)
        docker = optional(object({
          registry_url      = string
          image_name        = string
          image_tag         = string
          registry_username = optional(string)
          registry_password = optional(string)
        }))
        use_custom_runtime = optional(bool)
      }))

      # Keep the app loaded in memory (Dedicated plans; auto-managed on
      # serverless tiers; rejected on Free/Shared at apply time).
      always_on = optional(bool)

      # Custom startup command.
      app_command_line = optional(string)

      # ARM ID of the API Management API this app backs.
      api_management_api_id = optional(string)

      # URL of the OpenAPI/Swagger definition for this app's API.
      api_definition_url = optional(string)

      # Default documents served for directory requests, in order.
      default_documents = optional(list(string), [])

      # Health check endpoint path.
      health_check_path = optional(string)

      # Time (in minutes) after which an unhealthy instance is evicted. 2-10.
      health_check_eviction_time_in_min = optional(number)

      # Minimum TLS versions, as the spec enum's name strings (TLS_1_0 /
      # TLS_1_1 / TLS_1_2 / TLS_1_3). Absent means TLS_1_2.
      minimum_tls_version     = optional(string)
      scm_minimum_tls_version = optional(string)

      # Minimum accepted TLS cipher suite (Azure's own identifiers).
      # Absent accepts Azure's platform default.
      minimum_tls_cipher_suite = optional(string)

      # Serverless scaling dials (Consumption / Elastic Premium).
      app_scale_limit           = optional(number)
      elastic_instance_minimum  = optional(number)
      pre_warmed_instance_count = optional(number)

      # Number of worker instances on Dedicated plans.
      worker_count = optional(number)

      # Enable HTTP/2 protocol.
      http2_enabled = optional(bool, false)

      # Enable WebSocket connections.
      websockets_enabled = optional(bool, false)

      # Use a 32-bit worker process.
      use_32_bit_worker = optional(bool, false)

      # Route all outbound traffic through VNet.
      vnet_route_all_enabled = optional(bool, false)

      # FTPS state, as the spec enum's name string (ALL_ALLOWED /
      # FTPS_ONLY / DISABLED). Absent means DISABLED.
      ftps_state = optional(string)

      # Load balancing mode, as the spec enum's name string. Absent
      # means LEAST_REQUESTS.
      load_balancing_mode = optional(string)

      # Request pipeline mode, as the spec enum's name string
      # (INTEGRATED / CLASSIC). Absent means INTEGRATED.
      managed_pipeline_mode = optional(string)

      # Remote debugging (Visual Studio attach).
      remote_debugging_enabled = optional(bool, false)

      # Runtime scale monitoring for KEDA-based triggers.
      runtime_scale_monitoring_enabled = optional(bool)

      # CORS configuration.
      cors = optional(object({
        allowed_origins     = list(string)
        support_credentials = optional(bool, false)
      }))

      # IP restriction rules for the main site. action arrives as the
      # spec enum's name string (ALLOW / DENY).
      ip_restrictions = optional(list(object({
        name                      = optional(string)
        priority                  = optional(number)
        action                    = optional(string)
        ip_address                = optional(string)
        service_tag               = optional(string)
        virtual_network_subnet_id = optional(string)
        description               = optional(string)
        headers = optional(object({
          x_forwarded_for   = optional(list(string), [])
          x_forwarded_host  = optional(list(string), [])
          x_azure_fdid      = optional(list(string), [])
          x_fd_health_probe = optional(list(string), [])
        }))
      })), [])

      # Default action for main-site IP restrictions. Absent means ALLOW.
      ip_restriction_default_action = optional(string)

      # Use main site IP restrictions for SCM site.
      scm_use_main_ip_restriction = optional(bool, false)

      # IP restriction rules for the SCM (Kudu) site.
      scm_ip_restrictions = optional(list(object({
        name                      = optional(string)
        priority                  = optional(number)
        action                    = optional(string)
        ip_address                = optional(string)
        service_tag               = optional(string)
        virtual_network_subnet_id = optional(string)
        description               = optional(string)
        headers = optional(object({
          x_forwarded_for   = optional(list(string), [])
          x_forwarded_host  = optional(list(string), [])
          x_azure_fdid      = optional(list(string), [])
          x_fd_health_probe = optional(list(string), [])
        }))
      })), [])

      # Default action for SCM-site IP restrictions. Absent means ALLOW.
      scm_ip_restriction_default_action = optional(string)

      # App Service logs (disk quota + retention).
      app_service_logs = optional(object({
        disk_quota_mb         = optional(number, 35)
        retention_period_days = optional(number)
      }))

      # Use managed identity for ACR image pulls.
      container_registry_use_managed_identity = optional(bool, false)

      # Client ID of the managed identity for ACR pulls.
      container_registry_managed_identity_client_id = optional(string)

      # Classic Application Insights instrumentation key (resolved
      # secret literal); prefer the connection string on the spec.
      application_insights_key = optional(string)
    })

    # Application settings (environment variables) as key-value pairs.
    app_settings = optional(map(string), {})

    # Named connection strings. type arrives as the spec enum's name
    # string (e.g. SQL_AZURE, CUSTOM); values are resolved secrets.
    connection_strings = optional(list(object({
      name  = string
      type  = string
      value = string
    })), [])

    # Settings pinned to the production slot during slot swaps.
    sticky_settings = optional(object({
      app_setting_names       = optional(list(string), [])
      connection_string_names = optional(list(string), [])
    }))

    # Application Insights connection string for APM telemetry.
    application_insights_connection_string = optional(string)

    # Enforce HTTPS-only access.
    https_only = optional(bool, true)

    # Enable public network access.
    public_network_access_enabled = optional(bool, true)

    # Whether the Function App is enabled. false stops the app without
    # deleting it.
    enabled = optional(bool, true)

    # Built-in logging via AzureWebJobsDashboard.
    builtin_logging_enabled = optional(bool, true)

    # Force disable the auto-created Azure Files content share.
    content_share_force_disabled = optional(bool, false)

    # Enable client certificate authentication (mTLS).
    client_certificate_enabled = optional(bool, false)

    # Client certificate mode, as the spec enum's name string (REQUIRED /
    # OPTIONAL / OPTIONAL_INTERACTIVE_USER). Absent means OPTIONAL.
    client_certificate_mode = optional(string)

    # Paths excluded from client certificate validation (semicolon-separated).
    client_certificate_exclusion_paths = optional(string)

    # Subnet ID for VNet integration.
    virtual_network_subnet_id = optional(string)

    # Pull container images over the VNet integration.
    vnet_image_pull_enabled = optional(bool, false)

    # Route backup/restore traffic over the VNet integration.
    virtual_network_backup_restore_enabled = optional(bool, false)

    # Managed identity configuration. type arrives as the spec enum's
    # name string (SYSTEM_ASSIGNED / USER_ASSIGNED /
    # SYSTEM_AND_USER_ASSIGNED).
    identity = optional(object({
      type         = string
      identity_ids = optional(list(string), [])
    }))

    # User Assigned Identity ID for Key Vault references.
    key_vault_reference_identity_id = optional(string)

    # Basic-auth publishing toggles (Web Deploy and FTP). Flip both to
    # false to force identity-based deployment.
    webdeploy_publish_basic_authentication_enabled = optional(bool, true)
    ftp_publish_basic_authentication_enabled       = optional(bool, true)

    # Path to a local ZIP package to deploy on create/update.
    zip_deploy_file = optional(string)

    # Azure Storage Account mounts. type arrives as the spec enum's name
    # string (AZURE_FILES / AZURE_BLOB); access keys are resolved secrets.
    storage_mounts = optional(list(object({
      name         = string
      type         = string
      account_name = string
      share_name   = string
      access_key   = string
      mount_path   = optional(string)
    })), [])

    # Scheduled backups to a storage container (SAS URL is a resolved
    # secret). Standard tier or above.
    backup = optional(object({
      name                = string
      storage_account_url = string
      enabled             = optional(bool, true)
      schedule = object({
        frequency_interval       = number
        frequency_unit           = string
        keep_at_least_one_backup = optional(bool, false)
        retention_period_days    = optional(number, 30)
        start_time               = optional(string)
      })
    }))

    # App Service built-in authentication (Easy Auth v2). Enum-valued
    # fields arrive as the spec enum's name strings; provider secrets are
    # referenced by app setting NAME, never inline.
    auth_settings_v2 = optional(object({
      auth_enabled                            = optional(bool, false)
      runtime_version                         = optional(string, "~1")
      config_file_path                        = optional(string)
      require_authentication                  = optional(bool, false)
      unauthenticated_action                  = optional(string)
      default_provider                        = optional(string)
      excluded_paths                          = optional(list(string), [])
      require_https                           = optional(bool, true)
      http_route_api_prefix                   = optional(string, "/.auth")
      forward_proxy_convention                = optional(string)
      forward_proxy_custom_host_header_name   = optional(string)
      forward_proxy_custom_scheme_header_name = optional(string)

      login = object({
        logout_endpoint                   = optional(string)
        token_store_enabled               = optional(bool, false)
        token_refresh_extension_time      = optional(number, 72)
        token_store_path                  = optional(string)
        token_store_sas_setting_name      = optional(string)
        preserve_url_fragments_for_logins = optional(bool, false)
        allowed_external_redirect_urls    = optional(list(string), [])
        cookie_expiration_convention      = optional(string)
        cookie_expiration_time            = optional(string, "08:00:00")
        validate_nonce                    = optional(bool, true)
        nonce_expiration_time             = optional(string, "00:05:00")
      })

      apple_v2 = optional(object({
        client_id                  = string
        client_secret_setting_name = string
      }))

      active_directory_v2 = optional(object({
        client_id                            = string
        tenant_auth_endpoint                 = string
        client_secret_setting_name           = optional(string)
        client_secret_certificate_thumbprint = optional(string)
        login_parameters                     = optional(map(string), {})
        www_authentication_disabled          = optional(bool, false)
        jwt_allowed_groups                   = optional(list(string), [])
        jwt_allowed_client_applications      = optional(list(string), [])
        allowed_groups                       = optional(list(string), [])
        allowed_identities                   = optional(list(string), [])
        allowed_applications                 = optional(list(string), [])
        allowed_audiences                    = optional(list(string), [])
      }))

      azure_static_web_app_v2 = optional(object({
        client_id = string
      }))

      custom_oidc_v2 = optional(list(object({
        name                          = string
        client_id                     = string
        openid_configuration_endpoint = string
        name_claim_type               = optional(string)
        scopes                        = optional(list(string), [])
      })), [])

      facebook_v2 = optional(object({
        app_id                  = string
        app_secret_setting_name = string
        graph_api_version       = optional(string)
        login_scopes            = optional(list(string), [])
      }))

      github_v2 = optional(object({
        client_id                  = string
        client_secret_setting_name = string
        login_scopes               = optional(list(string), [])
      }))

      google_v2 = optional(object({
        client_id                  = string
        client_secret_setting_name = string
        allowed_audiences          = optional(list(string), [])
        login_scopes               = optional(list(string), [])
      }))

      microsoft_v2 = optional(object({
        client_id                  = string
        client_secret_setting_name = string
        allowed_audiences          = optional(list(string), [])
        login_scopes               = optional(list(string), [])
      }))

      twitter_v2 = optional(object({
        consumer_key                 = string
        consumer_secret_setting_name = string
      }))
    }))

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
