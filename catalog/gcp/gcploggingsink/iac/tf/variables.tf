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
  description = "GcpLoggingSink specification"
  type = object({
    scope = optional(object({
      project_id      = optional(string, "")
      folder_id       = optional(string, "")
      organization_id = optional(string, "")
      billing_account = optional(string, "")
    }))
    sink_name = optional(string, "")
    destination = object({
      gcs_bucket             = optional(string, "")
      bigquery_dataset       = optional(string, "")
      use_partitioned_tables = optional(bool, false)
      pubsub_topic           = optional(string, "")
      raw_uri                = optional(string, "")
    })
    filter      = optional(string, "")
    description = optional(string, "")
    disabled    = optional(bool, false)
    exclusions = optional(list(object({
      name        = string
      filter      = string
      description = optional(string, "")
      disabled    = optional(bool, false)
    })), [])
    include_children       = optional(bool, false)
    intercept_children     = optional(bool, false)
    unique_writer_identity = optional(bool)
    custom_writer_identity = optional(string, "")
    deletion_policy        = optional(string, "")
  })
}