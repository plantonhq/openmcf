# Usage limits cap what individual Redshift features may consume on
# this cluster (Spectrum scans, concurrency-scaling time, cross-region
# datasharing transfer). AWS generates the limit IDs at creation; the
# outputs export them keyed by the same triple used here so imports and
# out-of-band CLI operations can address each limit.
resource "aws_redshift_usage_limit" "this" {
  for_each = local.usage_limits_by_key

  cluster_identifier = aws_redshift_cluster.this.cluster_identifier

  feature_type = each.value.feature_type
  limit_type   = each.value.limit_type
  amount       = each.value.amount

  # Empty keeps the AWS defaults (monthly / log); the provider defaults
  # match, so sending only set values stays faithful.
  period        = each.value.period != "" ? each.value.period : null
  breach_action = each.value.breach_action != "" ? each.value.breach_action : null

  tags = local.aws_tags
}
