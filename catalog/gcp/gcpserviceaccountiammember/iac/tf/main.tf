# A single ADDITIVE IAM grant ON a service account: one role, to one member,
# on one service account resource. Additive means the provider merges this
# (role, member) pair into the account's IAM policy without touching any other
# member's bindings on the same role, and destroy subtracts only this exact
# pair — grants made by other charts, teams, or tools are never clobbered.
#
# A service account is both an identity and a resource; this grant is on the
# RESOURCE side — who may use or manage the account itself (impersonation via
# roles/iam.workloadIdentityUser or roles/iam.serviceAccountTokenCreator,
# deploy-as via roles/iam.serviceAccountUser). Account-scoped grants beat
# their project-level equivalents on least privilege: a project-level
# serviceAccountUser grant allows acting as EVERY account in the project.
#
# Every argument is immutable (ForceNew): IAM grants have no update — any
# change replaces the grant atomically, which is also how the API behaves.
resource "google_service_account_iam_member" "this" {
  # The account's project is embedded in the fully-qualified resource name,
  # so there is no project argument.
  service_account_id = var.spec.service_account_id
  role               = var.spec.role
  member             = var.spec.member

  # An IAM Condition is part of the grant's identity: the same role granted
  # with and without a condition are two independent bindings in the policy.
  dynamic "condition" {
    for_each = var.spec.condition != null ? [var.spec.condition] : []
    content {
      title       = condition.value.title
      expression  = condition.value.expression
      description = condition.value.description
    }
  }
}
