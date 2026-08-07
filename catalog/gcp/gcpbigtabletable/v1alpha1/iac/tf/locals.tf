locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Table name defaults to metadata.name when table_name is omitted.
  table_name = var.spec.table_name != "" ? var.spec.table_name : var.metadata.name

  # Families that carry a GC policy each get their own
  # google_bigtable_gc_policy resource — the API's own one-per-family
  # granularity, folded into this kind because a GC policy has no
  # independent life apart from its family.
  gc_policies = {
    for cf in var.spec.column_families :
    cf.family => cf.gc_policy
    if cf.gc_policy != null
  }
}
