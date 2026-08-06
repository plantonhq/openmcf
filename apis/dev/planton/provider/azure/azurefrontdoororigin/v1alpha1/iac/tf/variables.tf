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
  description = "Azure Front Door origin specification"
  type = object({
    # The origin group the origin belongs to, by ARM ID. References are
    # resolved to a literal ID by the platform before the module runs.
    # ForceNew.
    origin_group_id = string

    # The origin's name -- unique within the origin group. ForceNew.
    origin_name = string

    # The backend's address: hostname, IPv4, or IPv6.
    host_name = string

    # TLS certificate-name validation. Absent means true (tfvars drops
    # zero-valued proto fields; the provider requires the value, so the
    # module materializes the documented default). Azure requires true
    # with private_link.
    certificate_name_check_enabled = optional(bool)

    # Host header override toward the origin. Absent uses host_name
    # (Azure's behavior) -- correct for multi-tenant Azure backends.
    origin_host_header = optional(string)

    # Origin-leg ports. Azure defaults to 80/443 when omitted.
    http_port  = optional(number)
    https_port = optional(number)

    # Failover tier (1-5, lower serves first) and traffic share within
    # the tier (1-1000). Azure defaults: priority 1, weight 500.
    priority = optional(number)
    weight   = optional(number)

    # Whether the origin receives traffic. Azure defaults to true.
    enabled = optional(bool)

    # Private Link to the backend (Premium profiles only; Azure enforces
    # the SKU at apply because it lives on a different resource).
    private_link = optional(object({
      # The target resource's region -- private-link connections are
      # regional even though Front Door is global.
      location = string
      # The ARM ID of the target (site, storage account, Container Apps
      # environment, or Private Link Service).
      private_link_target_id = string
      # The sub-resource to attach to, as the spec enum's value name
      # (SITES / BLOB / BLOB_SECONDARY / WEB / WEB_SECONDARY /
      # MANAGED_ENVIRONMENTS / GATEWAY). Absent ONLY for Private Link
      # Service targets (the spec CEL enforces the pairing).
      target_type = optional(string)
      # The approval message shown to the target's owner.
      request_message = optional(string)
    }))
  })
}
