# One ADDITIVE IAM grant on a Google Service Account: the workload-identity
# principal (a Kubernetes ServiceAccount in the cluster project's implicit
# <project>.svc.id.goog pool) receives roles/iam.workloadIdentityUser on the
# GSA — the GCP half of GKE Workload Identity. The Kubernetes half (the
# iam.gke.io/gcp-service-account annotation on the KSA) belongs to the
# workload's own deployment.
#
# Additive means this grant merges into the GSA's IAM policy without touching
# any other principal's bindings, and destroy subtracts only this exact
# grant. Every argument is immutable (ForceNew): IAM grants have no update —
# any change replaces the grant atomically.
resource "google_service_account_iam_member" "workload_identity_binding" {
  service_account_id = local.service_account_id
  role               = "roles/iam.workloadIdentityUser"
  member             = local.workload_identity_member

  # An IAM Condition is part of the grant's identity: the same grant with
  # and without a condition are two independent bindings in the policy.
  dynamic "condition" {
    for_each = var.spec.condition != null ? [var.spec.condition] : []
    content {
      title       = condition.value.title
      expression  = condition.value.expression
      description = condition.value.description
    }
  }
}
