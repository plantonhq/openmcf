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
  description = "GcpSecretManagerSecret specification"
  type = object({
    project_id = optional(string, "")
    secret_id  = optional(string, "")
    region     = optional(string, "")
    replication = optional(object({
      auto = optional(object({
        customer_managed_encryption = optional(object({
          kms_key = string
        }))
      }))
      user_managed = optional(object({
        replicas = list(object({
          location = string
          customer_managed_encryption = optional(object({
            kms_key = string
          }))
        }))
      }))
    }))
    customer_managed_encryption = optional(object({
      kms_key = string
    }))
    labels              = optional(map(string), {})
    annotations         = optional(map(string), {})
    tags                = optional(map(string), {})
    expire_time         = optional(string, "")
    ttl                 = optional(string, "")
    version_aliases     = optional(map(string), {})
    version_destroy_ttl = optional(string, "")
    rotation = optional(object({
      rotation_period    = optional(string, "")
      next_rotation_time = optional(string, "")
    }))
    topics = optional(list(string), [])
    initial_version = optional(object({
      data            = string
      enabled         = optional(bool)
      is_base64       = optional(bool, false)
      deletion_policy = optional(string, "")
    }))
    iam_members = optional(list(object({
      role   = string
      member = string
      condition = optional(object({
        title       = string
        expression  = string
        description = optional(string, "")
      }))
    })), [])
    deletion_protection = optional(bool, false)
    deletion_policy     = optional(string, "")
  })
}