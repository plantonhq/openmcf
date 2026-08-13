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
  description = "Specification for the GCP Spanner database"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}).
    project_id = optional(string, "")
    instance   = string

    database_name            = optional(string, "")
    database_dialect         = optional(string, "")
    version_retention_period = optional(string, "")
    ddl                      = optional(list(string), [])
    enable_drop_protection   = optional(bool, false)

    encryption_config = optional(object({
      kms_key_name  = optional(string, "")
      kms_key_names = optional(list(string), [])
    }), null)

    default_time_zone   = optional(string, "")
    deletion_protection = optional(bool, true)
    deletion_policy     = optional(string, "")
  })
}
