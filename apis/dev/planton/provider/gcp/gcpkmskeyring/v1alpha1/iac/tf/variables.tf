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
  description = "Specification for the GCP KMS key ring"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}). If
    # project_id is empty, the provider's default project is used
    # (see locals.tf).
    project_id = optional(string, "")

    # Key ring name (the GCP resource name). Immutable (ForceNew) — and
    # since key rings can never be deleted from GCP, a name is permanently
    # consumed within its project+location once used.
    key_ring_name = string

    # Location (region, multi-region, or "global"). Immutable (ForceNew).
    location = string
  })

  validation {
    condition     = var.spec.key_ring_name != ""
    error_message = "key_ring_name is required."
  }

  validation {
    condition     = var.spec.location != ""
    error_message = "location is required."
  }
}
