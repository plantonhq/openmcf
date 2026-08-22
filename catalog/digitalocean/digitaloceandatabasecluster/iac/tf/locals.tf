locals {
  # Spec enum value names are exactly the DigitalOcean API slugs
  # (pg, mysql, redis, mongodb, kafka, opensearch, valkey / nyc3, sfo3, ...),
  # so they pass through unchanged.
  engine_slug = var.spec.engine
  region_slug = var.spec.region

  # Optional VPC UUID. References are resolved to the literal UUID before
  # the module runs, so the field arrives as a plain string.
  vpc_uuid = try(var.spec.vpc, "") != "" ? var.spec.vpc : null

  # Optional project placement.
  project_id = try(var.spec.project_id, "") != "" ? var.spec.project_id : null

  # The provider's storage_size_mib is a string holding a bare MiB count;
  # the spec carries GiB for ergonomics. Unset means "use the size slug's
  # default storage".
  storage_size_mib = var.spec.storage_gib != null && var.spec.storage_gib > 0 ? tostring(var.spec.storage_gib * 1024) : null

  # Engine-conditional arguments: null when unset so the provider's own
  # engine checks never fire for the engines they don't apply to.
  eviction_policy = try(var.spec.eviction_policy, "") != "" ? var.spec.eviction_policy : null
  sql_mode        = try(var.spec.sql_mode, "") != "" ? var.spec.sql_mode : null

  # Standard Planton labels rendered as DigitalOcean "key:value" tags —
  # the exact set and key spelling the Pulumi module applies, so both
  # provisioners tag identically.
  planton_tags = concat(
    [
      "planton-ai_resource:true",
      "planton-ai_name:${var.metadata.name}",
      "planton-ai_kind:DigitalOceanDatabaseCluster",
    ],
    try(var.metadata.org, "") != "" && var.metadata.org != null ? ["planton-ai_organization:${var.metadata.org}"] : [],
    try(var.metadata.env, "") != "" && var.metadata.env != null ? ["planton-ai_environment:${var.metadata.env}"] : [],
    try(var.metadata.id, "") != "" && var.metadata.id != null ? ["planton-ai_id:${var.metadata.id}"] : [],
  )

  tags = distinct(concat(coalesce(var.spec.tags, []), local.planton_tags))
}
