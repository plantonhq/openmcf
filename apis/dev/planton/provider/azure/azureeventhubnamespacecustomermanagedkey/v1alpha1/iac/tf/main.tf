# Customer-managed-key (BYOK) encryption applied ONTO an existing Event
# Hubs namespace -- Azure models CMK as a namespace property configured
# after creation, not a create-time block, and this resource mirrors
# that grain. The namespace must have single-tenant capacity (a
# dedicated cluster or PREMIUM) or Azure rejects the encryption patch.
#
# ADD-ONLY lifecycle (Azure's own contract): once CMK is enabled it can
# never be removed -- Azure has no decrypt-back path. The provider's
# Delete is deliberately a NO-OP: destroying this resource changes
# NOTHING on the namespace, and returning to Microsoft-managed keys
# requires replacing the namespace itself.
#
# Identity contract: a user-assigned identity named here must ALREADY be
# attached to the parent namespace's identity block, with wrap/unwrap
# access on the keys' vault -- Azure rejects the patch otherwise. Unset
# falls back to the namespace's system-assigned identity (grant IT the
# vault access).
resource "azurerm_eventhub_namespace_customer_managed_key" "main" {
  # ForceNew: the configuration is bound to its namespace for life.
  eventhub_namespace_id = var.spec.eventhub_namespace_id

  # Versionless key IDs make vault-side rotation propagate
  # automatically; pin versioned IDs only when a compliance regime
  # demands immutable key versions.
  key_vault_key_ids = var.spec.key_vault_key_ids

  # ForceNew: the second encryption layer is fixed the moment CMK is
  # first configured. Unset passes null, matching Azure's default
  # (false).
  infrastructure_encryption_enabled = var.spec.infrastructure_encryption_enabled

  user_assigned_identity_id = local.user_assigned_identity_id
}
