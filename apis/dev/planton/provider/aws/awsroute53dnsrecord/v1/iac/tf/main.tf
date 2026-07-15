resource "aws_route53_record" "this" {
  # The record's AWS identity is (zone, name, type, set_identifier) —
  # changing zone_id or name replaces the record.
  zone_id = var.spec.zone_id
  name    = var.spec.name
  type    = var.spec.type

  # Standard records carry values + ttl; alias records carry neither (the
  # alias target's TTL applies).
  ttl     = local.ttl
  records = local.records

  set_identifier = local.set_identifier

  # Health check gating this record's answers — valid with any non-simple
  # routing policy, most commonly the failover PRIMARY.
  health_check_id = local.health_check_id

  # Overwrite an existing record set with the same name and type instead of
  # failing — the adoption path for records created outside the graph.
  allow_overwrite = var.spec.allow_overwrite

  # Alias records point at an AWS resource's DNS name + that service's own
  # hosted zone ID (NOT this record's zone).
  dynamic "alias" {
    for_each = local.is_alias ? [var.spec.alias_target] : []
    content {
      name                   = alias.value.dns_name
      zone_id                = alias.value.zone_id
      evaluate_target_health = alias.value.evaluate_target_health
    }
  }

  # --- Routing policies (exactly one; mutually exclusive in the provider,
  # --- a proto oneof in the spec) -------------------------------------------

  dynamic "weighted_routing_policy" {
    for_each = local.has_weighted ? [local.routing_policy.weighted] : []
    content {
      weight = weighted_routing_policy.value.weight
    }
  }

  dynamic "latency_routing_policy" {
    for_each = local.has_latency ? [local.routing_policy.latency] : []
    content {
      region = latency_routing_policy.value.region
    }
  }

  dynamic "failover_routing_policy" {
    for_each = local.has_failover ? [local.routing_policy.failover] : []
    content {
      type = failover_routing_policy.value.failover_type
    }
  }

  dynamic "geolocation_routing_policy" {
    for_each = local.has_geolocation ? [local.routing_policy.geolocation] : []
    content {
      continent   = geolocation_routing_policy.value.continent != "" ? geolocation_routing_policy.value.continent : null
      country     = geolocation_routing_policy.value.country != "" ? geolocation_routing_policy.value.country : null
      subdivision = geolocation_routing_policy.value.subdivision != "" ? geolocation_routing_policy.value.subdivision : null
    }
  }

  # Geoproximity: exactly one location determinant (region, coordinates, or
  # Local Zone group — CEL-enforced) plus the optional bias dial.
  dynamic "geoproximity_routing_policy" {
    for_each = local.has_geoproximity ? [local.routing_policy.geoproximity] : []
    content {
      aws_region       = geoproximity_routing_policy.value.aws_region != "" ? geoproximity_routing_policy.value.aws_region : null
      local_zone_group = geoproximity_routing_policy.value.local_zone_group != "" ? geoproximity_routing_policy.value.local_zone_group : null
      bias             = geoproximity_routing_policy.value.bias != 0 ? geoproximity_routing_policy.value.bias : null

      dynamic "coordinates" {
        for_each = geoproximity_routing_policy.value.coordinates != null ? [geoproximity_routing_policy.value.coordinates] : []
        content {
          latitude  = coordinates.value.latitude
          longitude = coordinates.value.longitude
        }
      }
    }
  }

  dynamic "cidr_routing_policy" {
    for_each = local.has_cidr ? [local.routing_policy.cidr] : []
    content {
      collection_id = cidr_routing_policy.value.collection_id
      location_name = cidr_routing_policy.value.location_name
    }
  }

  # Multivalue answer is a bare flag in the provider, not a block; the spec
  # models it as an empty policy message so the oneof stays uniform.
  multivalue_answer_routing_policy = local.is_multivalue ? true : null
}
