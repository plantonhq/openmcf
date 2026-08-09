variable "metadata" {
  description = "Planton resource metadata"
  type = object({
    name    = string
    id      = optional(string, "")
    org     = optional(string, "")
    env     = optional(string, "")
    labels  = optional(map(string), {})
    tags    = optional(list(string), [])
    version = optional(string, "")
  })
}

variable "spec" {
  description = "GcpVertexAiEndpoint spec"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}). If
    # project_id is empty, the provider's default project is used
    # (see locals.tf).
    project_id = optional(string, "")

    # Region. Immutable (ForceNew).
    location = string

    display_name = string

    description = optional(string, "")

    # Numeric endpoint ID (max 10 digits). Empty means "derive a stable
    # ID from the resource identity" (see locals.tf). Immutable (ForceNew).
    endpoint_name = optional(string, "")

    # Resolved from a GcpVpcNetwork reference to the network self link.
    # VPC peering; mutually exclusive with PSC. Immutable (ForceNew).
    network = optional(string, "")

    # Resolved from a GcpKmsKey reference to the key path (CMEK).
    # Immutable (ForceNew).
    kms_key_name = optional(string, "")

    dedicated_endpoint_enabled = optional(bool, false)

    private_service_connect_config = optional(object({
      project_allowlist = optional(list(string), [])
    }), null)

    # User labels; merged beneath the platform attribution labels.
    labels = optional(map(string), {})

    request_response_logging_config = optional(object({
      enabled                  = optional(bool, false)
      sampling_rate            = optional(number, 0)
      bigquery_destination_uri = optional(string, "")
    }), null)
  })

  validation {
    condition     = var.spec.location != ""
    error_message = "location is required."
  }

  validation {
    condition     = var.spec.display_name != ""
    error_message = "display_name is required."
  }
}
