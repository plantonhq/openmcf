locals {
  # Client VPN endpoints have no AWS name argument (identity is the
  # generated cvpn-endpoint-... id), so the name lives on the Name tag,
  # identically on both engines.
  endpoint_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsClientVpn"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Authorization rules are keyed by (destination CIDR, grantee) -- the same
  # identity AWS uses -- so one CIDR can be granted to several IdP groups
  # and editing one grant never disturbs the others.
  authorization_rules = {
    for rule in var.spec.authorization_rules :
    "${rule.target_network_cidr}|${rule.access_group_id != "" ? rule.access_group_id : "all-groups"}" => rule
  }

  # Routes are keyed by (destination, subnet) -- the same pair AWS uses --
  # so a destination can be reached through several subnets.
  routes = {
    for route in var.spec.routes :
    "${route.destination_cidr_block}|${route.target_subnet_id}" => route
  }
}
