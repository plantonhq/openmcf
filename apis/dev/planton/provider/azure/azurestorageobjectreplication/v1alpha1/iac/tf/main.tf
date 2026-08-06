# The object replication policy. ONE logical policy, materialized by
# Azure on BOTH accounts under one server-assigned GUID -- the provider
# creates the destination side first (which assigns rule IDs), then
# mirrors it onto the source; destroy removes both sides. The policy
# carries no Azure tags (ARM does not support tags on
# objectReplicationPolicies); the platform's identity tags live on the
# accounts.
#
# Apply-time prerequisites the ACCOUNTS must carry (deliberately not
# mirrored as spec validation -- the accounts arrive as references):
# blob versioning + change feed on the source, blob versioning on the
# destination (the account spec's blob_properties). Replication is
# asynchronous with no default RPO guarantee.
resource "azurerm_storage_object_replication" "main" {
  source_storage_account_id      = var.spec.source_storage_account_id
  destination_storage_account_id = var.spec.destination_storage_account_id

  dynamic "rules" {
    for_each = var.spec.rules
    content {
      source_container_name      = rules.value.source_container_name
      destination_container_name = rules.value.destination_container_name

      # Unset lets the provider default (OnlyNewObjects -- no backfill)
      # apply; Everything backfills the whole container; an RFC 3339
      # instant backfills blobs created after that moment.
      copy_blobs_created_after = rules.value.copy_blobs_created_after == null || rules.value.copy_blobs_created_after == "" ? null : rules.value.copy_blobs_created_after

      # The spec names this prefix_match after ARM's own INCLUDE
      # semantics (prefixMatch); the provider attribute's
      # "filter_out" name is historical and means the same
      # include-only-these-prefixes behavior.
      filter_out_blobs_with_prefix = rules.value.prefix_match
    }
  }
}
