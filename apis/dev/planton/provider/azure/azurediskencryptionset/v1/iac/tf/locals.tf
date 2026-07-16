locals {
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # PARITY-EXCEPTION: resource_kind here is the family-wide snake-case
  # literal and resource_id falls back to metadata.name, while the Pulumi
  # module emits the lowered CloudResourceKind enum string and omits
  # resource_id when metadata.id is empty. Output-neutral (tags never feed
  # stack outputs); aligning the two shapes is a family-wide convention
  # change, not a per-kind fix.
  base_tags = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_disk_encryption_set"
    "resource_name" = var.metadata.name
  }

  org_tag = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "organization" = var.metadata.org } : {}

  env_tag = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "environment" = var.metadata.env } : {}

  final_tags = merge(local.base_tags, local.org_tag, local.env_tag, var.spec.tags)

  # Map the spec enums' name strings to ARM's values. null lets Azure apply
  # its default (EncryptionAtRestWithCustomerKey); only an explicit choice is
  # sent so an unspecified spec and Azure's default deploy identically on
  # both engines.
  encryption_type = (
    var.spec.encryption_type == "ENCRYPTION_AT_REST_WITH_CUSTOMER_KEY" ? "EncryptionAtRestWithCustomerKey" :
    var.spec.encryption_type == "ENCRYPTION_AT_REST_WITH_PLATFORM_AND_CUSTOMER_KEYS" ? "EncryptionAtRestWithPlatformAndCustomerKeys" :
    var.spec.encryption_type == "CONFIDENTIAL_VM_ENCRYPTED_WITH_CUSTOMER_KEY" ? "ConfidentialVmEncryptedWithCustomerKey" : null
  )

  identity_type = (
    var.spec.identity.type == "SYSTEM_ASSIGNED" ? "SystemAssigned" :
    var.spec.identity.type == "USER_ASSIGNED" ? "UserAssigned" :
    var.spec.identity.type == "SYSTEM_AND_USER_ASSIGNED" ? "SystemAssigned, UserAssigned" : null
  )
}
