# AWS Batch fair-share scheduling policy.
#
# The fair_share_policy block is emitted unconditionally: an empty block is
# valid AWS-side (all defaults), and keeping the block present makes every
# dial an in-place update as it is added or removed -- matching the Pulumi
# module's always-present FairSharePolicy.
resource "aws_batch_scheduling_policy" "this" {
  # The cloud name comes from metadata.name (the catalog naming basis) --
  # set explicitly so both engines create the same policy name.
  name = var.metadata.name

  fair_share_policy {
    compute_reservation = var.spec.compute_reservation
    share_decay_seconds = var.spec.share_decay_seconds

    dynamic "share_distribution" {
      for_each = var.spec.share_distributions
      content {
        share_identifier = share_distribution.value.share_identifier
        # weight_factor 0 means "unset" in the spec (AWS then defaults the
        # share's weight to 1.0) -- never send a literal zero, which is
        # below AWS's 0.0001 minimum.
        weight_factor = share_distribution.value.weight_factor > 0 ? share_distribution.value.weight_factor : null
      }
    }
  }

  tags = local.aws_tags
}
