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
  description = "Specification for the GCP Spanner backup schedule"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}).
    project_id = optional(string, "")
    instance   = string
    database   = string

    schedule_name      = optional(string, "")
    cron               = string
    retention_duration = string

    # Proto default "FULL" is materialized by the converter; guard anyway.
    backup_type = optional(string, "FULL")

    encryption_config = optional(object({
      encryption_type = string
      kms_key_name    = optional(string, "")
      kms_key_names   = optional(list(string), [])
    }), null)
  })
}
