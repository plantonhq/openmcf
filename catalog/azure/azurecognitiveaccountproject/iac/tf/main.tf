# Create the AI Foundry project on its Azure AI services account --
# the workspace a team's agents, evaluations and files live in,
# isolated from sibling projects on the same account. The parent
# account must be kind "AIServices" with project_management_enabled
# true (and therefore a managed identity); the FIRST project created
# on an account becomes the account's default (the is_default output).
#
# ARM cannot UPDATE description or display_name to an EMPTY value --
# the provider replaces the project when either is cleared (setting or
# changing them updates in place).
resource "azurerm_cognitive_account_project" "main" {
  name = var.spec.name

  # The parent account. ForceNew.
  cognitive_account_id = var.spec.cognitive_account_id

  location = var.spec.region

  # Required by the provider: every project carries an identity -- it
  # is what the project's agents and evaluations act as.
  identity {
    type         = local.identity_type_wire[var.spec.identity.type]
    identity_ids = length(var.spec.identity.identity_ids) > 0 ? var.spec.identity.identity_ids : null
  }

  description  = var.spec.description != "" ? var.spec.description : null
  display_name = var.spec.display_name != "" ? var.spec.display_name : null

  tags = local.final_tags
}
