locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-zero-trust-gateway-settings")

  # Labels
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))
}
