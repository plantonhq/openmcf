# The optional installation namespace. Created before the release;
# deleted with the resource. Pulumi twin: namespace.go.

resource "kubernetes_namespace_v1" "openfga" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}
