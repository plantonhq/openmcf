# API Gateway v2 VPC link.
#
# A VPC link is a set of AWS-managed ENIs provisioned into the referenced
# subnets; HTTP API private integrations route through it to reach internal
# ALBs, NLBs, or Cloud Map services. AWS has no update API for the network
# attachment -- changing subnets or security groups replaces the link (only
# the name mutates in place), which is why both lists are create-time inputs
# with no in-place drift story.
resource "aws_apigatewayv2_vpc_link" "this" {
  name = local.resource_name

  # Referenced subnet/SG IDs arrive pre-resolved as plain strings. Security
  # groups may be empty: AWS then applies no filtering on the link side and
  # reachability is governed solely by the target's security groups.
  subnet_ids         = var.spec.subnet_ids
  security_group_ids = var.spec.security_group_ids

  tags = local.aws_tags
}
