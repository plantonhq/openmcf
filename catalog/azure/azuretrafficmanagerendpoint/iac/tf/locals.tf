# Traffic Manager endpoints carry NO ARM tags on any engine (the
# provider exposes none) -- the platform's derived tags land on the
# owning profile instead, so this module derives no tag map.
locals {
  # Shared-field defaults, sent explicitly so plans stay deterministic:
  # weight 1 (the provider's default); enabled true. Priority is
  # deliberately NOT defaulted -- unset lets Azure assign the next free
  # value in creation order (the service owns that default).
  weight  = coalesce(var.spec.weight, 1)
  enabled = coalesce(var.spec.enabled, true)

  # Empty collections pass as null so the provider never sees an
  # empty-but-present argument.
  geo_mappings = length(var.spec.geo_mappings) > 0 ? var.spec.geo_mappings : null
}
