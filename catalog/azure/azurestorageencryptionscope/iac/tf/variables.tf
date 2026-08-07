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
  description = "Azure Storage Encryption Scope specification"
  type = object({
    # The storage account the scope lives in. References are resolved
    # to a literal ARM ID by the platform before the module runs.
    storage_account_id = string

    # The scope's name: 4-63 letters and digits (no hyphens); unique
    # within the account.
    scope_name = string

    # Who owns the scope's key, as the spec enum's name string
    # (MICROSOFT_STORAGE, MICROSOFT_KEY_VAULT).
    source = string

    # The Key Vault key the scope encrypts under -- required when source
    # is MICROSOFT_KEY_VAULT (enforced in the spec). References are
    # resolved to a literal key URI by the platform before the module
    # runs; the versionless URI lets rotation propagate.
    key_vault_key_id = optional(string)

    # Whether data under this scope is double-encrypted with a second
    # platform-managed layer. Fixed at creation.
    infrastructure_encryption_required = optional(bool)
  })
}
