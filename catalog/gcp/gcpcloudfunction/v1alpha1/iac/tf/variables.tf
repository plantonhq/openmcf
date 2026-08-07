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
  description = "Specification for the GCP Cloud Function (Gen 2)"
  type = object({
    project_id    = optional(string, "")
    region        = string
    function_name = optional(string, "")
    description   = optional(string, "")
    labels        = optional(map(string), {})
    kms_key_name  = optional(string, "")

    build_config = object({
      runtime     = string
      entry_point = string
      source = object({
        storage_source = optional(object({
          bucket     = string
          object     = string
          generation = optional(number, null)
        }), null)
        repo_source = optional(object({
          repo_name    = string
          branch_name  = optional(string, "")
          tag_name     = optional(string, "")
          commit_sha   = optional(string, "")
          dir          = optional(string, "")
          invert_regex = optional(bool, false)
          project_id   = optional(string, "")
        }), null)
      })
      build_environment_variables = optional(map(string), {})
      service_account             = optional(string, "")
      worker_pool                 = optional(string, "")
      docker_repository           = optional(string, "")
      update_policy               = optional(string, "")
    })

    service_config = optional(object({
      service_account_email             = optional(string, "")
      available_memory                  = optional(string, "")
      available_cpu                     = optional(string, "")
      timeout_seconds                   = optional(number, null)
      max_instance_request_concurrency  = optional(number, null)
      environment_variables             = optional(map(string), {})
      secret_environment_variables = optional(list(object({
        key        = string
        secret     = string
        version    = optional(string, "")
        project_id = optional(string, "")
      })), [])
      secret_volumes = optional(list(object({
        mount_path = string
        secret     = string
        project_id = optional(string, "")
        versions = optional(list(object({
          version = string
          path    = string
        })), [])
      })), [])
      vpc_connector                  = optional(string, "")
      vpc_connector_egress_settings  = optional(string, "")
      ingress_settings               = optional(string, "")
      scaling = optional(object({
        min_instance_count = optional(number, null)
        max_instance_count = optional(number, null)
      }), null)
      all_traffic_on_latest_revision = optional(bool, true)
      binary_authorization_policy    = optional(string, "")
      allow_unauthenticated          = optional(bool, false)
    }), null)

    trigger = optional(object({
      trigger_type = optional(string, "")
      event_trigger = optional(object({
        event_type   = string
        pubsub_topic = optional(string, "")
        event_filters = optional(list(object({
          attribute = string
          value     = string
          operator  = optional(string, "")
        })), [])
        trigger_region        = optional(string, "")
        retry_policy          = optional(string, "")
        service_account_email = optional(string, "")
      }), null)
    }), null)
  })
}
