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
  description = "GcpVertexAiDeployedIndex spec"
  type = object({
    # Region of the index endpoint — the Vertex AI API host is regional,
    # so the deployment must be addressed in the endpoint's region.
    # Immutable (ForceNew).
    location = string

    # User-chosen deployment handle (letter start; letters, numbers,
    # underscores; <=128 chars — CEL-enforced pre-deploy). Immutable.
    deployed_index_id = string

    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved). index is the full index resource
    # path; index_endpoint is the full endpoint resource path. Both
    # immutable (ForceNew).
    index          = string
    index_endpoint = string

    # Immutable on this resource class (unusually for a display name).
    display_name = optional(string, "")

    # Sizing arms: at most one (CEL-enforced). Omitting both lets GCP
    # deploy with automatic resources at its default bounds. Replica
    # bounds are the ONLY in-place-updatable fields on the deployment.
    automatic_resources = optional(object({
      min_replica_count = optional(number, 0)
      max_replica_count = optional(number, 0)
    }), null)
    dedicated_resources = optional(object({
      machine_type      = optional(string, "")
      min_replica_count = number
      max_replica_count = optional(number, 0)
    }), null)

    # IP-space partitioning group. Empty lets GCP default ("default").
    # Immutable.
    deployment_group = optional(string, "")

    # Send private-endpoint access logs to Cloud Logging. Immutable.
    enable_access_logging = optional(bool, false)

    # Names of reserved VPC_PEERING address ranges to deploy into
    # (GcpGlobalAddress references resolve to their names). Immutable.
    reserved_ip_ranges = optional(list(string), [])

    # JWT auth on the private query endpoint. allowed_issuers are
    # service-account emails (GcpServiceAccount references resolve to
    # them). Immutable.
    auth_config = optional(object({
      allowed_issuers = optional(list(string), [])
      audiences       = optional(list(string), [])
    }), null)
  })

  validation {
    condition     = var.spec.location != ""
    error_message = "location is required."
  }

  validation {
    condition     = var.spec.deployed_index_id != ""
    error_message = "deployed_index_id is required."
  }

  validation {
    condition     = var.spec.index != ""
    error_message = "index is required."
  }

  validation {
    condition     = var.spec.index_endpoint != ""
    error_message = "index_endpoint is required."
  }
}

variable "provider_config" {
  description = "GCP provider configuration"
  type = object({
    service_account_key = optional(string, "")
  })
  default = { service_account_key = "" }
}
