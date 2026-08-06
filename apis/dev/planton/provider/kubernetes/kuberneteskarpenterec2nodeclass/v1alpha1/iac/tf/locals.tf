locals {
  # Planton identity labels — the planton.ai/* convention, identical to the
  # Pulumi module's label set (twin discipline). Conditional entries use the
  # null-prune idiom: heterogeneous conditional merges fail HCL type
  # unification when sibling entries infer as different object types.
  labels = {
    for k, v in {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKarpenterEc2NodeClass"
      "planton.ai/resource-id"   = (var.metadata.id != null && var.metadata.id != "") ? var.metadata.id : null
      "planton.ai/organization"  = (var.metadata.org != null && var.metadata.org != "") ? var.metadata.org : null
      "planton.ai/environment"   = (var.metadata.env != null && var.metadata.env != "") ? var.metadata.env : null
    } : k => v if v != null
  }

  # Cluster-scoped CR: the converter already emits camelCase, null-pruned keys,
  # so the spec is passed through unchanged.
  manifest_spec = var.spec
}
