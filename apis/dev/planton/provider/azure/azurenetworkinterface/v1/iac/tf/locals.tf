locals {
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # PARITY-EXCEPTION: resource_kind here is the family-wide snake-case
  # literal and resource_id falls back to metadata.name, while the Pulumi
  # module emits the lowered CloudResourceKind enum string and omits
  # resource_id when metadata.id is empty. Output-neutral (tags never feed
  # stack outputs); aligning the two shapes is a family-wide convention
  # change, not a per-kind fix.
  base_tags = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_network_interface"
    "resource_name" = var.metadata.name
  }

  org_tag = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "organization" = var.metadata.org } : {}

  env_tag = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "environment" = var.metadata.env } : {}

  # Metadata-derived tags first, then the user's spec tags merged over them:
  # user tags deliberately win so an org's governance conventions (cost
  # center, owner) can override the derived values where they collide.
  final_tags = merge(local.base_tags, local.org_tag, local.env_tag, var.spec.tags)

  # Map the spec enums' name strings to ARM values. Unset applies Azure's
  # defaults (Dynamic allocation, IPv4), so an unspecified spec and Azure's
  # default deploy identically on both engines.
  ip_configurations = [
    for c in var.spec.ip_configurations : {
      name                 = c.name
      subnet_id            = c.subnet_id
      allocation           = c.private_ip_allocation == "STATIC" ? "Static" : "Dynamic"
      private_ip_address   = c.private_ip_address
      version              = c.private_ip_version == "IPV6" ? "IPv6" : "IPv4"
      public_ip_address_id = c.public_ip_address_id
      primary              = c.primary
      gateway_lb_fip_id    = c.gateway_load_balancer_frontend_ip_configuration_id
    }
  ]

  # ARM auxiliary values: AcceleratedConnections/Floating/MaxConnections
  # and A1/A2/A4/A8; null sends nothing (the non-appliance default).
  auxiliary_mode = (
    var.spec.auxiliary_mode == "ACCELERATED_CONNECTIONS" ? "AcceleratedConnections" :
    var.spec.auxiliary_mode == "FLOATING" ? "Floating" :
    var.spec.auxiliary_mode == "MAX_CONNECTIONS" ? "MaxConnections" : null
  )
  auxiliary_sku = (
    var.spec.auxiliary_sku == "A1" ? "A1" :
    var.spec.auxiliary_sku == "A2" ? "A2" :
    var.spec.auxiliary_sku == "A4" ? "A4" :
    var.spec.auxiliary_sku == "A8" ? "A8" : null
  )
}
