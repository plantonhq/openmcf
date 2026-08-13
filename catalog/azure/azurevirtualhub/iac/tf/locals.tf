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
    "resource_kind" = "azure_virtual_hub"
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

  # The spec's enum NAMES mapped onto ARM's vocabulary. Unset optionals
  # apply ARM's defaults explicitly so the rendered plan shows the real
  # value -- mirroring the Pulumi module's nil handling.
  sku = (
    var.spec.sku == null || var.spec.sku == ""
    ? "Standard"
    : lookup({ "BASIC" = "Basic", "STANDARD" = "Standard" }, var.spec.sku, var.spec.sku)
  )

  hub_routing_preference_wire = {
    "EXPRESS_ROUTE" = "ExpressRoute"
    "VPN_GATEWAY"   = "VpnGateway"
    "AS_PATH"       = "ASPath"
  }
  hub_routing_preference = (
    var.spec.hub_routing_preference == null || var.spec.hub_routing_preference == ""
    ? "ExpressRoute"
    : lookup(local.hub_routing_preference_wire, var.spec.hub_routing_preference, var.spec.hub_routing_preference)
  )

  destinations_type_wire = {
    "CIDR"        = "CIDR"
    "RESOURCE_ID" = "ResourceId"
    "SERVICE"     = "Service"
  }

  match_condition_wire = {
    "CONTAINS"     = "Contains"
    "EQUALS"       = "Equals"
    "NOT_CONTAINS" = "NotContains"
    "NOT_EQUALS"   = "NotEquals"
  }

  action_type_wire = {
    "ADD"     = "Add"
    "DROP"    = "Drop"
    "REMOVE"  = "Remove"
    "REPLACE" = "Replace"
  }

  # Unset leaves the provider's default ("Unknown": evaluation stops
  # after the match) -- the spec deliberately does not model ARM's
  # "Unknown" as a choosable value.
  next_step_wire = {
    "CONTINUE"  = "Continue"
    "TERMINATE" = "Terminate"
  }

  routing_policy_destination_wire = {
    "INTERNET"        = "Internet"
    "PRIVATE_TRAFFIC" = "PrivateTraffic"
  }
}
