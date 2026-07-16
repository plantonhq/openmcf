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
  description = "Specification for the GCP Firestore backup schedule"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}).
    project_id = optional(string, "")
    database   = string
    retention  = string

    # Exactly one recurrence shape: daily true, or weekly_recurrence.day set.
    daily = optional(bool, false)
    weekly_recurrence = optional(object({
      day = string
    }), null)
  })

  validation {
    condition     = var.spec.database != ""
    error_message = "database is required."
  }
}
