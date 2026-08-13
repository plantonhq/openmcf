output "machine_learning_compute_instance_id" {
  description = "The Azure Resource Manager ID of the compute instance"
  value       = azurerm_machine_learning_compute_instance.main.id
}

output "machine_learning_compute_instance_name" {
  description = "The instance's name -- what its owner selects as their compute in notebooks and the ML studio"
  value       = azurerm_machine_learning_compute_instance.main.name
}

output "system_assigned_identity_principal_id" {
  description = "The principal (object) ID of the instance's system-assigned identity, when one is enabled"
  value       = try(azurerm_machine_learning_compute_instance.main.identity[0].principal_id, "")
}

output "ssh_username" {
  description = "The admin username for SSH access, assigned by the service -- populated only when the ssh block is configured"
  value       = try(azurerm_machine_learning_compute_instance.main.ssh[0].username, "")
}

output "ssh_port" {
  description = "The port the instance answers SSH on, assigned by the service -- populated only when the ssh block is configured"
  value       = try(azurerm_machine_learning_compute_instance.main.ssh[0].port, 0)
}
