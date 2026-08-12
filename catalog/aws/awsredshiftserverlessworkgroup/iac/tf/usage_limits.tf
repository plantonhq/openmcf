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

  resource_arn = aws_redshiftserverless_workgroup.this.arn

  usage_type = each.value.usage_type
  amount     = each.value.amount

  # Empty keeps the AWS defaults (monthly / log); the provider defaults
  # match, so sending only set values stays faithful.
  period        = each.value.period != "" ? each.value.period : null
  breach_action = each.value.breach_action != "" ? each.value.breach_action : null
}
