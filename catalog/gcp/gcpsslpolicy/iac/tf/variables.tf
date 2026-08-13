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

    # Cipher-suite profile: COMPATIBLE (GCP default), MODERN, RESTRICTED,
    # CUSTOM, or FIPS_202205 (requires min_tls_version TLS_1_2). Empty falls
    # through to the API default (COMPATIBLE). Mutable.
    profile = optional(string, "")

    # Minimum TLS version clients may negotiate: TLS_1_0 (GCP default),
    # TLS_1_1, TLS_1_2, or TLS_1_3 (requires the RESTRICTED profile). Empty
    # falls through to the API default. Mutable.
    min_tls_version = optional(string, "")

    # Exact cipher suites to allow — required with (and only valid with) the
    # CUSTOM profile. Mutable.
    custom_features = optional(list(string), [])

    # Post-quantum key exchange (X25519MLKEM768) rollout stance: DEFAULT
    # (follow GCP's timeline), ENABLED, or DEFERRED. Empty falls through to
    # the API default (DEFAULT). Mutable.
    post_quantum_key_exchange = optional(string, "")

    # What happens to the policy when this resource is destroyed:
    # DELETE (default), PREVENT, or ABANDON.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.profile == "" || contains(["COMPATIBLE", "MODERN", "RESTRICTED", "CUSTOM", "FIPS_202205"], var.spec.profile)
    error_message = "profile must be COMPATIBLE, MODERN, RESTRICTED, CUSTOM, or FIPS_202205."
  }

  validation {
    condition     = var.spec.min_tls_version == "" || contains(["TLS_1_0", "TLS_1_1", "TLS_1_2", "TLS_1_3"], var.spec.min_tls_version)
    error_message = "min_tls_version must be TLS_1_0, TLS_1_1, TLS_1_2, or TLS_1_3."
  }

  # Mirrors the provider's own CustomizeDiff rule so the mismatch fails at
  # plan time here instead of at apply time in GCP.
  validation {
    condition     = var.spec.profile == "CUSTOM" ? length(var.spec.custom_features) > 0 : length(var.spec.custom_features) == 0
    error_message = "the CUSTOM profile requires custom_features, and custom_features is only valid with the CUSTOM profile."
  }

  # The two provider-documented pairings, failed at plan time instead of at
  # apply time in GCP.
  validation {
    condition     = var.spec.profile != "FIPS_202205" || var.spec.min_tls_version == "TLS_1_2"
    error_message = "the FIPS_202205 profile requires min_tls_version TLS_1_2."
  }

  validation {
    condition     = var.spec.min_tls_version != "TLS_1_3" || var.spec.profile == "RESTRICTED"
    error_message = "min_tls_version TLS_1_3 requires the RESTRICTED profile."
  }

  validation {
    condition     = var.spec.post_quantum_key_exchange == "" || contains(["DEFAULT", "ENABLED", "DEFERRED"], var.spec.post_quantum_key_exchange)
    error_message = "post_quantum_key_exchange must be DEFAULT, ENABLED, or DEFERRED."
  }
}
