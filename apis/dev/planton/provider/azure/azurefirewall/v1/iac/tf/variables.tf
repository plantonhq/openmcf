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
  description = "Azure Firewall specification"
  type = object({
    # The Azure region the firewall lives in; must match its virtual
    # network (or hub).
    region = string

    # The resource group the firewall is created in. References are
    # resolved to a literal name by the platform before the module runs.
    resource_group = string

    # The firewall's name, unique within the resource group. Renaming
    # replaces the firewall.
    name = string

    # The deployment model as the proto enum value name
    # (AZFW_VNET/AZFW_HUB). Absent means AZFW_VNET.
    sku_name = optional(string)

    # The tier as the proto enum value name (BASIC/STANDARD/PREMIUM).
    # Absent means STANDARD. Must match the attached policy's tier.
    sku_tier = optional(string)

    # Data-path IP configurations: exactly one carries subnet_id (the
    # "AzureFirewallSubnet" subnet, /26+); extra blocks add public IPs.
    # Ids arrive as resolved ARM-id literals.
    ip_configurations = optional(list(object({
      name                 = string
      subnet_id            = optional(string)
      public_ip_address_id = optional(string)
    })), [])

    # The management path (forced tunneling / BASIC tier): the dedicated
    # "AzureFirewallManagementSubnet" (/26+) and its REQUIRED public IP.
    # Fixed at creation.
    management_ip_configuration = optional(object({
      name                 = string
      subnet_id            = string
      public_ip_address_id = string
    }))

    # The firewall policy this instance enforces (resolved ARM id).
    firewall_policy_id = optional(string)

    # Threat-intelligence posture as the proto enum value name
    # (ALERT/DENY/OFF) -- only meaningful without a policy. Absent lets
    # Azure own its default (Alert).
    threat_intel_mode = optional(string)

    # Custom upstream DNS servers. Setting servers implicitly turns the
    # DNS proxy ON in Azure's wire encoding.
    dns_servers = optional(list(string), [])

    # Run the firewall as a DNS proxy.
    dns_proxy_enabled = optional(bool, false)

    # SNAT-exempt ranges: CIDRs or the literal "IANAPrivateRanges" token.
    private_ip_ranges = optional(list(string), [])

    # The Virtual WAN hub target (AZFW_HUB model only).
    virtual_hub = optional(object({
      virtual_hub_id  = string
      public_ip_count = optional(number)
    }))

    # Availability zones, e.g. ["1", "2", "3"]. Fixed at creation.
    zones = optional(list(string), [])

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
