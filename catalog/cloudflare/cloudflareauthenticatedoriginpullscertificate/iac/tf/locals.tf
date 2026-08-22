locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-authenticated-origin-pulls-certificate")

  # Labels
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))

  # The scope decides which of Cloudflare's two upload surfaces receives the
  # certificate. Exactly one of the two resources below is created.
  is_zone_scope = coalesce(var.spec.scope, "zone") == "zone"
}
