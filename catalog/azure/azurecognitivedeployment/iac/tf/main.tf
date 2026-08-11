# Create the model deployment on its Azure AI services account --
# which actual model applications call through the account's endpoint
# (the deployment's NAME is the model parameter they pass). The
# pay-per-token SKUs (Standard/GlobalStandard/DataZone*) carry no idle
# cost -- capacity is a rate limit in thousands of tokens-per-minute;
# the ProvisionedManaged SKUs bill their PTU capacity continuously
# while the deployment exists.
#
# Which models exist at which capacities differs per region and per
# subscription quota -- ARM rejects what the region or quota cannot
# host, so quota errors surface here, never as silent degradation.
resource "azurerm_cognitive_deployment" "main" {
  name = var.spec.name

  # The parent account (kind "OpenAI" or "AIServices"). ForceNew.
  cognitive_account_id = var.spec.cognitive_account_id

  dynamic_throttling_enabled = var.spec.dynamic_throttling_enabled

  # Format and name are ForceNew -- a different model is a new
  # deployment; version updates in place (omit it to track the
  # model's current default version).
  model {
    format  = var.spec.model.format
    name    = var.spec.model.name
    version = var.spec.model.version != "" ? var.spec.model.version : null
  }

  sku {
    # Already a wire value in the spec.
    name = var.spec.sku.name
    # Enum name -> wire value; unspecified ("") omits the property so
    # ARM derives the tier from the SKU name.
    tier   = lookup(local.sku_tier_wire, var.spec.sku.tier, null)
    size   = var.spec.sku.size != "" ? var.spec.sku.size : null
    family = var.spec.sku.family != "" ? var.spec.sku.family : null
    # Unset applies the provider default of 1. The in-place scale knob.
    capacity = var.spec.sku.capacity
  }

  # Optional+Computed on the provider: emit null when the spec leaves
  # it empty so ARM assigns its default policy and reads don't drift.
  rai_policy_name = var.spec.rai_policy_name != "" ? var.spec.rai_policy_name : null

  # Enum name -> wire value; unspecified ("") lets the provider apply
  # its default, "OnceNewDefaultVersionAvailable".
  version_upgrade_option = lookup(local.version_upgrade_option_wire, var.spec.version_upgrade_option, null)
}
