# The Front Door route -- connects an endpoint (entry hostname) to an
# origin group (backend pool) by URL pattern. The ARM parent is the
# ENDPOINT (ForceNew); the origin group is the updatable destination.
# No Azure tags: ARM does not support tags on routes.
resource "azurerm_cdn_frontdoor_route" "main" {
  name                          = var.spec.route_name
  cdn_frontdoor_endpoint_id     = var.spec.endpoint_id
  cdn_frontdoor_origin_group_id = var.spec.origin_group_id

  # Azure never receives these -- the provider uses them purely to
  # sequence route creation after the origins exist (ARM rejects a route
  # whose origin group has no origins yet).
  cdn_frontdoor_origin_ids = var.spec.origin_ids

  # Attached delivery policies and custom domains. Front Door treats an
  # EMPTY collection differently from an absent one (empty means
  # "disassociate", which only matters on update), so empty lists are
  # normalized to null -- absence and emptiness then agree.
  cdn_frontdoor_rule_set_ids      = var.spec.rule_set_ids != null && length(coalesce(var.spec.rule_set_ids, [])) > 0 ? var.spec.rule_set_ids : null
  cdn_frontdoor_custom_domain_ids = var.spec.custom_domain_ids != null && length(coalesce(var.spec.custom_domain_ids, [])) > 0 ? var.spec.custom_domain_ids : null

  patterns_to_match   = var.spec.patterns_to_match
  supported_protocols = local.supported_protocols
  forwarding_protocol = local.forwarding_protocol

  # Sent only when explicitly set: Azure's defaults (https redirect on,
  # linked to the default domain, enabled) apply when omitted.
  https_redirect_enabled = var.spec.https_redirect_enabled
  link_to_default_domain = var.spec.link_to_default_domain
  enabled                = var.spec.enabled

  cdn_frontdoor_origin_path = var.spec.origin_path

  # The cache block is sent only when configured: Front Door treats
  # ABSENT cache settings as caching disabled (the provider transmits an
  # explicit null), so omitting the block is a real behavior switch.
  dynamic "cache" {
    for_each = var.spec.cache != null ? [var.spec.cache] : []
    content {
      query_string_caching_behavior = local.query_string_caching_behavior_map[coalesce(cache.value.query_string_caching_behavior, "IGNORE_QUERY_STRING")]
      query_strings                 = cache.value.query_strings
      compression_enabled           = cache.value.compression_enabled
      content_types_to_compress     = cache.value.content_types_to_compress
    }
  }
}
