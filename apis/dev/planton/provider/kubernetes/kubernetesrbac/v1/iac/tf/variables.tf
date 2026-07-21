# Input variables for Kubernetes RBAC Terraform module.
# The spec mirrors spec.proto with every StringValueOrRef flattened to a plain
# string (references are resolved before the module is invoked).

variable "metadata" {
  description = "Metadata for the RBAC grant resource"
  type = object({
    name = string
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes RBAC grant"
  type = object({
    # Scope: exactly one of the following. namespace_scope creates Role/RoleBinding
    # objects; cluster_scope creates ClusterRole/ClusterRoleBinding objects.
    namespace_scope = optional(object({
      # The namespace the grant applies to; omitted means the cluster's "default".
      namespace = optional(string, "default")
    }))
    # Intentionally empty object: cluster scope carries no parameters.
    cluster_scope = optional(object({}))

    # Role source: exactly one of the following. create_role defines a new
    # Role/ClusterRole; existing_role binds to one already in the cluster
    # (e.g. the built-in "view"/"edit"/"admin"/"cluster-admin" ClusterRoles).
    create_role = optional(object({
      # Name of the created role; defaults to the component's metadata.name.
      name = optional(string)

      # Permission rules. Each rule independently grants verbs over resources or
      # non-resource URLs; there is no ordering and no deny semantics. May be
      # empty only when aggregation_rule is set.
      rules = optional(list(object({
        verbs             = list(string)
        api_groups        = optional(list(string), [])
        resources         = optional(list(string), [])
        resource_names    = optional(list(string), [])
        non_resource_urls = optional(list(string), [])
      })), [])

      # ClusterRole aggregation (cluster scope only): the controller composes the
      # role's rules from every ClusterRole matching ANY of these selectors.
      aggregation_rule = optional(object({
        cluster_role_selectors = list(object({
          match_labels = optional(map(string), {})
          match_expressions = optional(list(object({
            key      = string
            operator = string
            values   = optional(list(string), [])
          })), [])
        }))
      }))
    }))

    existing_role = optional(object({
      name = string
      # Only meaningful in namespace scope: a RoleBinding may reference either a
      # namespaced Role (false) or a ClusterRole (true). Ignored in cluster scope,
      # where the reference is always a ClusterRole.
      is_cluster_role = optional(bool, false)
    }))

    # Who receives the permissions. Each entry sets exactly one of
    # service_account / user / group. Empty means role-definition-only grant.
    subjects = optional(list(object({
      service_account = optional(object({
        name = string
        # Defaults to the grant's namespace in namespace scope; required by spec
        # validation in cluster scope (a ServiceAccount always has a namespace).
        namespace = optional(string)
      }))
      user  = optional(string)
      group = optional(string)
    })), [])

    # Extra labels/annotations applied to every created RBAC object.
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})
  })
}
