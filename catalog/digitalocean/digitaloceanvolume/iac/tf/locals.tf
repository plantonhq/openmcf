locals {
  # The proto enum value names ARE the strings the provider expects
  # (ext4/xfs); "unformatted" (the enum zero value, which may also arrive as
  # an omitted empty string) means "do not format" and must be sent as null.
  filesystem_type = contains(["ext4", "xfs"], var.spec.filesystem_type) ? var.spec.filesystem_type : null

  # The label only means anything when a filesystem is being formatted.
  filesystem_label = local.filesystem_type != null && var.spec.initial_filesystem_label != "" ? var.spec.initial_filesystem_label : null

  # Standard Planton labels rendered as DigitalOcean "key:value" tags — the
  # exact set and key spelling the Pulumi module applies, so both
  # provisioners tag identically.
  planton_tags = concat(
    [
      "planton-ai_resource:true",
      "planton-ai_name:${var.metadata.name}",
      "planton-ai_kind:DigitalOceanVolume",
    ],
    try(var.metadata.org, "") != "" && var.metadata.org != null ? ["planton-ai_organization:${var.metadata.org}"] : [],
    try(var.metadata.env, "") != "" && var.metadata.env != null ? ["planton-ai_environment:${var.metadata.env}"] : [],
    try(var.metadata.id, "") != "" && var.metadata.id != null ? ["planton-ai_id:${var.metadata.id}"] : [],
  )

  tags = distinct(concat(coalesce(var.spec.tags, []), local.planton_tags))
}
