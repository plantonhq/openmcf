# Kubernetes ConfigMap Terraform Module
# Creates a Kubernetes ConfigMap with UTF-8 data, base64 binary data, and immutability support

resource "kubernetes_config_map_v1" "configmap" {
  metadata {
    name        = var.spec.name
    namespace   = local.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  # data takes plain UTF-8 strings; binary_data takes base64-encoded strings
  # (Kubernetes stores binaryData as base64 on the wire, so values pass
  # through unchanged).
  data        = local.data
  binary_data = local.binary_data

  # When true, the cluster rejects any update to data/binary_data after
  # creation; changing content requires delete-and-recreate.
  immutable = var.spec.immutable
}
