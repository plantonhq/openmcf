# One Resolver query logging configuration with its VPC associations
# managed in-line.
#
# Lifecycle facts the render below depends on:
#   - name and destination_arn are both ForceNew - the configuration
#     is immutable except tags (existing log data survives replacement
#     in the destination);
#   - associations are pure joins (config, vpc) - every argument
#     ForceNew, no update path;
#   - an association can FAIL asynchronously after a clean apply when
#     the resolver cannot write to the destination (permissions/policy)
#     - the provider's waiter surfaces the association's error code.

resource "aws_route53_resolver_query_log_config" "this" {
  name            = var.metadata.name
  destination_arn = var.spec.destination_arn

  tags = local.aws_tags
}

resource "aws_route53_resolver_query_log_config_association" "this" {
  for_each = toset(var.spec.vpc_ids)

  resolver_query_log_config_id = aws_route53_resolver_query_log_config.this.id
  resource_id                  = each.value
}
