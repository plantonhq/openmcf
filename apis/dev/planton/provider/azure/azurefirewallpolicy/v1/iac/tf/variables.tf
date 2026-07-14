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
  description = "Azure Firewall Policy specification"
  type = object({
    # The Azure region the policy lives in. Regional resource, attachable
    # to firewalls in any region.
    region = string

    # The resource group the policy is created in. References are resolved
    # to a literal name by the platform before the module runs.
    resource_group = string

    # The policy's name, unique within the resource group. Renaming
    # replaces the policy.
    name = string

    # The tier as the proto enum value name (BASIC/STANDARD/PREMIUM).
    # Absent means STANDARD. Fixed at creation; must match the tier of
    # every attached firewall.
    sku = optional(string)

    # The parent policy's ARM id for inheritance (resolved literal).
    base_policy_id = optional(string)

    # Threat-intelligence posture as the proto enum value name
    # (ALERT/DENY/OFF). Absent means ALERT, Azure's default.
    threat_intelligence_mode = optional(string)

    # Traffic threat intelligence must never flag.
    threat_intelligence_allowlist = optional(object({
      ip_addresses = optional(list(string), [])
      fqdns        = optional(list(string), [])
    }))

    # DNS settings for attached firewalls. proxy_enabled is required for
    # FQDN network rules to resolve deterministically.
    dns = optional(object({
      servers       = optional(list(string), [])
      proxy_enabled = optional(bool, false)
    }))

    # Premium-only IDPS configuration. State/mode values arrive as proto
    # enum value names (IDPS_OFF/IDPS_ALERT/IDPS_DENY).
    intrusion_detection = optional(object({
      mode = optional(string)
      signature_overrides = optional(list(object({
        id    = optional(string)
        state = optional(string)
      })), [])
      private_ranges = optional(list(string), [])
      traffic_bypass = optional(list(object({
        name                  = string
        description           = optional(string)
        protocol              = string
        source_addresses      = optional(list(string), [])
        source_ip_groups      = optional(list(string), [])
        destination_addresses = optional(list(string), [])
        destination_ip_groups = optional(list(string), [])
        destination_ports     = optional(list(string), [])
      })), [])
    }))

    # The policy's managed identity. type arrives as the proto enum value
    # name (SYSTEM_ASSIGNED/USER_ASSIGNED/SYSTEM_AND_USER_ASSIGNED); ids
    # are resolved ARM ids. TLS inspection requires a user-assigned
    # identity with Key Vault secret read access.
    identity = optional(object({
      type                       = string
      user_assigned_identity_ids = optional(list(string), [])
    }))

    # Premium-only TLS inspection: the CA certificate's Key Vault SECRET
    # id (resolved literal; versionless follows renewals) and a display
    # name.
    tls_certificate = optional(object({
      key_vault_secret_id = string
      name                = string
    }))

    # Firewall Policy Analytics wiring into Log Analytics.
    insights = optional(object({
      enabled                            = bool
      default_log_analytics_workspace_id = string
      retention_in_days                  = optional(number)
      log_analytics_workspaces = optional(list(object({
        workspace_id      = string
        firewall_location = string
      })), [])
    }))

    # Explicit forward proxy settings. The azurerm provider caps proxy
    # ports at 35536 (its published validation bound).
    explicit_proxy = optional(object({
      enabled         = optional(bool, false)
      http_port       = optional(number)
      https_port      = optional(number)
      enable_pac_file = optional(bool, false)
      pac_file_port   = optional(number)
      pac_file        = optional(string)
    }))

    # Allow SQL redirect ports (11000-11999, 14000-14999) in FQDN network
    # rules.
    sql_redirect_allowed = optional(bool, false)

    # Ranges the firewall treats as private (never SNATed). CIDRs or
    # single IPv4 addresses.
    private_ip_ranges = optional(list(string), [])

    # Auto-learn private ranges. Azure only records "Enabled"; disabling
    # is by omission, so the module sends the flag only when true.
    auto_learn_private_ranges_enabled = optional(bool, false)

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
