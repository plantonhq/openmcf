locals {
  # The domain name is metadata.name -- create-only in AWS, constrained to
  # ^[a-z][0-9a-z\-]{2,27}$ (3-28 chars), and the basis both engines share so a
  # manifest deploys identically on either.
  domain_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsOpenSearchDomain"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # A VPC domain is selected by the presence of subnets; adding or removing
  # vpc_options after creation replaces the domain (provider ForceNew).
  has_vpc_options = var.spec.vpc_options != null ? length(var.spec.vpc_options.subnet_ids) > 0 : false

  # access_policies is a free-form JSON document (google.protobuf.Struct in the
  # spec); the provider wants it as a JSON string.
  access_policies_json = var.spec.access_policies != null ? jsonencode(var.spec.access_policies) : null

  # FGAC is emitted only when enabled: AWS treats advanced_security_options as
  # one-way (it cannot be disabled once on), so an explicit enabled=false block
  # adds nothing a missing block does not.
  has_advanced_security = var.spec.advanced_security_options != null ? var.spec.advanced_security_options.enabled : false

  has_cognito = var.spec.cognito_options != null ? var.spec.cognito_options.enabled : false
}
