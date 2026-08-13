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
  description = "Specification for the self-managed GCP SSL certificate"
  type = object({
    # The GCP project that owns the certificate. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string. Empty falls back to the
    # provider's default project (see locals.tf).
    project_id = optional(string, "")

    # Name of the certificate in GCP (RFC1035). Empty defaults to
    # metadata.name (see locals.tf). Self-managed and Google-managed
    # certificates share one namespace per scope. Immutable (ForceNew).
    certificate_name = optional(string, "")

    # Region for a REGIONAL certificate; empty means GLOBAL. The scope
    # selects which provider resource is created (see main.tf). Immutable.
    region = optional(string, "")

    # What this certificate secures and where it came from. Immutable.
    description = optional(string, "")

    # PEM certificate chain: leaf first, then intermediates (max 5 certs,
    # at least one intermediate). Public handshake material. Immutable.
    certificate = string

    # PEM unencrypted private key matching the certificate. Secret material:
    # write-only in GCP, never surfaced in outputs. Immutable.
    private_key = string

    # What happens to the certificate when this resource is destroyed:
    # DELETE (default), PREVENT, or ABANDON.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = startswith(var.spec.certificate, "-----BEGIN CERTIFICATE-----")
    error_message = "certificate must be a PEM-encoded certificate chain (-----BEGIN CERTIFICATE-----)."
  }

  validation {
    condition = anytrue([
      startswith(var.spec.private_key, "-----BEGIN PRIVATE KEY-----"),
      startswith(var.spec.private_key, "-----BEGIN RSA PRIVATE KEY-----"),
      startswith(var.spec.private_key, "-----BEGIN EC PRIVATE KEY-----"),
    ])
    error_message = "private_key must be a PEM-encoded unencrypted private key."
  }

  # NOTE: never guard optional strings with coalesce() here — HCL's coalesce
  # skips empty strings as well as nulls, so coalesce("", "") errors and the
  # validation fails on a legitimately-empty value.
  validation {
    condition     = try(var.spec.certificate_name, "") == "" || can(regex("^[a-z]([-a-z0-9]{0,61}[a-z0-9])?$", var.spec.certificate_name))
    error_message = "certificate_name must be RFC1035-compliant: 1-63 lowercase letters, digits, or hyphens."
  }
}
