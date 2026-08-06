# App Runner VPC connector.
#
# A VPC connector is a set of AWS-managed ENIs provisioned into the
# referenced subnets; App Runner services route their OUTBOUND traffic
# through it to reach private VPC resources (databases, caches, internal
# APIs). AWS has no update API for connectors -- changing subnets or
# security groups replaces the connector (registered as a new revision
# under the same name), which is why both lists are create-time inputs
# with no in-place drift story.
resource "aws_apprunner_vpc_connector" "this" {
  vpc_connector_name = local.resource_name

  # Referenced subnet/SG IDs arrive pre-resolved as plain strings. AWS
  # requires at least one security group on a connector -- the groups govern
  # what the connected services can reach, so the targets' groups must also
  # admit ingress from these.
  subnets         = var.spec.subnet_ids
  security_groups = var.spec.security_group_ids

  tags = local.aws_tags
}
