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
  description = "KubernetesOpenBao specification"
  type = object({
    namespace        = string
    create_namespace = optional(bool, false)
    chart_version    = optional(string)
    server = optional(object({
      dev        = optional(object({}))
      standalone = optional(object({}))
      ha = optional(object({
        replicas = optional(number)
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
      data_storage = optional(object({
        size          = optional(string)
        storage_class = optional(string, "")
      }))
      audit_storage = optional(object({
        size          = optional(string)
        storage_class = optional(string, "")
      }))
      log_level  = optional(string)
      log_format = optional(string)
      scheduling = optional(object({
        node_selector = optional(map(string), {})
        tolerations = optional(list(object({
          key                = optional(string, "")
          operator           = optional(string, "")
          value              = optional(string, "")
          effect             = optional(string, "")
          toleration_seconds = optional(number)
        })), [])
      }))
    }))
    tls = optional(object({
      enabled          = optional(bool, false)
      cert_secret_name = optional(string, "")
    }))
    auto_unseal = optional(object({
      aws_kms = optional(object({
        region            = string
        kms_key_id        = string
        access_key_id     = optional(string, "")
        secret_access_key = optional(string, "")
      }))
      gcp_kms = optional(object({
        project                           = string
        region                            = string
        key_ring                          = string
        crypto_key                        = string
        workload_identity_service_account = optional(string, "")
      }))
      azure_key_vault = optional(object({
        vault_name    = string
        key_name      = string
        tenant_id     = string
        client_id     = optional(string, "")
        client_secret = optional(string, "")
      }))
      transit = optional(object({
        address    = string
        key_name   = string
        mount_path = optional(string)
        token      = optional(string, "")
      }))
    }))
    injector = optional(object({
      enabled        = optional(bool, false)
      replicas       = optional(number)
      failure_policy = optional(string)
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
    }))
    ui_enabled             = optional(bool)
    network_policy_enabled = optional(bool, false)
    metrics = optional(object({
      enabled                 = optional(bool, false)
      service_monitor_enabled = optional(bool, false)
    }))
    snapshot_agent = optional(object({
      enabled                    = optional(bool, false)
      schedule                   = optional(string)
      s3_host                    = string
      s3_bucket                  = string
      s3_expire_days             = optional(number)
      s3_credentials_secret_name = string
      bao_role                   = optional(string)
      bao_auth_path              = optional(string)
    }))
    service_account = optional(object({
      annotations            = optional(map(string), {})
      auth_delegator_enabled = optional(bool)
    }))
    helm_values = optional(string, "")
  })
}
