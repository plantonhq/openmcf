# Create the Data Protection Resource Guard -- the approval gate
# behind Multi-User Authorization. The guard's protection comes from
# SCOPE SEPARATION: deploy it in a resource group a different
# administrator controls than the vaults it guards; a guard in the
# same scope as its vaults is a speed bump, not a control.
#
# The guard is a free configuration object; vaults opt in by
# referencing its ARM ID.
resource "azurerm_data_protection_resource_guard" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  # Operations EXCLUDED from the approval requirement. Empty guards
  # everything (the strongest posture). Updates in place.
  vault_critical_operation_exclusion_list = length(var.spec.vault_critical_operation_exclusion_list) > 0 ? var.spec.vault_critical_operation_exclusion_list : null

  tags = local.final_tags
}
