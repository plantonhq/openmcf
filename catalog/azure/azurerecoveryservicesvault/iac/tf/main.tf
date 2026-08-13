# Create the Recovery Services vault -- the safe that classic Azure
# Backup data and Site Recovery configuration live in. The vault is
# free at rest; cost follows the protected items and their backup
# storage.
#
# Destroy semantics kept deliberately at the engines' defaults:
# deleting the vault FAILS while protected items remain inside it (the
# provider's purge_protected_items_from_vault_on_destroy feature stays
# off in provider.tf) -- stop and delete protections first, then the
# vault.
resource "azurerm_recovery_services_vault" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # The platform default fills the sku ("Standard"); the null guard
  # keeps direct module invocations safe. "RS0" is the legacy
  # tier-style spelling ARM also accepts, priced the same.
  sku = var.spec.sku != null ? var.spec.sku : "Standard"

  # Optional-with-default on the provider: null lets the provider
  # default apply (GeoRedundant / true).
  storage_mode_type             = var.spec.storage_mode_type
  public_network_access_enabled = var.spec.public_network_access_enabled

  # Plain bool: false is the provider's own default, so passing the
  # zero value through is exact. Enabling is in-place; DISABLING
  # replaces the vault (the provider's one-way ForceNew, recorded on
  # the spec field).
  cross_region_restore_enabled = var.spec.cross_region_restore_enabled

  # Optional+Computed on the provider: unset leaves the service default
  # (Disabled). Transitions are limited (Locked <- Unlocked <->
  # Disabled) and Locked is permanent -- recorded on the spec field;
  # the provider stages Locked through Unlocked itself.
  immutability = var.spec.immutability != "" ? var.spec.immutability : null

  # ForceNew on the provider; false is its default.
  classic_vmware_replication_enabled = var.spec.classic_vmware_replication_enabled

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_wire[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  # Customer-managed-key encryption. Once enabled it can never be
  # disabled, infrastructure_encryption_enabled can never change, and
  # the sku freezes (the provider's own update guards -- recorded on
  # the spec fields). Versionless key URIs are accepted (the provider
  # validates VersionTypeAny) and rotate automatically -- the
  # reference's default target.
  dynamic "encryption" {
    for_each = var.spec.encryption != null ? [var.spec.encryption] : []
    content {
      key_id = encryption.value.key_id
      # ARM requires an explicit choice; the plain bool always ships.
      infrastructure_encryption_enabled = encryption.value.infrastructure_encryption_enabled
      use_system_assigned_identity      = encryption.value.use_system_assigned_identity
      user_assigned_identity_id         = encryption.value.user_assigned_identity_id != "" ? encryption.value.user_assigned_identity_id : null
    }
  }

  # The vault's built-in Azure Monitor alert settings. Every switch
  # defaults ON both provider- and service-side; null passes each
  # unset switch to the provider default. The three v5-new switches
  # here are the trio the Pulumi engine cannot express (its
  # PARITY-EXCEPTION) -- THIS module honors them fully.
  dynamic "monitoring" {
    for_each = var.spec.monitoring != null ? [var.spec.monitoring] : []
    content {
      alerts_for_all_job_failures_enabled            = monitoring.value.alerts_for_all_job_failures_enabled
      alerts_for_all_failover_issues_enabled         = monitoring.value.alerts_for_all_failover_issues_enabled
      alerts_for_all_replication_issues_enabled      = monitoring.value.alerts_for_all_replication_issues_enabled
      alerts_for_critical_operation_failures_enabled = monitoring.value.alerts_for_critical_operation_failures_enabled
      email_notifications_for_site_recovery_enabled  = monitoring.value.email_notifications_for_site_recovery_enabled
    }
  }

  tags = local.final_tags
}

# The composed Resource Guard association (Multi-User Authorization):
# privileged vault operations then require an approval through the
# guard. ARM pins the association's own name to the literal
# "VaultProxy" -- one guard per vault.
resource "azurerm_recovery_services_vault_resource_guard_association" "main" {
  count = var.spec.resource_guard_id != "" ? 1 : 0

  vault_id          = azurerm_recovery_services_vault.main.id
  resource_guard_id = var.spec.resource_guard_id
}
