# Kubernetes PodDisruptionBudget Terraform module.
#
# The selector block is always rendered (it is required in the spec): an
# empty selector block is the "all pods in the namespace" form —
# deliberately explicit, because a policy/v1 budget with a NULL selector
# matches no pods at all.

resource "kubernetes_pod_disruption_budget_v1" "pod_disruption_budget" {
  metadata {
    name        = var.spec.name
    namespace   = local.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  spec {
    # min_available/max_unavailable are IntOrString upstream: a numeric
    # string is an absolute count, a "%"-suffixed one a percentage. The spec
    # enforces exactly one is set.
    min_available   = var.spec.min_available != "" ? var.spec.min_available : null
    max_unavailable = var.spec.max_unavailable != "" ? var.spec.max_unavailable : null

    selector {
      match_labels = length(try(var.spec.selector.match_labels, {})) > 0 ? var.spec.selector.match_labels : null
      dynamic "match_expressions" {
        for_each = try(var.spec.selector.match_expressions, [])
        content {
          key      = match_expressions.value.key
          operator = match_expressions.value.operator
          values   = length(match_expressions.value.values) > 0 ? match_expressions.value.values : null
        }
      }
    }
  }

  lifecycle {
    # PARITY-EXCEPTION: the kubernetes provider's PDB resource cannot express
    # spec.unhealthyPodEvictionPolicy; the Pulumi module always sends it
    # (with the server default IfHealthyBudget). The server default is what a
    # non-default value would override, so failing loudly beats silently
    # deploying the default where always_allow was requested.
    precondition {
      condition     = local.unhealthy_pod_eviction_policy == "IfHealthyBudget"
      error_message = "The terraform kubernetes provider cannot express unhealthy_pod_eviction_policy: always_allow. Deploy this budget with the pulumi provisioner, or drop the field (the cluster then applies the default if_healthy_budget behavior)."
    }
  }
}
