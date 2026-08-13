output "availability_set_id" {
  description = "The Azure Resource Manager ID of the availability set -- what VMs reference to join it"
  value       = azurerm_availability_set.main.id
}

output "availability_set_name" {
  description = "The availability set's name"
  value       = azurerm_availability_set.main.name
}
