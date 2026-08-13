output "gallery_id" {
  description = "The Azure Resource Manager ID of the gallery"
  value       = azurerm_shared_image_gallery.main.id
}

output "gallery_name" {
  description = "The gallery's name -- what image definitions reference as their gallery"
  value       = azurerm_shared_image_gallery.main.name
}

output "unique_name" {
  description = "The globally unique name Azure assigns the gallery (used in cross-tenant and community addressing)"
  value       = azurerm_shared_image_gallery.main.unique_name
}

output "community_gallery_name" {
  description = "The public community-gallery name generated from the sharing prefix; empty unless Community-shared"
  value       = try(azurerm_shared_image_gallery.main.sharing[0].community_gallery[0].name, "")
}
