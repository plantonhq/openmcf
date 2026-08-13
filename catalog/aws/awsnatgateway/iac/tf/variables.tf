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
  description = "AwsNatGateway specification"
  type = object({
    region = string
    connectivity_type = optional(string, "")
    subnet_id = optional(string, "")
    allocation_id = optional(string, "")
    private_ip = optional(string, "")
    secondary_allocation_ids = optional(list(string), [])
    secondary_private_ip_addresses = optional(list(string), [])
    secondary_private_ip_address_count = optional(number, 0)
    availability_mode = optional(string, "")
    vpc_id = optional(string, "")
    availability_zone_addresses = optional(list(object({
      availability_zone = optional(string, "")
      availability_zone_id = optional(string, "")
      allocation_ids = optional(list(string), [])
    })), [])
  })
}
