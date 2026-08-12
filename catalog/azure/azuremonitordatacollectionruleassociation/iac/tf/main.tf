# Attach one machine to an Azure Monitor data collection rule or data
# collection endpoint. The association is an extension resource living
# ON the target machine ({target_id}/providers/Microsoft.Insights/
# dataCollectionRuleAssociations/{name}) -- creating it puts the
# machine under the rule, destroying it detaches the machine without
# touching the rule. Free at rest.
resource "azurerm_monitor_data_collection_rule_association" "main" {
  target_resource_id = var.spec.target_resource_id

  # Left unset for endpoint bindings so the provider applies Azure's
  # mandated fixed name ("configurationAccessEndpoint"); required (and
  # spec-enforced) for rule bindings.
  name = var.spec.name != "" ? var.spec.name : null

  # Exactly one of the two bindings (spec CEL mirrors the provider's
  # ExactlyOneOf).
  data_collection_rule_id     = var.spec.data_collection_rule_id
  data_collection_endpoint_id = var.spec.data_collection_endpoint_id

  # Sent only when non-empty for a clean plan; Azure treats an absent
  # and an empty description identically.
  description = var.spec.description != "" ? var.spec.description : null
}
