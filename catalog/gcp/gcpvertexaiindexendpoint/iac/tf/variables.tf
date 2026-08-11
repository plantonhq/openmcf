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
  description = "GcpVertexAiIndexEndpoint spec"
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

    # Public querying arm. Mutually exclusive with network and PSC
    # (CEL-enforced pre-deploy). Immutable (ForceNew).
    public_endpoint_enabled = optional(bool, false)

    # Resolved from a GcpVpcNetwork reference (canonical output is the
    # network self-link). Normalized to the API's relative form in
    # locals.tf. VPC peering; mutually exclusive with the other arms.
    # Immutable (ForceNew).
    network = optional(string, "")

    # Private Service Connect arm. Immutable (ForceNew).
    private_service_connect_config = optional(object({
      enable_private_service_connect = optional(bool, true)
      project_allowlist              = optional(list(string), [])

      # PSC endpoints Vertex AI creates automatically in consumer
      # projects/networks. network arrives resolved (a GcpVpcNetwork
      # reference's self-link), normalized to the API's relative form
      # in locals.tf.
      psc_automation_configs = optional(list(object({
        network    = string
        project_id = string
      })), [])
    }), null)

    # User labels; merged beneath the platform attribution labels.
    labels = optional(map(string), {})

    # CMEK key resource path (a GcpKmsKey reference resolves to it).
    # Empty means Google-managed encryption. Immutable (ForceNew).
    kms_key_name = optional(string, "")

    # Client-side destroy behavior: DELETE (default), PREVENT, ABANDON.
    deletion_policy = optional(string, "")
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
