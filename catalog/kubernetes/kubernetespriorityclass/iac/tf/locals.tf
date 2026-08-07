# Computed values shared across the module.

locals {
  # Resource-identity labels: the kuberneteslabelkeys set, identical to what
  # the Pulumi module stamps for the same manifest. User labels merge in
  # first so they can never override the identity keys.
  identity_labels_base = {
    "planton.ai/resource" = "true"
    "planton.ai/name"     = var.metadata.name
    "planton.ai/kind"     = "KubernetesPriorityClass"
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

  # The API string with the server default (PreemptLowerPriority) applied,
  # identical to the Pulumi module.
  preemption_policy_map = {
    "preempt_lower_priority" = "PreemptLowerPriority"
    "never"                  = "Never"
  }
  preemption_policy = lookup(local.preemption_policy_map, try(var.spec.preemption_policy, "preempt_lower_priority"), "PreemptLowerPriority")
}
