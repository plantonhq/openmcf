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
  description = "Azure Firewall Policy Rule Collection Group specification"
  type = object({
    # The parent policy's ARM id (resolved literal). The group nests under
    # this policy; changing it replaces the group.
    firewall_policy_id = string

    # The group's name, unique within the policy. Renaming replaces the
    # group.
    name = string

    # Group evaluation priority among the policy's groups: 100-65000,
    # lower first.
    priority = number

    # Application rule collections (evaluated after DNAT and network
    # rules). Actions and protocol types arrive as proto enum value names
    # (ALLOW/DENY; HTTP/HTTPS/MSSQL).
    application_rule_collections = optional(list(object({
      name     = string
      priority = number
      action   = string
      rules = list(object({
        name        = string
        description = optional(string)
        protocols = optional(list(object({
          type = string
          port = number
        })), [])
        http_headers = optional(list(object({
          name  = string
          value = string
        })), [])
        source_addresses      = optional(list(string), [])
        source_ip_groups      = optional(list(string), [])
        destination_addresses = optional(list(string), [])
        destination_fqdns     = optional(list(string), [])
        destination_urls      = optional(list(string), [])
        destination_fqdn_tags = optional(list(string), [])
        terminate_tls         = optional(bool, false)
        web_categories        = optional(list(string), [])
      }))
    })), [])

    # Network rule collections (evaluated after DNAT, before application
    # rules). Protocols arrive as proto enum value names
    # (ANY/TCP/UDP/ICMP).
    network_rule_collections = optional(list(object({
      name     = string
      priority = number
      action   = string
      rules = list(object({
        name                  = string
        description           = optional(string)
        protocols             = list(string)
        source_addresses      = optional(list(string), [])
        source_ip_groups      = optional(list(string), [])
        destination_addresses = optional(list(string), [])
        destination_ip_groups = optional(list(string), [])
        destination_fqdns     = optional(list(string), [])
        destination_ports     = list(string)
      }))
    })), [])

    # DNAT rule collections (evaluated first; a match implicitly allows
    # the translated flow). The collection action is always Azure's
    # "Dnat" -- a one-value vocabulary the module sends unconditionally.
    nat_rule_collections = optional(list(object({
      name     = string
      priority = number
      rules = list(object({
        name                = string
        description         = optional(string)
        protocols           = list(string)
        source_addresses    = optional(list(string), [])
        source_ip_groups    = optional(list(string), [])
        destination_address = optional(string)
        # ARM caps DNAT destination ports at ONE entry today; the list
        # shape mirrors ARM's own.
        destination_ports  = optional(list(string), [])
        translated_address = optional(string)
        translated_fqdn    = optional(string)
        translated_port    = number
      }))
    })), [])
  })
}
