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
  description = "Azure Linux Web App specification"
  type = object({
    # The Azure region where the Web App will be created. ForceNew.
    region = string

    # The Azure Resource Group name. References are resolved to a literal
    # name by the platform before the module runs. ForceNew.
    resource_group = string

    # The name of the Web App (globally unique -- it forms
    # {web_app_name}.azurewebsites.net). ForceNew.
    web_app_name = string

    # The App Service Plan resource ID. References are resolved before
    # the module runs.
    service_plan_id = string

    # Site configuration including runtime, scaling, security, networking,
    # and auto-heal.
    site_config = object({
      # Application stack (runtime selection). Exactly one runtime.
      application_stack = optional(object({
        # .NET runtime version.
        dotnet_version = optional(string)

        # Node.js runtime version (LTS identifiers, e.g. "20-lts").
        node_version = optional(string)

        # Python runtime version.
        python_version = optional(string)

        # Java runtime version; requires java_server + java_server_version.
        java_version = optional(string)

        # Java application server, as the spec enum's name string
        # (JAVA_SE / TOMCAT / JBOSSEAP).
        java_server = optional(string)

        # Java application server version.
        java_server_version = optional(string)

        # PHP runtime version.
        php_version = optional(string)

        # Ruby runtime version (legacy).
        ruby_version = optional(string)

        # Go runtime version (legacy).
        go_version = optional(string)

        # Custom container. The registry password is a resolved secret
        # literal by the time the module runs.
        docker = optional(object({
          registry_url      = string
          image_name        = string
          image_tag         = string
          registry_username = optional(string)
          registry_password = optional(string)
        }))
      }))

      # Keep the Web App always loaded in memory. Azure rejects true on
      # Free/Shared plans at apply time.
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

      # Minimum accepted TLS cipher suite (Azure's own identifiers, e.g.
      # TLS_AES_256_GCM_SHA384). Absent accepts Azure's platform default.
      minimum_tls_cipher_suite = optional(string)

      # Number of worker instances.
      worker_count = optional(number)

      # Enable HTTP/2 protocol.
      http2_enabled = optional(bool, false)

      # Enable WebSocket connections.
      websockets_enabled = optional(bool, false)

      # Use a 32-bit worker process (spec default false -- 64-bit for
      # production; the provider's own default is true).
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

      # MySQL in-app (dev/test convenience only).
      local_mysql_enabled = optional(bool, false)

      # Remote debugging (Visual Studio attach). Azure picks the
      # debugger version.
      remote_debugging_enabled = optional(bool, false)

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

      # Default action for main-site IP restrictions, as the spec enum's
      # name string (ALLOW / DENY). Absent means ALLOW.
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

      # Use managed identity for ACR image pulls.
      container_registry_use_managed_identity = optional(bool, false)

      # Client ID of the managed identity for ACR pulls.
      container_registry_managed_identity_client_id = optional(string)

      # Auto-heal: recycle on trigger conditions. The heal action is
      # implicitly Recycle (the only Linux action).
      auto_heal_setting = optional(object({
        trigger = object({
          requests = optional(object({
            count    = number
            interval = string
          }))
          status_codes = optional(list(object({
            status_code_range = string
            count             = number
            interval          = string
            sub_status        = optional(number)
            win32_status_code = optional(number)
            path              = optional(string)
          })), [])
          slow_request = optional(object({
            time_taken = string
            interval   = string
            count      = number
          }))
          slow_request_with_path = optional(list(object({
            time_taken = string
            interval   = string
            count      = number
            path       = optional(string)
          })), [])
        })
        minimum_process_execution_time = optional(string)
      }))
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

    # Application Insights connection string for APM telemetry (merged
    # into app_settings).
    application_insights_connection_string = optional(string)

    # Enforce HTTPS-only access.
    https_only = optional(bool, true)

    # Enable public network access.
    public_network_access_enabled = optional(bool, true)

    # Whether the Web App is enabled. false stops the app without
    # deleting it.
    enabled = optional(bool, true)

    # Enable client affinity (ARR session affinity cookies).
    client_affinity_enabled = optional(bool, false)

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

    # Logging configuration for the Web App.
    logs = optional(object({
      # Application-level log settings. Levels arrive as the spec enum's
      # name strings (OFF / ERROR / WARNING / INFORMATION / VERBOSE).
      application_logs = optional(object({
        file_system_level = optional(string)
        azure_blob_storage = optional(object({
          level             = string
          sas_url           = string
          retention_in_days = optional(number, 0)
        }))
      }))

      # HTTP request/response log settings -- exactly one destination.
      http_logs = optional(object({
        file_system = optional(object({
          retention_in_mb   = optional(number, 35)
          retention_in_days = optional(number, 0)
        }))
        azure_blob_storage = optional(object({
          sas_url           = string
          retention_in_days = optional(number, 0)
        }))
      }))

      # Enable failed request tracing.
      failed_request_tracing = optional(bool, false)

      # Enable detailed error messages in responses.
      detailed_error_messages = optional(bool, false)
    }))

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
