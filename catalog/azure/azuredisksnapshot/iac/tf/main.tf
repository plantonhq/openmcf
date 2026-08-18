# Create the managed disk snapshot. The source fields pair with
# create_option ("Copy" reads source_resource_id; "Import" reads
# source_uri + storage_account_id) -- the provider's own schema does
# not tie them together and Azure validates the pairing at create
# time, so the module sends each source field only when set.
#
# Removing encryption_settings from a snapshot that had them forces
# replacement (the provider's CustomizeDiff -- Azure cannot disable
# encryption in place).
resource "azurerm_snapshot" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  create_option = var.spec.create_option

  source_resource_id = (
    var.spec.source_resource_id != null && var.spec.source_resource_id != ""
    ? var.spec.source_resource_id
    : null
  )
  source_uri = (
    var.spec.source_uri != null && var.spec.source_uri != ""
    ? var.spec.source_uri
    : null
  )
  storage_account_id = (
    var.spec.storage_account_id != null && var.spec.storage_account_id != ""
    ? var.spec.storage_account_id
    : null
  )

  incremental_enabled = var.spec.incremental_enabled

  # Unset inherits the source's size (Azure computes it); the provider
  # sends the value only when > 0.
  disk_size_gb = var.spec.disk_size_gb

  # Unset rides the provider default, "AllowAll".
  network_access_policy = (
    var.spec.network_access_policy != null && var.spec.network_access_policy != ""
    ? var.spec.network_access_policy
    : null
  )
  disk_access_id = (
    var.spec.disk_access_id != null && var.spec.disk_access_id != ""
    ? var.spec.disk_access_id
    : null
  )

  # Unset rides the provider default, true.
  public_network_access_enabled = var.spec.public_network_access_enabled

  dynamic "encryption_settings" {
    for_each = var.spec.encryption_settings != null ? [var.spec.encryption_settings] : []
    content {
      disk_encryption_key {
        secret_url      = encryption_settings.value.disk_encryption_key.secret_url
        source_vault_id = encryption_settings.value.disk_encryption_key.source_vault_id
      }

      dynamic "key_encryption_key" {
        for_each = encryption_settings.value.key_encryption_key != null ? [encryption_settings.value.key_encryption_key] : []
        content {
          key_url         = key_encryption_key.value.key_url
          source_vault_id = key_encryption_key.value.source_vault_id
        }
      }
    }
  }

  # The source fields are create-time-only by contract: a snapshot's
  # creation data is immutable history, and the provider (v5 pin) never
  # reads source_resource_id/source_uri back from Azure -- an adopted
  # (imported) snapshot therefore holds a null source in state, and
  # without this guard every post-import plan proposes a destroy+create
  # that would delete the very artifact the user adopted. Ignoring the
  # pair also means an in-place source edit is a no-op rather than a
  # silent deletion of a backup artifact; capturing a different disk is
  # a NEW snapshot resource (the spec field comments teach this). The
  # Pulumi module carries the same guard -- keep the engines in step.
  lifecycle {
    ignore_changes = [source_resource_id, source_uri]
  }

  tags = local.final_tags
}
