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
  description = "Azure Event Hubs Namespace specification"
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

    # Throughput units on BASIC/STANDARD (1-40); processing units on
    # PREMIUM (1, 2, 4, 8, 16). Unset keeps Azure's default (1).
    capacity = optional(number)

    # STANDARD's elastic throughput scaling -- grows TUs up to the
    # ceiling below under load (never shrinks them back).
    auto_inflate_enabled = optional(bool)

    # The TU ceiling auto-inflate may reach (0-40). Azure validates the
    # auto-inflate pairing at apply time.
    maximum_throughput_units = optional(number)

    # The dedicated Event Hubs cluster to place the namespace on, by ARM
    # id (resolved to a literal before the module runs). ForceNew.
    dedicated_cluster_id = optional(string)

    # Managed identity for the namespace (required for identity-based
    # capture auth and CMK).
    identity = optional(object({
      # SYSTEM_ASSIGNED, USER_ASSIGNED, or SYSTEM_AND_USER_ASSIGNED.
      type = string
      # AzureUserAssignedIdentity ARM ids, resolved to literals before
      # the module runs.
      user_assigned_identity_ids = optional(list(string), [])
    }))

    # Whether SAS authentication is allowed (Azure's default true).
    # False = keyless posture: Entra-only data-plane auth.
    local_authentication_enabled = optional(bool, true)

    # Whether the namespace accepts public-internet traffic. Must agree
    # with the rule set's own dial when that block is present.
    public_network_access_enabled = optional(bool, true)

    # The namespace firewall (not available on BASIC).
    network_rule_sets = optional(object({
      # ALLOW or DENY for unmatched traffic -- required by Azure when
      # the rule set is declared (the spec enum enforces an explicit
      # choice).
      default_action = string
      # Whether admitted traffic may arrive over public IP space; must
      # equal the namespace-level dial (spec CEL front-loads the pair).
      public_network_access_enabled = optional(bool, true)
      # Whether trusted Microsoft services bypass the firewall.
      trusted_service_access_enabled = optional(bool)
      # Admitted public IPv4 addresses/CIDRs (each entry is an allow
      # rule -- Azure's per-rule action has exactly one legal value).
      ip_rules = optional(list(string), [])
      # Admitted VNet subnets (service endpoints).
      virtual_network_rules = optional(list(object({
        subnet_id                                       = string
        ignore_missing_virtual_network_service_endpoint = optional(bool)
      })), [])
    }))

    # User tags, merged over the Planton-derived identity tags (user
    # values win on key conflicts).
    tags = optional(map(string), {})
  })
}
