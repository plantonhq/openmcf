# Create the AI Foundry hub -- ARM-wise an ML workspace of kind "Hub":
# the shared foundation whose security/storage/network posture every
# AzureAiFoundryProject inside it inherits. The key vault and storage
# attachments are ForceNew; insights and registry update in place
# (unlike the classic ML workspace, where the registry is ForceNew).
#
# Deletion is a SOFT delete: the hub becomes a purgeable ghost that
# keeps holding its name (the ML workspace class). The provider purges
# it when the features flag
# `machine_learning.purge_soft_deleted_workspace_on_destroy` is enabled.
resource "azurerm_ai_foundry" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # The two required companion services (both ForceNew).
  key_vault_id       = var.spec.key_vault_id
  storage_account_id = var.spec.storage_account_id

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_wire[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  # Attachable and re-pointable in place.
  application_insights_id = var.spec.application_insights_id != "" ? var.spec.application_insights_id : null
  container_registry_id   = var.spec.container_registry_id != "" ? var.spec.container_registry_id : null

  primary_user_assigned_identity = var.spec.primary_user_assigned_identity != "" ? var.spec.primary_user_assigned_identity : null

  # The spec's bool onto the provider's two-value string. Unset (null)
  # means the proto default, true -- the provider's own default,
  # "Enabled".
  public_network_access = coalesce(var.spec.public_network_access_enabled, true) ? "Enabled" : "Disabled"

  # Customer-managed-key encryption; the whole block is ForceNew. The
  # key id is a VERSIONED Key Vault key URL -- the provider's hub
  # contract (versionless is rejected; rotation does not
  # auto-propagate, unlike the classic ML workspace).
  dynamic "encryption" {
    for_each = var.spec.encryption != null ? [var.spec.encryption] : []
    content {
      key_vault_id              = encryption.value.key_vault_id
      key_id                    = encryption.value.key_id
      user_assigned_identity_id = encryption.value.user_assigned_identity_id != "" ? encryption.value.user_assigned_identity_id : null
    }
  }

  # The managed virtual network. isolation_mode is Optional+Computed on
  # the provider -- unspecified omits it and the value is read back.
  dynamic "managed_network" {
    for_each = var.spec.managed_network != null ? [var.spec.managed_network] : []
    content {
      isolation_mode = lookup(local.isolation_mode_wire, managed_network.value.isolation_mode, null)
    }
  }

  # Sent only when true (both engines): the property is
  # Optional+Computed and the SERVICE flips it true when encryption is
  # enabled -- a pinned false would fight that read-back. ForceNew.
  high_business_impact_enabled = var.spec.high_business_impact_enabled ? true : null

  description   = var.spec.description != "" ? var.spec.description : null
  friendly_name = var.spec.friendly_name != "" ? var.spec.friendly_name : null

  tags = local.final_tags
}
