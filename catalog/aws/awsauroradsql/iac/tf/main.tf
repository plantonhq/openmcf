# One Aurora DSQL cluster with its multi-region pairing managed
# in-line.
#
# Lifecycle facts the render below depends on:
#   - only multi_region_properties.witness_region replaces the
#     cluster; deletion protection and the KMS key update in place (a
#     key change re-encrypts, no replacement);
#   - the PEERING resource is a disguised UpdateCluster call AWS
#     accepts only while the cluster is in PENDING_SETUP - creating
#     it right after the cluster (as here) is the one valid order;
#   - the peering has NO update path at the provider (changes error
#     at apply, they do not replace) and a no-op delete - changing
#     peers means recreating the CLUSTER;
#   - force_destroy makes the provider disable deletion protection
#     before deleting.

resource "aws_dsql_cluster" "this" {
  deletion_protection_enabled = var.spec.deletion_protection_enabled ? true : null
  force_destroy               = var.spec.force_destroy ? true : null
  kms_encryption_key          = var.spec.kms_encryption_key != "" ? var.spec.kms_encryption_key : null

  dynamic "multi_region_properties" {
    for_each = var.spec.multi_region != null ? [var.spec.multi_region] : []
    content {
      # The witness makes the CLUSTER multi-region at create; the peer
      # list lands via the peering resource below.
      witness_region = multi_region_properties.value.witness_region
    }
  }

  tags = local.aws_tags
}

resource "aws_dsql_cluster_peering" "this" {
  count = var.spec.multi_region != null ? 1 : 0

  identifier     = aws_dsql_cluster.this.identifier
  clusters       = var.spec.multi_region.peer_cluster_arns
  witness_region = var.spec.multi_region.witness_region
}
