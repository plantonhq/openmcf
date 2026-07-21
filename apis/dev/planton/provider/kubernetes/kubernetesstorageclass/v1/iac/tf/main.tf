# Kubernetes StorageClass Terraform module.
#
# provisioner and parameters are immutable upstream: the provider forces
# replacement on change, matching the Pulumi module's DeleteBeforeReplace
# semantics (StorageClass names are cluster-unique, so delete-then-create is
# the only safe order).

resource "kubernetes_storage_class_v1" "storage_class" {
  metadata {
    name        = var.spec.name
    labels      = local.labels
    annotations = local.annotations
  }

  storage_provisioner = var.spec.provisioner
  parameters          = length(try(var.spec.parameters, {})) > 0 ? var.spec.parameters : null

  # Always sent explicitly (with the API server's defaults applied in
  # locals) so both engines submit byte-identical objects.
  reclaim_policy      = local.reclaim_policy
  volume_binding_mode = local.volume_binding_mode

  allow_volume_expansion = var.spec.allow_volume_expansion

  mount_options = length(try(var.spec.mount_options, [])) > 0 ? var.spec.mount_options : null

  dynamic "allowed_topologies" {
    for_each = length(try(var.spec.allowed_topologies, [])) > 0 ? [1] : []
    content {
      dynamic "match_label_expressions" {
        for_each = var.spec.allowed_topologies[0].match_label_expressions
        content {
          key    = match_label_expressions.value.key
          values = match_label_expressions.value.values
        }
      }
    }
  }

  lifecycle {
    # PARITY-EXCEPTION: the kubernetes provider models allowed_topologies as
    # a SINGLE selector term (max_items = 1), while the API (and the Pulumi
    # module) accept multiple OR'd terms. One term passes through intact;
    # failing the plan loudly on several beats silently dropping the rest.
    precondition {
      condition     = length(try(var.spec.allowed_topologies, [])) <= 1
      error_message = "The terraform kubernetes provider supports at most ONE allowed_topologies term. Deploy this class with the pulumi provisioner, or combine the zone values into a single term (values within one requirement already OR together)."
    }
  }
}
