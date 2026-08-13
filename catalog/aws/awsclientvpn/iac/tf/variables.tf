variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name        = string
    id          = optional(string, "")
    org         = optional(string, "")
    env         = optional(string, "")
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})
    tags        = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsClientVpn specification"
  type = object({
    region      = string
    description = optional(string, "")
    authentication_options = list(object({
      type                           = string
      root_certificate_chain_arn     = optional(string, "")
      active_directory_id            = optional(string, "")
      saml_provider_arn              = optional(string, "")
      self_service_saml_provider_arn = optional(string, "")
    }))
    server_certificate_arn   = string
    client_cidr_block        = optional(string, "")
    split_tunnel             = optional(bool, false)
    transport_protocol       = optional(string, "")
    vpn_port                 = optional(number)
    endpoint_ip_address_type = optional(string, "")
    traffic_ip_address_type  = optional(string, "")
    vpc_id                   = optional(string, "")
    security_group_ids       = optional(list(string), [])
    subnet_ids               = optional(list(string), [])
    transit_gateway_configuration = optional(object({
      transit_gateway_id    = string
      availability_zones    = optional(list(string), [])
      availability_zone_ids = optional(list(string), [])
    }))
    authorization_rules = optional(list(object({
      target_network_cidr  = string
      access_group_id      = optional(string, "")
      authorize_all_groups = optional(bool, false)
      description          = optional(string, "")
    })), [])
    routes = optional(list(object({
      destination_cidr_block = string
      target_subnet_id       = string
      description            = optional(string, "")
    })), [])
    session_timeout_hours         = optional(number)
    disconnect_on_session_timeout = optional(bool, false)
    self_service_portal_enabled   = optional(bool, false)
    client_connect_options = optional(object({
      lambda_function_arn = string
    }))
    client_login_banner = optional(object({
      banner_text = string
    }))
    client_route_enforcement_enabled = optional(bool, false)
    dns_servers                      = optional(list(string), [])
    connection_log = optional(object({
      cloudwatch_log_group  = string
      cloudwatch_log_stream = optional(string, "")
    }))
  })
}
