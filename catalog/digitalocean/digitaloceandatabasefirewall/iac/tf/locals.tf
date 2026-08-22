locals {
  # The provider takes one polymorphic rule set of {type, value} rows; the
  # spec carries one TYPED list per source kind so a value can never be
  # paired with the wrong type. Fan the lists back out here. References
  # (droplets, clusters, apps) are resolved to literal ids before the
  # module runs, so every list arrives as plain strings.
  rules = concat(
    [for v in coalesce(var.spec.ip_rules, []) : { type = "ip_addr", value = v }],
    [for v in coalesce(var.spec.droplet_ids, []) : { type = "droplet", value = v }],
    [for v in coalesce(var.spec.kubernetes_cluster_ids, []) : { type = "k8s", value = v }],
    [for v in coalesce(var.spec.app_ids, []) : { type = "app", value = v }],
    [for v in coalesce(var.spec.tags, []) : { type = "tag", value = v }],
  )
}
