# AWS Client VPN: the managed OpenVPN front door for remote access into AWS
# networks. The endpoint carries everything decided at create time
# (authentication, client CIDR, transport, IP address types, transit-gateway
# attachment) plus the in-place dials (sessions, banner, connect handler,
# logging, DNS). Its three endpoint-scoped satellites -- target network
# associations, authorization rules, and routes -- are folded here because
# none has identity outside its endpoint; each is still its own provider
# resource so membership edits apply in place.

resource "aws_ec2_client_vpn_endpoint" "this" {
  description            = var.spec.description != "" ? var.spec.description : null
  server_certificate_arn = var.spec.server_certificate_arn

  # Client addressing: required except for pure-IPv6 tunnel traffic (CEL
  # enforces the coupling; AWS assigns addressing there).
  client_cidr_block = var.spec.client_cidr_block != "" ? var.spec.client_cidr_block : null

  # Split tunnel is sent explicitly: AWS defaults to full tunnel (false),
  # and both engines must agree on the sent value rather than each relying
  # on provider defaults.
  split_tunnel = var.spec.split_tunnel

  transport_protocol = var.spec.transport_protocol != "" ? var.spec.transport_protocol : null
  vpn_port           = var.spec.vpn_port

  endpoint_ip_address_type = var.spec.endpoint_ip_address_type != "" ? var.spec.endpoint_ip_address_type : null
  traffic_ip_address_type  = var.spec.traffic_ip_address_type != "" ? var.spec.traffic_ip_address_type : null

  # One block per authentication option (max 2; a client passes on ANY).
  # The per-type identity source is CEL-guaranteed. The client CA chain is
  # deliberately its own ref instead of silently reusing the server
  # certificate: the two play different roles even when a self-signed setup
  # points both at the same imported certificate.
  dynamic "authentication_options" {
    for_each = var.spec.authentication_options
    content {
      type                           = authentication_options.value.type
      root_certificate_chain_arn     = authentication_options.value.root_certificate_chain_arn != "" ? authentication_options.value.root_certificate_chain_arn : null
      active_directory_id            = authentication_options.value.active_directory_id != "" ? authentication_options.value.active_directory_id : null
      saml_provider_arn              = authentication_options.value.saml_provider_arn != "" ? authentication_options.value.saml_provider_arn : null
      self_service_saml_provider_arn = authentication_options.value.self_service_saml_provider_arn != "" ? authentication_options.value.self_service_saml_provider_arn : null
    }
  }

  # Connection logging: presence of the spec block is the switch. The
  # provider requires this block either way, so absence maps to
  # enabled=false -- there is no separate boolean for a manifest to
  # contradict.
  connection_log_options {
    enabled               = var.spec.connection_log != null
    cloudwatch_log_group  = var.spec.connection_log != null ? var.spec.connection_log.cloudwatch_log_group : null
    cloudwatch_log_stream = var.spec.connection_log != null && try(var.spec.connection_log.cloudwatch_log_stream, "") != "" ? var.spec.connection_log.cloudwatch_log_stream : null
  }

  # VPC attachment surface (CEL guarantees it never coexists with the
  # transit-gateway arm). When omitted, AWS infers the VPC from the first
  # associated subnet and applies that VPC's default security group.
  # Unresolved (empty-string) references are filtered out before the length
  # check — mirroring the Pulumi module — so a list of one unresolved ref
  # omits the argument instead of sending [""] to the provider.
  vpc_id             = var.spec.vpc_id != "" ? var.spec.vpc_id : null
  security_group_ids = length([for sg in var.spec.security_group_ids : sg if sg != ""]) > 0 ? [for sg in var.spec.security_group_ids : sg if sg != ""] : null

  # Transit-gateway attachment surface (ForceNew): clients reach every
  # network the gateway routes to, without per-subnet associations.
  dynamic "transit_gateway_configuration" {
    for_each = var.spec.transit_gateway_configuration != null ? [var.spec.transit_gateway_configuration] : []
    content {
      transit_gateway_id    = transit_gateway_configuration.value.transit_gateway_id
      availability_zones    = length(transit_gateway_configuration.value.availability_zones) > 0 ? transit_gateway_configuration.value.availability_zones : null
      availability_zone_ids = length(transit_gateway_configuration.value.availability_zone_ids) > 0 ? transit_gateway_configuration.value.availability_zone_ids : null
    }
  }

  # Sessions and client experience (all update in place).
  session_timeout_hours         = var.spec.session_timeout_hours
  disconnect_on_session_timeout = var.spec.disconnect_on_session_timeout
  self_service_portal           = var.spec.self_service_portal_enabled ? "enabled" : "disabled"

  # Posture hook: presence of the spec block is the switch. The block is
  # ALWAYS materialized with an explicit enabled value — the provider
  # declares it Optional+Computed and diff-suppresses a missing block, so
  # omitting it after the hook was enabled would leave the Lambda gate
  # active forever. Sending enabled=false (payload deliberately dropped,
  # matching the provider's own expander) is what makes removal a real
  # disable.
  client_connect_options {
    enabled             = var.spec.client_connect_options != null
    lambda_function_arn = var.spec.client_connect_options != null ? var.spec.client_connect_options.lambda_function_arn : null
  }

  # Login banner: presence of the spec block is the switch — always
  # materialized for the same removal-must-disable reason as above.
  client_login_banner_options {
    enabled     = var.spec.client_login_banner != null
    banner_text = var.spec.client_login_banner != null ? var.spec.client_login_banner.banner_text : null
  }

  # Route enforcement: a single hardening dial, always sent explicitly in
  # BOTH states — the provider's Optional+Computed block absorbs the
  # AWS-side value when omitted, so flipping the spec back to false would
  # otherwise be a silent no-op.
  client_route_enforcement_options {
    enforced = var.spec.client_route_enforcement_enabled
  }

  dns_servers = length(var.spec.dns_servers) > 0 ? var.spec.dns_servers : null

  tags = local.aws_tags
}

# Target network associations: one per subnet, keyed by subnet ID so adding
# an AZ never disturbs the others. Each attach/detach takes AWS several
# minutes -- these are the slow part of any Client VPN deploy.
resource "aws_ec2_client_vpn_network_association" "this" {
  for_each = toset(var.spec.subnet_ids)

  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.this.id
  subnet_id              = each.value
}

# Authorization rules: the endpoint's network ACL. Nothing is reachable
# until a rule authorizes it -- even with an association in place.
resource "aws_ec2_client_vpn_authorization_rule" "this" {
  for_each = local.authorization_rules

  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.this.id
  target_network_cidr    = each.value.target_network_cidr

  # Exactly one grantee arm (CEL-enforced): a specific IdP group or every
  # authenticated client.
  access_group_id      = each.value.access_group_id != "" ? each.value.access_group_id : null
  authorize_all_groups = each.value.access_group_id == "" ? true : null

  description = each.value.description != "" ? each.value.description : null
}

# Additional routes beyond the auto-created per-VPC ones. Each route depends
# on its target subnet's association: AWS rejects a route whose subnet is
# still associating, so the edge makes ordering explicit instead of leaning
# on provider-side retries.
resource "aws_ec2_client_vpn_route" "this" {
  for_each = local.routes

  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.this.id
  destination_cidr_block = each.value.destination_cidr_block
  target_vpc_subnet_id   = each.value.target_subnet_id
  description            = each.value.description != "" ? each.value.description : null

  depends_on = [aws_ec2_client_vpn_network_association.this]
}
