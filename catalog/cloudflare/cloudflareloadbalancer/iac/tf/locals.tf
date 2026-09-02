locals {
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-load-balancer")

  # Zone (StringValueOrRef flattened to a plain string by the converter).
  zone_id = try(var.spec.zone_id, "")

  # Enums flatten to their string names. session_affinity may be omitted when
  # "none": the provider's schema default IS "none" and the API echoes it, so
  # the omission stays convergent (measured live). steering_policy is DIFFERENT:
  # the provider defaults it to "" while the API canonicalizes and stores
  # "off"/"geo", so an omitted value re-plans after a blind import ("off" in
  # state vs "" in config -- measured live 2026-08-28). ALWAYS send the
  # canonical value, mirroring the API's documented "" mapping: geo when any
  # geo-pool map is set, otherwise off.
  session_affinity = (
    try(var.spec.session_affinity, "") == "" || var.spec.session_affinity == "none"
  ) ? null : var.spec.session_affinity
  steering_policy = (
    try(var.spec.steering_policy, "") == "" || var.spec.steering_policy == "off"
    ) ? (
    (length(local.region_pools) > 0 || length(local.country_pools) > 0 || length(local.pop_pools) > 0) ? "geo" : "off"
  ) : var.spec.steering_policy

  # Geo-pool maps: list-of-{code,pool_ids} -> provider's { code => pool_ids }.
  region_pools  = { for e in try(var.spec.region_pools, []) : e.code => e.pool_ids }
  country_pools = { for e in try(var.spec.country_pools, []) : e.code => e.pool_ids }
  pop_pools     = { for e in try(var.spec.pop_pools, []) : e.code => e.pool_ids }

  # require_all_headers is sent ONLY when true: the API never echoes a false
  # back (GET returns null), so an always-sent false re-plans on every
  # refresh-inclusive plan forever (measured live 2026-08-28 -- the class from
  # the harness's idempotency gate).
  saa = try(var.spec.session_affinity_attributes, null) == null ? null : {
    drain_duration         = try(var.spec.session_affinity_attributes.drain_duration, 0) > 0 ? var.spec.session_affinity_attributes.drain_duration : null
    headers                = length(try(var.spec.session_affinity_attributes.headers, [])) > 0 ? var.spec.session_affinity_attributes.headers : null
    require_all_headers    = try(var.spec.session_affinity_attributes.require_all_headers, false) ? true : null
    samesite               = try(var.spec.session_affinity_attributes.samesite, "") != "" ? var.spec.session_affinity_attributes.samesite : null
    secure                 = try(var.spec.session_affinity_attributes.secure, "") != "" ? var.spec.session_affinity_attributes.secure : null
    zero_downtime_failover = try(var.spec.session_affinity_attributes.zero_downtime_failover, "") != "" ? var.spec.session_affinity_attributes.zero_downtime_failover : null
  }

  adaptive_routing = try(var.spec.adaptive_routing, null) == null ? null : {
    failover_across_pools = try(var.spec.adaptive_routing.failover_across_pools, false)
  }

  location_strategy = (
    try(var.spec.location_strategy, null) == null ||
    (try(var.spec.location_strategy.mode, "") == "" && try(var.spec.location_strategy.prefer_ecs, "") == "")
    ) ? null : {
    mode       = try(var.spec.location_strategy.mode, "") != "" ? var.spec.location_strategy.mode : null
    prefer_ecs = try(var.spec.location_strategy.prefer_ecs, "") != "" ? var.spec.location_strategy.prefer_ecs : null
  }

  random_steering = try(var.spec.random_steering, null) == null ? null : {
    default_weight = try(var.spec.random_steering.default_weight, 0) > 0 ? var.spec.random_steering.default_weight : null
    pool_weights   = length(try(var.spec.random_steering.pool_weights, {})) > 0 ? var.spec.random_steering.pool_weights : null
  }

  # Traffic rules: the spec's typed list -> the provider's rules[] shape.
  # Override semantics need real presence: an unset override inherits the load
  # balancer's setting, while an explicit value -- INCLUDING "none"/"off"/0 --
  # overrides it. The spec carries presence on priority, session_affinity, and
  # steering_policy (proto optional), so those pass through as-is (null when
  # unset). disabled is sent only when true (false is the provider default).
  # terminates is sent as TRUE whenever the rule carries a fixed_response:
  # Cloudflare auto-marks fixed-response rules terminating and echoes
  # terminates=true on every read, so a null send re-plans the rules list
  # after a blind import (measured live 2026-08-28); an explicit false is
  # never sent -- it would fight the same echo.
  rules = length(try(var.spec.rules, [])) == 0 ? null : [
    for r in var.spec.rules : {
      name       = try(r.name, "") != "" ? r.name : null
      condition  = try(r.condition, "") != "" ? r.condition : null
      priority   = try(r.priority, null)
      disabled   = try(r.disabled, false) ? true : null
      terminates = (try(r.terminates, false) || try(r.fixed_response, null) != null) ? true : null
      fixed_response = try(r.fixed_response, null) == null ? null : {
        content_type = try(r.fixed_response.content_type, "") != "" ? r.fixed_response.content_type : null
        location     = try(r.fixed_response.location, "") != "" ? r.fixed_response.location : null
        message_body = try(r.fixed_response.message_body, "") != "" ? r.fixed_response.message_body : null
        status_code  = try(r.fixed_response.status_code, 0) != 0 ? r.fixed_response.status_code : null
      }
      overrides = try(r.overrides, null) == null ? null : {
        adaptive_routing = try(r.overrides.adaptive_routing, null) == null ? null : {
          failover_across_pools = try(r.overrides.adaptive_routing.failover_across_pools, false)
        }
        country_pools = length(try(r.overrides.country_pools, [])) > 0 ? { for e in r.overrides.country_pools : e.code => e.pool_ids } : null
        default_pools = length(try(r.overrides.default_pools, [])) > 0 ? r.overrides.default_pools : null
        fallback_pool = try(r.overrides.fallback_pool, "") != "" ? r.overrides.fallback_pool : null
        location_strategy = try(r.overrides.location_strategy, null) == null ? null : {
          mode       = try(r.overrides.location_strategy.mode, "") != "" ? r.overrides.location_strategy.mode : null
          prefer_ecs = try(r.overrides.location_strategy.prefer_ecs, "") != "" ? r.overrides.location_strategy.prefer_ecs : null
        }
        pop_pools = length(try(r.overrides.pop_pools, [])) > 0 ? { for e in r.overrides.pop_pools : e.code => e.pool_ids } : null
        random_steering = try(r.overrides.random_steering, null) == null ? null : {
          default_weight = try(r.overrides.random_steering.default_weight, 0) > 0 ? r.overrides.random_steering.default_weight : null
          pool_weights   = length(try(r.overrides.random_steering.pool_weights, {})) > 0 ? r.overrides.random_steering.pool_weights : null
        }
        region_pools     = length(try(r.overrides.region_pools, [])) > 0 ? { for e in r.overrides.region_pools : e.code => e.pool_ids } : null
        session_affinity = try(r.overrides.session_affinity, null)
        session_affinity_attributes = try(r.overrides.session_affinity_attributes, null) == null ? null : {
          drain_duration         = try(r.overrides.session_affinity_attributes.drain_duration, 0) > 0 ? r.overrides.session_affinity_attributes.drain_duration : null
          headers                = length(try(r.overrides.session_affinity_attributes.headers, [])) > 0 ? r.overrides.session_affinity_attributes.headers : null
          require_all_headers    = try(r.overrides.session_affinity_attributes.require_all_headers, false) ? true : null
          samesite               = try(r.overrides.session_affinity_attributes.samesite, "") != "" ? r.overrides.session_affinity_attributes.samesite : null
          secure                 = try(r.overrides.session_affinity_attributes.secure, "") != "" ? r.overrides.session_affinity_attributes.secure : null
          zero_downtime_failover = try(r.overrides.session_affinity_attributes.zero_downtime_failover, "") != "" ? r.overrides.session_affinity_attributes.zero_downtime_failover : null
        }
        session_affinity_ttl = try(r.overrides.session_affinity_ttl, 0) > 0 ? r.overrides.session_affinity_ttl : null
        steering_policy      = try(r.overrides.steering_policy, null)
        ttl                  = try(r.overrides.ttl, 0) > 0 ? r.overrides.ttl : null
      }
    }
  ]
}
