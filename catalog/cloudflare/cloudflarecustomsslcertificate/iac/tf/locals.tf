locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-custom-ssl-certificate")

  # Labels
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))

  # Empty strings mean "not set" for plain proto3 string fields -- drop them
  # rather than sending empty values the API would reject or record.
  policy        = var.spec.policy != "" ? var.spec.policy : null
  custom_csr_id = var.spec.custom_csr_id != "" ? var.spec.custom_csr_id : null

  # The nested geo restriction is sent only when the label is present.
  geo_restrictions = var.spec.geo_restrictions != null && try(var.spec.geo_restrictions.label, null) != null ? {
    label = var.spec.geo_restrictions.label
  } : null
}
