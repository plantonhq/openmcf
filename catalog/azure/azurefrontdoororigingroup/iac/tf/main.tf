# The Front Door origin group -- the load-balanced pool of backends a
# route sends traffic to. The backends are AzureFrontDoorOrigin
# components referencing this group. The parent is addressed by the
# profile's full ARM id; the provider derives the resource group and
# profile name from it. No Azure tags: ARM does not support tags on
# origin groups.
resource "azurerm_cdn_frontdoor_origin_group" "main" {
  name                     = var.spec.origin_group_name
  cdn_frontdoor_profile_id = var.spec.profile_id

  # Azure requires load-balancing settings on every origin group, so the
  # block is always sent -- with the spec's values when present,
  # otherwise Azure's own defaults (sample size 4, 3 successful samples
  # required, 50 ms additional latency), exactly what an unset spec
  # block documents. tfvars drops zero-valued proto fields, so each
  # field materializes its documented default here.
  load_balancing {
    sample_size                        = try(coalesce(var.spec.load_balancing.sample_size, 4), 4)
    successful_samples_required        = try(coalesce(var.spec.load_balancing.successful_samples_required, 3), 3)
    additional_latency_in_milliseconds = try(coalesce(var.spec.load_balancing.additional_latency_in_milliseconds, 50), 50)
  }

  # The health probe is sent only when configured: Front Door treats
  # ABSENT probe settings as probing disabled (all origins assumed
  # healthy). Note for maintainers: Azure's PATCH would silently null
  # probe settings on unrelated updates; azurerm ships a workaround
  # client for this, and both engines inherit it from the provider layer.
  dynamic "health_probe" {
    for_each = var.spec.health_probe != null ? [var.spec.health_probe] : []
    content {
      protocol            = local.health_probe_protocol_map[health_probe.value.protocol]
      interval_in_seconds = health_probe.value.interval_in_seconds
      request_type        = local.health_probe_request_type_map[coalesce(health_probe.value.request_type, "HEAD")]
      path                = coalesce(health_probe.value.path, "/")
    }
  }

  # Sent only when explicitly set: Azure's defaults are session affinity
  # on and a 10-minute traffic-restore ramp.
  session_affinity_enabled                                   = var.spec.session_affinity_enabled
  restore_traffic_time_to_healed_or_new_endpoint_in_minutes = var.spec.restore_traffic_time_to_healed_or_new_endpoint_in_minutes
}
