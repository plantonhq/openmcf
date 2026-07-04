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
    "resource_kind" = "azure_container_registry"
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

  # Map the spec enum's name string to ARM's SKU value. The unspecified
  # spec applies the STANDARD baseline on both engines (azurerm requires
  # an explicit sku, so the default is materialized here rather than sent
  # as null).
  sku = (
    var.spec.sku == "BASIC" ? "Basic" :
    var.spec.sku == "PREMIUM" ? "Premium" : "Standard"
  )

  # Map the bypass enum's name string to ARM's value. null lets Azure
  # apply its default (AzureServices); only an explicit choice is sent, so
  # an unspecified spec and Azure's default deploy identically on both
  # engines.
  network_rule_bypass_option = (
    var.spec.network_rule_bypass_option == "AZURE_SERVICES" ? "AzureServices" :
    var.spec.network_rule_bypass_option == "NONE" ? "None" : null
  )

  # Map the identity type enum's name string to ARM's comma-separated
  # value.
  identity_type = (
    var.spec.identity == null ? null :
    var.spec.identity.type == "SYSTEM_ASSIGNED" ? "SystemAssigned" :
    var.spec.identity.type == "USER_ASSIGNED" ? "UserAssigned" :
    var.spec.identity.type == "SYSTEM_AND_USER_ASSIGNED" ? "SystemAssigned, UserAssigned" : null
  )

  # Map the network-rule default action enum's name string to ARM's value.
  # Azure's default (Allow) applies when unspecified.
  network_rule_default_action = (
    var.spec.network_rule_set == null ? null :
    var.spec.network_rule_set.default_action == "DENY" ? "Deny" : "Allow"
  )
}
