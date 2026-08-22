locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-ip-access-rule")

  # Labels
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))
}
