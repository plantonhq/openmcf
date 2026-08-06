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

    fields = list(object({
      field_path = string
      order      = optional(string, "")
      array_config = optional(string, "")
      vector_config = optional(object({
        dimension = number
      }), null)
    }))
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
