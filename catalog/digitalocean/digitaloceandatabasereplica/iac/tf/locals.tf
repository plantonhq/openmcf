locals {
  # Spec enum value names are exactly the DigitalOcean region slugs
  # (nyc3, sfo3, ...), so they pass through unchanged.
  region_slug = var.spec.region

  # Optional VPC UUID. References are resolved to the literal UUID before
  # the module runs, so the field arrives as a plain string. Null when
  # unset: private_network_uuid is Optional+Computed upstream, so the
  # region's default VPC applies without drift.
  vpc_uuid = try(var.spec.vpc, "") != "" ? var.spec.vpc : null

  # The provider's storage_size_mib is a string holding a bare MiB count;
  # the spec carries the number. Unset means "use the size slug's default
  # storage" (Optional+Computed upstream -- drift-safe to omit).
  storage_size_mib = var.spec.storage_size_mib != null && var.spec.storage_size_mib > 0 ? tostring(var.spec.storage_size_mib) : null

  # Standard Planton labels rendered as DigitalOcean "key:value" tags —
  # the exact set and key spelling the Pulumi module applies, so both
  # provisioners tag identically. NOTE: replica tags are CREATE-ONLY
  # upstream -- changing the final set REPLACES the replica.
  planton_tags = concat(
    [
      "planton-ai_resource:true",
      "planton-ai_name:${var.metadata.name}",
      "planton-ai_kind:DigitalOceanDatabaseReplica",
    ],
    try(var.metadata.org, "") != "" && var.metadata.org != null ? ["planton-ai_organization:${var.metadata.org}"] : [],
    try(var.metadata.env, "") != "" && var.metadata.env != null ? ["planton-ai_environment:${var.metadata.env}"] : [],
    try(var.metadata.id, "") != "" && var.metadata.id != null ? ["planton-ai_id:${var.metadata.id}"] : [],
  )

  tags = distinct(concat(coalesce(var.spec.tags, []), local.planton_tags))
}
