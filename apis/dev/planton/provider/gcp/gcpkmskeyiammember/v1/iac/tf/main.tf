# A single ADDITIVE IAM grant ON a KMS crypto key: one role, to one member,
# on one key. Additive means the provider merges this (role, member) pair into
# the key's IAM policy without touching any other member's bindings on the
# same role, and destroy subtracts only this exact pair — grants made by other
# charts, teams, or tools are never clobbered.
#
# The canonical use is CMEK: granting a consuming service's agent
# roles/cloudkms.cryptoKeyEncrypterDecrypter on exactly the key it needs.
# Key-scoped beats project-scoped on least privilege AND gives orchestration a
# real dependency edge — the encrypted resource can be ordered after the grant
# it requires, closing the first-deploy IAM-propagation race.
#
# Every argument is immutable (ForceNew): IAM grants have no update — any
# change replaces the grant atomically, which is also how the API behaves.
resource "google_kms_crypto_key_iam_member" "this" {
  # The key's project and location are embedded in its resource path, so
  # there is no project or location argument.
  crypto_key_id = var.spec.crypto_key_id
  role          = var.spec.role
  member        = var.spec.member

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
