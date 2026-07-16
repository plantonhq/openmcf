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
  description = "Specification for the GCP Bigtable table"
  type = object({
    # The GCP project owning the instance. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Short name of the parent Bigtable instance (a GcpBigtableInstance
    # reference resolves to it). Immutable (ForceNew).
    instance = string

    # Table name (1-50 chars); empty falls back to metadata.name.
    # Immutable (ForceNew).
    table_name = optional(string, "")

    # Column families with optional per-family GC policies. Mutable.
    column_families = optional(list(object({
      family = string
      type   = optional(string, "")
      gc_policy = optional(object({
        mode            = optional(string, "")
        max_age         = optional(string, "")
        max_versions    = optional(number, 0)
        gc_rules        = optional(string, "")
        ignore_warnings = optional(bool, false)
      }), null)
    })), [])

    # Row keys to pre-split the table at. Immutable — changing this
    # REPLACES the table and its data.
    split_keys = optional(list(string), [])

    # Change stream retention (24h-168h as a duration string; "0"
    # disables). Empty leaves change streams off.
    change_stream_retention = optional(string, "")

    # Built-in automated backups (both duration strings).
    automated_backup_policy = optional(object({
      retention_period = string
      frequency        = string
    }), null)

    # API-side deletion guard. The spec defaults this to PROTECTED
    # (Planton middleware materializes the default); the module sends it
    # explicitly so destroy behavior is identical on both engines.
    deletion_protection = optional(string, "PROTECTED")

    # Structured row key schema as the API's Type JSON. In-place update
    # is not supported: clear, apply, then set the new schema.
    row_key_schema = optional(string, "")
  })

  validation {
    condition     = var.spec.instance != ""
    error_message = "instance is required."
  }
}
