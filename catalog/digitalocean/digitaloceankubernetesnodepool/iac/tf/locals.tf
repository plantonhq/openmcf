locals {
  # Standard Planton labels rendered as DigitalOcean "key:value" tags —
  # the exact set and key spelling the Pulumi module applies, so both
  # provisioners tag identically.
  planton_tags = concat(
    [
      "planton-ai_resource:true",
      "planton-ai_name:${var.metadata.name}",
      "planton-ai_kind:DigitalOceanKubernetesNodePool",
    ],
    try(var.metadata.org, "") != "" && var.metadata.org != null ? ["planton-ai_organization:${var.metadata.org}"] : [],
    try(var.metadata.env, "") != "" && var.metadata.env != null ? ["planton-ai_environment:${var.metadata.env}"] : [],
    try(var.metadata.id, "") != "" && var.metadata.id != null ? ["planton-ai_id:${var.metadata.id}"] : [],
  )

  tags = distinct(concat(coalesce(var.spec.tags, []), local.planton_tags))

  # The same Planton identity as Kubernetes node labels, with user labels
  # winning on key collisions — the exact map the Pulumi module applies.
  planton_labels = merge(
    {
      "planton-ai_resource" = "true"
      "planton-ai_name"     = var.metadata.name
      "planton-ai_kind"     = "DigitalOceanKubernetesNodePool"
    },
    try(var.metadata.org, "") != "" && var.metadata.org != null ? { "planton-ai_organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" && var.metadata.env != null ? { "planton-ai_environment" = var.metadata.env } : {},
    try(var.metadata.id, "") != "" && var.metadata.id != null ? { "planton-ai_id" = var.metadata.id } : {},
  )

  node_labels = merge(local.planton_labels, coalesce(var.spec.labels, {}))
}
