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
  description = "Azure Container Registry specification"
  type = object({
    # The Azure region the registry's home replica lives in.
    region = string

    # The resource group the registry lives in. References are resolved to
    # a literal name by the platform before the module runs.
    resource_group = string

    # The registry's name: 5-50 lowercase alphanumerics, globally unique
    # (it becomes the {name}.azurecr.io login server).
    registry_name = string

    # The pricing tier, as the spec enum's name string
    # (BASIC/STANDARD/PREMIUM). Unset applies the STANDARD baseline.
    sku = optional(string)

    # Whether the built-in admin account (static username/password) is
    # enabled. Azure defaults to false; Entra-based auth is the production
    # path.
    admin_user_enabled = optional(bool, false)

    # Whether the registry accepts public internet connections (Azure
    # defaults to true; false requires Premium and means
    # private-endpoints-only).
    public_network_access_enabled = optional(bool, true)

    # Whether the home replica is spread across availability zones
    # (Premium; fixed at creation).
    zone_redundancy_enabled = optional(bool, false)

    # Whether unauthenticated pulls are allowed (Standard/Premium; makes
    # every repository publicly readable).
    anonymous_pull_enabled = optional(bool, false)

    # Whether the registry gets dedicated regional data endpoints for
    # exact firewall allowlisting (Premium).
    data_endpoint_enabled = optional(bool, false)

    # Whether newly pushed images are quarantined until scanning tooling
    # passes them (Premium).
    quarantine_policy_enabled = optional(bool, false)

    # Days after which untagged manifests are purged (Premium; unset keeps
    # them forever, Azure's default).
    retention_policy_in_days = optional(number)

    # Whether Docker Content Trust (image signing) is enabled (Premium;
    # Azure defaults to false).
    trust_policy_enabled = optional(bool, false)

    # Whether artifacts can be exported out of the registry (Azure
    # defaults to true; disabling requires Premium + public access off).
    export_policy_enabled = optional(bool, true)

    # Whether trusted Azure services bypass network restrictions, as the
    # spec enum's name string (AZURE_SERVICES/NONE). Unset applies Azure's
    # default (AzureServices).
    network_rule_bypass_option = optional(string)

    # Network access rules for a public registry (Premium): default action
    # (ALLOW/DENY enum name string) plus an IPv4 CIDR allowlist.
    network_rule_set = optional(object({
      default_action = optional(string)
      ip_rules = optional(list(object({
        ip_range = string
      })), [])
    }))

    # Additional regions the registry replicates to (Premium; must not
    # contain the home region).
    georeplications = optional(list(object({
      location                  = string
      zone_redundancy_enabled   = optional(bool, false)
      regional_endpoint_enabled = optional(bool, false)
      tags                      = optional(map(string), {})
    })), [])

    # The registry's managed identity: type is the spec enum's name string
    # (SYSTEM_ASSIGNED/USER_ASSIGNED/SYSTEM_AND_USER_ASSIGNED);
    # identity_ids are resolved user-assigned-identity ARM IDs.
    identity = optional(object({
      type         = string
      identity_ids = optional(list(string), [])
    }))

    # Customer-managed-key encryption (Premium; fixed at creation):
    # identity_client_id is the resolved client id of the unwrapping
    # user-assigned identity; key_vault_key_id is the Key Vault key ID.
    encryption = optional(object({
      identity_client_id = string
      key_vault_key_id   = string
    }))

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
