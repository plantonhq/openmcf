# Kubernetes RBAC Terraform Module
#
# Deploys one RBAC grant: a role (created or existing) plus, when subjects are
# present, a binding pointing every subject at that role. Count guards select
# which of the four Kubernetes RBAC object kinds are actually created:
#
#   namespace_scope + create_role  -> kubernetes_role_v1
#   cluster_scope   + create_role  -> kubernetes_cluster_role_v1
#   namespace_scope + subjects     -> kubernetes_role_binding_v1
#   cluster_scope   + subjects     -> kubernetes_cluster_role_binding_v1
#
# When existing_role is used, no role resource is created — the binding
# references the existing role by name through its immutable role_ref.

resource "kubernetes_role_v1" "role" {
  count = local.is_namespace_scoped && local.has_create_role ? 1 : 0

  metadata {
    name        = local.role_name
    namespace   = local.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  # Each rule independently grants verbs over resources; there is no ordering and
  # no deny semantics. non_resource_urls never appear here: they cannot be
  # namespaced, and spec validation rejects them outside cluster scope.
  dynamic "rule" {
    for_each = local.rules
    content {
      verbs          = rule.value.verbs
      api_groups     = length(rule.value.api_groups) > 0 ? rule.value.api_groups : null
      resources      = length(rule.value.resources) > 0 ? rule.value.resources : null
      resource_names = length(rule.value.resource_names) > 0 ? rule.value.resource_names : null
    }
  }
}

resource "kubernetes_cluster_role_v1" "cluster_role" {
  count = !local.is_namespace_scoped && local.has_create_role ? 1 : 0

  metadata {
    name        = local.role_name
    labels      = local.labels
    annotations = local.annotations
  }

  dynamic "rule" {
    for_each = local.rules
    content {
      verbs             = rule.value.verbs
      api_groups        = length(rule.value.api_groups) > 0 ? rule.value.api_groups : null
      resources         = length(rule.value.resources) > 0 ? rule.value.resources : null
      resource_names    = length(rule.value.resource_names) > 0 ? rule.value.resource_names : null
      non_resource_urls = length(rule.value.non_resource_urls) > 0 ? rule.value.non_resource_urls : null
    }
  }

  # ClusterRole aggregation: the controller continuously composes this role's
  # rules from every ClusterRole matching ANY selector; directly listed rules
  # become controller-managed. Only expressible on ClusterRoles — the Kubernetes
  # Role type has no aggregationRule field.
  dynamic "aggregation_rule" {
    for_each = local.aggregation_rule != null ? [local.aggregation_rule] : []
    content {
      dynamic "cluster_role_selectors" {
        for_each = aggregation_rule.value.cluster_role_selectors
        content {
          match_labels = cluster_role_selectors.value.match_labels

          dynamic "match_expressions" {
            for_each = cluster_role_selectors.value.match_expressions
            content {
              key      = match_expressions.value.key
              operator = match_expressions.value.operator
              values   = match_expressions.value.values
            }
          }
        }
      }
    }
  }
}

# PARITY-EXCEPTION: this provider's subject schema defaults namespace to
# "default" for every subject kind, so User/Group subjects carry namespace
# "default" in the stored object while the Pulumi module omits it. The RBAC
# authorizer ignores namespace on User/Group subjects, so authorization behavior
# is identical across both engines.
resource "kubernetes_role_binding_v1" "binding" {
  count = local.is_namespace_scoped && local.has_binding ? 1 : 0

  metadata {
    name        = local.binding_name
    namespace   = local.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  # roleRef is immutable in Kubernetes and always addresses the role through the
  # rbac.authorization.k8s.io API group. Its kind is the grant's resolved role
  # kind: the created Role, or for existing roles, Role-or-ClusterRole per
  # is_cluster_role.
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = local.role_kind
    name      = local.role_name
  }

  dynamic "subject" {
    for_each = local.subjects
    content {
      kind      = subject.value.kind
      name      = subject.value.name
      namespace = subject.value.namespace
      api_group = subject.value.api_group
    }
  }

  # role_ref references the role by name only, so the dependency on a role
  # created in the same run must be explicit for correct apply ordering.
  depends_on = [kubernetes_role_v1.role, kubernetes_cluster_role_v1.cluster_role]
}

resource "kubernetes_cluster_role_binding_v1" "cluster_binding" {
  count = !local.is_namespace_scoped && local.has_binding ? 1 : 0

  metadata {
    name        = local.binding_name
    labels      = local.labels
    annotations = local.annotations
  }

  # In cluster scope the reference is always a ClusterRole (local.role_kind).
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = local.role_kind
    name      = local.role_name
  }

  dynamic "subject" {
    for_each = local.subjects
    content {
      kind      = subject.value.kind
      name      = subject.value.name
      namespace = subject.value.namespace
      api_group = subject.value.api_group
    }
  }

  depends_on = [kubernetes_cluster_role_v1.cluster_role]
}
