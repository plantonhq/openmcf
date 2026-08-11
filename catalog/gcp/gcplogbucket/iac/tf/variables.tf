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
  description = "GcpLogBucket specification"
  type = object({
    scope = optional(object({
      project_id      = optional(string, "")
      folder_id       = optional(string, "")
      organization_id = optional(string, "")
      billing_account = optional(string, "")
    }))
    bucket_id        = string
    location         = optional(string, "")
    description      = optional(string, "")
    retention_days   = optional(number, 0)
    locked           = optional(bool, false)
    enable_analytics = optional(bool)
    cmek_kms_key     = optional(string, "")
    index_configs = optional(list(object({
      field_path = string
      type       = string
    })), [])
    log_views = optional(list(object({
      view_id     = string
      filter      = optional(string, "")
      description = optional(string, "")
    })), [])
    linked_bigquery_dataset = optional(object({
      link_id     = string
      description = optional(string, "")
    }))
    scope_settings = optional(object({
      disable_default_sink = optional(bool, false)
      kms_key              = optional(string, "")
      storage_location     = optional(string, "")
    }))
    deletion_policy = optional(string, "")
  })
}