locals {
  # Determine if this is an alias record
  is_alias = var.spec.alias_target != null

  # TTL is only applicable for non-alias records
  ttl = local.is_alias ? null : var.spec.ttl

  # Records are only for non-alias records
  records = local.is_alias ? null : var.spec.values

  # Determine routing policy type
  has_weighted    = var.spec.routing_policy != null ? var.spec.routing_policy.weighted != null : false
  has_latency     = var.spec.routing_policy != null ? var.spec.routing_policy.latency != null : false
  has_failover    = var.spec.routing_policy != null ? var.spec.routing_policy.failover != null : false
  has_geolocation = var.spec.routing_policy != null ? var.spec.routing_policy.geolocation != null : false
}
