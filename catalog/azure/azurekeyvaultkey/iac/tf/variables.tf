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
  description = "Azure Key Vault key specification"
  type = object({
    # The key's name within the vault: 1-127 letters/digits/hyphens.
    name = string

    # The vault the key lives in. References are resolved to a literal
    # ARM ID by the platform before the module runs.
    key_vault_id = string

    # The algorithm family, as the spec enum's name string
    # (RSA/RSA_HSM/EC/EC_HSM); mapped to ARM's hyphenated values in
    # locals.
    key_type = string

    # For RSA keys: the modulus size in bits (2048/3072/4096). Unset for
    # EC keys (spec validation enforces the pairing).
    key_size = optional(number)

    # For EC keys: the curve, as the spec enum's name string
    # (P_256/P_256K/P_384/P_521). Unset lets Azure default to P-256.
    curve = optional(string)

    # The permitted operations, as the spec enum's name strings
    # (DECRYPT/ENCRYPT/SIGN/UNWRAP_KEY/VERIFY/WRAP_KEY); mapped to
    # Azure's camelCase values in locals.
    key_opts = list(string)

    # Activation / expiry instants, RFC 3339 UTC.
    not_before_date = optional(string)
    expiration_date = optional(string)

    # Automatic rotation policy (ISO 8601 durations, e.g. "P90D").
    rotation_policy = optional(object({
      expire_after         = optional(string)
      notify_before_expiry = optional(string)
      automatic = optional(object({
        time_after_creation = optional(string)
        time_before_expiry  = optional(string)
      }))
    }))

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
