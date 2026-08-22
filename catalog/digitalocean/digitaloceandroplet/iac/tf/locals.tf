locals {
  # Region is optional: unset lets DigitalOcean choose a region with
  # available capacity. The unspecified enum name must never be sent as a
  # slug (the enum's value names ARE the region slugs).
  region = (
    try(var.spec.region, "") != "" &&
    var.spec.region != "digital_ocean_region_unspecified"
  ) ? var.spec.region : null

  # Optional VPC UUID. References are resolved to the literal UUID before
  # the module runs, so the field arrives as a plain string. Unset means
  # the droplet lands in the region's default VPC.
  vpc_uuid = try(var.spec.vpc, "") != "" ? var.spec.vpc : null

  # SSH keys are create-only; an empty list is sent as null so a keyless
  # droplet never diffs against an empty set.
  ssh_keys = length(coalesce(var.spec.ssh_keys, [])) > 0 ? var.spec.ssh_keys : null

  # Flattened StringValueOrRef list: each entry is already the volume UUID.
  volume_ids = compact([
    for vol in coalesce(var.spec.volume_ids, []) : vol
  ])

  # Cloud-init user data is create-only and hash-stored by DigitalOcean;
  # empty is sent as null, never as an empty string.
  user_data = try(var.spec.user_data, "") != "" ? var.spec.user_data : null

  gpu_partition_mode = try(var.spec.gpu_partition_mode, "") != "" ? var.spec.gpu_partition_mode : null

  # Standard Planton labels rendered as DigitalOcean "key:value" tags —
  # the exact set and key spelling the Pulumi module applies, so both
  # provisioners tag identically.
  planton_tags = concat(
    [
      "planton-ai_resource:true",
      "planton-ai_name:${var.metadata.name}",
      "planton-ai_kind:DigitalOceanDroplet",
    ],
    try(var.metadata.org, "") != "" && var.metadata.org != null ? ["planton-ai_organization:${var.metadata.org}"] : [],
    try(var.metadata.env, "") != "" && var.metadata.env != null ? ["planton-ai_environment:${var.metadata.env}"] : [],
    try(var.metadata.id, "") != "" && var.metadata.id != null ? ["planton-ai_id:${var.metadata.id}"] : [],
  )

  tags = distinct(concat(coalesce(var.spec.tags, []), local.planton_tags))
}
