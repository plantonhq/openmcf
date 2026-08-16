locals {
  # Create-only network placement. The provider validates the CIDRs and
  # rejects empty strings, so unset must arrive as null, not "".
  cluster_subnet     = try(var.spec.cluster_subnet, "") != "" ? var.spec.cluster_subnet : null
  service_subnet     = try(var.spec.service_subnet, "") != "" ? var.spec.service_subnet : null
  worker_subnet_uuid = try(var.spec.worker_subnet_uuid, "") != "" ? var.spec.worker_subnet_uuid : null

  # 0 (proto3 unset) means DigitalOcean's 7-day default; keep it out of
  # state rather than pinning an explicit zero.
  kubeconfig_expire_seconds = var.spec.kubeconfig_expire_seconds > 0 ? var.spec.kubeconfig_expire_seconds : null

  # Standard Planton labels rendered as DigitalOcean "key:value" tags —
  # the exact set and key spelling the Pulumi module applies, so both
  # provisioners tag identically.
  planton_tags = concat(
    [
      "planton-ai_resource:true",
      "planton-ai_name:${var.metadata.name}",
      "planton-ai_kind:DigitalOceanKubernetesCluster",
    ],
    try(var.metadata.org, "") != "" && var.metadata.org != null ? ["planton-ai_organization:${var.metadata.org}"] : [],
    try(var.metadata.env, "") != "" && var.metadata.env != null ? ["planton-ai_environment:${var.metadata.env}"] : [],
    try(var.metadata.id, "") != "" && var.metadata.id != null ? ["planton-ai_id:${var.metadata.id}"] : [],
  )

  tags = distinct(concat(coalesce(var.spec.tags, []), local.planton_tags))

  # The same Planton identity as Kubernetes node labels on the default
  # pool, with user labels winning on key collisions — the exact map the
  # Pulumi module applies.
  planton_labels = merge(
    {
      "planton-ai_resource" = "true"
      "planton-ai_name"     = var.metadata.name
      "planton-ai_kind"     = "DigitalOceanKubernetesCluster"
    },
    try(var.metadata.org, "") != "" && var.metadata.org != null ? { "planton-ai_organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" && var.metadata.env != null ? { "planton-ai_environment" = var.metadata.env } : {},
    try(var.metadata.id, "") != "" && var.metadata.id != null ? { "planton-ai_id" = var.metadata.id } : {},
  )

  default_node_pool_labels = merge(local.planton_labels, coalesce(var.spec.default_node_pool.labels, {}))
}
