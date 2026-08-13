# Create one Traffic Manager endpoint -- a destination the referenced
# profile steers traffic to.
#
# The endpoint type is whichever variant block the spec carries
# (validation guarantees exactly one), so exactly one of the count-gated
# resources below materializes. Shared fields (weight, priority,
# enabled, geo/subnet claims, probe headers) feed whichever resource is
# created; which of them MATTER depends on the profile's routing method
# (Azure evaluates them there).
#
# Endpoints carry no ARM tags on any engine. Name, profile, type, and
# subnet claims are fixed at creation; everything else updates in place.
# Priority is sent only when set -- unset lets Azure assign the next
# free value in creation order.

# A public Azure resource as the destination: Azure reads the target's
# address (and its region, for Performance routing) from the resource
# itself, so there is no endpoint_location argument on this type.
resource "azurerm_traffic_manager_azure_endpoint" "main" {
  count = var.spec.azure != null ? 1 : 0

  name                 = var.spec.name
  profile_id           = var.spec.profile_id
  target_resource_id   = var.spec.azure.target_resource_id
  weight               = local.weight
  priority             = var.spec.priority
  enabled              = local.enabled
  always_serve_enabled = coalesce(var.spec.azure.always_serve_enabled, false)
  geo_mappings         = local.geo_mappings

  dynamic "subnet" {
    for_each = var.spec.subnets
    content {
      first = subnet.value.first
      last  = subnet.value.last != "" ? subnet.value.last : null
      scope = subnet.value.scope
    }
  }

  dynamic "custom_header" {
    for_each = var.spec.custom_headers
    content {
      name  = custom_header.value.name
      value = custom_header.value.value
    }
  }
}

# A DNS name or IP address as the destination. endpoint_location is
# REQUIRED by the service when the profile routes by Performance
# (external targets carry no discoverable region) -- enforced
# apply-time, the shape only the service can check.
resource "azurerm_traffic_manager_external_endpoint" "main" {
  count = var.spec.external != null ? 1 : 0

  name                 = var.spec.name
  profile_id           = var.spec.profile_id
  target               = var.spec.external.target
  endpoint_location    = var.spec.external.endpoint_location != "" ? var.spec.external.endpoint_location : null
  weight               = local.weight
  priority             = var.spec.priority
  enabled              = local.enabled
  always_serve_enabled = coalesce(var.spec.external.always_serve_enabled, false)
  geo_mappings         = local.geo_mappings

  dynamic "subnet" {
    for_each = var.spec.subnets
    content {
      first = subnet.value.first
      last  = subnet.value.last != "" ? subnet.value.last : null
      scope = subnet.value.scope
    }
  }

  dynamic "custom_header" {
    for_each = var.spec.custom_headers
    content {
      name  = custom_header.value.name
      value = custom_header.value.value
    }
  }
}

# Another Traffic Manager profile as the destination, composing routing
# methods into trees. This type carries no always-serve switch (the
# provider exposes none -- child health IS the point of nesting); the
# IPv4/IPv6 health floors pass through only when positive, mirroring
# the provider's own send-when-set behavior.
resource "azurerm_traffic_manager_nested_endpoint" "main" {
  count = var.spec.nested != null ? 1 : 0

  name                                  = var.spec.name
  profile_id                            = var.spec.profile_id
  target_resource_id                    = var.spec.nested.target_profile_id
  minimum_child_endpoints               = var.spec.nested.minimum_child_endpoints
  minimum_required_child_endpoints_ipv4 = var.spec.nested.minimum_required_child_endpoints_ipv4
  minimum_required_child_endpoints_ipv6 = var.spec.nested.minimum_required_child_endpoints_ipv6
  endpoint_location                     = var.spec.nested.endpoint_location != "" ? var.spec.nested.endpoint_location : null
  weight                                = local.weight
  priority                              = var.spec.priority
  enabled                               = local.enabled
  geo_mappings                          = local.geo_mappings

  dynamic "subnet" {
    for_each = var.spec.subnets
    content {
      first = subnet.value.first
      last  = subnet.value.last != "" ? subnet.value.last : null
      scope = subnet.value.scope
    }
  }

  dynamic "custom_header" {
    for_each = var.spec.custom_headers
    content {
      name  = custom_header.value.name
      value = custom_header.value.value
    }
  }
}
