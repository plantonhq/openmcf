# Computed values shared across the module.

locals {
  # Resource-identity labels: the kuberneteslabelkeys set, identical to what
  # the Pulumi module stamps for the same manifest. User labels merge in
  # first so they can never override the identity keys.
  identity_labels_base = {
    "planton.ai/resource" = "true"
    "planton.ai/name"     = var.metadata.name
    "planton.ai/kind"     = "KubernetesHorizontalPodAutoscaler"
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

  # The resolved scale target with the spec's defaults (apps/v1 Deployment)
  # applied, identical to the Pulumi module.
  target_api_version = try(var.spec.scale_target.api_version, "") != "" ? var.spec.scale_target.api_version : "apps/v1"
  target_kind        = try(var.spec.scale_target.kind, "") != "" ? var.spec.scale_target.kind : "Deployment"

  # The replica floor with the Kubernetes default (1) applied — sent
  # explicitly by both engines.
  min_replicas = try(var.spec.min_replicas, null) != null ? var.spec.min_replicas : 1

  # Enum value names -> Kubernetes API strings (identical to the Pulumi
  # module's mappings).
  metric_type_map = {
    "resource"           = "Resource"
    "container_resource" = "ContainerResource"
    "pods"               = "Pods"
    "object"             = "Object"
    "external"           = "External"
  }

  target_type_map = {
    "utilization"   = "Utilization"
    "raw_value"     = "Value"
    "average_value" = "AverageValue"
  }

  select_policy_map = {
    "max_change" = "Max"
    "min_change" = "Min"
    "disabled"   = "Disabled"
  }

  policy_type_map = {
    "pods"    = "Pods"
    "percent" = "Percent"
  }

  # Rendered for the outputs contract ("Kind/name").
  scale_target_string = "${local.target_kind}/${var.spec.scale_target.name}"
}
