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
  description = "Azure Private Endpoint specification"
  type = object({
    # The region and resource group the endpoint is created in, and the
    # subnet it draws its private IP from. References are resolved to literal
    # names/IDs by the platform before the module runs.
    region         = string
    resource_group = string
    name           = string
    subnet_id      = string

    # The single private link connection this endpoint establishes.
    private_service_connection = object({
      # Exactly one of resource id or alias is set (spec-guaranteed).
      private_connection_resource_id = optional(string, "")
      connection_alias               = optional(string, "")
      subresource_names              = optional(list(string), [])
      # Auto-approved by default; true requires request_message.
      is_manual_connection = optional(bool, false)
      request_message      = optional(string, "")
    })

    # Private DNS zones the endpoint registers its IP into (resolved to
    # literal zone ARM IDs). Empty means no DNS zone group.
    private_dns_zone_ids = optional(list(string), [])

    # Static IP assignments; empty means dynamic allocation from the subnet.
    ip_configurations = optional(list(object({
      name               = string
      private_ip_address = string
      subresource_name   = optional(string, "")
      member_name        = optional(string, "")
    })), [])

    # Application security groups the endpoint's NIC joins (resolved to
    # literal ARM IDs), realized as association resources.
    application_security_group_ids = optional(list(string), [])

    # A custom name for the auto-created network interface; empty lets Azure
    # name it.
    custom_network_interface_name = optional(string, "")

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}
