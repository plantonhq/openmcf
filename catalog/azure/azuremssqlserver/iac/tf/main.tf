# The deploying credential's context -- the tenant fallback for the
# Microsoft Entra administrator grant when the spec does not pin one.
data "azurerm_client_config" "current" {}

# The Azure SQL logical server: the administrative container carrying the
# login endpoint, authentication, TDE key, and networking posture.
# Databases (AzureMssqlDatabase) and elastic pools (AzureMssqlElasticPool)
# are first-class resources referencing this server's ARM id. The azurerm
# resource internally orchestrates the connection-policy and Entra-admin
# ARM sub-APIs, so both are plain arguments here.
resource "azurerm_mssql_server" "main" {
  name                = var.spec.server_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  version             = var.spec.version

  # SQL-auth credentials -- omitted (null) on an Entra-only server. The
  # login is fixed once set; ARM rejects a password change while
  # azuread_authentication_only is true.
  administrator_login          = local.administrator_login
  administrator_login_password = local.administrator_password

  # The Microsoft Entra administrator. With azuread_authentication_only
  # SQL logins are disabled server-wide.
  dynamic "azuread_administrator" {
    for_each = var.spec.azuread_administrator != null ? [var.spec.azuread_administrator] : []
    content {
      login_username              = azuread_administrator.value.login_username
      object_id                   = azuread_administrator.value.object_id
      tenant_id                   = local.aad_tenant_id
      azuread_authentication_only = azuread_administrator.value.azuread_authentication_only
    }
  }

  # The server's managed identity -- unwraps the TDE customer-managed
  # key. The primary user-assigned identity is the one ARM uses for Key
  # Vault access when several are attached.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type = local.identity_type_map[identity.value.type]
      # An empty list must become null: azurerm rejects identity_ids
      # alongside a pure SystemAssigned identity (the Pulumi module
      # omits the field identically).
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }
  primary_user_assigned_identity_id = var.spec.primary_user_assigned_identity_id

  # Server-level TDE customer-managed key (VERSIONED Key Vault key id --
  # ARM pins the exact version at the server level).
  transparent_data_encryption_key_vault_key_id = var.spec.transparent_data_encryption_key_vault_key_id

  # Unset maps to null so Azure applies its Default policy (Redirect
  # inside Azure, Proxy outside).
  connection_policy = local.connection_policy

  minimum_tls_version           = var.spec.minimum_tls_version
  public_network_access_enabled = var.spec.public_network_access_enabled

  # Outbound restriction: the allowlist itself is the
  # azurerm_mssql_outbound_firewall_rule resources below.
  outbound_network_restriction_enabled = var.spec.outbound_network_restriction_enabled

  # Microsoft Defender's agentless SQL scanning (no storage account
  # needed, unlike the classic vulnerability assessment).
  express_vulnerability_assessment_enabled = var.spec.express_vulnerability_assessment_enabled

  tags = local.final_tags
}

# Public-endpoint firewall allowlist. Only meaningful while public
# network access is enabled. 0.0.0.0-0.0.0.0 admits Azure-internal
# services only.
resource "azurerm_mssql_firewall_rule" "main" {
  for_each = { for rule in var.spec.firewall_rules : rule.name => rule }

  name             = each.value.name
  server_id        = azurerm_mssql_server.main.id
  start_ip_address = each.value.start_ip_address
  end_ip_address   = each.value.end_ip_address
}

# Subnet allowlist through Microsoft.Sql service endpoints -- the classic
# (non-Private-Link) way to keep traffic on the Azure backbone.
resource "azurerm_mssql_virtual_network_rule" "main" {
  for_each = { for rule in var.spec.virtual_network_rules : rule.name => rule }

  name                                 = each.value.name
  server_id                            = azurerm_mssql_server.main.id
  subnet_id                            = each.value.subnet_id
  ignore_missing_vnet_service_endpoint = each.value.ignore_missing_vnet_service_endpoint
}

# The FQDNs the server may reach OUT to while outbound restriction is
# enabled (elastic queries, linked external tables).
resource "azurerm_mssql_outbound_firewall_rule" "main" {
  for_each = toset(var.spec.outbound_firewall_rules)

  name      = each.value
  server_id = azurerm_mssql_server.main.id
}

# Server-level SQL auditing: audit events for every database on the
# server go to blob storage (blob_storage_endpoint + key) and/or Azure
# Monitor (log_monitoring_enabled, consumed through diagnostic settings).
resource "azurerm_mssql_server_extended_auditing_policy" "main" {
  count = var.spec.extended_auditing != null ? 1 : 0

  server_id                               = azurerm_mssql_server.main.id
  blob_storage_endpoint                   = var.spec.extended_auditing.storage_endpoint
  storage_account_access_key              = var.spec.extended_auditing.storage_account_access_key
  storage_account_access_key_is_secondary = var.spec.extended_auditing.storage_account_access_key_is_secondary
  retention_in_days                       = var.spec.extended_auditing.retention_in_days
  log_monitoring_enabled                  = var.spec.extended_auditing.log_monitoring_enabled
  storage_account_subscription_id         = var.spec.extended_auditing.storage_account_subscription_id
  predicate_expression                    = var.spec.extended_auditing.predicate_expression
  audit_actions_and_groups                = length(var.spec.extended_auditing.audit_actions_and_groups) > 0 ? var.spec.extended_auditing.audit_actions_and_groups : null
}

# Microsoft Defender for SQL threat detection at the server scope. The
# azurerm resource addresses the server by name + resource group rather
# than by id -- its own contract, not a choice here.
resource "azurerm_mssql_server_security_alert_policy" "main" {
  count = var.spec.security_alert_policy != null ? 1 : 0

  resource_group_name          = var.spec.resource_group
  server_name                  = azurerm_mssql_server.main.name
  state                        = local.alert_policy_state_map[var.spec.security_alert_policy.state]
  disabled_alerts              = [for alert in var.spec.security_alert_policy.disabled_alerts : local.alert_type_map[alert]]
  email_account_admins_enabled = var.spec.security_alert_policy.email_account_admins
  email_addresses              = var.spec.security_alert_policy.email_addresses
  retention_days               = var.spec.security_alert_policy.retention_days
  storage_endpoint             = var.spec.security_alert_policy.storage_endpoint
  storage_account_access_key   = var.spec.security_alert_policy.storage_account_access_key
}
