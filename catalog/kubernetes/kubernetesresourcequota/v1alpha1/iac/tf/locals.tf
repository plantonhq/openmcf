# Computed values shared across the module.

locals {
  # Resource-identity labels: the kuberneteslabelkeys set, identical to what
  # the Pulumi module stamps for the same manifest. User labels merge in
  # first so they can never override the identity keys.
  identity_labels_base = {
    "planton.ai/resource" = "true"
    "planton.ai/name"     = var.metadata.name
    "planton.ai/kind"     = "KubernetesResourceQuota"
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

  annotations = try(var.spec.annotations, {})

  # Fall back to the cluster's "default" namespace when the field arrives
  # null or empty — the same behavior as kubectl without a namespace flag.
  namespace = (
    try(var.spec.namespace, null) == null || try(var.spec.namespace, "") == ""
    ? "default"
    : var.spec.namespace
  )

  # Enum value names -> Kubernetes API strings (identical to the Pulumi
  # module's mapping).
  scope_map = {
    "terminating"                  = "Terminating"
    "not_terminating"              = "NotTerminating"
    "best_effort"                  = "BestEffort"
    "not_best_effort"              = "NotBestEffort"
    "priority_class"               = "PriorityClass"
    "cross_namespace_pod_affinity" = "CrossNamespacePodAffinity"
    "volume_attributes_class"      = "VolumeAttributesClass"
  }
  scopes = [for s in try(var.spec.scopes, []) : lookup(local.scope_map, s, s)]

  limit_type_map = {
    "container"               = "Container"
    "pod"                     = "Pod"
    "persistent_volume_claim" = "PersistentVolumeClaim"
  }

  # The companion LimitRange exists only when limit defaults are configured;
  # it shares the quota's name — one governance pair, one identity.
  create_limit_range = length(try(var.spec.limit_defaults, [])) > 0
  limit_range_name   = local.create_limit_range ? var.spec.name : ""
}
