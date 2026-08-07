# Create the disk encryption set -- the binding of a customer-managed Key
# Vault key to managed disks for encryption at rest. The set's identity must
# be granted crypto access on the key (out of band, or via a role assignment
# fixture) before disks can use it, and the referenced vault must have purge
# protection enabled (Azure requires it for disk encryption).
resource "azurerm_disk_encryption_set" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Versionless key URL when auto-rotation is on, versioned when off; the
  # provider validates the pairing at apply.
  key_vault_key_id = var.spec.key_vault_key_id

  # Azure defaults auto-rotation to false and encryption_type to
  # EncryptionAtRestWithCustomerKey; null lets the provider apply the
  # defaults so an unspecified spec deploys identically on both engines.
  auto_key_rotation_enabled = var.spec.auto_key_rotation_enabled
  encryption_type           = local.encryption_type
  federated_client_id       = var.spec.federated_client_id != "" ? var.spec.federated_client_id : null

  identity {
    type = local.identity_type
    # Required for the user-assigned flavors, empty for system-assigned
    # (spec-guaranteed by the identity CEL).
    identity_ids = length(var.spec.identity.identity_ids) > 0 ? var.spec.identity.identity_ids : null
  }

  tags = local.final_tags
}
