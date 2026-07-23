# Kubernetes PersistentVolumeClaim Terraform module.
#
# wait_until_bound is FALSE deliberately (mirroring the Pulumi module's
# skipAwait annotation): under a WaitForFirstConsumer StorageClass a claim is
# correctly Pending until a pod consumes it, and the provider's default wait
# would hang every such deploy. The attribute exists only in configuration —
# it is declared config-only in the provider import catalog.

resource "kubernetes_persistent_volume_claim_v1" "persistent_volume_claim" {
  metadata {
    name        = var.spec.name
    namespace   = local.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  spec {
    access_modes = local.access_modes

    resources {
      requests = {
        storage = var.spec.storage_request
      }
      limits = var.spec.storage_limit != "" ? { storage = var.spec.storage_limit } : null
    }

    # null = absent (cluster default class); "" = explicitly no dynamic
    # provisioning. The distinction is load-bearing upstream (see locals).
    storage_class_name = local.storage_class_name

    # Filesystem is the server default, but sending it explicitly keeps
    # both engines' submitted objects identical.
    volume_mode = local.volume_mode

    volume_name = var.spec.volume_name != "" ? var.spec.volume_name : null

    dynamic "selector" {
      for_each = try(var.spec.selector, null) != null ? [var.spec.selector] : []
      content {
        match_labels = length(try(selector.value.match_labels, {})) > 0 ? selector.value.match_labels : null
        dynamic "match_expressions" {
          for_each = try(selector.value.match_expressions, [])
          content {
            key      = match_expressions.value.key
            operator = match_expressions.value.operator
            values   = length(match_expressions.value.values) > 0 ? match_expressions.value.values : null
          }
        }
      }
    }
  }

  wait_until_bound = false

  lifecycle {
    # PARITY-EXCEPTION: the kubernetes provider's PVC resource cannot express
    # spec.dataSource/dataSourceRef (clone a PVC / restore a VolumeSnapshot);
    # the Pulumi module sends them. Failing the plan loudly beats silently
    # provisioning an EMPTY volume where the user asked for restored data.
    precondition {
      condition     = try(var.spec.data_source, null) == null
      error_message = "The terraform kubernetes provider cannot express a PVC data source (clone/snapshot-restore). Deploy this claim with the pulumi provisioner, or drop spec.data_source."
    }
  }
}
