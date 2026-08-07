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
  description = "Specification for the GCP Compute Engine SSL policy"
  type = object({
    # The GCP project that owns the SSL policy. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the SSL policy in GCP (RFC1035). Empty defaults to
    # metadata.name (see locals.tf). Immutable (ForceNew).
    ssl_policy_name = optional(string, "")

    # Region for a REGIONAL SSL policy; empty means GLOBAL. The scope selects
    # which provider resource is created (see main.tf). Immutable.
    region = optional(string, "")

    # Why this policy exists and which proxies should use it. Immutable on
    # this resource (a GCP API quirk — most descriptions are mutable).
    description = optional(string, "")

    # Cipher-suite profile: COMPATIBLE (GCP default), MODERN, RESTRICTED, or
    # CUSTOM. Empty falls through to the API default (COMPATIBLE). Mutable.
    profile = optional(string, "")

    # Minimum TLS version clients may negotiate: TLS_1_0 (GCP default),
    # TLS_1_1, or TLS_1_2. Empty falls through to the API default. Mutable.
    min_tls_version = optional(string, "")

    # Exact cipher suites to allow — required with (and only valid with) the
    # CUSTOM profile. Mutable.
    custom_features = optional(list(string), [])
  })

  validation {
    condition     = var.spec.profile == "" || contains(["COMPATIBLE", "MODERN", "RESTRICTED", "CUSTOM"], var.spec.profile)
    error_message = "profile must be COMPATIBLE, MODERN, RESTRICTED, or CUSTOM."
  }

  validation {
    condition     = var.spec.min_tls_version == "" || contains(["TLS_1_0", "TLS_1_1", "TLS_1_2"], var.spec.min_tls_version)
    error_message = "min_tls_version must be TLS_1_0, TLS_1_1, or TLS_1_2."
  }

  # Mirrors the provider's own CustomizeDiff rule so the mismatch fails at
  # plan time here instead of at apply time in GCP.
  validation {
    condition     = var.spec.profile == "CUSTOM" ? length(var.spec.custom_features) > 0 : length(var.spec.custom_features) == 0
    error_message = "the CUSTOM profile requires custom_features, and custom_features is only valid with the CUSTOM profile."
  }
}
