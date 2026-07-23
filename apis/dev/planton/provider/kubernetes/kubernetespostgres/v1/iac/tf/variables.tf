# Typed mirror of KubernetesPostgresSpec (spec.proto). The spec arrives from
# the proto->tfvars converter in snake_case with every StringValueOrRef
# foreign key -- `namespace` (KubernetesNamespace), `storage.storage_class`
# (KubernetesStorageClass), `certificates.server_tls_secret`
# (KubernetesCertificate), and the workload-identity references -- resolved
# to a literal string before Terraform runs.
#
# optional() defaults mirror the proto's (dev.planton.shared.options.default)
# annotations, so the module renders the same resource whether or not the
# platform's defaulting middleware ran. Fields whose ABSENCE is meaningful
# (enable_pdb, resize_in_use_volumes, connection_limit) carry NO default:
# the module renders them only on explicit divergence from the upstream
# default, exactly like the Pulumi module's *OrNil helpers.

variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name        = string
    id          = optional(string, "")
    org         = optional(string, "")
    env         = optional(string, "")
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})
    tags        = optional(list(string), [])
  })
}

variable "spec" {
  description = "KubernetesPostgres specification"
  type = object({
    namespace        = string
    create_namespace = optional(bool, false)
    instances        = optional(number, 1)
    image_name       = optional(string, "")

    storage = object({
      size                  = string
      storage_class         = optional(string, "")
      resize_in_use_volumes = optional(bool)
    })

    wal_storage = optional(object({
      size                  = string
      storage_class         = optional(string, "")
      resize_in_use_volumes = optional(bool)
    }))

    resources = optional(object({
      limits = optional(object({
        cpu    = optional(string, "")
        memory = optional(string, "")
      }))
      requests = optional(object({
        cpu    = optional(string, "")
        memory = optional(string, "")
      }))
    }))

    postgresql = optional(object({
      parameters               = optional(map(string), {})
      pg_hba                   = optional(list(string), [])
      pg_ident                 = optional(list(string), [])
      shared_preload_libraries = optional(list(string), [])
      synchronous = optional(object({
        method          = optional(string, "any")
        number          = number
        data_durability = optional(string, "required")
      }))
      enable_alter_system = optional(bool, false)
    }))

    bootstrap = optional(object({
      initdb = optional(object({
        database                  = optional(string, "app")
        owner                     = optional(string, "")
        owner_password            = optional(string, "")
        data_checksums            = optional(bool, false)
        encoding                  = optional(string, "UTF8")
        locale_collate            = optional(string, "")
        locale_ctype              = optional(string, "")
        post_init_sql             = optional(list(string), [])
        post_init_application_sql = optional(list(string), [])
        import = optional(object({
          type                    = string
          source_external_cluster = string
          databases               = list(string)
          roles                   = optional(list(string), [])
          schema_only             = optional(bool, false)
        }))
      }))
      recovery = optional(object({
        object_store = object({
          destination_path = string
          s3 = optional(object({
            region          = optional(string, "")
            endpoint_url    = optional(string, "")
            endpoint_ca_pem = optional(string, "")
            keyless         = optional(bool, false)
            access_keys = optional(object({
              access_key_id     = string
              secret_access_key = string
            }))
          }))
          gcs = optional(object({
            keyless                  = optional(bool, false)
            service_account_key_json = optional(string, "")
          }))
          azure_blob = optional(object({
            keyless           = optional(bool, false)
            connection_string = optional(string, "")
            storage_account   = optional(string, "")
            storage_key       = optional(string, "")
          }))
          wal = optional(object({
            compression  = optional(string, "")
            max_parallel = optional(number)
          }))
          data = optional(object({
            compression          = optional(string, "")
            jobs                 = optional(number)
            immediate_checkpoint = optional(bool, false)
          }))
        })
        source_server_name = string
        recovery_target = optional(object({
          target_time      = optional(string, "")
          target_lsn       = optional(string, "")
          target_name      = optional(string, "")
          target_immediate = optional(bool, false)
          backup_id        = optional(string, "")
        }))
      }))
      pg_basebackup = optional(object({
        source = string
      }))
    }))

    external_clusters = optional(list(object({
      name                  = string
      connection_parameters = optional(map(string), {})
      password              = optional(string, "")
    })), [])

    superuser = optional(object({
      enabled  = optional(bool, false)
      password = optional(string, "")
    }))

    roles = optional(list(object({
      name             = string
      comment          = optional(string, "")
      ensure           = optional(string, "present")
      password         = optional(string, "")
      disable_password = optional(bool, false)
      login            = optional(bool, false)
      superuser        = optional(bool, false)
      createdb         = optional(bool, false)
      createrole       = optional(bool, false)
      replication      = optional(bool, false)
      bypassrls        = optional(bool, false)
      in_roles         = optional(list(string), [])
      connection_limit = optional(number)
    })), [])

    backup = optional(object({
      object_store = object({
        destination_path = string
        s3 = optional(object({
          region          = optional(string, "")
          endpoint_url    = optional(string, "")
          endpoint_ca_pem = optional(string, "")
          keyless         = optional(bool, false)
          access_keys = optional(object({
            access_key_id     = string
            secret_access_key = string
          }))
        }))
        gcs = optional(object({
          keyless                  = optional(bool, false)
          service_account_key_json = optional(string, "")
        }))
        azure_blob = optional(object({
          keyless           = optional(bool, false)
          connection_string = optional(string, "")
          storage_account   = optional(string, "")
          storage_key       = optional(string, "")
        }))
        wal = optional(object({
          compression  = optional(string, "")
          max_parallel = optional(number)
        }))
        data = optional(object({
          compression          = optional(string, "")
          jobs                 = optional(number)
          immediate_checkpoint = optional(bool, false)
        }))
      })
      retention_policy = optional(string, "")
      schedules = optional(list(object({
        name      = string
        schedule  = string
        immediate = optional(bool, false)
        suspend   = optional(bool, false)
        target    = optional(string, "prefer-standby")
      })), [])
    }))

    workload_identity = optional(object({
      gke = optional(object({
        service_account_email = string
      }))
      eks = optional(object({
        role_arn = string
      }))
      aks = optional(object({
        client_id = string
        tenant_id = optional(string, "")
      }))
    }))

    certificates = optional(object({
      server_tls_secret    = optional(string, "")
      server_ca_secret     = optional(string, "")
      server_alt_dns_names = optional(list(string), [])
    }))

    monitoring = optional(object({
      tls_enabled             = optional(bool, false)
      disable_default_queries = optional(bool, false)
    }))

    scheduling = optional(object({
      anti_affinity_type = optional(string, "preferred")
      topology_key       = optional(string, "")
      node_selector      = optional(map(string), {})
      tolerations = optional(list(object({
        key                = optional(string, "")
        operator           = optional(string, "")
        value              = optional(string, "")
        effect             = optional(string, "")
        toleration_seconds = optional(number)
      })), [])
      priority_class_name = optional(string, "")
    }))

    update_strategy = optional(object({
      primary_update_strategy = optional(string, "unsupervised")
      primary_update_method   = optional(string, "restart")
    }))

    enable_pdb         = optional(bool)
    image_pull_secrets = optional(list(string), [])
  })
}
