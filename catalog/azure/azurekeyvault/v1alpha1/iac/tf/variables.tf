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
  description = "Azure Key Vault specification"
  type = object({
    # The Azure region the vault lives in.
    region = string

    # The resource group the vault lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The vault's name: 3-24 letters/digits/hyphens, globally unique (it
    # becomes the {name}.vault.azure.net endpoint).
    vault_name = string

    # The pricing tier, as the spec enum's name string (STANDARD/PREMIUM).
    # Unset applies the STANDARD baseline.
    sku = optional(string)

    # Whether data-plane authorization is Azure RBAC (true, the
    # recommended posture and the spec default) or legacy access policies
    # (false).
    rbac_authorization_enabled = optional(bool, true)

    # Legacy access-policy grants -- only honored by Azure when
    # rbac_authorization_enabled is false. Permission lists arrive as the
    # spec enums' FULL name strings (e.g. "KEY_GET", "SECRET_SET") and are
    # mapped to ARM's values in locals.
    access_policies = optional(list(object({
      object_id               = string
      tenant_id               = optional(string)
      application_id          = optional(string)
      key_permissions         = optional(list(string), [])
      secret_permissions      = optional(list(string), [])
      certificate_permissions = optional(list(string), [])
      storage_permissions     = optional(list(string), [])
    })), [])

    # The three resource-manager integration switches (Azure defaults all
    # to false).
    enabled_for_deployment          = optional(bool, false)
    enabled_for_disk_encryption     = optional(bool, false)
    enabled_for_template_deployment = optional(bool, false)

    # Whether the vault accepts public internet connections (Azure
    # defaults to true; false means private-endpoints-only).
    public_network_access_enabled = optional(bool, true)

    # Whether purge protection is on (Azure defaults to false;
    # irreversible once enabled).
    purge_protection_enabled = optional(bool, false)

    # Soft-delete retention in days, 7-90 (Azure defaults to 90; fixed at
    # creation).
    soft_delete_retention_days = optional(number, 90)

    # Network access rules for the public endpoint: default_action and
    # bypass arrive as the spec enums' name strings (ALLOW/DENY,
    # AZURE_SERVICES/NONE); virtual_network_subnet_ids are resolved subnet
    # ARM IDs.
    network_acls = optional(object({
      default_action             = string
      bypass                     = optional(string)
      ip_rules                   = optional(list(string), [])
      virtual_network_subnet_ids = optional(list(string), [])
    }))

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
