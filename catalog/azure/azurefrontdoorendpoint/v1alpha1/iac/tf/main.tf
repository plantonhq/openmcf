# The Front Door endpoint -- the public entry point client traffic
# arrives at. Azure generates a globally unique {name}-{hash}.z01.azurefd.net
# hostname (surfaced as the host_name output), which is why the endpoint
# name only needs per-profile uniqueness. The parent is addressed by the
# profile's full ARM id; the provider derives the resource group and
# profile name from it.
resource "azurerm_cdn_frontdoor_endpoint" "main" {
  name                     = var.spec.endpoint_name
  cdn_frontdoor_profile_id = var.spec.profile_id

  # Sent only when explicitly disabled: Azure's default is enabled, and
  # tfvars drops zero-valued proto fields (the platform materializes the
  # documented default centrally, so absent means true).
  enabled = var.spec.enabled

  tags = local.final_tags
}
