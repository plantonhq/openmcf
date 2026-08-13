# Usage limits cap the workgroup's consumption -- RPU-hours of
# serverless compute or terabytes of cross-region datasharing transfer.
# AWS generates the limit IDs at creation; the outputs export them keyed
# by the same pair used here so imports and out-of-band CLI operations
# can address each limit. The resource is untagged in AWS (unlike the
# provisioned cluster's usage limits).
resource "aws_redshiftserverless_usage_limit" "this" {
  # Keyed by usage_type/period with an unset period rendered as monthly
  # -- the same normalization the spec's uniqueness CEL applies, so
  # validate-time uniqueness IS plan-time key uniqueness.
  for_each = {
    for l in var.spec.usage_limits :
    "${l.usage_type}/${l.period != "" ? l.period : "monthly"}" => l
  }

  # Usage-limit calls are conflict-IMMUNE (live-probed: create and
  # delete both succeed against a busy workgroup) but each call flips
  # the workgroup to MODIFYING for ~15-30s after returning -- while the
  # endpoint-access create AND delete are conflict-SENSITIVE with no
  # provider retry. Limits therefore apply LAST (endpoint accesses
  # first, on the idle workgroup), and the destroy path crosses back
  # over the settle sleep so the endpoint delete never lands in a
  # limit-delete's busy window. Full contract in satellite_settle.tf.
  depends_on = [
    time_sleep.endpoint_access_settle,
    aws_redshiftserverless_endpoint_access.this,
    aws_redshiftserverless_custom_domain_association.this,
  ]

  resource_arn = aws_redshiftserverless_workgroup.this.arn

  usage_type = each.value.usage_type
  amount     = each.value.amount

  # Empty keeps the AWS defaults (monthly / log); the provider defaults
  # match, so sending only set values stays faithful.
  period        = each.value.period != "" ? each.value.period : null
  breach_action = each.value.breach_action != "" ? each.value.breach_action : null
}
