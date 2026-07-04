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
    "resource_kind" = "azure_key_vault"
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

  # Map the spec enum's name string to ARM's SKU value. The unspecified
  # spec applies the STANDARD baseline on both engines (azurerm requires
  # an explicit sku_name, so the default is materialized here rather than
  # sent as null).
  sku_name = var.spec.sku == "PREMIUM" ? "premium" : "standard"

  # Map the network default-action enum's name string to ARM's value.
  # default_action is required whenever the network_acls block is present
  # (spec validation enforces it), so there is no null fallback to invent.
  network_acls_default_action = (
    var.spec.network_acls == null ? null :
    var.spec.network_acls.default_action == "ALLOW" ? "Allow" : "Deny"
  )

  # Map the bypass enum's name string to ARM's value. azurerm requires a
  # bypass value whenever network_acls is set, so an unspecified spec
  # materializes Azure's own default (AzureServices) -- identical on both
  # engines.
  network_acls_bypass = (
    var.spec.network_acls == null ? null :
    var.spec.network_acls.bypass == "NONE" ? "None" : "AzureServices"
  )

  # The spec's permission enums arrive as FULL proto value names (the
  # tfvars wire format never strips prefixes); each map below carries the
  # complete verbatim vocabulary for its enum, mapped to ARM's data-plane
  # permission strings. A missing entry would silently drop a grant, so
  # the maps are exhaustive by construction.
  key_permission_map = {
    "KEY_GET"                 = "Get"
    "KEY_LIST"                = "List"
    "KEY_UPDATE"              = "Update"
    "KEY_CREATE"              = "Create"
    "KEY_IMPORT"              = "Import"
    "KEY_DELETE"              = "Delete"
    "KEY_RECOVER"             = "Recover"
    "KEY_BACKUP"              = "Backup"
    "KEY_RESTORE"             = "Restore"
    "KEY_DECRYPT"             = "Decrypt"
    "KEY_ENCRYPT"             = "Encrypt"
    "KEY_UNWRAP_KEY"          = "UnwrapKey"
    "KEY_WRAP_KEY"            = "WrapKey"
    "KEY_VERIFY"              = "Verify"
    "KEY_SIGN"                = "Sign"
    "KEY_PURGE"               = "Purge"
    "KEY_RELEASE"             = "Release"
    "KEY_ROTATE"              = "Rotate"
    "KEY_GET_ROTATION_POLICY" = "GetRotationPolicy"
    "KEY_SET_ROTATION_POLICY" = "SetRotationPolicy"
  }

  secret_permission_map = {
    "SECRET_GET"     = "Get"
    "SECRET_LIST"    = "List"
    "SECRET_SET"     = "Set"
    "SECRET_DELETE"  = "Delete"
    "SECRET_RECOVER" = "Recover"
    "SECRET_BACKUP"  = "Backup"
    "SECRET_RESTORE" = "Restore"
    "SECRET_PURGE"   = "Purge"
  }

  certificate_permission_map = {
    "CERTIFICATE_GET"             = "Get"
    "CERTIFICATE_LIST"            = "List"
    "CERTIFICATE_UPDATE"          = "Update"
    "CERTIFICATE_CREATE"          = "Create"
    "CERTIFICATE_IMPORT"          = "Import"
    "CERTIFICATE_DELETE"          = "Delete"
    "CERTIFICATE_RECOVER"         = "Recover"
    "CERTIFICATE_BACKUP"          = "Backup"
    "CERTIFICATE_RESTORE"         = "Restore"
    "CERTIFICATE_MANAGE_CONTACTS" = "ManageContacts"
    "CERTIFICATE_MANAGE_ISSUERS"  = "ManageIssuers"
    "CERTIFICATE_GET_ISSUERS"     = "GetIssuers"
    "CERTIFICATE_LIST_ISSUERS"    = "ListIssuers"
    "CERTIFICATE_SET_ISSUERS"     = "SetIssuers"
    "CERTIFICATE_DELETE_ISSUERS"  = "DeleteIssuers"
    "CERTIFICATE_PURGE"           = "Purge"
  }

  storage_permission_map = {
    "STORAGE_GET"            = "Get"
    "STORAGE_LIST"           = "List"
    "STORAGE_SET"            = "Set"
    "STORAGE_UPDATE"         = "Update"
    "STORAGE_DELETE"         = "Delete"
    "STORAGE_RECOVER"        = "Recover"
    "STORAGE_BACKUP"         = "Backup"
    "STORAGE_RESTORE"        = "Restore"
    "STORAGE_PURGE"          = "Purge"
    "STORAGE_REGENERATE_KEY" = "RegenerateKey"
    "STORAGE_GET_SAS"        = "GetSAS"
    "STORAGE_LIST_SAS"       = "ListSAS"
    "STORAGE_SET_SAS"        = "SetSAS"
    "STORAGE_DELETE_SAS"     = "DeleteSAS"
  }
}
