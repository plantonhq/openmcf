# Kubernetes RBAC

Deploys one Kubernetes RBAC grant as a single resource: a scope (one namespace or the whole cluster), a role (custom policy rules or an existing built-in role by name), and the subjects that receive it (ServiceAccounts, users, or groups). The module derives the right object pair automatically -- Role + RoleBinding at namespace scope, ClusterRole + ClusterRoleBinding at cluster scope. Manages authorization declaratively through a Kubernetes Provider Connection with full audit trail and versioning.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Role or ClusterRole** -- created when the grant defines a role from policy rules (or an aggregation rule at cluster scope); skipped when binding an existing role by name
- **RoleBinding or ClusterRoleBinding** -- created when subjects are listed; a grant with a defined role and no subjects creates the role only (a permission set published for later grants)
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. The connection's own identity must hold the permissions being granted (Kubernetes forbids privilege escalation through RBAC).
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- For namespace-scoped grants, the target namespace must already exist. Use the Kubernetes Namespace component to manage namespaces declaratively.
- When binding an existing role, that Role or ClusterRole must already exist (the built-ins -- view, edit, admin, cluster-admin -- always do).

## Deploy

### Console

Open the deployment store, find **RBAC on Kubernetes**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Namespace App Reader** preset for a least-privilege app grant or **Grant Builtin View** to bind the built-in view role in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesRbac
metadata:
  name: checkout-reader
  org: acme-corp
  env: prod
spec:
  namespaceScope:
    namespace:
      value: backend-services
  createRole:
    rules:
      - verbs: ["get", "list", "watch"]
        apiGroups: [""]
        resources: ["configmaps", "secrets"]
  subjects:
    - serviceAccount:
        name:
          value: checkout-identity
```

```shell
planton apply -f rbac.yaml
```

This creates a Role permitting read access to ConfigMaps and Secrets in `backend-services`, bound to the `checkout-identity` ServiceAccount through a RoleBinding.

## Key Configuration

These are the most important decisions when configuring a Kubernetes RBAC grant. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Scope is the blast radius** -- Namespace scope confines every permission to one namespace (the least-privilege default). Cluster scope reaches every namespace plus cluster-wide objects (nodes, CRDs, namespaces themselves) -- reserve it for operators and platform tooling.

**Define or bind** -- `createRole` defines a new role from policy rules (verbs x resources per rule; `resources` and `nonResourceUrls` are mutually exclusive per rule, and non-resource URLs need cluster scope). `existingRole` binds a role that already exists -- usually a built-in (`view`, `edit`, `admin`, `cluster-admin`; set `isClusterRole: true` for those).

**Subjects complete the grant** -- ServiceAccounts (reference the KubernetesServiceAccount for a graph edge), users, or groups from the cluster's authentication layer. At cluster scope, every ServiceAccount subject must name its namespace. Binding an existing role requires at least one subject; a defined role without subjects publishes the role only.

**Aggregation (cluster scope, advanced)** -- Instead of (or alongside) explicit rules, label selectors combine the rules of every matching ClusterRole -- how Kubernetes' own view role grows with CRDs. `In`/`NotIn` selector expressions require values; `Exists`/`DoesNotExist` forbid them.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespaceScope.namespace` | KubernetesNamespace (`spec.name`) | The namespace a namespace-scoped grant confines permissions to |
| `spec.subjects[].serviceAccount.name` | KubernetesServiceAccount (`spec.name`) | The workload identity receiving the permissions |
| `spec.subjects[].serviceAccount.namespace` | KubernetesNamespace (`spec.name`) | The subject identity's namespace (required at cluster scope) |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `role_name` | The name of the created or bound role | Audit tooling and cross-grant references |
| `role_kind` | `Role` or `ClusterRole` | Conditional downstream configuration |
| `binding_name` | The binding's name; empty for role-only grants | Operational verification |
| `binding_kind` | `RoleBinding` or `ClusterRoleBinding`; empty when no binding exists | Operational verification |
| `namespace` | The grant's namespace; empty for cluster-scoped grants | Compliance reporting |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Least-privilege app grant** -- A custom namespace-scoped role granting exactly the verbs the workload calls, bound to its ServiceAccount. Start from the **Namespace App Reader** preset.

**Human read access** -- The built-in `view` ClusterRole bound at namespace scope to a group from your identity provider. Start from the **Grant Builtin View** preset.

**Platform operator** -- A cluster-scoped custom role for tooling that watches resources across namespaces. Start from the **Cluster Operator** preset.

**Aggregated ClusterRole** -- A label-selector role that grows as matching ClusterRoles are added. Start from the **Aggregated ClusterRole** preset.

## Works With

- **Kubernetes ServiceAccount** -- the usual subject: identity there, permissions here.
- **Kubernetes Namespace** -- namespace-scoped grants reference their namespace so charts deploy in dependency order.
- **Kubernetes Deployment and the other workload kinds** -- run as the granted ServiceAccount to exercise the permissions.
