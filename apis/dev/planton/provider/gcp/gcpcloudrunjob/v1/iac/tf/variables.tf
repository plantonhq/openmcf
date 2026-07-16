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
  description = "Specification for the GCP Cloud Run job"
  type = object({
    project_id = optional(string, "")
    region     = string
    job_name   = optional(string, "")
    labels     = optional(map(string), {})
    annotations = optional(map(string), {})

    template = object({
      containers = list(object({
        name    = optional(string, "")
        image   = string
        command = optional(list(string), [])
        args    = optional(list(string), [])
        env = optional(list(object({
          name  = string
          value = optional(string, "")
          value_from_secret = optional(object({
            secret  = string
            version = optional(string, "")
          }), null)
        })), [])
        resources = optional(object({
          cpu    = optional(string, "")
          memory = optional(string, "")
        }), null)
        volume_mounts = optional(list(object({
          name       = string
          mount_path = string
        })), [])
        working_dir = optional(string, "")
        depends_on  = optional(list(string), [])
      }))

      volumes = optional(list(object({
        name = string
        cloud_sql_instance = optional(object({
          instances = list(string)
        }), null)
        secret = optional(object({
          secret       = string
          default_mode = optional(number, null)
          items = optional(list(object({
            path    = string
            version = optional(string, "")
            mode    = optional(number, null)
          })), [])
        }), null)
        empty_dir = optional(object({
          medium     = optional(string, "")
          size_limit = optional(string, "")
        }), null)
        gcs = optional(object({
          bucket    = string
          read_only = optional(bool, false)
        }), null)
        nfs = optional(object({
          server    = string
          path      = string
          read_only = optional(bool, false)
        }), null)
      })), [])

      service_account       = optional(string, "")
      execution_environment = optional(string, "")
      encryption_key        = optional(string, "")
      timeout_seconds       = optional(number, null)
      max_retries           = optional(number, null)

      vpc_access = optional(object({
        connector = optional(string, "")
        network_interfaces = optional(list(object({
          network    = optional(string, "")
          subnetwork = optional(string, "")
          tags       = optional(list(string), [])
        })), [])
        egress = optional(string, "")
      }), null)

      node_selector = optional(object({
        accelerator = string
      }), null)
    })

    task_count   = optional(number, null)
    parallelism  = optional(number, null)
    launch_stage = optional(string, "")

    binary_authorization = optional(object({
      use_default              = optional(bool, false)
      policy                   = optional(string, "")
      breakglass_justification = optional(string, "")
    }), null)

    gpu_zonal_redundancy_disabled = optional(bool, false)
    deletion_protection           = optional(bool, true)
  })
}
