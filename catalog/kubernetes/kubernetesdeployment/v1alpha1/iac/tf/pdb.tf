# PodDisruptionBudget guarding availability during voluntary disruptions.
# The PDB selects on the workload's SELECTOR labels — the same immutable set
# the Deployment selects pods with; matching the full governance label set
# would silently guard zero pods if a non-selector label ever drifts.
resource "kubernetes_pod_disruption_budget_v1" "this" {
  count = local.pdb_enabled ? 1 : 0

  metadata {
    name      = var.metadata.name
    namespace = local.namespace
    labels    = local.final_labels
  }

  spec {
    selector {
      match_labels = local.selector_labels
    }

    # Exactly one bound may be set (spec validation enforces it); when neither
    # is set, default to min_available 1 — at least one pod always survives.
    min_available = (
      try(var.spec.availability.pod_disruption_budget.min_available, "") != ""
      ? var.spec.availability.pod_disruption_budget.min_available
      : (try(var.spec.availability.pod_disruption_budget.max_unavailable, "") == "" ? "1" : null)
    )
    max_unavailable = (
      try(var.spec.availability.pod_disruption_budget.max_unavailable, "") != ""
      ? var.spec.availability.pod_disruption_budget.max_unavailable
      : null
    )
  }

  depends_on = [kubernetes_deployment_v1.this]
}
