locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-zt-service-token")

  # Labels
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))

  # Scope: exactly one of account_id or zone_id is set (enforced by the spec).
  account_id = try(var.spec.account_id, "")
  zone_id    = try(var.spec.zone_id, "")
}
