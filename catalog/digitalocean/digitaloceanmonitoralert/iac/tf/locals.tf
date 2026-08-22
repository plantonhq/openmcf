locals {
  # The provider takes one untyped entities list; the spec splits it into
  # three typed reference lists (droplets, load balancers, database
  # clusters). References arrive flattened as plain id strings -- the
  # Planton orchestrator resolves valueFrom references before Terraform
  # runs -- so the module merges them back into the provider's shape. Spec
  # validation already guarantees only the metric family's own list is
  # populated.
  entities = concat(
    var.spec.droplet_ids,
    var.spec.load_balancer_ids,
    var.spec.database_cluster_ids,
  )
}
