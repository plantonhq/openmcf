# A REST API VPC link: the network attachment REST API integrations
# route through to reach an internal Network Load Balancer.
#
# Lifecycle facts the render below depends on:
#   - the target NLB is immutable in AWS (the provider replaces the
#     link when it changes) -- expected: a link IS its network
#     attachment;
#   - AWS takes exactly one balancer per link (target_arns caps at one
#     item at the pin) -- the spec is singular by design;
#   - creation waits for the attachment to reach AVAILABLE (up to ~20
#     minutes upstream) before integrations can reference the link.

resource "aws_api_gateway_vpc_link" "this" {
  # metadata.name is the naming basis on both engines.
  name        = var.metadata.name
  description = var.spec.description != "" ? var.spec.description : null

  target_arns = [var.spec.target_arn]

  tags = local.aws_tags
}
