# Create the Data Protection backup vault -- the safe that modern
# Azure Backup data (disks, blobs, AKS clusters, flexible-server
# databases, Data Lake storage) lives in. The vault is free at rest;
# cost follows the protected instances and their backup storage.
#
# Destroy note: Azure's delete call returns before the vault is fully
# gone; the provider polls until the name is actually free (its own
# workaround for the service bug), so destroy runs a little longer
# than the API suggests.
resource "azurerm_data_protection_backup_vault" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  # Both ForceNew on the provider; required in the spec.
  datastore_type = var.spec.datastore_type
  redundancy     = var.spec.redundancy

  # Sent ONLY when true: the provider errors when this argument is
  # EXPLICITLY present (even as false) on a non-GeoRedundant vault, so
  # an unset spec value must reach the provider as null, never false.
  # Enabling is in-place; DISABLING replaces the vault (the provider's
  # one-way ForceNew, recorded on the spec field).
  cross_region_restore_enabled = var.spec.cross_region_restore_enabled ? true : null

  # Soft-delete retention window (days). Null lets the provider
  # default (14) apply.
  retention_duration_in_days = var.spec.retention_duration_in_days

  # Null lets the provider defaults apply (On / Disabled). Both carry
  # one-way doors -- AlwaysOn and Locked are permanent (leaving either
  # replaces the vault; the provider's ForceNew transitions, recorded
  # on the spec fields).
  soft_delete  = var.spec.soft_delete
  immutability = var.spec.immutability

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_wire[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  tags = local.final_tags
}

# Customer-managed-key encryption, composed from the provider's
# sibling resource (which rewrites the vault's own security settings
# -- one ARM object, one authoritative shape). Azure unwraps the key
# with the vault's SYSTEM-assigned identity (the provider hardcodes
# it; the spec's CEL requires that identity flavor).
#
# ONE-WAY DOOR: once enabled, CMK can never be removed -- the
# provider's delete for this resource is a documented no-op; only
# deleting the vault removes the encryption. The KEY itself rotates in
# place (the one updatable part), and versionless key URIs (the
# reference's default) make rotation automatic.
resource "azurerm_data_protection_backup_vault_customer_managed_key" "main" {
  count = var.spec.encryption != null ? 1 : 0

  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.main.id
  key_vault_key_id                = var.spec.encryption.key_id
}
