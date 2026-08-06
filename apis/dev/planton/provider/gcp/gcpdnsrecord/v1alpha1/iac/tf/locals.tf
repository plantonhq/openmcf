locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project (ambient credentials decide).
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # DNS record sets carry no labels in GCP — attribution lives on the zone.

  # Cloud DNS rejects a literal 0 TTL only at the resolver level, not the
  # API; pass the value through untouched so the spec default (300) and any
  # explicit value behave identically on both engines.
  ttl_seconds = var.spec.ttl_seconds

  routing_policy = var.spec.routing_policy

  # The routing policy's optional health check: empty string means unset.
  routing_health_check = (
    local.routing_policy != null && local.routing_policy.health_check != ""
  ) ? local.routing_policy.health_check : null
}
