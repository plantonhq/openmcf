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
    "resource_kind" = "azure_machine_learning_online_endpoint"
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

  # The spec's identity flavors to the provider's wire values. azapi
  # accepts the azurerm-style comma+space dual value and normalizes it to
  # ARM's own "SystemAssigned,UserAssigned" on the wire -- the same ARM
  # identity type the Pulumi module sends.
  identity_type_wire = {
    "SYSTEM_ASSIGNED"          = "SystemAssigned"
    "USER_ASSIGNED"            = "UserAssigned"
    "SYSTEM_AND_USER_ASSIGNED" = "SystemAssigned, UserAssigned"
  }

  # The ARM properties object, assembled key-by-key so unset optionals are
  # OMITTED (ARM applies its own defaults) instead of sent as nulls.
  endpoint_properties = merge(
    {
      authMode = var.spec.auth_mode
    },
    var.spec.description != "" ? { description = var.spec.description } : {},
    var.spec.public_network_access_enabled != null ? {
      publicNetworkAccess = var.spec.public_network_access_enabled ? "Enabled" : "Disabled"
    } : {},
    length(var.spec.traffic) > 0 ? { traffic = var.spec.traffic } : {},
    length(var.spec.mirror_traffic) > 0 ? { mirrorTraffic = var.spec.mirror_traffic } : {},
    length(var.spec.properties) > 0 ? { properties = var.spec.properties } : {},
  )

  # Bring-your-own auth keys ride azapi's write-only sensitive_body overlay:
  # they merge into the ARM request but never enter Terraform state --
  # exactly right for a property ARM treats as create-time input and never
  # returns on any read (retrieval is the separate listKeys action).
  initial_keys = var.spec.initial_auth_keys != null ? merge(
    var.spec.initial_auth_keys.primary_key != "" ? { primaryKey = var.spec.initial_auth_keys.primary_key } : {},
    var.spec.initial_auth_keys.secondary_key != "" ? { secondaryKey = var.spec.initial_auth_keys.secondary_key } : {},
  ) : {}
}
