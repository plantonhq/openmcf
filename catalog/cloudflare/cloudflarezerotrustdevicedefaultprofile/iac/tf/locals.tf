locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-zero-trust-device-default-profile")

  # Labels
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))
}
