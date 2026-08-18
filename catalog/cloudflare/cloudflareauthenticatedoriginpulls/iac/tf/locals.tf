locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-authenticated-origin-pulls")

  # Labels
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))

  # One provider resource per association row, keyed by hostname (unique per
  # zone). The provider hard-fails any config list that does not hold exactly
  # one element -- fan-out is the only correct shape.
  hostname_associations = { for a in var.spec.hostname_associations : a.hostname => a }
}
