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
  description = "Azure Network Interface specification"
  type = object({
    # The Azure region the NIC lives in (must match its virtual network
    # and the VM that attaches it).
    region = string

    # The resource group the NIC lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The NIC's name, unique within the resource group. Renaming replaces
    # the NIC and detaches it from any VM.
    name = string

    # The NIC's IP configurations (at least one). When multiple are
    # declared, the first must be marked primary (ARM's contract;
    # spec-level validation enforces it).
    ip_configurations = list(object({
      # A configuration-level label, unique within the NIC.
      name = string

      # The subnet the private address lives in, as a resolved ARM ID.
      # Required for IPv4 configurations; IPv6 configurations inherit the
      # NIC's subnet placement.
      subnet_id = optional(string)

      # DYNAMIC (unset) lets Azure pick a free address; STATIC pins
      # private_ip_address.
      private_ip_allocation = optional(string)

      # For STATIC allocation: the exact address to pin.
      private_ip_address = optional(string)

      # The address family, as the spec enum's name string (IPV4/IPV6).
      # Unset applies Azure's default (IPv4).
      private_ip_version = optional(string)

      # The fronting public IP, as a resolved ARM ID (omit for
      # private-only NICs).
      public_ip_address_id = optional(string)

      # Whether this is the NIC's primary configuration.
      primary = optional(bool, false)

      # Gateway-SKU load balancer frontend for service chaining (niche).
      gateway_load_balancer_frontend_ip_configuration_id = optional(string)

      # Load-balancer backend pools this configuration joins, as resolved
      # ARM IDs; each membership is realized as an association resource.
      load_balancer_backend_address_pool_ids = optional(list(string), [])

      # Single-target inbound NAT rules this configuration completes, as
      # resolved ARM IDs; realized as association resources.
      load_balancer_inbound_nat_rule_ids = optional(list(string), [])

      # Application Gateway backend pools this configuration joins, as
      # ARM IDs; realized as association resources.
      application_gateway_backend_address_pool_ids = optional(list(string), [])
    }))

    # DNS servers overriding the virtual network's DNS for this NIC only.
    dns_servers = optional(list(string), [])

    # The VNet-internal DNS label other VMs can resolve this NIC by.
    internal_dns_name_label = optional(string)

    # SR-IOV acceleration (Azure defaults to false; enable on supported VM
    # sizes).
    accelerated_networking_enabled = optional(bool, false)

    # Whether the NIC forwards traffic not addressed to it (network
    # virtual appliances only).
    ip_forwarding_enabled = optional(bool, false)

    # NVA-acceleration auxiliary mode/SKU (preview), as the spec enums'
    # name strings. Both or neither (spec-level validation enforces the
    # pairing).
    auxiliary_mode = optional(string)
    auxiliary_sku  = optional(string)

    # Edge Zone pinning for edge-computing workloads (fixed at creation).
    edge_zone = optional(string)

    # The NSG filtering this NIC's traffic, as a resolved ARM ID; realized
    # as an association resource.
    network_security_group_id = optional(string)

    # Application security groups this NIC joins, as ARM IDs; each is
    # realized as an association resource.
    application_security_group_ids = optional(list(string), [])

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
