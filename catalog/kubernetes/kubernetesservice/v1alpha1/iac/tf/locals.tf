# Computed values shared across the module.
#
# The enum maps translate proto enum value names (what the manifest carries)
# into the exact Kubernetes API strings — resolved once here so the resource
# and the outputs agree on wire values, mirroring the Pulumi module's
# resolve* helpers.

locals {
  # Resource-identity labels: the kuberneteslabelkeys set, identical to what
  # the Pulumi module stamps for the same manifest. User labels merge in
  # first so they can never override the identity keys.
  identity_labels_base = {
    "planton.ai/resource" = "true"
    "planton.ai/name"     = var.metadata.name
    "planton.ai/kind"     = "KubernetesService"
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

  service_type_map = {
    "cluster_ip"    = "ClusterIP"
    "node_port"     = "NodePort"
    "load_balancer" = "LoadBalancer"
    "external_name" = "ExternalName"
  }
  service_type = lookup(local.service_type_map, try(var.spec.type, "cluster_ip"), "ClusterIP")

  is_external_name = local.service_type == "ExternalName"
  is_load_balancer = local.service_type == "LoadBalancer"
  # Types with externally-facing addresses, where external_traffic_policy applies.
  is_externally_reachable = local.service_type == "NodePort" || local.is_load_balancer

  external_traffic_policy_map = {
    "cluster" = "Cluster"
    "local"   = "Local"
  }
  external_traffic_policy = lookup(local.external_traffic_policy_map, try(var.spec.external_traffic_policy, "cluster"), "Cluster")

  internal_traffic_policy_map = {
    "internal_cluster" = "Cluster"
    "internal_local"   = "Local"
  }
  # null when unset — the field is only sent when the user chose a policy.
  internal_traffic_policy = (
    try(var.spec.internal_traffic_policy, null) == null
    ? null
    : lookup(local.internal_traffic_policy_map, var.spec.internal_traffic_policy, "Cluster")
  )

  session_affinity_map = {
    "none"      = "None"
    "client_ip" = "ClientIP"
  }
  session_affinity = lookup(local.session_affinity_map, try(var.spec.session_affinity, "none"), "None")

  ip_family_map = {
    "ipv4" = "IPv4"
    "ipv6" = "IPv6"
  }
  ip_families = [for f in try(var.spec.ip_families, []) : lookup(local.ip_family_map, f, "IPv4")]

  ip_family_policy_map = {
    "single_stack"       = "SingleStack"
    "prefer_dual_stack"  = "PreferDualStack"
    "require_dual_stack" = "RequireDualStack"
  }
  ip_family_policy = (
    try(var.spec.ip_family_policy, null) == null
    ? null
    : lookup(local.ip_family_policy_map, var.spec.ip_family_policy, null)
  )

  # In-cluster DNS name — resolves for every type; for ExternalName it is the
  # CNAME alias the service exists to provide.
  kube_endpoint = "${var.spec.name}.${local.namespace}.svc.cluster.local"

  # Port-forward needs pods behind the service; an ExternalName alias has none.
  port_forward_command = (
    local.is_external_name || length(try(var.spec.ports, [])) == 0
    ? ""
    : "kubectl port-forward -n ${local.namespace} service/${var.spec.name} ${var.spec.ports[0].port}:${var.spec.ports[0].port}"
  )
}
