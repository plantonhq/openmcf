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
  description = "Azure SQL Database logical server specification"
  type = object({
    # The Azure region the server lives in.
    region = string

    # The resource group the server lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The server's name: 1-63 lowercase letters/digits/hyphens, globally
    # unique (it becomes {server_name}.database.windows.net).
    server_name = string

    # The SQL Server version identifier ("2.0" legacy, "12.0" current).
    version = optional(string, "12.0")

    # SQL-auth admin credentials. Empty on an Entra-only server
    # (azuread_administrator.azuread_authentication_only = true). The
    # login is fixed once set; ARM rejects a password change while
    # Entra-only auth is on.
    administrator_login    = optional(string)
    administrator_password = optional(string)

    # The server's Microsoft Entra administrator; with
    # azuread_authentication_only it becomes the ONLY auth mechanism.
    # tenant_id falls back to the deploying credential's tenant.
    azuread_administrator = optional(object({
      login_username              = string
      object_id                   = string
      tenant_id                   = optional(string)
      azuread_authentication_only = optional(bool, false)
    }))

    # The server's managed identity: type arrives as the spec enum's name
    # string (SYSTEM_ASSIGNED / USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED)
    # with the user-assigned identity ARM IDs.
    identity = optional(object({
      type         = string
      identity_ids = optional(list(string), [])
    }))

    # Which attached user-assigned identity ARM uses for Key Vault
    # access; required when the identity type includes USER_ASSIGNED.
    primary_user_assigned_identity_id = optional(string)

    # The server-level TDE customer-managed key, by VERSIONED Key Vault
    # key ID.
    transparent_data_encryption_key_vault_key_id = optional(string)

    # Connection policy, as the spec enum's name string
    # (DEFAULT/PROXY/REDIRECT). Unset is not sent (Azure applies
    # Default).
    connection_policy = optional(string)

    # The TLS floor ("1.2" is the only accepted value on current API
    # versions). ARM rejects removing it once set.
    minimum_tls_version = optional(string, "1.2")

    # Whether the server is reachable over the public internet.
    public_network_access_enabled = optional(bool, true)

    # Outbound restriction: when enabled, the server may only reach OUT
    # to the FQDNs below (each its own ARM sub-resource).
    outbound_network_restriction_enabled = optional(bool, false)
    outbound_firewall_rules              = optional(list(string), [])

    # Microsoft Defender's agentless SQL scanning.
    express_vulnerability_assessment_enabled = optional(bool, false)

    # Public-endpoint firewall allowlist, one Azure sub-resource each.
    firewall_rules = optional(list(object({
      name             = string
      start_ip_address = string
      end_ip_address   = string
    })), [])

    # Subnet allowlist through Microsoft.Sql service endpoints, one Azure
    # sub-resource each.
    virtual_network_rules = optional(list(object({
      name                                 = string
      subnet_id                            = string
      ignore_missing_vnet_service_endpoint = optional(bool, false)
    })), [])

    # Server-level SQL auditing to blob storage and/or Azure Monitor.
    extended_auditing = optional(object({
      storage_endpoint                        = optional(string)
      storage_account_access_key              = optional(string)
      storage_account_access_key_is_secondary = optional(bool, false)
      retention_in_days                       = optional(number, 0)
      log_monitoring_enabled                  = optional(bool, true)
      storage_account_subscription_id         = optional(string)
      predicate_expression                    = optional(string)
      audit_actions_and_groups                = optional(list(string), [])
    }))

    # Microsoft Defender threat detection at the server scope. state and
    # disabled_alerts arrive as spec enum name strings.
    security_alert_policy = optional(object({
      state                      = string
      disabled_alerts            = optional(list(string), [])
      email_account_admins       = optional(bool, false)
      email_addresses            = optional(list(string), [])
      retention_days             = optional(number, 0)
      storage_endpoint           = optional(string)
      storage_account_access_key = optional(string)
    }))

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
