# The Data Lake Gen2 filesystem, addressed by the parent account's ARM
# ID. This is a DATA-PLANE resource: the provider talks to the account's
# dfs endpoint (with the account's shared key by default), so the
# account must be reachable from where the deploy runs -- a data-plane
# firewall that blocks the runner blocks the create even though ARM
# would allow it. Filesystems carry no Azure tags (the properties map is
# the filesystem-level metadata surface); the platform's identity tags
# live on the account.
#
# POSIX access control (owner/group/ace) requires hierarchical namespace
# on the ACCOUNT -- Azure rejects it on flat-namespace accounts at apply
# time, deliberately not mirrored as spec validation because the account
# arrives as a reference.
resource "azurerm_storage_data_lake_gen2_filesystem" "main" {
  name               = var.spec.filesystem_name
  storage_account_id = var.spec.storage_account_id

  # Sub-account key isolation: data that doesn't name its own scope
  # encrypts under this one. Fixed at creation.
  default_encryption_scope = local.default_encryption_scope

  # Root-path ownership. Unset leaves Azure's defaults ($superuser --
  # the shared-key principal); both accept an Entra object ID or the
  # literal $superuser.
  owner = local.owner
  group = local.group

  # The root path's POSIX ACL. Access entries gate the root itself;
  # default entries are the template newly created children inherit --
  # how a zone's permission posture propagates to files landing in it.
  # The scope's unset default is "access", matching the spec enum's
  # unspecified value.
  dynamic "ace" {
    for_each = var.spec.aces
    content {
      scope       = ace.value.scope == null || ace.value.scope == "" ? null : local.ace_scope_map[ace.value.scope]
      type        = local.ace_type_map[ace.value.type]
      id          = ace.value.object_id == null || ace.value.object_id == "" ? null : ace.value.object_id
      permissions = ace.value.permissions
    }
  }

  # Azure requires the VALUES to be base64-encoded; keys stay plain.
  properties = var.spec.properties
}
