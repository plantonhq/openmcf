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
  description = "AwsSubnet specification"
  type = object({
    region = string
    vpc_id = string
    availability_zone = optional(string, "")
    cidr_block = optional(string, "")
    map_public_ip_on_launch = optional(bool, false)
    assign_ipv6_address_on_creation = optional(bool, false)
    ipv6_cidr_block = optional(string, "")
    enable_dns64 = optional(bool, false)
    enable_resource_name_dns_a_record_on_launch = optional(bool, false)
    enable_resource_name_dns_aaaa_record_on_launch = optional(bool, false)
    private_dns_hostname_type_on_launch = optional(string)
    route_table_id = optional(string, "")
    routes = optional(list(object({
      destination_cidr_block = optional(string, "")
      destination_ipv6_cidr_block = optional(string, "")
      destination_prefix_list_id = optional(string, "")
      target_type = optional(string, "")
      target_id = string
    })), [])
    availability_zone_id = optional(string, "")
    ipv4_ipam_pool_id = optional(string, "")
    ipv4_netmask_length = optional(number)
    ipv6_ipam_pool_id = optional(string, "")
    ipv6_netmask_length = optional(number)
    ipv6_native = optional(bool, false)
    propagating_vgws = optional(list(string), [])
  })
}