locals {
  # The resource server's display name comes from the spec; metadata.name is
  # the graph identity. No aws_tags map here: the aws_cognito_resource_server
  # resource is not taggable (identity tagging lives on the pool).
  resource_name = var.metadata.name

  # The fully-qualified scope identifiers ("{identifier}/{scope_name}") are
  # computed from the spec rather than read back from the provider, so the
  # export order matches the spec order deterministically on both engines.
  scope_identifiers = [for s in var.spec.scopes : "${var.spec.identifier}/${s.scope_name}"]
}
