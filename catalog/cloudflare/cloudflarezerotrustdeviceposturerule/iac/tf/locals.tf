locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-zero-trust-device-posture-rule")

  # Labels
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))
}
