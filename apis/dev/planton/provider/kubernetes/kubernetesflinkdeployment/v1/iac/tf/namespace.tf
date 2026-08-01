# The optional namespace. Created before the CR (the CR and every
# operator-created child is namespaced); deleted with the resource.
# Pre-existing-namespace deployments leave create_namespace false.
#
# NOTE the operator watch contract: the KubernetesFlinkOperator's watch
# scope must cover this namespace — a freshly created namespace needs an
# operator watching it before the CR reconciles.
resource "kubernetes_namespace_v1" "namespace" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}
