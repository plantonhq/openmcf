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
  description = "Azure Disk Encryption Set specification"
  type = object({
    # The region and resource group the set is created in. References are
    # resolved to literal names by the platform before the module runs.
    region         = string
    resource_group = string
    name           = string

    # The Key Vault key URL that encrypts disks. Versionless when
    # auto_key_rotation_enabled is true, versioned when false; the provider
    # validates the pairing at apply.
    key_vault_key_id = string

    # Follow the key's latest version automatically (Azure default false).
    auto_key_rotation_enabled = optional(bool)

    # What the set encrypts, as the spec enum's name string
    # (ENCRYPTION_AT_REST_WITH_CUSTOMER_KEY /
    # ENCRYPTION_AT_REST_WITH_PLATFORM_AND_CUSTOMER_KEYS /
    # CONFIDENTIAL_VM_ENCRYPTED_WITH_CUSTOMER_KEY). Unset lets Azure default.
    encryption_type = optional(string)

    # Multi-tenant app client id for cross-tenant key access; empty for the
    # same-tenant case.
    federated_client_id = optional(string, "")

    # The set's managed identity: type is the spec enum name string
    # (SYSTEM_ASSIGNED / USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED); ids are
    # the resolved user-assigned identity ARM IDs.
    identity = object({
      type         = string
      identity_ids = optional(list(string), [])
    })

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}
