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
  description = "Specification for the GCP Firestore database"
  type = object({
    # The GCP project for the database. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Location (multi-region nam5/eur3 or a single region). Immutable
    # (ForceNew).
    location_id = string

    # Database ID ("(default)" or a custom name). Immutable (ForceNew).
    database_name = string

    # FIRESTORE_NATIVE or DATASTORE_MODE. Mutable (a significant
    # operational change).
    type = string

    # OPTIMISTIC / PESSIMISTIC / OPTIMISTIC_WITH_ENTITY_GROUPS; empty
    # leaves GCP's per-type default.
    concurrency_mode = optional(string, "")

    # POINT_IN_TIME_RECOVERY_ENABLED extends version retention to 7 days.
    point_in_time_recovery_enablement = optional(string, "")

    # GCP-side delete guard; deletion fails while ENABLED.
    delete_protection_state = optional(string, "DELETE_PROTECTION_DISABLED")

    # STANDARD or ENTERPRISE (ENTERPRISE requires FIRESTORE_NATIVE).
    # Immutable (ForceNew).
    database_edition = optional(string, "")

    # CMEK key resource ID; empty means Google-managed encryption.
    # Immutable (ForceNew).
    kms_key_name = optional(string, "")

    # ENABLED couples the database's lifecycle to the project's App
    # Engine app (legacy); DISABLED keeps it independent. Immutable
    # (ForceNew).
    app_engine_integration_mode = optional(string, "")
  })

  validation {
    condition     = var.spec.location_id != ""
    error_message = "location_id is required."
  }

  validation {
    condition     = var.spec.database_name != ""
    error_message = "database_name is required."
  }

  validation {
    condition     = var.spec.type != ""
    error_message = "type is required."
  }
}
