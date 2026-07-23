# Local values: the spec's three orthogonal choices (scope, role source, subjects)
# resolved into concrete Kubernetes object names and kinds, so main.tf contains no
# decision logic of its own. Mirrors module/locals.go in the Pulumi module.

locals {
  # Standard Planton labels merged with user labels; applied to every created object.
  standard_labels = {
    "managed-by"    = "planton"
    "resource"      = var.metadata.name
    "resource-kind" = "KubernetesRbac"
  }

  labels      = merge(local.standard_labels, var.spec.labels)
  annotations = var.spec.annotations

  # Scope: namespace_scope produces Role/RoleBinding, cluster_scope produces
  # ClusterRole/ClusterRoleBinding. Spec validation guarantees exactly one is set.
  is_namespace_scoped = var.spec.namespace_scope != null

  # An omitted namespace lands the grant in the cluster's "default" namespace
  # (the variable default handles that); cluster scope has no namespace at all.
  namespace = local.is_namespace_scoped ? try(var.spec.namespace_scope.namespace, "default") : ""

  # Role source: either we create the role, or we bind to an existing one.
  has_create_role = var.spec.create_role != null

  # The created role defaults its name to the component's own metadata.name so
  # simple grants need no explicit role name; an existing role brings its own.
  role_name = local.has_create_role ? coalesce(try(var.spec.create_role.name, null), var.metadata.name) : try(var.spec.existing_role.name, "")

  # The role's kind — also the kind used in the binding's roleRef, so the two can
  # never diverge. In namespace scope a binding may point at either a namespaced
  # Role or a ClusterRole (how built-in roles like "view" are granted
  # per-namespace, per existing_role.is_cluster_role); in cluster scope the role
  # is always a ClusterRole.
  role_kind = local.is_namespace_scoped ? (
    local.has_create_role ? "Role" : (try(var.spec.existing_role.is_cluster_role, false) ? "ClusterRole" : "Role")
  ) : "ClusterRole"

  # A binding exists only when there are subjects to bind. Subjects absent means
  # the grant only publishes a role definition for later bindings.
  has_binding  = length(var.spec.subjects) > 0
  binding_name = local.has_binding ? var.metadata.name : ""
  binding_kind = local.has_binding ? (local.is_namespace_scoped ? "RoleBinding" : "ClusterRoleBinding") : ""

  # Rules and aggregation of the created role (empty/null when binding to an
  # existing role).
  rules            = try(var.spec.create_role.rules, [])
  aggregation_rule = try(var.spec.create_role.aggregation_rule, null)

  # Subjects normalized to the Kubernetes rbac/v1 Subject shape:
  #   - service_account -> kind ServiceAccount with a namespace: the subject's own
  #     when set, otherwise the grant's namespace (namespace scope only — spec
  #     validation guarantees cluster-scoped grants set it explicitly, because a
  #     ServiceAccount always lives in some namespace).
  #   - user / group -> kind User/Group under the rbac.authorization.k8s.io API
  #     group. Plain strings matched against what the cluster's authenticator
  #     asserts; Kubernetes has no User or Group objects.
  subjects = [for s in var.spec.subjects :
    s.service_account != null ? {
      kind      = "ServiceAccount"
      name      = try(s.service_account.name, "")
      namespace = try(coalesce(s.service_account.namespace, local.namespace), local.namespace)
      api_group = null
    } : {
      kind      = s.user != null ? "User" : "Group"
      name      = s.user != null ? s.user : s.group
      namespace = null
      api_group = "rbac.authorization.k8s.io"
    }
  ]
}
