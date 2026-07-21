# Kubernetes ServiceAccount Terraform Module
# Creates a Kubernetes ServiceAccount with image-pull secrets, the token-automount
# setting, and cloud workload-identity annotations.

resource "kubernetes_service_account_v1" "service_account" {
  metadata {
    name        = var.spec.name
    namespace   = var.spec.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  # Name-only references to dockerconfigjson secrets in the same namespace.
  dynamic "image_pull_secret" {
    for_each = var.spec.image_pull_secrets
    content {
      name = image_pull_secret.value
    }
  }

  # Tri-state automount: the Kubernetes API distinguishes unset (cluster default
  # applies) from an explicit true/false, but this provider's schema defaults the
  # attribute to true and cannot leave it off the object. That is behaviorally
  # equivalent: the Kubernetes cluster default for an absent field IS "mount the
  # token" (true), so applying true when the variable is null produces the exact
  # ServiceAccount behavior the Pulumi module gets by omitting the field. Not a
  # parity exception — observable behavior matches across both engines.
  automount_service_account_token = coalesce(var.spec.automount_service_account_token, true)
}
