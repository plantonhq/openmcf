locals {
  # CloudWatch dashboards are untaggable at AWS (the resource has no tags
  # argument), so this module carries no tag map - the one deliberate
  # absence against the catalog's tag convention (mirrored in the Pulumi
  # module).

  # The dashboard body arrives from the tfvars layer as a nested object;
  # the provider wants the document as a JSON string. AWS normalizes the
  # JSON server-side and the provider diffs it semantically, so key order
  # never causes drift.
  dashboard_body = jsonencode(var.spec.dashboard_body)
}
