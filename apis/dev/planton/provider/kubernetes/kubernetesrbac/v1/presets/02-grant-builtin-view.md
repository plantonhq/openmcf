# Grant Built-in `view` to a Group

This preset binds the Kubernetes built-in `view` ClusterRole to a group, confined to one namespace. No role is created — only a RoleBinding — giving everyone in the group read-only access to most namespaced objects with zero rules to maintain.

## When to Use

- Team read access to a namespace: developers inspecting a production namespace, support staff browsing workloads, auditors reviewing state
- Any human-access need the built-ins already cover — prefer this over hand-written reader roles, because `view` is aggregated and automatically absorbs read permissions for CRDs that operators register

## Key Configuration Choices

- **`existingRole: { name: view, isClusterRole: true }`** — binds a role that already exists instead of creating one. `isClusterRole: true` is essential: `view` is a ClusterRole, and a namespace-scoped RoleBinding pointing at a ClusterRole is exactly how built-ins are granted per-namespace (the binding confines the cluster-wide definition to this one namespace)
- **`view` excludes Secrets and RBAC objects** — reading a Secret is effectively holding the credential, so `view` deliberately omits it; grant Secret access separately and narrowly if truly needed. Sibling built-ins: `edit` (read/write, no RBAC) and `admin` (edit plus namespace-local RBAC management)
- **Group subject** — a plain string matched against what the cluster's authenticator asserts (an OIDC groups claim, a cloud IAM mapping). Kubernetes has no Group objects; membership is managed in the identity provider, which is exactly why binding groups beats binding individual users
- **Subjects are required with `existingRole`** — binding an existing role to nobody would deploy nothing; the schema rejects it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Namespace the read access is confined to | Your namespace management |
| `<your-group>` | Group name as asserted by your cluster's authenticator, e.g. an OIDC groups claim like `dev-team` | Your identity provider / cluster OIDC configuration |

## Related Presets

- **01-namespace-app-reader** — a custom, narrower reader role for a workload identity
- **03-cluster-operator** — cluster-scoped permissions including nodes and `/metrics`
- **04-aggregated-clusterrole** — a label-composed ClusterRole published without a binding
