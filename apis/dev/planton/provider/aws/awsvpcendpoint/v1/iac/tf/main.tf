# The VPC endpoint composes onto its neighbors instead of embedding
# them: the VPC attaches by reference, gateway endpoints reference route
# tables (a subnet-owned table or the VPC's main/default table), and
# interface endpoints reference AwsSubnet and AwsSecurityGroup nodes.
# This module never modifies a resource it merely references -- in
# particular it never edits the referenced route tables' own routes; AWS
# manages the endpoint's prefix-list route as part of the endpoint
# itself.
#
# Create-only in AWS: the endpoint type, the service target
# (service_name / resource_configuration_arn / service_network_arn), and
# service_region. Attachments (route tables, subnets, security groups),
# policy, DNS options, and IP address type update in place.
resource "aws_vpc_endpoint" "this" {
  vpc_id = var.spec.vpc_id

  # Empty keeps AWS's own default (Gateway) -- passing nothing and
  # passing "Gateway" are equivalent, so only a non-empty spec value is
  # forwarded, keeping the diff surface minimal.
  vpc_endpoint_type = var.spec.endpoint_type != "" ? var.spec.endpoint_type : null

  # Exactly one of the three service targets is set (CEL-enforced);
  # forward whichever one carries the target.
  service_name               = var.spec.service_name != "" ? var.spec.service_name : null
  resource_configuration_arn = var.spec.resource_configuration_arn != "" ? var.spec.resource_configuration_arn : null
  service_network_arn        = var.spec.service_network_arn != "" ? var.spec.service_network_arn : null

  # Gateway endpoints attach through route tables; ENI-based types
  # attach through subnets (one ENI per subnet) and, for Interface,
  # security groups. The spec's CEL gating guarantees only the arms
  # matching the endpoint type are populated. Empty lists become null so
  # AWS's own defaults apply (e.g. the VPC default security group on an
  # interface endpoint with no groups given).
  route_table_ids    = length(var.spec.route_table_ids) > 0 ? var.spec.route_table_ids : null
  subnet_ids         = length(var.spec.subnet_ids) > 0 ? var.spec.subnet_ids : null
  security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null

  # private_dns_enabled is only expressible on Interface endpoints
  # (CEL), where it updates in place. Forwarded only when true so a
  # gateway endpoint's create call carries no DNS argument at all --
  # matching the Pulumi module, which prunes the false boolean.
  private_dns_enabled = var.spec.private_dns_enabled ? true : null

  dynamic "dns_options" {
    for_each = var.spec.dns_options != null ? [var.spec.dns_options] : []
    content {
      dns_record_ip_type = dns_options.value.dns_record_ip_type != "" ? dns_options.value.dns_record_ip_type : null
      # The S3 dual-stack pattern: in-VPC traffic rides the free gateway
      # endpoint while on-premises resolver traffic reaches the
      # interface endpoint. Forwarded only when true -- AWS rejects the
      # flag on services without a gateway+interface pair.
      private_dns_only_for_inbound_resolver_endpoint = dns_options.value.private_dns_only_for_inbound_resolver_endpoint ? true : null
      # Preference + domain list are the Lattice types' private-DNS
      # controls; the spec couples them (domains exactly when the
      # preference says "specified domains").
      private_dns_preference        = dns_options.value.private_dns_preference != "" ? dns_options.value.private_dns_preference : null
      private_dns_specified_domains = length(dns_options.value.private_dns_specified_domains) > 0 ? dns_options.value.private_dns_specified_domains : null
    }
  }

  ip_address_type = var.spec.ip_address_type != "" ? var.spec.ip_address_type : null

  # Empty policy means AWS's full-access default; forwarding nothing
  # (instead of an empty string) lets AWS attach its default document.
  policy = var.spec.policy != "" ? var.spec.policy : null

  # Pin specific ENI addresses per subnet -- for appliances or firewall
  # rules that need stable endpoint IPs. Each listed subnet must also
  # appear in subnet_ids.
  dynamic "subnet_configuration" {
    for_each = var.spec.subnet_configurations
    content {
      subnet_id = subnet_configuration.value.subnet_id
      ipv4      = subnet_configuration.value.ipv4 != "" ? subnet_configuration.value.ipv4 : null
      ipv6      = subnet_configuration.value.ipv6 != "" ? subnet_configuration.value.ipv6 : null
    }
  }

  # Cross-region interface endpoints: reach a service in another region
  # without cross-region networking of your own. Create-only in AWS.
  service_region = var.spec.service_region != "" ? var.spec.service_region : null

  # Accept the connection automatically when the PrivateLink service
  # requires acceptance and lives in the same account; cross-account
  # services must accept on their side regardless.
  auto_accept = var.spec.auto_accept ? true : null

  tags = local.aws_tags
}
