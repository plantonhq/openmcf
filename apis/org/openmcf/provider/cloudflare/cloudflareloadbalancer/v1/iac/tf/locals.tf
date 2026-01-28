locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-load-balancer")
  
  # Tags/labels
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))

  # Zone ID - resolve from value or ref
  zone_id = coalesce(try(var.spec.zone_id.value, null), try(var.spec.zone_id.ref, null))

  # Health probe path with default
  health_probe_path = coalesce(try(var.spec.health_probe_path, null), "/")

  # Proxied setting with default
  proxied = coalesce(try(var.spec.proxied, null), true)

  # Session affinity mapping (enum to string)
  session_affinity_map = {
    0 = "none"
    1 = "cookie"
  }
  session_affinity = lookup(local.session_affinity_map, try(var.spec.session_affinity, 0), "none")

  # Steering policy mapping (enum to string)
  steering_policy_map = {
    0 = "off"     # STEERING_OFF - active-passive failover
    1 = "geo"     # STEERING_GEO - geographic routing
    2 = "random"  # STEERING_RANDOM - weighted distribution
  }
  steering_policy = lookup(local.steering_policy_map, try(var.spec.steering_policy, 0), "off")

  # Origins list
  origins = try(var.spec.origins, [])
}

