output "machine_learning_compute_cluster_id" {
  description = "The Azure Resource Manager ID of the compute cluster"
  value       = azurerm_machine_learning_compute_cluster.main.id
}

output "machine_learning_compute_cluster_name" {
  description = "The cluster's name -- what jobs and pipelines reference as their compute target within the workspace"
  value       = azurerm_machine_learning_compute_cluster.main.name
}

output "system_assigned_identity_principal_id" {
  description = "The principal (object) ID of the cluster's system-assigned identity, when one is enabled"
  value       = try(azurerm_machine_learning_compute_cluster.main.identity[0].principal_id, "")
}
