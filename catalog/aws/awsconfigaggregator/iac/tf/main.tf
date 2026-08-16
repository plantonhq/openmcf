# AWS Config aggregation, both sides:
#   - the AGGREGATOR (spec.aggregation) - deployed in the collecting
#     account; references no Config recorder (aggregation works in an
#     account with zero recorders);
#   - the reciprocal GRANTS (spec.authorizations) - deployed in each
#     source account, naming the aggregator account+region allowed to
#     collect from it.
#
# Lifecycle facts the renders below depend on:
#   - the provider replaces the aggregator only when a source block
#     APPEARS on an existing aggregator (absent -> present); content
#     changes and block removal update in place;
#   - the spec CEL guarantees exactly one source shape arrives here;
#   - grants are keyed "{account_id}:{authorized_aws_region}" (the
#     provider's own import ID), so reordering the spec list never
#     churns them.

resource "aws_config_configuration_aggregator" "this" {
  count = var.spec.aggregation != null ? 1 : 0

  # metadata.name is the aggregator name on both engines.
  name = var.metadata.name

  dynamic "account_aggregation_source" {
    for_each = var.spec.aggregation.account_source != null ? [var.spec.aggregation.account_source] : []
    content {
      account_ids = account_aggregation_source.value.account_ids
      all_regions = account_aggregation_source.value.all_regions
      regions     = length(account_aggregation_source.value.regions) > 0 ? account_aggregation_source.value.regions : null
    }
  }

  dynamic "organization_aggregation_source" {
    for_each = var.spec.aggregation.organization_source != null ? [var.spec.aggregation.organization_source] : []
    content {
      role_arn    = organization_aggregation_source.value.role_arn
      all_regions = organization_aggregation_source.value.all_regions
      regions     = length(organization_aggregation_source.value.regions) > 0 ? organization_aggregation_source.value.regions : null
    }
  }

  tags = local.aws_tags
}

# The source-account side: each grant authorizes ONE aggregator
# (account+region) to collect this account's Config data. The
# provider's deprecated "region" alias is deliberately not rendered -
# authorized_aws_region is the surviving argument.
resource "aws_config_aggregate_authorization" "grants" {
  for_each = {
    for grant in var.spec.authorizations :
    "${grant.account_id}:${grant.authorized_aws_region}" => grant
  }

  account_id            = each.value.account_id
  authorized_aws_region = each.value.authorized_aws_region

  tags = local.aws_tags
}
