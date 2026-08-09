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
  description = "Specification for the Cloud Composer user workloads Secret"
  type = object({
    # The GCP project of the Composer environment. The CLI's tfvars
    # converter resolves StringValueOrRef fields to their literal string
    # before the module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Region of the Composer environment. Immutable (ForceNew).
    region = string

    # The Composer environment name (resolved ref). Immutable (ForceNew).
    environment = string

    # Kubernetes Secret name. Immutable (ForceNew).
    secret_name = string

    # Key-value entries; values are base64-encoded secret material (the
    # API's Kubernetes Secret contract). The provider schema marks this
    # attribute sensitive, so plans redact it; it still lands in IaC
    # state, which is the engine's secret boundary.
    data = map(string)

    # Client-side destroy behavior: DELETE (default), PREVENT, ABANDON.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.region != ""
    error_message = "region is required."
  }

  validation {
    condition     = var.spec.environment != ""
    error_message = "environment is required."
  }

  validation {
    condition     = var.spec.secret_name != ""
    error_message = "secret_name is required."
  }

  validation {
    condition     = length(var.spec.data) > 0
    error_message = "data must have at least one entry."
  }
}
