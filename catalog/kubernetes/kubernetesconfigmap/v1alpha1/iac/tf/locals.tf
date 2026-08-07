# Local values and computed configuration

locals {
  # Build combined labels
  standard_labels = {
    "managed-by"    = "planton"
    "resource"      = var.metadata.name
    "resource-kind" = "KubernetesConfigMap"
  }

  labels = merge(local.standard_labels, try(var.spec.labels, {}))

  # Build annotations
  annotations = try(var.spec.annotations, {})

  # Fall back to the cluster's "default" namespace when the field arrives
  # null or empty — the same behavior as kubectl without a namespace flag.
  namespace = (
    try(var.spec.namespace, null) == null || try(var.spec.namespace, "") == ""
    ? "default"
    : var.spec.namespace
  )

  # UTF-8 entries pass through as-is; binary entries are already
  # base64-encoded strings (the wire form Kubernetes stores for binaryData)
  # and also pass through unchanged.
  data        = try(var.spec.data, {})
  binary_data = try(var.spec.binary_data, {})
}
