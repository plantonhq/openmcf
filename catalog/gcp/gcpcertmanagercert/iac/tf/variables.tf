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
  description = "Specification for the Certificate Manager certificate"
  type = object({
    # StringValueOrRef fields arrive as PLAIN STRINGS: the tfvars converter
    # flattens refs before the module ever sees them.
    project_id = optional(string, "")

    cert_name   = optional(string, "")
    description = optional(string, "")

    # Certificate Manager location; empty means "global".
    location = optional(string, "")

    # DEFAULT, EDGE_CACHE, ALL_REGIONS, or CLIENT_AUTH; empty means DEFAULT.
    scope = optional(string, "")

    # Google-managed arm: provisioned and renewed automatically. Exactly
    # one of managed or self_managed is set (enforced pre-deploy).
    managed = optional(object({
      domains = list(string)
      # Fully-qualified DNS authorization IDs (flattened refs).
      dns_authorizations = optional(list(string), [])
      issuance_config    = optional(string, "")
    }), null)

    # Self-managed arm: bring-your-own PEM certificate and key.
    self_managed = optional(object({
      pem_certificate = string
      pem_private_key = string
    }), null)

    labels = optional(map(string), {})

    # What happens to the certificate when this resource is destroyed:
    # DELETE (default), PREVENT, or ABANDON.
    deletion_policy = optional(string, "")
  })
}
