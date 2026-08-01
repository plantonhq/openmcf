# The optional namespace. Created before the CR (the CR and every
# operator-created child is namespaced); deleted with the resource.
# Pre-existing-namespace deployments leave create_namespace false.
#
# NOTE the operator watch contract: the namespace must be inside the
# KubernetesKubeRayOperator's watch scope (cluster-wide with the
# operator's defaults) before the CR reconciles. And the GCS
# fault-tolerance credential Secret rides a secretKeyRef, readable only
# from this same namespace — co-locate the cluster with its Valkey or
# replicate the Secret.
resource "kubernetes_namespace_v1" "namespace" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}
