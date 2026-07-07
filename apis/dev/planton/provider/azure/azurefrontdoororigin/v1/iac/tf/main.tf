# The Front Door origin -- one backend inside an origin group. The
# parent is addressed by the origin group's full ARM id; the provider
# derives the resource group, profile, and group names from it. No Azure
# tags: ARM does not support tags on origins.
resource "azurerm_cdn_frontdoor_origin" "main" {
  name                          = var.spec.origin_name
  cdn_frontdoor_origin_group_id = var.spec.origin_group_id
  host_name                     = var.spec.host_name

  # Always sent (the provider requires it): keeping certificate-name
  # validation on is the secure posture, and Azure requires it with
  # Private Link.
  certificate_name_check_enabled = local.certificate_name_check_enabled

  # Optional dials sent only when set: Azure's own defaults (ports
  # 80/443, priority 1, weight 500, enabled) apply when omitted.
  origin_host_header = var.spec.origin_host_header
  http_port          = var.spec.http_port
  https_port         = var.spec.https_port
  priority           = var.spec.priority
  weight             = var.spec.weight
  enabled            = var.spec.enabled

  # Private Link keeps origin traffic off the public internet.
  # PREMIUM-profile only -- Azure rejects it at apply on STANDARD (the
  # SKU lives on a different resource, so it cannot be checked
  # statically). target_type is omitted for Private Link Service
  # targets, whose ARM id is itself the attachment point.
  dynamic "private_link" {
    for_each = var.spec.private_link != null ? [var.spec.private_link] : []
    content {
      location               = private_link.value.location
      private_link_target_id = private_link.value.private_link_target_id
      target_type            = private_link.value.target_type != null ? local.private_link_target_type_map[private_link.value.target_type] : null
      request_message        = private_link.value.request_message
    }
  }
}
