# Create the Traffic Manager profile -- Azure's DNS-based traffic
# director. The profile owns {relative_name}.trafficmanager.net and
# answers each lookup with the address of one of its endpoints
# (AzureTrafficManagerEndpoint resources referencing this profile),
# chosen by routing method and endpoint health.
#
# Traffic Manager is GLOBAL: the provider pins the ARM location to
# "global" itself, which is why the spec carries no region. Everything
# except the resource name, the resource group, and the DNS relative
# name updates in place.
resource "azurerm_traffic_manager_profile" "main" {
  name                   = var.spec.name
  resource_group_name    = var.spec.resource_group
  traffic_routing_method = var.spec.routing_method
  profile_status         = local.profile_status
  # Present only for MultiValue routing (spec validation requires it
  # there); the provider sends it only when set.
  max_return           = var.spec.max_return
  traffic_view_enabled = coalesce(var.spec.traffic_view_enabled, false)
  tags                 = local.final_tags

  dns_config {
    # Globally unique across ALL of Azure (the trafficmanager.net
    # namespace is shared) -- Azure rejects a taken name at apply time.
    relative_name = var.spec.dns_config.relative_name
    # The TTL is always sent explicitly -- 60 (the Azure portal's own
    # default) when the spec leaves it unset -- so plans stay
    # deterministic.
    ttl = coalesce(var.spec.dns_config.ttl_seconds, 60)
  }

  monitor_config {
    protocol = var.spec.monitor_config.protocol
    port     = var.spec.monitor_config.port
    path     = var.spec.monitor_config.path != "" ? var.spec.monitor_config.path : null
    # Probe cadence defaults (30s interval / 10s timeout / 3 tolerated
    # failures) are always sent explicitly so both engines send
    # identical wire shapes. Spec validation already enforces the
    # provider's fast-interval contract (interval 10 narrows timeout
    # to an explicit 5-9).
    interval_in_seconds          = coalesce(var.spec.monitor_config.interval_in_seconds, 30)
    timeout_in_seconds           = coalesce(var.spec.monitor_config.timeout_in_seconds, 10)
    tolerated_number_of_failures = coalesce(var.spec.monitor_config.tolerated_number_of_failures, 3)
    # "min-max" strings (e.g. "200-299"); unset expects exactly 200.
    expected_status_code_ranges = length(var.spec.monitor_config.expected_status_code_ranges) > 0 ? var.spec.monitor_config.expected_status_code_ranges : null

    # Headers sent with every HTTP/HTTPS probe (e.g. a Host header for
    # name-based virtual hosting); endpoints can override per-endpoint.
    dynamic "custom_header" {
      for_each = var.spec.monitor_config.custom_headers
      content {
        name  = custom_header.value.name
        value = custom_header.value.value
      }
    }
  }
}
