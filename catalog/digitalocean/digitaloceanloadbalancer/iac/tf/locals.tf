locals {
  # Optional VPC UUID. References are resolved to the literal UUID before
  # the module runs, so the field arrives as a plain string. Unset means
  # DigitalOcean places a regional balancer in the region's default VPC
  # (GLOBAL balancers take no VPC at all).
  vpc_uuid = try(var.spec.vpc, "") != "" ? var.spec.vpc : null

  # Region is required for REGIONAL / REGIONAL_NETWORK and forbidden for
  # GLOBAL. The unspecified enum name must never be sent as a slug.
  region = (
    try(var.spec.region, "") != "" &&
    var.spec.region != "digital_ocean_region_unspecified"
  ) ? var.spec.region : null

  type          = try(var.spec.type, "") != "" ? var.spec.type : null
  size          = try(var.spec.size, "") != "" ? var.spec.size : null
  size_unit     = try(var.spec.size_unit, 0) > 0 ? var.spec.size_unit : null
  project_id    = try(var.spec.project_id, "") != "" ? var.spec.project_id : null
  subnet_uuid   = try(var.spec.subnet_uuid, "") != "" ? var.spec.subnet_uuid : null
  ip            = try(var.spec.ip, "") != "" ? var.spec.ip : null
  network       = try(var.spec.network, "") != "" ? var.spec.network : null
  network_stack = try(var.spec.network_stack, "") != "" ? var.spec.network_stack : null
  tls_cipher_policy = try(var.spec.tls_cipher_policy, "") != "" ? var.spec.tls_cipher_policy : null
  droplet_tag   = try(var.spec.droplet_tag, "") != "" ? var.spec.droplet_tag : null

  # 0 (proto3 unset) defers to DigitalOcean's 60-second default.
  http_idle_timeout_seconds = try(var.spec.http_idle_timeout_seconds, 0) > 0 ? var.spec.http_idle_timeout_seconds : null

  # Flattened StringValueOrRef list: each entry is already the numeric
  # Droplet ID string. Empty means tag-managed membership (or none).
  droplet_ids = [
    for id in coalesce(var.spec.droplet_ids, []) : tonumber(id)
  ]
  droplet_ids_or_null = length(local.droplet_ids) > 0 ? local.droplet_ids : null

  target_load_balancer_ids = length(coalesce(var.spec.target_load_balancer_ids, [])) > 0 ? var.spec.target_load_balancer_ids : null
}
