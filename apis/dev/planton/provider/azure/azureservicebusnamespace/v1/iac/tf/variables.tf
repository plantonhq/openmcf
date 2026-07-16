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
  description = "Azure Service Bus Namespace specification"
  type = object({
    # The Azure region for the namespace (e.g. "eastus").
    region = string

    # The resource group name. References are resolved to a literal by
    # the platform before the module runs.
    resource_group = string

    # The globally-unique namespace name -- becomes the endpoint
    # {name}.servicebus.windows.net. ForceNew.
    namespace_name = string

    # The pricing tier, as the spec enum's value name (BASIC, STANDARD,
    # PREMIUM). Unset deploys STANDARD.
    sku = optional(string)

    # PREMIUM messaging units (1, 2, 4, 8, 16). Required with PREMIUM;
    # must be absent otherwise (the spec CELs enforce the pairing).
    capacity = optional(number)

    # PREMIUM namespace partitions (1, 2, 4). Required with PREMIUM;
    # ForceNew -- the layout is fixed at creation.
    premium_messaging_partitions = optional(number)

    # Managed identity for the namespace (required for CMK).
    identity = optional(object({
      # SYSTEM_ASSIGNED, USER_ASSIGNED, or SYSTEM_AND_USER_ASSIGNED.
      type = string
      # AzureUserAssignedIdentity ARM ids, resolved to literals before
      # the module runs.
      user_assigned_identity_ids = optional(list(string), [])
    }))

    # Customer-managed-key encryption (PREMIUM only). Once set, removal
    # replaces the namespace (Azure's own contract).
    customer_managed_key = optional(object({
      # The Key Vault key data-plane ID (versionless for
      # rotation-follows-latest, or a pinned version).
      key_vault_key_id = string
      # The user-assigned identity that unwraps the key -- must also be
      # attached via the identity block.
      user_assigned_identity_id = string
      # A second infrastructure-encryption layer. ForceNew.
      infrastructure_encryption_enabled = optional(bool)
    }))

    # Whether SAS authentication is allowed (Azure's default true).
    # False = keyless posture: Entra-only data-plane auth.
    local_auth_enabled = optional(bool, true)

    # Whether the namespace accepts public-internet traffic.
    public_network_access_enabled = optional(bool, true)

    # The namespace firewall (PREMIUM only).
    network_rule_set = optional(object({
      # ALLOW or DENY for unmatched traffic; unset keeps Azure's open
      # default (ALLOW).
      default_action = optional(string)
      # Whether admitted traffic may arrive over public IP space.
      public_network_access_enabled = optional(bool, true)
      # Whether trusted Microsoft services bypass the firewall.
      trusted_services_allowed = optional(bool)
      # Admitted public IPv4 addresses/CIDRs.
      ip_rules = optional(list(string), [])
      # Admitted VNet subnets (service endpoints).
      network_rules = optional(list(object({
        subnet_id                            = string
        ignore_missing_vnet_service_endpoint = optional(bool)
      })), [])
    }))

    # User tags, merged over the Planton-derived identity tags (user
    # values win on key conflicts).
    tags = optional(map(string), {})
  })
}
