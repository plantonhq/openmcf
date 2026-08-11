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
  description = "Specification for the GCP Pub/Sub schema"
  type = object({
    # The GCP project for the schema. The CLI's tfvars converter resolves
    # StringValueOrRef fields to their literal string before the module
    # runs, so this arrives as a plain string. If empty, the provider's
    # default project is used (see locals.tf).
    project_id = optional(string, "")

    # Schema name (the GCP resource name). Immutable (ForceNew) — renaming
    # replaces the schema, unlike definition changes which commit revisions.
    schema_name = string

    # The schema definition language: AVRO or PROTOCOL_BUFFER.
    type = string

    # The schema definition text (Avro JSON or a protobuf message
    # definition, matching the declared type).
    definition = string

    # Deletion policy: "", "DELETE" (default), "PREVENT" (destroy fails),
    # or "ABANDON" (remove from management, leave serving in GCP).
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.schema_name != ""
    error_message = "schema_name is required."
  }

  validation {
    condition     = contains(["AVRO", "PROTOCOL_BUFFER"], var.spec.type)
    error_message = "type must be AVRO or PROTOCOL_BUFFER."
  }

  validation {
    condition     = var.spec.definition != ""
    error_message = "definition is required."
  }

  validation {
    condition     = contains(["", "DELETE", "PREVENT", "ABANDON"], var.spec.deletion_policy)
    error_message = "deletion_policy must be one of: DELETE, PREVENT, ABANDON."
  }
}
