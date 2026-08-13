output "machine_learning_workspace_id" {
  description = "The Azure Resource Manager ID of the workspace -- what datastores, compute, and outbound rules reference as their workspace_id"
  value       = azurerm_machine_learning_workspace.main.id
}

output "machine_learning_workspace_name" {
  description = "The name of the workspace. ARM addresses datastores, compute and outbound rules as children of this name"
  value       = azurerm_machine_learning_workspace.main.name
}

output "workspace_guid" {
  description = "The workspace's immutable GUID (distinct from the ARM ID) -- what some data-plane SDKs and diagnostic settings identify the workspace by"
  value       = azurerm_machine_learning_workspace.main.workspace_id
}

output "discovery_url" {
  description = "The workspace's regional discovery URL -- where SDKs resolve the workspace's data-plane service endpoints from"
  value       = azurerm_machine_learning_workspace.main.discovery_url
}

output "system_assigned_identity_principal_id" {
  description = "The principal (object) ID of the workspace's system-assigned identity, when one is enabled"
  value       = try(azurerm_machine_learning_workspace.main.identity[0].principal_id, "")
}

output "fqdn_outbound_rule_ids" {
  description = "The ARM ID of each FQDN outbound rule on the workspace, keyed by the rule's name from the spec"
  value       = { for name, rule in azurerm_machine_learning_workspace_network_outbound_rule_fqdn.fqdn_rules : name => rule.id }
}

output "private_endpoint_outbound_rule_ids" {
  description = "The ARM ID of each private-endpoint outbound rule on the workspace, keyed by the rule's name from the spec"
  value       = { for name, rule in azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint.private_endpoint_rules : name => rule.id }
}

output "service_tag_outbound_rule_ids" {
  description = "The ARM ID of each service-tag outbound rule on the workspace, keyed by the rule's name from the spec"
  value       = { for name, rule in azurerm_machine_learning_workspace_network_outbound_rule_service_tag.service_tag_rules : name => rule.id }
}
