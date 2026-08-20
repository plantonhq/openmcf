locals {
  # CloudWatch Logs resource policies are untaggable at AWS (the resource
  # has no tags argument), so this module carries no tag map - the one
  # deliberate absence against the catalog's tag convention (mirrored in
  # the Pulumi module).

  # The policy document arrives from the tfvars layer as a nested object;
  # the provider wants a JSON string and diffs IAM documents
  # semantically.
  policy_document = jsonencode(var.spec.policy_document)
}
