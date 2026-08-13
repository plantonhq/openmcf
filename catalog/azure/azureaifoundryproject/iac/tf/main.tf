# Create the AI Foundry project -- ARM-wise an ML workspace of kind
# "Project" linked to its hub. The project INHERITS the hub's posture
# (vault, storage, insights, registry, network, encryption) and
# deploys into the HUB's resource group -- the provider derives the
# group from the hub reference, which is why no resource-group
# argument exists here. The hub linkage is ForceNew.
resource "azurerm_ai_foundry_project" "main" {
  name               = var.spec.name
  location           = var.spec.region
  ai_services_hub_id = var.spec.ai_services_hub_id

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_wire[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  # Only legal alongside the identity block (spec CEL mirrors the
  # provider's RequiredWith).
  primary_user_assigned_identity = var.spec.primary_user_assigned_identity != "" ? var.spec.primary_user_assigned_identity : null

  # Sent only when true (both engines): the property is
  # Optional+Computed and the SERVICE flips it true when hub
  # encryption applies -- a pinned false would fight that read-back.
  # ForceNew.
  high_business_impact_enabled = var.spec.high_business_impact_enabled ? true : null

  description   = var.spec.description != "" ? var.spec.description : null
  friendly_name = var.spec.friendly_name != "" ? var.spec.friendly_name : null

  tags = local.final_tags
}
