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
  description = "Specification for the GCP Firestore composite index"
  type = object({
    project_id  = optional(string, "")
    database    = optional(string, "")
    collection  = string
    query_scope = optional(string, "COLLECTION")
    api_scope   = optional(string, "ANY_API")
    density     = optional(string, "")

    # Each field plays exactly one role: order, array_config, vector_config,
    # or search_config (enforced by the spec's CEL rules before the module
    # runs).
    fields = list(object({
      field_path   = string
      order        = optional(string, "")
      array_config = optional(string, "")
      vector_config = optional(object({
        dimension = number
      }), null)
      # Firestore Enterprise search surface: text and/or geo indexing.
      search_config = optional(object({
        text_spec = optional(object({
          index_specs = list(object({
            index_type = optional(string, "")
            match_type = optional(string, "")
          }))
        }), null)
        geo_spec = optional(object({
          geo_json_indexing_disabled = optional(bool, false)
        }), null)
      }), null)
    }))

    # MongoDB-style array indexing (api_scope MONGODB_COMPATIBLE_API only).
    multikey = optional(bool, false)

    # Enforce uniqueness of the indexed field values across documents.
    unique = optional(bool, false)

    # Return as soon as index creation is requested instead of waiting for
    # the background build.
    skip_wait = optional(bool, false)

    # Destroy-time guard: "" / "DELETE" deletes, "PREVENT" fails the
    # destroy, "ABANDON" unmanages without deleting.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.collection != ""
    error_message = "collection is required."
  }

  validation {
    condition     = length(var.spec.fields) >= 1
    error_message = "at least one field is required."
  }
}
