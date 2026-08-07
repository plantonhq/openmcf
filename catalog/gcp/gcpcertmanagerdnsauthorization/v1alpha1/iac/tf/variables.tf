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
  description = "Specification for the Certificate Manager DNS authorization"
  type = object({
    # StringValueOrRef fields arrive as PLAIN STRINGS: the tfvars converter
    # flattens refs before the module ever sees them.
    project_id = optional(string, "")

    authorization_name = optional(string, "")

    # The domain being authorized (covers the domain and its wildcard).
    domain = string

    description = optional(string, "")

    # Certificate Manager location; empty means "global".
    location = optional(string, "")

    # FIXED_RECORD or PER_PROJECT_RECORD; empty lets GCP pick the
    # location-appropriate default.
    type = optional(string, "")

    labels = optional(map(string), {})
  })
}
