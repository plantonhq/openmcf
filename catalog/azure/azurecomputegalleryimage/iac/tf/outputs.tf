output "image_id" {
  description = "The Azure Resource Manager ID of the image definition -- deploying from it gets the latest (non-excluded) version"
  value       = azurerm_shared_image.main.id
}

output "image_name" {
  description = "The image definition's name within its gallery"
  value       = azurerm_shared_image.main.name
}

output "version_ids" {
  description = "The ARM IDs of the image's published versions, keyed by version name -- VMs pin to an exact release through these"
  value       = { for name, version in azurerm_shared_image_version.main : name => version.id }
}
