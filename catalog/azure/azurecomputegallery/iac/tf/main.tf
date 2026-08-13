# Create the Azure Compute Gallery -- the shared library an
# organization keeps its approved VM images in. A gallery is free at
# rest. The ENTIRE sharing tree is create-only in the provider
# (changing it forces replacement); Community sharing requires the
# community_gallery block (the spec's CEL front-loads the provider's
# expand-time check).
resource "azurerm_shared_image_gallery" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  description = (
    var.spec.description != null && var.spec.description != ""
    ? var.spec.description
    : null
  )

  dynamic "sharing" {
    for_each = var.spec.sharing != null ? [var.spec.sharing] : []
    content {
      permission = sharing.value.permission

      dynamic "community_gallery" {
        for_each = sharing.value.community_gallery != null ? [sharing.value.community_gallery] : []
        content {
          eula            = community_gallery.value.eula
          prefix          = community_gallery.value.prefix
          publisher_email = community_gallery.value.publisher_email
          publisher_uri   = community_gallery.value.publisher_uri
        }
      }
    }
  }

  tags = local.final_tags
}
