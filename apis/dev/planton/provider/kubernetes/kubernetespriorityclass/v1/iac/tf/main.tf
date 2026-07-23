# Kubernetes PriorityClass Terraform module.
#
# value is immutable upstream: the provider forces replacement on change,
# matching the Pulumi module's DeleteBeforeReplace semantics (PriorityClass
# names are cluster-unique, so delete-then-create is the only safe order).

resource "kubernetes_priority_class_v1" "priority_class" {
  metadata {
    name        = var.spec.name
    labels      = local.labels
    annotations = local.annotations
  }

  value = var.spec.value

  # The Kubernetes default is false; sending it explicitly keeps both
  # engines' submitted objects identical.
  global_default = var.spec.global_default

  description = var.spec.description

  # Always sent explicitly (with the API server's default applied in locals)
  # so both engines submit byte-identical objects.
  preemption_policy = local.preemption_policy
}
