locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-zero-trust-list")

  # Labels
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))

  # Items: Cloudflare treats them as a SET (order-insignificant). Empty
  # descriptions are dropped rather than sent as "".
  items = length(var.spec.items) > 0 ? [
    for item in var.spec.items : {
      value       = item.value
      description = item.description != "" ? item.description : null
    }
  ] : null
}
