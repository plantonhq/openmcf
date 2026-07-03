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
  description = "Specification for the GCP Cloud SQL database"
  type = object({
    # The GCP project that owns the Cloud SQL instance. The CLI's tfvars
    # converter resolves StringValueOrRef fields to their literal string
    # before the module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the Cloud SQL instance hosting this database; arrives as a
    # plain string after ref resolution. Immutable.
    instance = string

    # Name of the database inside the instance. Immutable.
    database_name = string

    # Character set (engine-specific; empty uses the engine default).
    charset = optional(string, "")

    # Collation (engine-specific; empty uses the engine default).
    collation = optional(string, "")
  })
}
