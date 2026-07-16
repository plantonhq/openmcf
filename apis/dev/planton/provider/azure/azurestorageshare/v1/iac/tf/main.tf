# The Azure Files share, addressed by the parent account's ARM ID (the
# control-plane path -- the account-name form is the provider's legacy
# data-plane path, removed in azurerm v5). Shares carry no Azure tags:
# ARM does not support tags on fileServices/shares, so the platform's
# identity tags live on the account.
resource "azurerm_storage_share" "main" {
  name               = var.spec.share_name
  storage_account_id = var.spec.storage_account_id

  # The provisioned quota in GB -- what SMB clients see as the drive
  # size and Azure enforces on writes. Grows in place; shrinking below
  # used capacity fails. Premium FileStorage bills this whether used or
  # not.
  quota = var.spec.quota_gb

  # SMB (the default) or NFS -- NFS requires a premium FileStorage
  # account and is reachable only over private network paths. Fixed at
  # creation.
  enabled_protocol = local.enabled_protocol

  # Sent only when the spec chooses a tier, so Azure's per-account-kind
  # default (TransactionOptimized on standard, Premium on FileStorage)
  # applies when unset. Premium is required -- and the only legal tier
  # -- on FileStorage accounts.
  access_tier = local.access_tier

  # Stored access policies (signed identifiers): revoking or shortening
  # a policy immediately revokes every SAS token anchored to it. At most
  # five per share (Azure's limit, enforced in the spec).
  dynamic "acl" {
    for_each = var.spec.acls
    content {
      id = acl.value.id

      dynamic "access_policy" {
        for_each = acl.value.access_policies
        content {
          permissions = access_policy.value.permissions
          start       = access_policy.value.start
          expiry      = access_policy.value.expiry
        }
      }
    }
  }

  metadata = var.spec.metadata
}
