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
  description = "Specification for the GCP managed SSL certificate"
  type = object({
    # The GCP project that owns the certificate. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string. Empty falls back to the
    # provider's default project (see locals.tf).
    project_id = optional(string, "")

    # Name of the certificate in GCP (RFC1035). Empty defaults to metadata.name
    # (see locals.tf). Immutable (ForceNew).
    certificate_name = optional(string, "")

    # What this certificate secures — immutable.
    description = optional(string, "")

    # Domains the certificate is valid for (1-100). Each must be a fully-
    # qualified domain name; wildcards are not supported. Immutable (ForceNew).
    domains = list(string)

    # What happens to the certificate when this resource is destroyed:
    # DELETE (default), PREVENT, or ABANDON.
    deletion_policy = optional(string, "")
  })

  # NOTE: never guard optional strings with coalesce() here — HCL's coalesce
  # skips empty strings as well as nulls, so coalesce("", "") errors and the
  # validation fails on a legitimately-empty value.
  validation {
    condition     = try(var.spec.certificate_name, "") == "" || can(regex("^[a-z]([-a-z0-9]{0,61}[a-z0-9])?$", var.spec.certificate_name))
    error_message = "certificate_name must be RFC1035-compliant: 1-63 lowercase letters, digits, or hyphens."
  }

  validation {
    condition     = length(var.spec.domains) >= 1
    error_message = "domains must contain at least one fully-qualified domain name."
  }

  validation {
    condition     = length(var.spec.domains) <= 100
    error_message = "domains must contain at most 100 entries."
  }
}
