variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Azure Subnet specification"
  type = object({
    # The parent virtual network's ARM ID. The subnet's resource group and
    # network name are derived from it (the subnet is an ARM child of the
    # network). References are resolved to a literal ID by the platform
    # before the module runs.
    virtual_network_id = string

    # The subnet's name, unique within the virtual network. Renaming
    # replaces the subnet and everything deployed into it.
    name = string

    # Self-managed CIDR blocks. Exactly one of address_prefixes or
    # ip_address_pool is set (spec-level validation enforces the XOR).
    address_prefixes = optional(list(string), [])

    # Delegated allocation from an Azure Network Manager IPAM pool -- the
    # alternative to self-managed address_prefixes.
    ip_address_pool = optional(object({
      id                     = string
      number_of_ip_addresses = string
    }))

    # Azure service endpoints to enable on the subnet.
    service_endpoints = optional(list(string), [])

    # ARM IDs of service endpoint policies narrowing the endpoints' reach.
    service_endpoint_policy_ids = optional(list(string), [])

    # Service delegations handing the subnet to a PaaS service. actions
    # omitted lets Azure apply the service's default action set.
    delegations = optional(list(object({
      name         = string
      service_name = string
      actions      = optional(list(string), [])
    })), [])

    # Private-endpoint network-policy mode, as the spec enum's name string
    # (ENABLED / NETWORK_SECURITY_GROUP_ENABLED / ROUTE_TABLE_ENABLED).
    # Unset lets Azure apply its default (Disabled).
    private_endpoint_network_policies = optional(string)

    # Whether standard network policies apply to Private Link Service
    # resources in the subnet. Azure defaults to true.
    private_link_service_network_policies_enabled = optional(bool, true)

    # Whether workloads get Azure's implicit default outbound access.
    # Azure's historical default is true; production subnets set false and
    # route egress explicitly (NAT gateway, LB outbound rules, firewall).
    default_outbound_access_enabled = optional(bool, true)

    # Cross-tenant sharing scope ("TENANT"), requiring
    # default_outbound_access_enabled = false. Unset means not shared.
    sharing_scope = optional(string)

    # Optional attachments, as resolved ARM IDs. Each drives one
    # subnet-side association resource.
    route_table_id            = optional(string)
    network_security_group_id = optional(string)
    nat_gateway_id            = optional(string)
  })
}
