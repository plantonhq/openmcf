# Computed values shared across the module.

locals {
  # Resource-identity labels: the kuberneteslabelkeys set, identical to what
  # the Pulumi module stamps for the same manifest. User labels merge in
  # first so they can never override the identity keys.
  identity_labels_base = {
    "planton.ai/resource" = "true"
    "planton.ai/name"     = var.metadata.name
    "planton.ai/kind"     = "KubernetesStorageClass"
  }

  id_label = (
    var.metadata.id != null && try(var.metadata.id, "") != ""
  ) ? { "planton.ai/id" = var.metadata.id } : {}

  org_label = (
    var.metadata.org != null && try(var.metadata.org, "") != ""
  ) ? { "planton.ai/organization" = var.metadata.org } : {}

  env_label = (
    var.metadata.env != null && try(var.metadata.env, "") != ""
  ) ? { "planton.ai/environment" = var.metadata.env } : {}

  labels = merge(
    try(var.spec.labels, {}),
    local.identity_labels_base,
    local.id_label,
    local.org_label,
    local.env_label,
  )

  # The default-class marker is a first-class spec field; the annotation is
  # only the wire form (identical to the Pulumi module).
  annotations = merge(
    try(var.spec.annotations, {}),
    var.spec.is_default_class ? { "storageclass.kubernetes.io/is-default-class" = "true" } : {},
  )

  # The Kubernetes API strings for the enum-typed policies, resolved with
  # the API server's own defaults so both engines submit identical objects
  # whether or not the spec set the optional fields.
  reclaim_policy_map = {
    "delete" = "Delete"
    "retain" = "Retain"
  }
  reclaim_policy = lookup(local.reclaim_policy_map, try(var.spec.reclaim_policy, "delete"), "Delete")

  volume_binding_mode_map = {
    "immediate"               = "Immediate"
    "wait_for_first_consumer" = "WaitForFirstConsumer"
  }
  volume_binding_mode = lookup(local.volume_binding_mode_map, try(var.spec.volume_binding_mode, "immediate"), "Immediate")
}
