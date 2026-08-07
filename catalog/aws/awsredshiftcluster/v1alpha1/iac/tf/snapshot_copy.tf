# Cross-region snapshot copy is a cluster setting keyed by the cluster
# identifier (AWS EnableSnapshotCopy/DisableSnapshotCopy), not a
# resource with its own identity -- folded for the same reason as
# logging. Changing the destination region tears the configuration down
# and re-enables it against the new region (provider replace semantics).
resource "aws_redshift_snapshot_copy" "this" {
  count = var.spec.snapshot_copy != null ? 1 : 0

  cluster_identifier = aws_redshift_cluster.this.cluster_identifier
  destination_region = var.spec.snapshot_copy.destination_region

  # 0 keeps the AWS defaults: 7 days for copied automated snapshots,
  # indefinite (-1) for copied manual snapshots.
  retention_period                 = var.spec.snapshot_copy.retention_period != 0 ? var.spec.snapshot_copy.retention_period : null
  manual_snapshot_retention_period = var.spec.snapshot_copy.manual_snapshot_retention_period != 0 ? var.spec.snapshot_copy.manual_snapshot_retention_period : null

  # Required by AWS when the cluster is KMS-encrypted: the grant lets
  # Redshift encrypt copied snapshots with a key in the destination
  # region.
  snapshot_copy_grant_name = var.spec.snapshot_copy.snapshot_copy_grant_name != "" ? var.spec.snapshot_copy.snapshot_copy_grant_name : null
}
