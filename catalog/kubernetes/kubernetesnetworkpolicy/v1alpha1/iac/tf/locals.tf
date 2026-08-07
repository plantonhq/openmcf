# Computed values shared across the module.

locals {
  # Resource-identity labels: the kuberneteslabelkeys set, identical to what
  # the Pulumi module stamps for the same manifest. User labels merge in
  # first so they can never override the identity keys.
  identity_labels_base = {
    "planton.ai/resource" = "true"
    "planton.ai/name"     = var.metadata.name
    "planton.ai/kind"     = "KubernetesNetworkPolicy"
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

  # The governed directions as Kubernetes API strings. ALWAYS sent explicitly
  # (mirroring the Pulumi module): when the spec omits policy_types, the API
  # server's own inference is applied here — ingress always, egress only when
  # egress rules exist — so both engines submit byte-identical direction sets
  # and the deployed object never depends on which engine applied it.
  policy_type_map = {
    "ingress" = "Ingress"
    "egress"  = "Egress"
  }
  policy_types = (
    length(try(var.spec.policy_types, [])) > 0
    ? [for t in var.spec.policy_types : lookup(local.policy_type_map, t, "Ingress")]
    : concat(["Ingress"], length(try(var.spec.egress_rules, [])) > 0 ? ["Egress"] : [])
  )

  # Rendered for the outputs contract ("Ingress", "Egress", or "Ingress,Egress").
  policy_types_string = join(",", local.policy_types)
}
