locals {
  # Planton identity labels — the planton.ai/* convention, identical to the
  # Pulumi module's label set (twin discipline). Conditional entries use the
  # null-prune idiom: heterogeneous conditional merges fail HCL type
  # unification when sibling entries infer as different object types.
  labels = {
    for k, v in {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKarpenterNodePool"
      "planton.ai/resource-id"   = (var.metadata.id != null && var.metadata.id != "") ? var.metadata.id : null
      "planton.ai/organization"  = (var.metadata.org != null && var.metadata.org != "") ? var.metadata.org : null
      "planton.ai/environment"   = (var.metadata.env != null && var.metadata.env != "") ? var.metadata.env : null
    } : k => v if v != null
  }

  # The spec proto deliberately FLATTENS the CRD's template nesting: users
  # write template.{labels,annotations,nodeClassRef,requirements,...} while
  # the CRD wants template.metadata.{labels,annotations} and
  # template.spec.{nodeClassRef,requirements,taints,expireAfter,...} (the
  # NodeClaim template is a full object template). The converter's camelCase
  # spec therefore CANNOT pass through verbatim for this kind — the module
  # rebuilds the nesting, exactly like the Pulumi twin's template builder.
  # Everything OUTSIDE template (disruption, limits, weight, replicas)
  # already matches the CRD shape and passes through unchanged.

  # group/kind fall back to the proto-declared AWS defaults when unset: the
  # CRD requires all three keys non-empty (twin of the Pulumi module's
  # buildNodeClassRef).
  node_class_ref = {
    group = try(coalesce(try(var.spec.template.nodeClassRef.group, null)), "karpenter.k8s.aws")
    kind  = try(coalesce(try(var.spec.template.nodeClassRef.kind, null)), "EC2NodeClass")
    name  = var.spec.template.nodeClassRef.name
  }

  # template.metadata renders only when labels or annotations exist (the
  # CRD tolerates an absent metadata; an empty object is noise).
  template_metadata = {
    for k, v in {
      labels      = try(var.spec.template.labels, null)
      annotations = try(var.spec.template.annotations, null)
    } : k => v if v != null
  }

  # The NodeClaim template spec: converter-emitted keys pass through where
  # present (they are already camelCase and null-pruned); expireAfter and
  # terminationGracePeriod render only when set so the apiserver applies
  # the CRD defaults (expireAfter: "720h") otherwise.
  template_spec = {
    for k, v in {
      nodeClassRef           = local.node_class_ref
      requirements           = try(var.spec.template.requirements, null)
      taints                 = try(var.spec.template.taints, null)
      startupTaints          = try(var.spec.template.startupTaints, null)
      expireAfter            = try(var.spec.template.expireAfter, null)
      terminationGracePeriod = try(var.spec.template.terminationGracePeriod, null)
    } : k => v if v != null
  }

  manifest_spec = merge(
    { for k, v in var.spec : k => v if k != "template" },
    {
      template = {
        for k, v in {
          metadata = length(local.template_metadata) > 0 ? local.template_metadata : null
          spec     = local.template_spec
        } : k => v if v != null
      }
    }
  )
}
