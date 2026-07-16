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
  description = "Azure Event Hubs Namespace Customer Managed Key specification"
  type = object({
    # The namespace to encrypt, by ARM ID. References are resolved to a
    # literal by the platform before the module runs. The namespace must
    # have single-tenant capacity (a dedicated cluster or PREMIUM) and
    # already carry the unwrapping identity. ForceNew.
    eventhub_namespace_id = string

    # The Key Vault keys that encrypt the namespace's data, by
    # data-plane key ID (1-10; references resolved to literals before
    # the module runs). Versionless IDs make vault-side rotation
    # propagate automatically.
    key_vault_key_ids = list(string)

    # A second infrastructure-encryption layer beneath the customer
    # keys. ForceNew -- fixed the moment CMK is first configured.
    infrastructure_encryption_enabled = optional(bool)

    # The user-assigned identity that unwraps the keys, by ARM ID. Must
    # already be attached via the namespace's identity block with
    # wrap/unwrap vault access. Unset uses the namespace's
    # system-assigned identity.
    user_assigned_identity_id = optional(string)
  })
}
