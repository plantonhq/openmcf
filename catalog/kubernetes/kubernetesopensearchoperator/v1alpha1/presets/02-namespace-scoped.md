# Namespace scoped preset

An operator fenced to a single namespace with namespace-scoped RBAC:
`watch_namespace` limits what it reconciles, `use_role_bindings`
swaps ClusterRoleBindings for RoleBindings in the watched namespace.
This is the posture for shared clusters — a team that owns one
namespace can run its own operator without cluster-wide permissions,
and multiple installs coexist because each watches only its own slice.

Do not use this when you administer the whole cluster and want one
operator serving every team: the standard preset's cluster-wide watch
is simpler and avoids per-namespace operator sprawl. Also note the
watch fence is per-install, not per-cluster — a KubernetesOpenSearch
declared outside the watched namespace is silently ignored, which
looks like a hang if you forget the fence exists.

Change `watch_namespace` (and the install namespace — keeping them
the same is the simplest mental model) to the namespace your search
clusters will live in.

See [02-namespace-scoped.yaml](./02-namespace-scoped.yaml) for the
manifest.
