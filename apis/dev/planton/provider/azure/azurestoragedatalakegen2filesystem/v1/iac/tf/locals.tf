locals {
  # The spec's ace scope/type enums arrive as the FULL proto value
  # names; the maps carry the complete vocabularies, translated to the
  # data plane's lowercase wire values.
  ace_scope_map = {
    "ACCESS"  = "access"
    "DEFAULT" = "default"
  }

  ace_type_map = {
    "USER"  = "user"
    "GROUP" = "group"
    "MASK"  = "mask"
    "OTHER" = "other"
  }

  # Optional strings are sent only when non-empty so unset stays unset
  # on the wire (owner/group are Computed -- always sending an empty
  # string would fight Azure's server-side defaults).
  default_encryption_scope = (
    var.spec.default_encryption_scope == null || var.spec.default_encryption_scope == "" ? null :
    var.spec.default_encryption_scope
  )

  owner = var.spec.owner == null || var.spec.owner == "" ? null : var.spec.owner
  group = var.spec.group == null || var.spec.group == "" ? null : var.spec.group

  # The account name, parsed from the resolved account ARM ID -- used
  # for the storage_account_name output and to construct the
  # filesystem's ARM container-proxy ID. The named-group regex fails the
  # plan loudly if the ID is not a storage-account ARM ID.
  storage_account_name = regex("/storageAccounts/(?P<name>[^/]+)$", var.spec.storage_account_id)["name"]

  # ADLS filesystems surface in ARM as blob containers, so the
  # management/RBAC identity of the filesystem is the container-proxy
  # ID under the account -- constructed (identically on both engines)
  # rather than read back, because the provider's own resource id is a
  # data-plane dfs URL that nothing management-grain can consume.
  filesystem_arm_id = "${var.spec.storage_account_id}/blobServices/default/containers/${var.spec.filesystem_name}"
}
