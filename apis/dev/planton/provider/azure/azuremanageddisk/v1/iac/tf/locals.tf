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
    "resource_kind" = "azure_managed_disk"
    "resource_name" = var.metadata.name
  }

  org_tag = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "organization" = var.metadata.org } : {}

  env_tag = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "environment" = var.metadata.env } : {}

  # Metadata-derived tags first, then the user's spec tags merged over them:
  # user tags deliberately win so an org's governance conventions (cost
  # center, owner) can override the derived values where they collide.
  final_tags = merge(local.base_tags, local.org_tag, local.env_tag, var.spec.tags)

  # Map the spec enums' name strings to ARM values. Only explicit choices
  # are ever sent (null otherwise), so an unspecified spec and Azure's
  # defaults deploy identically on both engines.
  storage_account_type = {
    "STANDARD_LRS"     = "Standard_LRS"
    "STANDARD_SSD_LRS" = "StandardSSD_LRS"
    "STANDARD_SSD_ZRS" = "StandardSSD_ZRS"
    "PREMIUM_LRS"      = "Premium_LRS"
    "PREMIUM_ZRS"      = "Premium_ZRS"
    "PREMIUM_V2_LRS"   = "PremiumV2_LRS"
    "ULTRA_SSD_LRS"    = "UltraSSD_LRS"
  }[var.spec.storage_account_type]

  create_option = {
    "EMPTY"         = "Empty"
    "COPY"          = "Copy"
    "FROM_IMAGE"    = "FromImage"
    "IMPORT"        = "Import"
    "IMPORT_SECURE" = "ImportSecure"
    "RESTORE"       = "Restore"
    "UPLOAD"        = "Upload"
  }[var.spec.create_option]

  os_type = (
    var.spec.os_type == "LINUX" ? "Linux" :
    var.spec.os_type == "WINDOWS" ? "Windows" : null
  )

  security_type = (
    var.spec.security_type == "CONFIDENTIAL_VM_VMGUEST_STATE_ONLY_ENCRYPTED_WITH_PLATFORM_KEY" ? "ConfidentialVM_VMGuestStateOnlyEncryptedWithPlatformKey" :
    var.spec.security_type == "CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_PLATFORM_KEY" ? "ConfidentialVM_DiskEncryptedWithPlatformKey" :
    var.spec.security_type == "CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_CUSTOMER_KEY" ? "ConfidentialVM_DiskEncryptedWithCustomerKey" : null
  )

  network_access_policy = (
    var.spec.network_access_policy == "ALLOW_ALL" ? "AllowAll" :
    var.spec.network_access_policy == "ALLOW_PRIVATE" ? "AllowPrivate" :
    var.spec.network_access_policy == "DENY_ALL" ? "DenyAll" : null
  )
}
