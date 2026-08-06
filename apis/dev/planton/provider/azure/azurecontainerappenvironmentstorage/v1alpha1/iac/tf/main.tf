# The storage registration makes an Azure Files share mountable by apps
# and jobs in the environment (their volumes reference it by name).
# Everything except the SMB access_key is ForceNew -- key rotation is the
# one in-place update; every other change recreates the registration and
# briefly breaks volume mounts that reference it.
resource "azurerm_container_app_environment_storage" "main" {
  name                         = var.spec.storage_name
  container_app_environment_id = var.spec.container_app_environment_id
  share_name                   = var.spec.share_name
  access_mode                  = local.access_mode_map[var.spec.access_mode]

  # Exactly one protocol (spec-enforced): SMB addresses the share by
  # account name + access key; NFS addresses the account's file endpoint
  # and requires a VNet-injected environment.
  account_name   = var.spec.account_name
  access_key     = var.spec.access_key
  nfs_server_url = var.spec.nfs_server_url
}
