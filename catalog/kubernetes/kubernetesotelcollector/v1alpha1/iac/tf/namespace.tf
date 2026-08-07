# The optional deployment namespace. Created before the CR; deleted with
# the resource (pre-existing-namespace deploys leave create_namespace
# false).
resource "kubernetes_namespace_v1" "namespace" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}
