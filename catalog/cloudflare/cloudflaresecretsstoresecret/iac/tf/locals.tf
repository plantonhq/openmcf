locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-secrets-store-secret")

  # Labels
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))

  # Empty string means "not set" for the optional comment -- drop it rather
  # than sending an empty value.
  comment = var.spec.comment != "" ? var.spec.comment : null
}
