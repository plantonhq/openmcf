locals {
  # A record is standard (values + ttl) or alias (alias_target) — never both
  # (CEL-enforced). The generator flattens StringValueOrRef fields to plain
  # strings, so alias presence is the block itself.
  is_alias = var.spec.alias_target != null

  # TTL and values only exist on standard records; alias records take the
  # target's TTL (the provider rejects ttl alongside alias).
  ttl     = local.is_alias ? null : var.spec.ttl
  records = local.is_alias ? null : var.spec.values

  # Null-when-unset scalars so absent optional inputs stay absent.
  set_identifier  = var.spec.set_identifier != "" ? var.spec.set_identifier : null
  health_check_id = var.spec.health_check_id != "" ? var.spec.health_check_id : null

  # Exactly one routing policy per record (proto oneof); records with the
  # same name/type but different set_identifier values form a routing group.
  routing_policy   = var.spec.routing_policy
  has_weighted     = local.routing_policy != null ? local.routing_policy.weighted != null : false
  has_latency      = local.routing_policy != null ? local.routing_policy.latency != null : false
  has_failover     = local.routing_policy != null ? local.routing_policy.failover != null : false
  has_geolocation  = local.routing_policy != null ? local.routing_policy.geolocation != null : false
  has_geoproximity = local.routing_policy != null ? local.routing_policy.geoproximity != null : false
  has_cidr         = local.routing_policy != null ? local.routing_policy.cidr != null : false
  is_multivalue    = local.routing_policy != null ? local.routing_policy.multivalue_answer != null : false
}
