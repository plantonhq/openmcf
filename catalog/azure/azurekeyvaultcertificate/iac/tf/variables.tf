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
  description = "Azure Key Vault certificate specification"
  type = object({
    # The certificate's name within the vault: 1-127
    # letters/digits/hyphens.
    name = string

    # The vault the certificate lives in. References are resolved to a
    # literal ARM ID by the platform before the module runs.
    key_vault_id = string

    # An existing bundle to IMPORT: base64 PFX/PEM contents (carries the
    # private key -- sensitive) and its optional password. Omit to have
    # the vault generate from certificate_policy (spec validation requires
    # at least one of the two).
    certificate = optional(object({
      contents = string
      password = optional(string)
    }))

    # How the vault generates and manages the certificate. Enum-typed
    # sub-fields arrive as the spec enums' FULL name strings and are
    # mapped to Azure's case-sensitive values in locals.
    certificate_policy = optional(object({
      # "Self" (self-signed), "Unknown" (out-of-band CA via CSR), or the
      # name of a CA issuer configured on the vault.
      issuer_name = string

      # The private key's shape: key_type as RSA/RSA_HSM/EC/EC_HSM/OCT,
      # curve as P_256/P_256K/P_384/P_521.
      key_properties = object({
        exportable = optional(bool, false)
        key_type   = string
        key_size   = optional(number)
        curve      = optional(string)
        reuse_key  = optional(bool, false)
      })

      # Renewal/notification actions: action_type as
      # AUTO_RENEW/EMAIL_CONTACTS; exactly one trigger field per action
      # (spec validation enforces it).
      lifetime_actions = optional(list(object({
        action_type = string
        trigger = object({
          days_before_expiry  = optional(number)
          lifetime_percentage = optional(number)
        })
      })), [])

      # The secret face's media type, as PKCS12 or PEM.
      secret_properties = object({
        content_type = string
      })

      # X.509 content -- required when generating (spec validation
      # enforces it); omitted for imports. key_usage entries arrive as
      # the spec enum's name strings (e.g. "DIGITAL_SIGNATURE").
      x509_certificate_properties = optional(object({
        subject = string
        subject_alternative_names = optional(object({
          dns_names = optional(list(string), [])
          emails    = optional(list(string), [])
          upns      = optional(list(string), [])
        }))
        key_usage          = list(string)
        extended_key_usage = optional(list(string), [])
        validity_in_months = number
      }))
    }))

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
