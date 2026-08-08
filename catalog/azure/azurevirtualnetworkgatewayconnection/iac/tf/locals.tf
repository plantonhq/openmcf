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
    "resource_kind" = "azure_virtual_network_gateway_connection"
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

  # Enum wire maps. tfvars carries FULL proto enum value names; the maps
  # translate them to azurerm's exact (case-sensitive) vocabulary. The
  # connection type is spec-required (never empty for a valid manifest).
  type_wire = {
    "IPSEC"         = "IPsec"
    "VNET_TO_VNET"  = "Vnet2Vnet"
    "EXPRESS_ROUTE" = "ExpressRoute"
  }
  connection_type = lookup(local.type_wire, var.spec.type, null)

  # Sent only when specified -- the provider treats the protocol as
  # Computed, so omission lets Azure apply its default (IKEv2).
  protocol_wire = {
    "IKE_V1" = "IKEv1"
    "IKE_V2" = "IKEv2"
  }
  connection_protocol = (
    var.spec.connection_protocol != null && var.spec.connection_protocol != ""
  ) ? lookup(local.protocol_wire, var.spec.connection_protocol, null) : null

  # Sent explicitly (Default when unspecified) -- deterministic payloads
  # on both engines.
  mode_wire = {
    "DEFAULT"        = "Default"
    "INITIATOR_ONLY" = "InitiatorOnly"
    "RESPONDER_ONLY" = "ResponderOnly"
  }
  connection_mode = lookup(local.mode_wire, coalesce(var.spec.connection_mode, "DEFAULT"), "Default")
}
