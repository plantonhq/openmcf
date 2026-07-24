# Namespace-scoped preset

The tenancy posture: an operator that watches exactly one namespace
and holds only namespace-scoped Roles instead of cluster-wide RBAC.
On a shared cluster this is what lets a team run their own ClickHouse
engine without a cluster administrator granting — or trusting — 
anything beyond their namespace. Each team that needs ClickHouse
installs its own copy of this preset in its own namespace.

The two settings travel together by design: a wide watch list with
namespace-scoped RBAC would let the operator SEE clusters it cannot
manage, which presents as ClickHouse resources sitting silently
unreconciled. Keep `watch_namespaces` matching where the RBAC
actually lives — here, the install namespace itself.

One caveat worth knowing before choosing this shape: the four CRDs
are cluster-scoped no matter what (all Kubernetes CRDs are), so the
FIRST install on a cluster still needs permissions to create them,
and CRD upgrades ride whichever install's hook runs first. The
per-team isolation is about workloads and RBAC, not the API types.

The first thing to change is the namespace pair and the password.

See [02-namespace-scoped.yaml](./02-namespace-scoped.yaml) for the
manifest.
