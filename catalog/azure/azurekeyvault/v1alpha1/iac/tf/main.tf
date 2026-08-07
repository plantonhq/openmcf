# The vault authenticates against the deploying credential's Azure AD
# tenant -- a vault cannot be managed cross-tenant, so the tenant is read
# from the ambient client configuration instead of being modeled as a
# contradictable spec field.
data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "main" {
  name                = var.spec.vault_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = local.sku_name

  # RBAC is the spec's recommended default; access_policy blocks below are
  # only honored by Azure when this is false (ARM stores but ignores
  # policies on an RBAC-mode vault).
  rbac_authorization_enabled = var.spec.rbac_authorization_enabled

  # Legacy access-policy grants. Declared inline (not via the standalone
  # access-policy resource) so the vault owns its complete grant list --
  # azurerm cannot mix the two shapes on one vault without perpetual
  # drift.
  dynamic "access_policy" {
    for_each = var.spec.access_policies
    content {
      # An unset tenant falls back to the vault's own tenant -- access
      # policies cannot span tenants in practice.
      tenant_id      = access_policy.value.tenant_id != null ? access_policy.value.tenant_id : data.azurerm_client_config.current.tenant_id
      object_id      = access_policy.value.object_id
      application_id = access_policy.value.application_id

      # Permission lists arrive as FULL proto enum names; the exhaustive
      # maps in locals translate them to ARM's data-plane strings.
      key_permissions         = [for p in access_policy.value.key_permissions : local.key_permission_map[p]]
      secret_permissions      = [for p in access_policy.value.secret_permissions : local.secret_permission_map[p]]
      certificate_permissions = [for p in access_policy.value.certificate_permissions : local.certificate_permission_map[p]]
      storage_permissions     = [for p in access_policy.value.storage_permissions : local.storage_permission_map[p]]
    }
  }

  # Resource-manager integration switches (Azure defaults: all false).
  enabled_for_deployment          = var.spec.enabled_for_deployment
  enabled_for_disk_encryption     = var.spec.enabled_for_disk_encryption
  enabled_for_template_deployment = var.spec.enabled_for_template_deployment

  public_network_access_enabled = var.spec.public_network_access_enabled

  # Purge protection is irreversible once enabled; with it on, destroying
  # the vault schedules deletion for the end of the soft-delete retention
  # window instead of purging (the provider's default behavior purges
  # soft-deleted vaults on destroy when purge protection is off).
  purge_protection_enabled   = var.spec.purge_protection_enabled
  soft_delete_retention_days = var.spec.soft_delete_retention_days

  # Public-endpoint firewall. azurerm requires default_action and bypass
  # whenever the block is present; the bypass fallback to AzureServices
  # (Azure's own default) is materialized in locals.
  dynamic "network_acls" {
    for_each = var.spec.network_acls != null ? [var.spec.network_acls] : []
    content {
      default_action             = local.network_acls_default_action
      bypass                     = local.network_acls_bypass
      ip_rules                   = network_acls.value.ip_rules
      virtual_network_subnet_ids = network_acls.value.virtual_network_subnet_ids
    }
  }

  tags = local.final_tags
}
