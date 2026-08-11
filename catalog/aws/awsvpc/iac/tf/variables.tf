variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name = string
    id = optional(string, "")
    org = optional(string, "")
    env = optional(string, "")
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    tags = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsVpc specification"
  type = object({
    region = string
    cidr_block = optional(string, "")
    secondary_ipv4_cidrs = optional(list(object({
      cidr_block = optional(string, "")
      ipam_pool_id = optional(string, "")
      netmask_length = optional(number, 0)
    })), [])
    ipv4_ipam_pool_id = optional(string, "")
    ipv4_netmask_length = optional(number, 0)
    instance_tenancy = optional(string, "")
    enable_dns_support = optional(bool)
    enable_dns_hostnames = optional(bool, false)
    enable_network_address_usage_metrics = optional(bool, false)
    assign_generated_ipv6_cidr_block = optional(bool, false)
    ipv6_cidr_block = optional(string, "")
    ipv6_cidr_block_network_border_group = optional(string, "")
    ipv6_ipam_pool_id = optional(string, "")
    ipv6_netmask_length = optional(number, 0)
    secondary_ipv6_cidrs = optional(list(object({
      assign_generated = optional(bool, false)
      ipv6_pool = optional(string, "")
      ipam_pool_id = optional(string, "")
      cidr_block = optional(string, "")
      netmask_length = optional(number, 0)
    })), [])
    encryption_control = optional(object({
      mode = string
      exclude_internet_gateway = optional(bool, false)
      exclude_egress_only_internet_gateway = optional(bool, false)
      exclude_nat_gateway = optional(bool, false)
      exclude_virtual_private_gateway = optional(bool, false)
      exclude_vpc_peering = optional(bool, false)
      exclude_vpc_lattice = optional(bool, false)
      exclude_lambda = optional(bool, false)
      exclude_elastic_file_system = optional(bool, false)
    }))
  })
}
